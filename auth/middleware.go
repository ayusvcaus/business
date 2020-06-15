package auth

import (
	"context"
	"net/http"

	//"strconv"
	"strings"

	//"github.com/alecthomas/log4go"

	//"github.com/ayusvcaus/business/graph/model"
	"github.com/ayusvcaus/business/persist"
	//"github.com/ayusvcaus/business/graph/resolvers"
)

const BEARER_SCHEMA = "Bearer "

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			//Allow client without bearer token incase login
			if !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
				next.ServeHTTP(w, r)
				return
			}
			token := authHeader[len(BEARER_SCHEMA):]
			//Jwt
			username, err := ParseToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}
			// create user and check if user exists in db
			user := persist.User{Username: username}
			//id := user.GetUserIdByUsername(username)
			//if id < 1 {
			//next.ServeHTTP(w, r)
			//http.Error(w, "Invalid User ID", http.StatusForbidden)
			//return
			//}
			//user.ID = id
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *persist.User {
	raw, _ := ctx.Value(userCtxKey).(*persist.User)
	return raw
}
