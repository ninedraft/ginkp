package command

import "time"

type Code byte

const (
	M_NUL = Code(iota)
	M_ADR
	M_PWD
	M_FILE
	M_OK
	M_EOB
	M_GOT
	M_ERR
	M_BSY
	M_GET
	M_SKIP
)

type FileConfig struct {
	Name   string
	Size   int
	Time   time.Time
	Offset int
}

type Command interface {
	String() string
	CommandCode() Code
	Unmarshall(p []byte) error
	Marhshal() ([]byte, error)
}

type MNUL struct {
	Message string
}

type MADR struct {
	Addresses []string
}

type MPWD struct {
	Password string
}

type MFILE struct {
	FileConfig
}
