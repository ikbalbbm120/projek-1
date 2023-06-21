package controller

import (
	"net/http"
	"projek-1/helper"
	"projek-1/model/web"
	"projek-1/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ControllerImpl struct {
	Service service.Service
}

func NewController(service service.Service) Controller {
	return &ControllerImpl {
		Service: service,
	}
}

func (Controller *ControllerImpl) Create (writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryCreateRequest := web.CategoryCreateRequest{}
	helper.ReadFromRequestBody(request, &categoryCreateRequest)

	categroyResponse := Controller.Service.Create(request.Context(), categoryCreateRequest)
	webResponse := web.WebResponse {
		Code: 200,
		Status: "ok",
		Data: categroyResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func(controller *ControllerImpl) Update (writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	CategoryUpdateRequest := web.CategoryUpdateRequest{}
	helper.ReadFromRequestBody(request, &CategoryUpdateRequest)

	CategoryUpdateRequest.Id = id

	categoryResponse := controller.Service.Update(request.Context(), CategoryUpdateRequest)
	webResponse := web.WebResponse {
		Code: 200,
		Status: "ok",
		Data: categoryResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (Controller *ControllerImpl) Delete (writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)

	Controller.Service.Delete(request.Context(), id)
	webResponse := web.WebResponse {
		Code: 200,
		Status: "ok",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) FindById (writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)

	categoryResponse := controller.Service.FindById(request.Context(), id)
	webResponse := web.WebResponse {
		Code: 200,
		Status: "ok",
		Data: categoryResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) FindAll (writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryResponses := controller.Service.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

