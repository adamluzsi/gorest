package gorest

import "context"

type ContextHandler interface {
	ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error)
}

type ContextHandlerFunc func(context.Context, string) (context.Context, bool, error)

func (fn ContextHandlerFunc) ContextWithResource(ctx context.Context, resourceID string) (context.Context, bool, error) {
	return fn(ctx, resourceID)
}
