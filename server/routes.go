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
	router.Route("/rms", func(api chi.Router) {
		api.Post("/sign-in", handlers.SignIn)
		api.Post("/sign-up", handlers.SignUp)

		api.Route("/task", func(r chi.Router) {
			r.Use(helpers.AuthMiddleware)

			r.Route("/admins", func(admin chi.Router) {
				admin.Use(helpers.AdminPermissionMiddleWare)
				admin.Route("/user", func(user chi.Router) {
					user.Post("/", handlers.AddNewUser)
					user.Get("/", handlers.FetchUsersByAdmin)
				})
				admin.Route("/sub-admin", func(sub chi.Router) {
					sub.Get("/", handlers.FetchAllSubAdmins)
					sub.Post("/", handlers.AddSubAdmins)
				})
				admin.Route("/restaurant", func(res chi.Router) {
					res.Post("/", handlers.AddRestaurant)
				})
				admin.Route("/dish", func(dish chi.Router) {
					dish.Put("/", handlers.UpdateDish)
					dish.Post("/", handlers.AddDish)
				})
			})

			r.Route("/sub-admins", func(sub chi.Router) {
				sub.Use(helpers.SubAdminPermissionMiddleWare)
				sub.Route("/user", func(user chi.Router) {
					user.Post("/", handlers.AddNewUser)
					user.Get("/", handlers.FetchUsersBySubAdmin)
				})
				sub.Route("/restaurant", func(res chi.Router) {
					res.Post("/", handlers.AddRestaurant)
				})
				sub.Route("/dish", func(dish chi.Router) {
					dish.Put("/", handlers.UpdateDish)
					dish.Post("/", handlers.AddDish)
				})
			})

			r.Post("/add-user-location", handlers.AddNewLocation)
			r.Put("/update-user", handlers.UpdateUser)
			r.Post("/add-role", handlers.AddRole)
			r.Get("/distance", handlers.GetDistance)

			r.Route("/restaurant", func(restaurant chi.Router) {
				restaurant.Get("/restaurant-list", handlers.GetRestaurantList)
				restaurant.Get("/dish-list", handlers.GetRestaurantDishList)
			})
		})
	})

	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
