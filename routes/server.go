package routes

import (
	"github.com/go-chi/chi/v5"
	"storex/handlers"
	"storex/middleware"
	//"storex/middleware"
)

func Routes(r chi.Router) {
	r.Route("/api", func(api chi.Router) {
		AuthRoutes(api)
		//RegisterUserRoutes(api)
		//RegisterAssetRoutes(api)
		//RegisterStatusRoutes(api)
		//RegisterServiceRoutes(api)

		UsersRoutes(api)
	})
}

func AuthRoutes(r chi.Router) {
	r = r.Route("/auth", func(auth chi.Router) {
		auth.Post("/login", handlers.Login)
	})
}

func UsersRoutes(r chi.Router) { //routes for reading here
	r.Route("/users", func(users chi.Router) {
		users.Use(middleware.AuthMiddleware())
		users.Use(middleware.RequireRoles("admin", "employee_manager"))
		users.Get("/", handlers.ListUsers)
		users.Post("/", handlers.CreateUser)

	})

}
