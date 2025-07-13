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
		AssetsRoutes(api)
	})
}

func AuthRoutes(r chi.Router) {
	r = r.Route("/auth", func(auth chi.Router) {
		auth.Post("/login", handlers.Login)
		auth.Get("/refresh_token", handlers.RefreshToken)

	})
}

func UsersRoutes(r chi.Router) { //routes for reading here
	r.Route("/users", func(users chi.Router) {
		users.Use(middleware.AuthMiddleware())
		users.Use(middleware.RequireRoles("admin", "employee_manager"))
		users.Get("/", handlers.ListUsers)
		users.Post("/", handlers.CreateUser)
		users.Patch("/{user_id}", handlers.UpdateUser)

	})

}

func AssetsRoutes(r chi.Router) { //routes for reading here
	r.Route("/asset", func(asset chi.Router) {
		asset.Use(middleware.AuthMiddleware())
		asset.Use(middleware.RequireRoles("admin", "asset_manager"))
		//asset.Get("/", handlers.ListUsers)
		asset.Post("/", handlers.CreateAsset)
		asset.Get("/", handlers.ListAssets)
		//asset.Patch("/{user_id}", handlers.UpdateUser)

	})

}
