package main

import (
	
	"fmt"
	"net/http"
	"subscribly/database"
	"subscribly/handlers"
	"subscribly/middleware"
	
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {

	database.InitializePool()
	r := chi.NewRouter()
	fmt.Println("server is starting......")

	r.Use(chimiddleware.Logger)

	r.Get("/healthz", handlers.CheckHealth)

	r.With(middleware.SignUpValidator).Post("/signup", handlers.UserSignUp)
	r.With(middleware.LoginValidator).Post("/login", handlers.UserLogin)
	r.With(middleware.JwtMiddleware).Post("/logout", handlers.UserLogOut)

	r.With(middleware.Pagnation).Get("/users", handlers.GetAllUsers)

	r.With(middleware.JwtMiddleware).Post("/orgs", handlers.CreateNewOrg)

	r.With(middleware.JwtMiddleware).With(middleware.Pagnation).Get("/orgs", handlers.GetAllOrgs)

	r.With(middleware.JwtMiddleware).Post("/orgs/join/{orgName}", handlers.JoinOrg)

	r.With(middleware.JwtMiddleware).With(middleware.OwnerMiddleware).Get("/plans", handlers.GetAllPlans)

	r.With(middleware.JwtMiddleware).With(middleware.NewOrgValidator).Post("/plan/upgrade/{planName}", handlers.UpgradePlan)

	r.Post("/refresh", handlers.GenerateNewAccessToken)

	r.With(middleware.JwtMiddleware).With(middleware.AdminMiddleware).With(middleware.SearchPagnation).Get("/subscriptions", handlers.GetAllSubscriptions)

	r.With(middleware.JwtMiddleware).Post("/use", handlers.Use)
	r.With(middleware.JwtMiddleware).Get("/use/left", handlers.UsesLeft)

	r.With(middleware.JwtMiddleware).With(middleware.OwnerMiddleware).With(middleware.Pagnation).Get("/orgs/members", handlers.GetAllMembersInOrg)

	http.ListenAndServe(":5000", r)
}