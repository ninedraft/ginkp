package mailer

import "io"
import "net"
import "github.com/Sirupsen/logrus"
import "context"
import (
	. "ginkp/ctxutil"
)

type Tunnel interface {
	io.ReadWriter
}

type cState byte

const (
	cConnInit = cState(iota)
	cWaitConn
	cSendPasswd
	cWaitAddr
	cAuthRemote
	cIfSecure
	cWaitOK
	cOpts
	cPwdAck
	cInitTransfer
	cSwitch
	cRecive
	cTransmit
)

type Client struct {
	Address  string
	Log      *logrus.Logger
	Listener net.Listener
	Context  context.Context
	state    cState
}

func (client *Client) RunWithContext(ctx context.Context) {
	go func() {
		for {
			conn, err := client.Listener.Accept()
			if err != nil {
				client.Log.Error(err)
			}
			go newRX(CtxBuilderFromCtx(client.Context).
				With(ctxConn, conn).
				With(ctxClient, client).
				With(ctxLog, client.Log).
				Ctx())
		}
	}()
}
