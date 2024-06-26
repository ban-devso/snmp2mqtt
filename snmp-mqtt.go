package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ban-devso/snmp-mqtt/config"
	"github.com/ban-devso/snmp-mqtt/snmp"

	"github.com/docopt/docopt-go"
)

var exitChan = make(chan int)

// VERSION because...
const VERSION = "0.0.3"

func cliArguments() {
	usage := `
Usage: snmp-mqtt [options]

Options:
  --endpoints_map=<endpoints_map>     SNMP Endpoints Map File [default: ./endpoints.json]
  --config=<config>                   Config File [default: ./config.json]
  --server=<server>                   MQTT server host/IP [default: 127.0.0.1]
  --port=<port>                       MQTT server port [default: 1883]
  --portsnmp=<portsnmp>               SNMP port [default: 161]
  --clientid=<clientid>               MQTT client identifier [default: snmp]
  --interval=<interval>               Poll interval (seconds) [default: 5]
  -h, --help                          Show this screen.
  -v, --version                       Show version.
`
	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	configFile, _ := args.String("--config")
	err := config.LoadMap(configFile)
	if err != nil {
		log.Println(err)
		log.Fatal("error opening " + configFile)
	}

	config.Conf.Server, _ = args.String("--server")
	config.Conf.Port, _ = args.Int("--port")
	config.Conf.ClientID, _ = args.String("--clientid")
	config.Conf.Interval, _ = args.Int("--interval")
	config.Conf.SNMPPort, _ = args.Int("--portsnmp")

	log.Printf("server: %s, port: %d, client identifier: %s, poll interval: %d, SNMP port: %d", config.Conf.Server, config.Conf.Port, config.Conf.ClientID, config.Conf.Interval, config.Conf.SNMPPort)

	if config.Conf.SNMPMap != nil {
		for _, endpoint := range config.Conf.SNMPMap.SNMPEndpoints {
			log.Printf("SNMP Endpoint: %s, Port: %d, Community: %s", endpoint.Endpoint, endpoint.Port, endpoint.Community)
			for _, oidTopic := range endpoint.OIDTopics {
				log.Printf("OID: %s, Topic: %s", oidTopic.OID, oidTopic.Topic)
			}
		}
	}
}

// sigChannelListen basic handlers for inbound signals
func sigChannelListen() {
	signalChan := make(chan os.Signal, 1)
	code := 0

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)
	signal.Notify(signalChan, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		log.Printf("Received signal %s. shutting down", sig)
	case code = <-exitChan:
		switch code {
		case 0:
			log.Println("Shutting down")
		default:
			log.Println("*Shutting down")
		}
	}

	os.Exit(code)
}

func main() {
	cliArguments()

	// catch signals
	go sigChannelListen()

	// run sensor poll loop
	snmp.Init()

	os.Exit(0)
}
