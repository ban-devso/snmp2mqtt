package snmp

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ban-devso/snmp-mqtt/config"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gosnmp/gosnmp"
)

// Init contains the generic read/publish loop
func Init() {
	opts := mqtt.NewClientOptions().AddBroker(config.ConnectionString()).SetClientID(config.Conf.ClientID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var wg sync.WaitGroup

	for {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, endpoint := range config.Conf.SNMPMap.SNMPEndpoints {
				log.Println("Polling endpoint " + endpoint.Endpoint)

				snmp := gosnmp.GoSNMP{}

				snmp.Target = endpoint.Endpoint
				// port SNMP
				if endpoint.Port != 0 {
					snmp.Port = uint16(endpoint.Port)
				} else if config.Conf.SNMPPort != 0 {
					snmp.Port = uint16(config.Conf.SNMPPort) // Использование порта из конфигурации по умолчанию
				} else {
					snmp.Port = 161 // Использование порта по умолчанию
				}
				snmp.Version = gosnmp.Version2c
				snmp.Community = endpoint.Community

				snmp.Timeout = time.Duration(5 * time.Second)
				err := snmp.Connect()
				if err != nil {
					log.Fatal("SNMP Connect error\n")
				}

				oids := []string{}

				for _, oidTopic := range endpoint.OIDTopics {
					oids = append(oids, oidTopic.OID)
				}

				result, err := snmp.Get(oids)
				if err != nil {
					log.Printf("error in Get: %s", err)
				} else {
					for _, variable := range result.Variables {
						for _, oidTopic := range endpoint.OIDTopics {
							if strings.Compare(oidTopic.OID, variable.Name) == 0 {
								convertedValue := ""

								switch variable.Type {
								case gosnmp.OctetString:
									convertedValue = string(variable.Value.([]byte))
								default:
									convertedValue = fmt.Sprintf("%d", gosnmp.ToBigInt(variable.Value))
								}

								log.Printf("%s = %s", oidTopic.Topic, convertedValue)
								token := client.Publish(oidTopic.Topic, 0, false, convertedValue)

								token.Wait()
								if token.Error() != nil {
									log.Fatal(token.Error())
								}
							}
						}
					}
				}
				snmp.Conn.Close()
			}

		}()

		time.Sleep(time.Duration(config.Conf.Interval) * time.Second)
	}

	wg.Wait()

	client.Disconnect(250)
	time.Sleep(1 * time.Second)
}
