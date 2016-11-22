package mailer

type ctxKey string

const (
	ctxClient = ctxKey("client")
	ctxLog    = ctxKey("log")
	ctxConn   = ctxKey("connection")
)
