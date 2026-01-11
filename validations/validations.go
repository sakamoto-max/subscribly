package validations

import (
	"subscribly/models"
)

type validationErrors struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}



var (
	EmailReq = validationErrors{Path: "email", Message: "required"}
	NameReq = validationErrors{Path: "name", Message: "required"}
	PasswordReq = validationErrors{Path: "password", Message: "required"}
	OrgNameReq = validationErrors{Path: "org_name", Message: "required"}
	PlanDoesNotExist = validationErrors{Path: "plan_Name", Message: "doesn't exist"}
	PlanNameReq = validationErrors{Path: "plan_name", Message: "required"}
)


func SignUpValidator(userSentDetails models.UserLoginCreds) ([]validationErrors, bool) {

	var allErrors []validationErrors

	if userSentDetails.Email == ""{
		allErrors = append(allErrors, EmailReq)
	}

	if userSentDetails.Name == "" {
		allErrors = append(allErrors, NameReq)
	}

	if userSentDetails.Password == "" {
		allErrors = append(allErrors, PasswordReq)
		return allErrors, false
	}

	return allErrors, true
}

func LoginValidator(userSentDetails models.UserLoginCreds) ([]validationErrors, bool) {

	var allErrors []validationErrors

	if userSentDetails.Email == ""{
		allErrors = append(allErrors, EmailReq)
	}

	if userSentDetails.Password == "" {
		allErrors = append(allErrors, PasswordReq)
		return allErrors, false
	}

	return allErrors, true

}

func NewOrgValidator(userSentDetails models.Org) (validationErrors, bool) {

	if userSentDetails.OrgName == "" {
		return OrgNameReq, false
	}

	return OrgNameReq, true

}

func UpgradePlanValidator(planName string) (validationErrors, bool) {

	if planName == "" {
		return PlanNameReq, false
	}
	
	if planName != "free" {
		if planName != "pro" {
			if planName != "ultra" {
				return PlanDoesNotExist, false
			}
		}
	}

	return PlanDoesNotExist, true
}