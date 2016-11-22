package mailer

type txState byte

const (
	txGNF txState = iota
	txTryR
	txReadS
	txWLA
	txDone
)

type tX struct {
	state txState
}
