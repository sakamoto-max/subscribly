package main

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"
// )

// func main() {
// 	fmt.Println("server is startin..........")
// 	// Create an HTTP server that listens on port 8000
// 	http.ListenAndServe(":5000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		// This prints to STDOUT to show that processing has started
// 		fmt.Fprint(os.Stdout, "processing request\n")
// 		// We use `select` to execute a piece of code depending on which
// 		// channel receives a message first
// 		select {
// 		case <-time.After(2 * time.Second):
// 			// If we receive a message after 2 seconds
// 			// that means the request has been processed
// 			// We then write this as the response
// 			w.Write([]byte("request processed"))
// 		case <-ctx.Done():
// 			// If the request gets cancelled, log it
// 			// to STDERR
// 			fmt.Fprint(os.Stderr, "request cancelled\n")
// 		}
// 	}))
// }

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

	r.With(middleware.JwtMiddleware).With(middleware.RoleMiddleware).Get("/plans", handlers.GetAllPlans)

	r.With(middleware.JwtMiddleware).With(middleware.NewOrgValidator).Post("/plan/upgrade/{planName}", handlers.UpgradePlan)

	r.Post("/refresh", handlers.GenerateNewAccessToken)

	// r.Get("/subscriptions", handlers.GetAllSubscriptions)
	r.With(middleware.JwtMiddleware).Post("/use", handlers.Use)
	r.With(middleware.JwtMiddleware).Get("/use/left", handlers.UsesLeft)

	r.With(middleware.JwtMiddleware).With(middleware.RoleMiddleware).With(middleware.Pagnation).Get("/orgs/members", handlers.GetAllMembersInOrg)

	http.ListenAndServe(":5000", r)
}
