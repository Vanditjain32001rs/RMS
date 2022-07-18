package helpers

import (
	"RMS/models"
	"net/http"
)

func SubAdminPermissionMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context().Value("User").(models.ContextMap)
		signedUserRole := ctx.UserRole
		if signedUserRole != "subadmin" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
