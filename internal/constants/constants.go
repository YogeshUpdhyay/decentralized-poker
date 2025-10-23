package constants

const (
	Empty                       = ""
	PeerID                      = "peerID"
	ApplicationDataDir          = ".yoker"
	ApplicationIdentityFileName = "identitypi.key"
	ApplicationConfigFileName   = "config.yml"
	ApplicationDBFilleName      = "yoker.db"
	DatabasePathDefault         = ApplicationDataDir + "/" + ApplicationDBFilleName

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

	// request status
	RequestStatusSent             = "sent"
	RequestStatusAwaitingDecision = "awaiting_decision"
	RequestStatusAccepted         = "accepted"
	RequestStatusRejected         = "rejected"

	// message statuses
	ToBeSent  = "toBeSent"
	Sending   = "sending"
	Sent      = "sent"
	Delivered = "delivered"
	Read      = "read"
	Received  = "received"

	// events
	EventNewConnectionRequest = "new_connection_request"
	EventThreadListUpdated    = "thread_list_upadted"
	EventNewMessage           = "new_message"

	// dummy
	DummyAvatarUrl = "https://avatar.iran.liara.run/username?username=dummy&bold=false&length=1"
)
