package models

// TrafficMessage represents kafka topic message for
// nac-passive-scanner-traffic-topic topic.
type TrafficMessage struct {
	// Customer universally unique identifier.
	CustomerUUID string `json:"customerUuid"`

	// Le identification.
	LeID string `json:"leId"`

	// Passive scanner traffic flow
	PassiveScannerTrafficFlows []passiveScannerTrafficFlow `json:"passiveScannerTrafficFlows"`
}

type passiveScannerTrafficFlow struct {
	// Client identification
	ClientID string `json:"clientId"`

	// Server identification
	ServerID string `json:"serverId"`

	// Client IP address
	ClientIP string `json:"clientIp"`

	// Server IP address
	ServerIP string `json:"ClientIp"`

	// Server Port
	ServerPort int `json:"serverPort"`

	// Protocol
	Protocol string `json:"protocol"`

	// Server application
	ServerApplication string `json:"ServerApplication"`

	// Server name
	ServerName string `json:"serverName"`

	// TxBytes
	TxBytes int `json:"txBytes"`

	// TxPkts
	TxPkts int `json:"txPkts"`

	// RxBytes
	RxBytes int `json:"rxBytes"`

	// RxPkts
	RxPkts int `json:"rxPkts"`

	// Timestamp
	Timestamp int `json:"timestamp"`
}
