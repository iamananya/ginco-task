package middlewares

import (
	"context"
	"net/http"

	"github.com/iamananya/ginco-task/pkg/models"
)

// AuthenticationMiddleware is a middleware function that authenticates the user using the X-Token header
// and adds the user to the context of the request.
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		users, err := models.GetUserByToken(token)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(users) == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user := users[0]
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
