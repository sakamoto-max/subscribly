package middleware

import (
	"net/http"
	"subscribly/customerrors"
	"subscribly/utils"
)

func OwnerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, ok := GetClaimsFromContext(r.Context())
		if !ok {
			utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
			return 
		}


		if claims.Role != "owner" {
			utils.ErrorWriter(w, customerrors.ErrNoAccess, http.StatusBadRequest)
			return 
		}

		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, ok := GetClaimsFromContext(r.Context())
		if !ok {
			utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
			return 
		}

		if claims.Role != "admin" {
			utils.ErrorWriter(w, customerrors.ErrOnlyAdminAccess, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)		
	})

}