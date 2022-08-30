package customcontext

import "context"

type (
	Context struct {
		Ctx       context.Context
		CancelCtx context.CancelFunc
	}
)

func NewContext() (result *Context) {
	result = &Context{}

	result.Ctx, result.CancelCtx = context.WithCancel(context.Background())
	return
}

func (ctx *Context) Close() {
	ctx.CancelCtx()
}

func (ctx *Context) AddToContext(key, value interface{}) {
	ctx.Ctx = context.WithValue(ctx.Ctx, key, value)
	return
}

func (ctx *Context) GetFromContext(key interface{}) (result interface{}) {
	result = ctx.Ctx.Value(key)
	return
}
