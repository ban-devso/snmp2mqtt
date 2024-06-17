package config

import (
	"encoding/json"
	"os"
	"strconv"
)

// OIDTopicObject maps OIDs to MQTT topics
type OIDTopicObject struct {
	OID   string `json:"oid"`
	Topic string `json:"topic"`
}

// SNMPEndpointObject is the SNMP Endpoint definition
type SNMPEndpointObject struct {
	Endpoint  string           `json:"endpoint"`
	Community string           `json:"community"`
	Port      int              `json:"port"` // порт SNMP
	OIDTopics []OIDTopicObject `json:"oidTopics"`
}

// SNMPMapObject basic map of endpoints
type SNMPMapObject struct {
	SNMPEndpoints []SNMPEndpointObject `json:"snmpEndpoints"`
}

// Config represents the configuration structure
type Config struct {
	SNMPMap   *SNMPMapObject `json:"snmpMap"`
	Server    string         `json:"server"`
	Port      int            `json:"port"`
	ClientID  string         `json:"clientid"`
	Interval  int            `json:"interval"`
	SNMPPort  int            `json:"snmpPort"`
}

var (
	// Global configuration
	Conf Config
)

// ConnectionString returns the MQTT connection string
func ConnectionString() string {
	return "tcp://" + Conf.Server + ":" + strconv.Itoa(Conf.Port)
}

// LoadMap loads the configuration from the file
func LoadMap(file string) error {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Conf)

	if err != nil {
		return err
	}

	return nil
}
