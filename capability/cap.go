package capability

// A Capabiltity is a capability that the server implements
type Cap struct {
	Name, Value string
}

func (c Cap) String() string { return c.Name }

var (
	AwayNotify  = Cap{Name: "away-notify"}
	CapNotify   = Cap{Name: "cap-notify"}
	Chghost     = Cap{Name: "chghost"}
	EchoMessage = Cap{Name: "echo-message"}
	MessageTags = Cap{Name: "message-tags"}
	MultiPrefix = Cap{Name: "multi-prefix"}
	SASL        = Cap{Name: "sasl", Value: "PLAIN,EXTERNAL,SCRAM-SHA-256"}
	ServerTime  = Cap{Name: "server-time"}
	Setname     = Cap{Name: "setname"}
	STS         = Cap{Name: "sts", Value: "port=%s,duration=%.f"}
)
