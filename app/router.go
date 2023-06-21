package app

import (
	"projek-1/controller"
	"projek-1/exception"
	

	"github.com/julienschmidt/httprouter"
)

func NewRouter (Controller controller.Controller) httprouter.Router {
	router := httprouter.New()

	router.GET("/api/categories", Controller.FindAll)
	router.GET("/api/categories/:categoryId", Controller.FindById)
	router.POST("/api/categories", Controller.Create)
	router.PUT("/api/categories/:categoryId", Controller.Update)
	router.DELETE("/api/categories/:categoryId", Controller.Delete)

	router.PanicHandler = exception.ErrorHandhler
	return *router
}