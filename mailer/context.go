package mailer

// +gen ring
type ctxKey string

const (
	ctxClient = ctxKey("client")
	ctxLog    = ctxKey("log")
	ctxConn   = ctxKey("connection")
)
