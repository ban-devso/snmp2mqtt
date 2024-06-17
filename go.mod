module github.com/ban-devso/snmp-mqtt

go 1.22.4

replace github.com/soniah/gosnmp => github.com/gosnmp/gosnmp v1.36.0

require (
	github.com/dchote/snmp-mqtt v0.0.0-20191120133650-5021794990f3
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/gosnmp/gosnmp v1.37.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/soniah/gosnmp v1.37.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)
