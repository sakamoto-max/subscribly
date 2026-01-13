package handlers

import (
	// "context"
	"encoding/json"
	"fmt"
	"net/http"
	"subscribly/auth"
	"subscribly/customerrors"
	"subscribly/middleware"
	"subscribly/services"
	"subscribly/utils"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "app is alive",
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func UserSignUp(w http.ResponseWriter, r *http.Request) {

	userSentDetails, ok := middleware.GetSignUpDetailsFromctx(r.Context())
	if !ok {
		fmt.Println("error getting details form context")
		utils.ErrorWriter(w, customerrors.ErrGettingDetailsFromContext, http.StatusBadRequest)
		return
	}	

	response, err := services.UserSignUpService(userSentDetails)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, response, http.StatusCreated)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	userSentDetails, ok := middleware.GetLoginDetailsFromctx(r.Context())
	if !ok {
		fmt.Println("error getting details form context")
		utils.ErrorWriter(w, customerrors.ErrGettingDetailsFromContext, http.StatusBadRequest)
		return

	}

	userId,role,  err := services.UserLoginService(userSentDetails)
	if err != nil {

		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	AccessToken, RefreshToken, err := auth.GenerateAccesndRefresh(userId, role)

	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    RefreshToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   600,
	}

	http.SetCookie(w, &cookie)

	response := map[string]string{
		"message":      "login successful",
		"Access_Token": AccessToken,
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func UserLogOut(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("jwt")
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	ExpiredCookie, err := middleware.ExpireTheCookie(cookie)

	http.SetCookie(w, ExpiredCookie)

	response := map[string]string{
		"message": "logout successful",
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func CreateNewOrg(w http.ResponseWriter, r *http.Request) {

	userSentDetails, ok := middleware.GetNewOrgDetailsFromCtx(r.Context())
	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingDetailsFromContext, http.StatusBadRequest)
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		fmt.Println("error occured")
		return
	}
	userSentDetails.OwnerId = claims.UserId

	response, err := services.CreateNewOrgService(userSentDetails)
	if err != nil {
		fmt.Printf("error occured : %v\n", err)

		response := map[string]string{
			"message": "error occured",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	utils.ResponseWriter(w, response, http.StatusCreated)
}

func GetAllOrgs(w http.ResponseWriter, r *http.Request) {

	response, err := services.GetAllOrgsService()
	if err != nil {
		fmt.Printf("error occured : %v\n", err)

		response := map[string]string{
			"message": "error occured",
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func JoinOrg(w http.ResponseWriter, r *http.Request) {

	orgName := r.PathValue("orgName")

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		return
	}

	err := services.JoinOrgService(claims.UserId, orgName)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": "joined successfully",
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func GetAllPlans(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		return
	}

	response, err := services.GetAllPlansService(claims.UserId)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func Use(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetClaimsFromContext(r.Context())

	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		return
	}

	response, err := services.UseService(claims.UserId)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func GetAllMembersInOrg(w http.ResponseWriter, r *http.Request) {
	
	pageLimit, ok := middleware.GetPageAndLimitFromCtx(r.Context())
	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingPage, http.StatusBadRequest)
		return
	}
	// fmt.Println("page limit fetched")

	claims, ok := middleware.GetClaimsFromContext(r.Context())

	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		return
	}

	response, err := services.GetAllMembersService(claims.UserId, pageLimit.PageNumber, pageLimit.Limit)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, response, http.StatusOK)

}

func UpgradePlan(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetClaimsFromContext(r.Context())

	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		return
	}

	planName, ok := middleware.GetNewPlanFromCtx(r.Context())
	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingDetailsFromContext, http.StatusBadRequest)
		return
	}

	response, err := services.UpgradePlanService(claims.UserId, planName)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func UsesLeft(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r.Context())

	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingClaims, http.StatusBadRequest)
		return
	}

	response, err := services.UsesLeftService(claims.UserId)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, response, http.StatusOK)
}

func GenerateNewAccessToken(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("jwt")
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	refreshToken := cookie.Value

	claims, err := auth.ValidateJwtToken(refreshToken)
	if err != nil {
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	newAccessToken, err := auth.GenerateNewToken(claims.UserId, claims.Role)

	response := map[string]string{
		"new_access_token": newAccessToken,
	}

	utils.ResponseWriter(w, response, http.StatusCreated)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	pageLimit, ok := middleware.GetPageAndLimitFromCtx(r.Context())
	if !ok {
		utils.ErrorWriter(w, customerrors.ErrGettingPage, http.StatusBadRequest)
		return
	}

	repsonse, err := services.GetAllUsersService(pageLimit.PageNumber, pageLimit.Limit) 
	if err != nil{
		utils.ErrorWriter(w, err, http.StatusBadRequest)
		return
	}

	utils.ResponseWriter(w, repsonse, http.StatusOK)
}

func GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {

	search, ok := middleware.GetSearchFromCts(r.Context())
	if !ok {
		//
	}

	fmt.Println(search)

	if search == "" {
		response, err := services.GetAllSubscriptionsService()
		if err != nil{
			utils.ErrorWriter(w, err, http.StatusBadRequest)
			return
		}
		utils.ResponseWriter(w, response, http.StatusOK)
	}else {
		response, err := services.GetAllSubscriptionsServiceWithFilter(search)
		if err != nil{
			utils.ErrorWriter(w, err, http.StatusBadRequest)
			return	
		}
		utils.ResponseWriter(w, response, http.StatusOK)

	}
}