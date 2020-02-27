package gorest

import "context"

// ContextHandler responsible to setup the request context with the requested resource based on the resource id.
type ContextHandler interface {
	// ContextWithResource responsible to validate the received resource id, and confirm if the requester is authorized to access it.
	// In case everything align, it should store the resource object that was fetched from some sort of external resource into the context.
	// In case an error occurs, an error object expected to be received back from it.
	// In case the resource is not found or the requester is not authorized to access it, a not found is expected as a return value.
	//
	// The signature is inspirited by the combination of os.LookupEnv and context.WithValue
	ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error)
}

type ContextHandlerFunc func(context.Context, string) (context.Context, bool, error)

func (fn ContextHandlerFunc) ContextWithResource(ctx context.Context, resourceID string) (context.Context, bool, error) {
	return fn(ctx, resourceID)
}
