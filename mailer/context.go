package mailer

type ctxKey byte

const (
	ctxClient ctxKey = iota
	ctxLog
	ctxConn
)
