package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserLoginCreds struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserDetails struct {
	Id        int       `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	JoinedOrg bool 	`json:"joined_org,omitempty"`
}


type Org struct {
	Id        int       `json:"id"`
	OrgName   string    `json:"org_name"`
	OwnerId   int       `json:"Owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Claims struct {
	UserId int `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Plans struct {
	Id        int    `json:"id"`
	PlanName  string `json:"plan_name"`
	Price     int    `json:"price"`
	TotalUses int    `json:"total_uses"`
}

type OrgWithPlans struct {
	OrgId       int       `json:"org_id"`
	OrgName     string    `json:"org_name"`
	OwnerId     int       `json:"Owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	PlanId      int       `json:"plan_id"`
	PlanName    string    `json:"plan_name"`
	Price       int       `json:"price"`
	PurchasedAt time.Time `json:"purchased_at"`
	TotalUses   int       `json:"total_uses"`
	UsesLeft    int       `json:"uses_left"`
}

type ReminingCounts struct {
	UsesLeft int `json:"uses_left"`
}

type Members struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type NewPlan struct {
	Message string `json:"message"`
	PlanName string `json:"plan_name"`
	TotalUses int `json:"total_uses"`
}

type UsesLeft struct {
	NoOfUsesLeft int `json:"uses_left"`
}

type PageLimit struct{
	PageNumber int `json:"page"`
	Limit int `json:"limit"`
	Search string `json:"search"`
}

type Subscriptions struct {
	ID int `json:"id"`
	OrgID int `json:"org_id"`
	PlanName string `json:"Plan_name"`
	PurchasedAt time.Time `json:"purchased_at"`
	TotalUses int `json:"total_uses"`
	UsesLeft int `json:"uses_left"`
}