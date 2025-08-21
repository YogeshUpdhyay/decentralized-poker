package constants

const (
	Empty                       = ""
	ApplicationDataDir          = ".yoker"
	ApplicationIdentityFileName = "identity.key"
	ApplicationConfigFileName   = "config.yml"

	Ping = "PING"
	Pong = "PONG"

	ServerName = "serverName"

	ServerNameDefault     = "yoker alpha"
	ServerPortDefault     = "3000"
	ServerVersionDefault  = "yoker1.0.0"
	StreamProtocolDefault = "/ypoker/1.0.0"

	// connection states
	ConnectionStateActive   = "active"
	ConnectionStateInactive = "inactive"
)
