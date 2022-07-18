package helpers

import (
	"RMS/models"
	"net/http"
)

func AdminPermissionMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context().Value("User").(models.ContextMap)
		signedUserRole := ctx.UserRole
		if signedUserRole != "admin" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
