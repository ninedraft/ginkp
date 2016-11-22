package mailer

type ctxKey int

const (
	ctxClient ctxKey = iota
	ctxLog
	ctxConn
)
