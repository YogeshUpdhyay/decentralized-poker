package constants

const (
	Empty                       = ""
	ApplicationDataDir          = ".yoker"
	ApplicationIdentityFileName = "identitygamma.key"
	ApplicationConfigFileName   = "config.yml"

	Ping = "PING"
	Pong = "PONG"

	ServerName = "serverName"

	ServerNameDefault     = "yoker gamma"
	ServerPortDefault     = "9000"
	ServerVersionDefault  = "yoker1.0.0"
	StreamProtocolDefault = "/ypoker/1.0.0"

	// connection states
	ConnectionStateActive   = "active"
	ConnectionStateInactive = "inactive"

	// message statuses
	ToBeSent  = "toBeSent"
	Sending   = "sending"
	Sent      = "sent"
	Delivered = "delivered"
	Read      = "read"
	Received  = "received"
)
