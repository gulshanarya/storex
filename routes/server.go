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

func UsersRoutes(r chi.Router) {
	r.Route("/users", func(users chi.Router) {
		users.Use(middleware.AuthMiddleware())
		//users.Get("/{user_id}", handlers.GetAssetsByUser)
		users.Use(middleware.RequireRoles("admin", "employee_manager"))
		users.Get("/", handlers.ListUsers)
		users.Post("/", handlers.CreateUser)
		users.Patch("/{user_id}", handlers.UpdateUser)
		users.Delete("/{user_id}", handlers.DeleteUser)
	})
}

func AssetsRoutes(r chi.Router) {
	r.Route("/asset", func(asset chi.Router) {
		asset.Use(middleware.AuthMiddleware())
		asset.Use(middleware.RequireRoles("admin", "asset_manager"))
		asset.Post("/", handlers.CreateAsset)
		asset.Get("/", handlers.ListAssets)
		asset.Patch("/{id}", handlers.UpdateAsset)
		asset.Post("/assign", handlers.AssignAsset)
		asset.Patch("/retrieve/{asset_id}", handlers.RetrieveAsset)
		asset.Get("/timeline", handlers.AssetTimeline)
		asset.Get("/user/timeline", handlers.UserAssetTimeline)
	})

}
