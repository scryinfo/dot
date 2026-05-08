package contextex

import "context"

type ContextEx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewContextEx() *ContextEx {
	ctx, cansel := context.WithCancel(context.Background())
	return &ContextEx{ctx: ctx, cancel: cansel}
}

func (p *ContextEx) Cancel() {
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *ContextEx) Context() context.Context {
	return p.ctx
}
