package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"todo/api/responses"

	jwt "github.com/dgrijalva/jwt-go"
)

func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = map[string]interface{}{"status": "failed", "detail": "Missing authorization token"}

		var header = r.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			responses.JsonResponse(w, http.StatusForbidden, resp)
			return
		}

		accessTokenValue := header[7:len(header)]

		token, err := jwt.Parse(accessTokenValue, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			resp["status"] = "failed"
			resp["detail"] = "Invalid token, please login"
			fmt.Println("err is", err)
			responses.JsonResponse(w, http.StatusForbidden, resp)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "userID", claims["userID"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
