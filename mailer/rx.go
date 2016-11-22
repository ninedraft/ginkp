package mailer

import (
	"context"
	"ginkp/frame"
	"io"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
)

type rxState byte

const (
	rxWaitF rxState = iota
	rxAccF
	rxReceD
	rxWriteD
	rxEOB
	rxDone
)

type rX struct {
	state rxState
}

type rxConfig struct {
	Conn     net.Conn
	Deadline time.Duration
	Client   *Client
}

func newRX(ctx context.Context) {
	//client := ctx.Value(ctxClient).(*Client)
	log := ctx.Value(ctxLog).(*logrus.Logger)
	conn := ctx.Value(ctxConn).(net.Conn)
	buf := getChunk()
	defer func() { returnChunk(buf) }()
	fr := frame.NewFrame()
	defer func() { frame.DeleteFrame(fr) }()

	select {
	case <-ctx.Done():
		return
	default:
	}
	_, err := io.CopyBuffer(fr, conn, buf)
	if err != nil && err != io.EOF {
		log.Error(err)
		return
	}
}
