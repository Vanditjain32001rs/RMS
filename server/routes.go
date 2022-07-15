package server

import (
	"RMS/handlers"
	"RMS/helpers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetUpRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/user", func(api chi.Router) {
		api.Post("/sign-in", handlers.SignIn)
		api.Post("/sign-up", handlers.SignUp)

		api.Route("/task", func(r chi.Router) {
			r.Use(helpers.AuthMiddleware)
			r.Post("/add-user", handlers.AddNewUser)
			r.Post("/add-subadmin", handlers.AddSubAdmins)
			r.Post("/add-user-location", handlers.AddNewLocation)
			r.Get("/fetch-users", handlers.FetchUsers)
			r.Get("/fetch-all-subadmin", handlers.FetchAllSubAdmins)
			r.Post("/add-role", handlers.AddRole)
			r.Get("/distance", handlers.GetDistance)

			r.Route("/restaurant", func(rest chi.Router) {
				r.Post("/add-restaurant", handlers.AddRestaurant)
				r.Post("/add-dish", handlers.AddDish)
				r.Get("/fetch-restaurant-list", handlers.GetRestaurantList)
				r.Get("/fetch-restaurant-dish", handlers.GetRestaurantDishList)
			})
		})
	})

	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
