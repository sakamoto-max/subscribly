package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"subscribly/models"
	"subscribly/utils"
)

type validationErrors struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

var (
	EmailReq         = validationErrors{Path: "email", Message: "required"}
	NameReq          = validationErrors{Path: "name", Message: "required"}
	PasswordReq      = validationErrors{Path: "password", Message: "required"}
	OrgNameReq       = validationErrors{Path: "org_name", Message: "required"}
	PlanDoesNotExist = validationErrors{Path: "plan_Name", Message: "doesn't exist"}
	PlanNameReq      = validationErrors{Path: "plan_name", Message: "required"}
)

var UserSingUpKey string = "userSignUp"
var UserLoginKey string = "userLogIn"
var NewOrgKey string = "newOrg"
var UpgradePlanKey string = "upgradePlan"

func SignUpValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userSentDetails models.UserLoginCreds
		var errOccured bool = false

		json.NewDecoder(r.Body).Decode(&userSentDetails)

		var allErrors []validationErrors

		if userSentDetails.Name == "" {
			allErrors = append(allErrors, NameReq)
			errOccured = true
		}

		if userSentDetails.Email == "" {
			allErrors = append(allErrors, EmailReq)
			errOccured = true
		}

		if userSentDetails.Password == "" {
			allErrors = append(allErrors, PasswordReq)
			errOccured = true
		}

		if errOccured {
			utils.ValidationErrWriter(w, allErrors)
			return
		}

		ctx := context.WithValue(r.Context(), UserSingUpKey, &userSentDetails)

		next.ServeHTTP(w, r.WithContext(ctx))
		// set the req to the context
	})
}

func LoginValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userSentDetails models.UserLoginCreds
		var allErrors []validationErrors
		var errOccured bool

		json.NewDecoder(r.Body).Decode(&userSentDetails)

		if userSentDetails.Email == "" {
			allErrors = append(allErrors, EmailReq)
			errOccured = true
		}

		if userSentDetails.Password == "" {
			allErrors = append(allErrors, PasswordReq)
			errOccured = true
		}

		if errOccured {
			json.NewEncoder(w).Encode(allErrors)
			return
		}

		ctx := context.WithValue(r.Context(), UserLoginKey, &userSentDetails)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewOrgValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userSentDetails models.Org
		var allErrors []validationErrors
		var errOccured bool

		json.NewDecoder(r.Body).Decode(&userSentDetails)

		if userSentDetails.OrgName == "" {
			allErrors = append(allErrors, OrgNameReq)
			errOccured = true
		}

		if errOccured {
			json.NewEncoder(w).Encode(allErrors)
			return
		}

		ctx := context.WithValue(r.Context(), NewOrgKey, &userSentDetails)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func UpgradePlanValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		planName := r.PathValue("planName")

		if planName == "" {
			json.NewEncoder(w).Encode(PlanNameReq)
			return
		}

		if planName != "free" {
			if planName != "pro" {
				if planName != "ultra" {
					json.NewEncoder(w).Encode(PlanDoesNotExist)
					return 
				}
			}
		}

		ctx := context.WithValue(r.Context(), UpgradePlanKey, planName)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func GetSignUpDetailsFromctx(ctx context.Context) (*models.UserLoginCreds, bool) {
	creds, ok := ctx.Value(UserSingUpKey).(*models.UserLoginCreds)

	return creds, ok
}

func GetLoginDetailsFromctx(ctx context.Context) (*models.UserLoginCreds, bool) {
	creds, ok := ctx.Value(UserLoginKey).(*models.UserLoginCreds)

	return creds, ok
}

func GetNewOrgDetailsFromCtx(ctx context.Context) (*models.Org, bool) {
	details, ok := ctx.Value(NewOrgKey).(*models.Org)

	return details, ok
}

func GetNewPlanFromCtx(ctx context.Context) (string, bool) {

	planName, ok := ctx.Value(NewOrgKey).(string)

	return planName, ok

}