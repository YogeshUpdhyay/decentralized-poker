package constants

const (
	Empty                       = ""
	ApplicationDataDir          = ".yoker"
	ApplicationIdentityFileName = "identitygamma.key"
	ApplicationConfigFileName   = "config.yml"
	ApplicationDBFilleName      = "yoker.db"
	DefaultDatabasePath         = ApplicationDataDir + "/" + ApplicationDBFilleName

	Ping = "PING"
	Pong = "PONG"

	ServerName = "serverName"

	ServerNameDefault     = "yoker alpha"
	ServerPortDefault     = "9000"
	ServerVersionDefault  = "yoker1.0.0"
	StreamProtocolDefault = "/ypoker/1.0.0"

	// connection states
	ConnectionStateActive   = "active"
	ConnectionStateInactive = "inactive"
	ConnectionStatePending  = "pending"

	// message statuses
	ToBeSent  = "toBeSent"
	Sending   = "sending"
	Sent      = "sent"
	Delivered = "delivered"
	Read      = "read"
	Received  = "received"
)
