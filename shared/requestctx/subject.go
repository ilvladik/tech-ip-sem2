package requestctx

import "context"

type subjectKey struct{}

func Subject(ctx context.Context) string {
	if rid, ok := ctx.Value(subjectKey{}).(string); ok {
		return rid
	}
	return ""
}

func WithSubject(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, subjectKey{}, requestID)
}
