package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/orka-org/orkacore/pkg/tokens"
)

func AuthMiddleware(secret, name string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			userID := ""
			if tr, ok := transport.FromServerContext(ctx); ok {
				header := tr.RequestHeader()
				if header.Get("Authorization") != "" {
					token, err := tokens.ValidateAuthorizationHeader(header.Get("Authorization"))
					if err != nil {
						return nil, err
					}
					p, err := tokens.NewTokenFactory(name, secret).NewTokenPayload().Parse(token)
					if err != nil {
						return nil, err
					}
					userID = p.GetID()
				} else {
					return nil, errors.Unauthorized("No Authorization Header", "The Authorization header is missing")
				}

				defer func() {
					// Do something on exiting
				}()
			}
			ctx = context.WithValue(ctx, "userID", userID)
			return handler(ctx, req)
		}
	}
}
