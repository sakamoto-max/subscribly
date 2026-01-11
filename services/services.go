package services

import (
	"strings"
	"subscribly/customerrors"
	"subscribly/models"
	"subscribly/repository"
	"subscribly/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

func UserSignUpService(Credentials models.UserLoginCreds) (models.UserDetails, error) {

	var response models.UserDetails
	Credentials.Email = strings.ToLower(Credentials.Email)

	// check if user already signedUp

	exists, err := repository.EmailExistsInDB(Credentials.Email)
	if err != nil {
		return response, err
	}

	if exists {
		return response, customerrors.ErrUserAlreadySignedUp
	}

	hashedPassword, err := utils.PasswordHasher(Credentials.Password)
	if err != nil {
		return response, err
	}
	Credentials.Password = hashedPassword

	if Credentials.Role == "" {
		Credentials.Role = "user"
	}

	response, err = repository.UserSignUpDB(Credentials)
	response.Role = Credentials.Role
	if err != nil {
		return response, err
	}

	return response, nil

}

func UserLoginService(Credential models.UserLoginCreds) (int, string, error) {

	var userId int
	var role string

	Credential.Email = strings.ToLower(Credential.Email)

	exists, err := repository.EmailExistsInDB(Credential.Email)
	if err != nil {
		return userId, role, err
	}

	if !exists{
		return userId, role, customerrors.ErrNotSignedUp
	}

	hashedPassword, err := repository.GetHashedPassFromDb(Credential.Email)
	if err != nil {
		return userId, role, err
	}

	err = utils.ComparePassword(Credential.Password, hashedPassword)
	if err != nil {
		return userId, role, customerrors.ErrPasswordInvalid
	}

	userId, err = repository.GetUserId(Credential.Email)
	if err != nil {
		return userId, role, err
	}

	role, err = repository.GetUserRoleFromDb(userId)
	if err != nil {
		return userId, role, err
	}

	return userId, role, nil
}

func CreateNewOrgService(orgName models.Org) (models.OrgWithPlans, error) {

	response, err := repository.CreateNewOrgInDB(orgName.OrgName, orgName.OwnerId)
	if err != nil {
		return response, err
	}

	return response, nil
}

func GetAllOrgsService() ([]models.Org, error) {
	var allOrgs []models.Org
	rows, err := repository.GetAllOrgsFromDB()
	if err != nil {
		return allOrgs, err
	}

	defer rows.Close()

	var id int
	var orgName string
	var ownerId int
	var createdAT time.Time

	for rows.Next() {
		err := rows.Scan(&id, &orgName, &ownerId, &createdAT)
		if err != nil {
			return allOrgs, err
		}
		a := models.Org{Id: id, OrgName: orgName, OwnerId: ownerId, CreatedAt: createdAT}

		allOrgs = append(allOrgs, a)
	}

	return allOrgs, nil
}

func JoinOrgService(userID int, orgName string) error {

	orgId, err := repository.GetOrgId(orgName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customerrors.ErrOrgNotExists
		}
		return err
	}

	joined, err := repository.UserJoinedOrg(userID)
	if err != nil {
		return err
	}

	if joined {
		return customerrors.ErrAlreadyJoinedOrg
	}

	err = repository.JoinOrgInDB(userID, orgId)
	if err != nil {
		return err
	}

	return nil

}

func GetAllPlansService(userId int) ([]models.Plans, error) {

	// check if the user is owner or not
	var allPlans []models.Plans

	rows, err := repository.GetAllPlansFormDB()
	if err != nil {
		return allPlans, err
	}

	defer rows.Close()

	var id int
	var planName string
	var price int
	var totalUses int

	for rows.Next() {
		err := rows.Scan(&id, &planName, &price, &totalUses)
		if err != nil {
			return allPlans, err
		}

		newPlan := models.Plans{Id: id, PlanName: planName, Price: price, TotalUses: totalUses}

		allPlans = append(allPlans, newPlan)
	}

	return allPlans, nil

}

func UseService(userId int) (models.ReminingCounts, error) {

	var a models.ReminingCounts

	orgId, err := repository.FindOrgIDByuserID(userId)

	if err != nil {
		return a, err
	}

	// check if the count is zero
	usesLeft, err := repository.GetNumberOfUsesLeft(userId)
	if err != nil {
		return a, err
	}
	if usesLeft == 0 {
		return a, customerrors.ErrUsesOver
	}

	usesLeft, err = repository.ReductTheUseCount(orgId)
	if err != nil {
		return a, err
	}

	a = models.ReminingCounts{UsesLeft: usesLeft}

	return a, nil
}

func GetAllMembersService(userId int, page int, limit int) ([]models.Members, error) {

	var allMembers []models.Members

	offSet := (page - 1) * limit

	orgId, err := repository.GetOrgIdByOwnerId(userId)
	if err != nil {
		return allMembers, err
	}

	rows, err := repository.GetAllMembersInOrgDB(orgId, limit, offSet)
	if err != nil {
		return allMembers, err
	}

	defer rows.Close()

	var id int
	var name string
	var email string

	for rows.Next() {
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			return allMembers, err
		}

		a := models.Members{Id: id, Name: name, Email: email}

		allMembers = append(allMembers, a)
	}

	return allMembers, nil

}

func UpgradePlanService(userId int, planName string) (models.NewPlan, error) {

	var response models.NewPlan

	planName = strings.ToLower(planName)

	currentPlan, err := repository.GetCurrentPlan(userId)
	if err != nil {
		return response, err
	}

	if planName == currentPlan {
		return response, customerrors.ErrSamePlan
		// do this
	}

	newPlanId, err := repository.GetPlanId(planName)
	if err != nil {
		return response, err
	}

	orgId, err := repository.GetOrgIdByOwnerId(userId)
	if err != nil {
		return response, err
	}

	response, err = repository.ChangePlan(orgId, newPlanId)
	if err != nil {
		return response, err
	}

	response.PlanName = planName

	response.Message = "successfully upgraded"

	return response, nil
}

func UsesLeftService(userId int) (models.UsesLeft, error) {

	var uLeft models.UsesLeft

	usesLeft, err := repository.GetNumberOfUsesLeft(userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uLeft, customerrors.ErrNotJoinedOrg
		}
		return uLeft, err
	}

	uLeft.NoOfUsesLeft = usesLeft

	return uLeft, nil
}

func GetAllUsersService(page int, limit int) ([]models.Members, error) {

	var allUsers []models.Members

	offset := (page - 1) * limit

	rows, err := repository.GetAllUsers(limit, offset)
	if err != nil {
		return allUsers, err
	}

	defer rows.Close()

	var id int
	var name string
	var email string

	for rows.Next() {
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			return allUsers, err
		}

		a := models.Members{Id: id, Name: name, Email: email}

		allUsers = append(allUsers, a)
	}

	return allUsers, nil

}
