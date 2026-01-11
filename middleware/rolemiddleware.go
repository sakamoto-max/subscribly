package middleware

import (
	"net/http"
	"subscribly/customerrors"
	"subscribly/utils"
)

func RoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, ok := GetClaimsFromContext(r.Context())
		if !ok {
			utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		}


		if claims.Role != "owner" {
			utils.ErrorWriter(w, customerrors.ErrNoAccess, http.StatusBadRequest)
			return 
		}

		next.ServeHTTP(w, r)
	})
}