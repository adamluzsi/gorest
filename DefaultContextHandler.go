package gorest

import "context"

type DefaultContextHandler struct{ ContextKey interface{} }

func (d DefaultContextHandler) WithResource(ctx context.Context, resourceID string) (context.Context, bool, error) {
	if resourceID == `` {
		return ctx, false, nil
	}
	return context.WithValue(ctx, d.ContextKey, resourceID), true, nil
}

func (d DefaultContextHandler) GetResourceID(ctx context.Context) interface{} {
	return ctx.Value(d.ContextKey)
}
