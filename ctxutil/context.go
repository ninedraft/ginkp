package ctxutil

import "context"
import "time"

type ContextBuilder interface {
	context.Context
	With(interface{}, interface{}) ContextBuilder
	Ctx() context.Context
}

type contextBuilder struct {
	ctx context.Context
}

func (cbuld *contextBuilder) Deadline() (deadline time.Time, ok bool) {
	deadline, ok = cbuld.ctx.Deadline()
	return
}

func (cbuld *contextBuilder) Done() <-chan struct{} {
	return cbuld.ctx.Done()
}

func (cbuild *contextBuilder) Err() error {
	return cbuild.ctx.Err()
}

func (cbuild *contextBuilder) Value(k interface{}) interface{} {
	return cbuild.ctx.Value(k)
}

func (cbuild *contextBuilder) Ctx() context.Context {
	return cbuild.ctx
}

func (cbuild *contextBuilder) With(k, v interface{}) ContextBuilder {
	cbuild.ctx = context.WithValue(cbuild.ctx, k, v)
	return cbuild
}

func CtxBuilderFromCtx(ctx context.Context) ContextBuilder {
	return &contextBuilder{
		ctx: ctx,
	}
}

func CtxBuilder() ContextBuilder {
	return &contextBuilder{
		ctx: context.Background(),
	}
}
