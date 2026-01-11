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

	r.Post("/signup", handlers.UserSignUp)
	r.Post("/login", handlers.UserLogin)
	r.With(middleware.JwtMiddleware).Post("/logout", handlers.UserLogOut)

	r.With(middleware.Pagnation).Get("/users", handlers.GetAllUsers)

	r.With(middleware.JwtMiddleware).Post("/orgs", handlers.CreateNewOrg)

	r.With(middleware.JwtMiddleware).With(middleware.Pagnation).Get("/orgs", handlers.GetAllOrgs)

	r.With(middleware.JwtMiddleware).Post("/orgs/join/{orgName}", handlers.JoinOrg)

	r.With(middleware.JwtMiddleware).With(middleware.RoleMiddleware).Get("/plans", handlers.GetAllPlans)

	r.With(middleware.JwtMiddleware).Post("/plan/upgrade/{planName}", handlers.UpgradePlan)

	r.Post("/refresh", handlers.GenerateNewAccessToken)

	// r.Get("/subscriptions", handlers.GetAllSubscriptions)
	r.With(middleware.JwtMiddleware).Post("/use", handlers.Use)
	r.With(middleware.JwtMiddleware).Get("/use/left", handlers.UsesLeft)

	r.With(middleware.JwtMiddleware).With(middleware.RoleMiddleware).With(middleware.Pagnation).Get("/orgs/members", handlers.GetAllMembersInOrg)

	http.ListenAndServe(":5000", r)
}
