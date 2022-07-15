package helpers

import (
	"RMS/models"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header["Token"]
		if token == nil {
			log.Printf("AuthMiddleWare : Token empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		authToken, err := jwt.Parse(token[0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error in parsing")
			}
			return []byte(GetSecretKey()), nil
		})
		log.Printf("AuthMiddleware : verified token")

		if claims, ok := authToken.Claims.(jwt.MapClaims); ok && authToken.Valid {

			username := fmt.Sprint(claims["username"])
			userRole := claims["userRole"].(string)
			userID := fmt.Sprint(claims["userID"])

			ctxMap := models.ContextMap{
				UserID:   userID,
				Username: username,
				UserRole: userRole,
			}

			ctx := context.WithValue(r.Context(), "User", ctxMap)
			log.Printf("AuthMiddleWare : sent value to context")
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			if err != nil {
				log.Printf("AuthMiddleWare : error in parsing the token.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

	})
}
