package middleware

import (
	"context"
	"net/http"
	"strconv"
	"subscribly/models"
)

var PageKey ContextKey = "page"

func Pagnation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var pageStr string
		var limitStr string
		
		var Page int
		var Limit int

		pageAndLimit := &models.PageLimit{}


		pageStr = r.URL.Query().Get("page")
		limitStr = r.URL.Query().Get("limit")

		Page, _ = strconv.Atoi(pageStr)
		Limit, _ = strconv.Atoi(limitStr)

		if Page == 0 || Page < 0 {
			Page = 1
		}

		if Limit == 0 || Limit < 0 {
			Limit = 20
		}

		pageAndLimit.PageNumber = Page
		pageAndLimit.Limit = Limit
		
		ctx := context.WithValue(r.Context(), PageKey, pageAndLimit)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}


func GetPageAndLimitFromCtx(ctx context.Context) (*models.PageLimit, bool) {
	page, ok := ctx.Value(PageKey).(*models.PageLimit)

	return page, ok

}
