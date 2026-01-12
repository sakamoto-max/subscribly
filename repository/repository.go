package repository

import (
	"context"
	"subscribly/customerrors"
	"subscribly/database"
	"subscribly/models"
	"github.com/jackc/pgx/v5"
)

func UserSignUpDB(C *models.UserLoginCreds) (models.UserDetails, error) {

	var D models.UserDetails

	trnx, err := database.DBConn.Begin(context.TODO())
	if err != nil {
		return D, err
	}

	err = trnx.QueryRow(context.TODO(), `
		INSERT INTO USERS(NAME, EMAIL, HASHED_PASSWORD, CREATED_AT, UPDATED_AT, JOINED_ORG)
		VALUES($1, $2, $3, NOW(), NOW(), false)
		RETURNING ID, NAME, EMAIL, CREATED_AT, UPDATED_AT
	`, C.Name, C.Email, C.Password).Scan(&D.Id, &D.Name, &D.Email, &D.CreatedAt, &D.UpdatedAt)
	if err != nil {
		return D, err
	}

	roleId, err := getRoleId(C.Role)
	if err != nil {
		return D, err
	}

	_, err = trnx.Exec(context.TODO(), `
		INSERT INTO USER_ROLES(USER_ID, ROLE_ID)
		VALUES($1, $2)
	`, D.Id, roleId)
	if err != nil {
		return D, err
	}

	err = trnx.Commit(context.TODO())
	if err != nil {
		return D, err
	}

	return D, nil
}

func getRoleId(roleName string) (int, error) {

	var id int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM ROLES
		WHERE ROLE_NAME = $1	
	`, roleName).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func EmailExistsInDB(email string) (bool, error) {
	var id int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM USERS
		WHERE EMAIL = $1
	`, email).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows{
			return false, nil
		}
		return false, err
	}

	return true, nil
}


func EmailExistsInDBV2(email string) (error) {
	var id int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM USERS
		WHERE EMAIL = $1
	`, email).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows{
			return customerrors.ErrEmailNotFound
		}
		return err
	}

	return nil
}

func GetHashedPassFromDb(email string) (string, error) {

	var hashedPassword string

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT HASHED_PASSWORD FROM USERS
		WHERE EMAIL = $1	
	`, email).Scan(&hashedPassword)

	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func CreateNewOrgInDB(orgName string, userId int) (models.OrgWithPlans, error) {

	var newOrg models.OrgWithPlans

	planDetails, err := GetPlanDetails("free")
	if err != nil {
		return newOrg, err
	}

	newOrg.PlanName = "free"

	trnx, err := database.DBConn.Begin(context.TODO())
	if err != nil {
		return newOrg, err
	}

	err = trnx.QueryRow(context.TODO(), `
		INSERT INTO ORGS(ORG_NAME, OWNER_ID, CREATED_AT)
		VALUES($1, $2, NOW())	
		RETURNING ID, ORG_NAME, OWNER_ID, CREATED_AT
	`, orgName, userId).Scan(&newOrg.OrgId, &newOrg.OrgName, &newOrg.OwnerId, &newOrg.CreatedAt)
	if err != nil {
		return newOrg, err
	}

	_, err = trnx.Exec(context.TODO(), `
		UPDATE USER_ROLES
		SET ROLE_ID = 2
		WHERE USER_ID = $1
	`, userId)
	if err != nil {
		return newOrg, err
	}

	err = trnx.QueryRow(context.TODO(), `
		INSERT INTO SUBSCRIPTIONS(ORG_ID, PLAN_ID, PURCHASED_AT, TOTAL_USES, NO_OF_TIMES_USED)
		VALUES($1, $2, NOW(), $3, $4)
		RETURNING PLAN_ID, PURCHASED_AT, TOTAL_USES, NO_OF_TIMES_USED
	`, newOrg.OrgId, planDetails.Id, planDetails.TotalUses, planDetails.TotalUses).Scan(
		&newOrg.PlanId, &newOrg.PurchasedAt, &newOrg.TotalUses, &newOrg.UsesLeft)

	_, err = trnx.Exec(context.TODO(), `
		INSERT INTO ORG_MEMBERS(org_id, user_id)
		VALUES($1, $2)	
	`, newOrg.OrgId, newOrg.OwnerId)
	if err != nil {
		return newOrg, err
	}

	err = trnx.Commit(context.TODO())
	if err != nil {
		return newOrg, err
	}

	return newOrg, nil

}

func GetUserId(email string) (int, error) {
	var userID int
	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM USERS
		WHERE EMAIL = $1	
	`, email).Scan(&userID)

	if err != nil {
		return userID, err
	}

	return userID, nil
}

func GetAllOrgsFromDB() (pgx.Rows, error) {
	rows, err := database.DBConn.Query(context.TODO(), `
		SELECT * FROM ORGS	
	`)

	if err != nil {
		return rows, err
	}

	return rows, nil
}

func JoinOrgInDB(userId int, OrgId int) error {

	trnx, err := database.DBConn.Begin(context.TODO())
	if err != nil {
		return err
	}

	_, err = trnx.Exec(context.TODO(), `
		INSERT INTO ORG_MEMBERS(ORG_ID, USER_ID)
		VALUES($1, $2)	
	`, OrgId, userId)
	if err != nil {
		return err
	}

	_, err = trnx.Exec(context.TODO(), `
		UPDATE USERS
		SET JOINED_ORG = $1
		WHERE ID = $2
	`, true, userId)
	if err != nil {
		return err
	}

	err = trnx.Commit(context.TODO())
	if err != nil {
		return err
	}

	return nil

}

func GetOrgId(orgName string) (int, error) {
	var id int
	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM ORGS
		WHERE ORG_NAME = $1	
	`, orgName).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func GerOrgIdByOwnerId(ownerID int) (int, error) {
	var id int
	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM ORGS
		WHERE OWNER_ID = $1	
	`, ownerID).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func GetAllPlansFormDB() (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.TODO(), `
		SELECT * FROM PLAN	
	`)

	if err != nil {
		return rows, err
	}

	return rows, nil
}

func UserRoleFromDB(userId int) (string, error) {
	var roleName string

	err := database.DBConn.QueryRow(context.TODO(), `
		Select role_name From user_roles
		inner join roles
		on user_roles.role_id = roles.id
		where user_id = $1
	`, userId).Scan(&roleName)

	if err != nil {
		return roleName, err
	}

	return roleName, nil
}

func GetPlanDetails(planName string) (models.Plans, error) {

	var planDetails models.Plans

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT * FROM PLAN	
		WHERE PLAN_NAME = $1
	`, planName).Scan(&planDetails.Id, &planDetails.PlanName, &planDetails.Price, &planDetails.TotalUses)

	if err != nil {
		return planDetails, err
	}

	return planDetails, nil
}
func GetPlanDetailsByplanId(planId int) (models.Plans, error) {

	var planDetails models.Plans

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT * FROM PLAN	
		WHERE id = $1
	`, planId).Scan(&planDetails.Id, &planDetails.PlanName, &planDetails.Price, &planDetails.TotalUses)

	if err != nil {
		return planDetails, err
	}

	return planDetails, nil
}

func FindOrgIDByuserID(userId int) (int, error) {

	var orgId int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ORG_ID FROM ORG_MEMBERS
		WHERE USER_ID = $1	
	`, userId).Scan(&orgId)

	if err != nil {
		return orgId, err
	}

	return orgId, nil
}

func ReductTheUseCount(orgId int) (int, error) {
	var usesLeft int
	err := database.DBConn.QueryRow(context.TODO(), `
		UPDATE SUBSCRIPTIONS
		SET NO_OF_TIMES_USED = NO_OF_TIMES_USED - 1
		WHERE ORG_ID = $1
		
		RETURNING NO_OF_TIMES_USED
	`, orgId).Scan(&usesLeft)

	if err != nil {
		return usesLeft, err
	}

	return usesLeft, nil
}

func UserJoinedOrg(userId int) (bool, error) {

	var joined bool

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT JOINED_ORG FROM USERS
		WHERE ID = $1	
	`, userId).Scan(&joined)

	if err != nil {
		return joined, err
	}

	return joined, nil

}

func GetUserRoleFromDb(userId int) (string, error) {
	var role string

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT role_name FROM USER_ROLES
		INNER JOIN ROLES
		ON USER_ROLES.ROLE_ID = ROLES.ID
		WHERE USER_ID = $1;	
	`, userId).Scan(&role)

	if err != nil {
		return role, err
	}

	return role, nil

}



func GetAllMembersInOrgDB(orgId int, limit int, offSet int) (pgx.Rows, error) {

	rows, err := database.DBConn.Query(context.TODO(), `
		SELECT USER_ID, NAME, EMAIL FROM ORG_MEMBERS
		INNER JOIN USERS
		ON ORG_MEMBERS.USER_ID = USERS.ID
		WHERE ORG_ID = $1
		LIMIT $2 OFFSET $3;
	`, orgId, limit, offSet)

	if err != nil {
		return rows, err
	}

	return rows, nil
}

func GetOrgIdByOwnerId(ownerId int) (int, error) {

	var orgId int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM ORGS
		WHERE OWNER_ID = $1	
	`, ownerId).Scan(&orgId)

	if err != nil {
		return orgId, err
	}

	return orgId, nil
}

func GetCurrentPlan(ownerId int) (string, error) {

	var currentPlan string

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT PLAN_NAME FROM SUBSCRIPTIONS
		INNER JOIN ORGS
		ON SUBSCRIPTIONS.ORG_ID = ORGS.ID
		INNER JOIN PLAN
		ON SUBSCRIPTIONS.PLAN_iD = PLAN.ID
		WHERE OWNER_ID = $1
	`, ownerId).Scan(&currentPlan)

	if err != nil {
		return currentPlan, err
	}

	return currentPlan, nil
}

func ChangePlan(orgId int, newPlanID int) (models.NewPlan, error) {

	var newPlan models.NewPlan

	planDetails, err := GetPlanDetailsByplanId(newPlanID)
	if err != nil {
		return newPlan, err
	}

	_, err = database.DBConn.Exec(context.TODO(), `
		UPDATE SUBSCRIPTIONS
		SET PLAN_ID = $1, PURCHASED_AT = NOW(), TOTAL_USES = $2, NO_OF_TIMES_USED = $3
		WHERE org_Id = $4
	`, planDetails.Id, planDetails.TotalUses, planDetails.TotalUses, orgId)

	if err != nil {
		return newPlan, err
	}

	newPlan.TotalUses = planDetails.TotalUses

	return newPlan, nil
}

func GetPlanId(planName string) (int, error) {
	var id int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT ID FROM PLAN
		WHERE PLAN_NAME = $1	
	`, planName).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func GetNumberOfUsesLeft(userId int) (int, error) {

	var usesLeft int

	err := database.DBConn.QueryRow(context.TODO(), `
		SELECT NO_OF_TIMES_USED FROM ORG_MEMBERS
		INNER JOIN SUBSCRIPTIONS
		ON ORG_MEMBERS.ORG_ID = SUBSCRIPTIONS.ORG_ID
		WHERE USER_ID = $1
	`, userId).Scan(&usesLeft)

	if err != nil {
		return usesLeft, err
	}

	return usesLeft, nil
}


func GetAllUsers(limit int, offSet int) (pgx.Rows, error) {
	rows, err := database.DBConn.Query(context.TODO(), `
		SELECT ID, NAME, EMAIL FROM USERS
		ORDER BY ID
		LIMIT $1 OFFSET $2;
	`, limit, offSet)

	if err != nil{
		return rows, err
	}

	return rows, nil
}

func GetSubscriptionsFromDb() (pgx.Rows, error) {
	rows, err := database.DBConn.Query(context.TODO(), `
		SELECT SUBSCRIPTIONS.ID, ORG_ID, PLAN_NAME, PURCHASED_AT, SUBSCRIPTIONS.TOTAL_USES, NO_OF_TIMES_USED FROM SUBSCRIPTIONS
		INNER JOIN PLAN
		ON SUBSCRIPTIONS.PLAN_ID = PLAN.ID
	`)

	if err != nil {
		return rows, err
	}

	return rows, nil
}
func GetSubscriptionsFromDbWithFilter(search string) (pgx.Rows, error) {
	rows, err := database.DBConn.Query(context.TODO(), `
		SELECT SUBSCRIPTIONS.ID, ORG_ID, PLAN_NAME, PURCHASED_AT, SUBSCRIPTIONS.TOTAL_USES, NO_OF_TIMES_USED FROM SUBSCRIPTIONS
		INNER JOIN PLAN
		ON SUBSCRIPTIONS.PLAN_ID = PLAN.ID
		WHERE PLAN_NAME = $1
	`, search)

	if err != nil {
		return rows, err
	}

	return rows, nil
}