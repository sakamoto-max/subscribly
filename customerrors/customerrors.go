package customerrors

import "errors"


type MyErrors struct {
	Message string `json:"message"`
}
var (

	// bad request
	ErrEmailNotFound = errors.New("email is not found")
	ErrNotSignedUp = errors.New("please signup first")
	ErrPasswordInvalid = errors.New("password is invalid")
	ErrInvalidToken = errors.New("jwt token is invalid")
	ErrPlanNotExists = errors.New("this plan does not exist")
	ErrCookieNotPresent = errors.New("cookie is not present")
	ErrGettingPage = errors.New("error getting page number")
	ErrUserAlreadySignedUp = errors.New("user already exists")


	// internal server errors

	ErrInternalServerError = errors.New("some internal server error occured")
	ErrValidationError = errors.New("validation error occured")
	ErrOrgNotExists = errors.New("the org doesn't exist")
	ErrGettingClaims = errors.New("Error getting claims")
	ErrTokenExpired = errors.New("token is expired, please login again")
	ErrRefreshTokenExpired = errors.New("your refresh token is expired please login again")
	ErrAccessTokenExpired = errors.New("your access token is expired please go to /refresh and get a new access token")
	ErrNoAccess = errors.New("only owners can access this")
	ErrSamePlan = errors.New("your org is currently using this plan")

	ErrNotJoinedOrg = errors.New("you haven't joined any org yet. please join an org")
	ErrUsesOver = errors.New("U have completed all the uses. please upgrade ur plan or wait for the next month")
	ErrRTokenNotPresent = errors.New("refresh token is not present")
	
	
	ErrDummyError = errors.New("this is a dummy error")

	ErrAlreadyJoinedOrg = errors.New("you have already joined an org")

	ErrNotLoggedIn = MyErrors{Message: "please login first"}
	ErrDummyMyError = MyErrors{Message: "this is a dummy error"}
)

// 2 errors
// bad request
// server error

