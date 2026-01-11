package middleware

import (
	"context"
	"fmt"
	"net/http"
	"subscribly/auth"
	"subscribly/customerrors"
	"subscribly/models"
	"subscribly/utils"

	"github.com/golang-jwt/jwt/v5"
)

func GlobalErrorHandlingMiddleware(err error)  {
	
}

type ContextKey string

// middlewares can send back response to the client
var ClaimsKey ContextKey = "claims"

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		AccessToken := r.Header.Get("token")

		claims, err := auth.ValidateJwtToken(AccessToken)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				cookie, err := r.Cookie("jwt")
				if err != nil {
					fmt.Printf("error in finding the cookie : %v\n", err)
					return
				}

				refreshToken := cookie.Value

				if refreshToken == "" {
					fmt.Println("refresh token is not found")
					return
				}

				err = auth.ValidateRefreshToken(refreshToken)
				if err != nil {
					if err == jwt.ErrTokenExpired {
						utils.ErrorWriter(w, customerrors.ErrRefreshTokenExpired, http.StatusBadRequest)
						return
					}
					fmt.Printf("error occured : %v\n", err)
					return
				}
				utils.ErrorWriter(w, customerrors.ErrAccessTokenExpired, http.StatusBadRequest)
				return

			}
			utils.ErrorWriter(w, err, http.StatusBadRequest)
			fmt.Println("token is invalid")
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaimsFromContext(ctx context.Context) (*models.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*models.Claims)

	return claims, ok
}



func ExpireTheCookie(cookie *http.Cookie) (*http.Cookie, error) {

	cookie.MaxAge = -1
	cookie.Value = ""

	return cookie, nil
}
