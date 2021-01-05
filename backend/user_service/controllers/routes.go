package controllers

import (
	"github.com/khoatxp/filmchon/backend/user_service/middlewares"
)

func (s *Server) initializeRoutes() {

	router := s.Router.Group("/user_service")
	{
		// Login Route
		router.POST("/login", s.Login)

		// Reset password
		router.POST("/password/forgot", s.ForgotPassword)
		router.POST("/password/reset", s.ResetPassword)

		//Users routes
		router.POST("/users", s.CreateUser)
		//FOR ADMIN
		// router.GET("/users", s.GetUsers)
		// router.GET("/users/:id", s.GetUser)
		router.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		router.PUT("/avatar/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateAvatar)
		router.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)
	}
}
