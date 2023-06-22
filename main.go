package main

import (
	"net/http"
	"projek-1/controller"
	"projek-1/repository"
	"projek-1/service"
	"projek-1/middleware"
	"projek-1/helper"
	"projek-1/app"
	"github.com/go-playground/validator/v10"
)

func main () {
	db := app.NewDB()
	validate := validator.New()
	Repository := repository.NewRepository()
	Service := service.NewService(Repository, db, validate)
	Controller := controller.NewController(Service)

	router := app.NewRouter(Controller)

	server := http.Server {
		Addr: "localhost:3000",

		Handler: middleware.NewAuthMiddleware(&router),
	}
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}