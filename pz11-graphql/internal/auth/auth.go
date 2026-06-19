package auth

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey struct{}

const DemoToken = "demo-token"

func WithAuthorization(ctx context.Context, header string) context.Context {
	return context.WithValue(ctx, ctxKey{}, header)
}

func AuthorizationFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKey{}).(string)
	return v
}

func RequireMutationToken(ctx context.Context) error {
	auth := AuthorizationFromContext(ctx)
	if strings.TrimPrefix(auth, "Bearer ") != DemoToken {
		return fmt.Errorf("unauthorized: mutations require Bearer %s", DemoToken)
	}
	return nil
}
