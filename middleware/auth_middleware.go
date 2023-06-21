package middleware

import (
	"net/http"
	"projek-1/model/web"
	"projek-1/helper"
)

type AuthMiddleware struct {
	Hadhler http.Handler
}

func NewAuthMiddleware(handhler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if "rahasia" == request.Header.Get("x-api-key") {
		//ok
		middleware.Hadhler.ServeHTTP(writer, request)
	} else {
		//error
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:	http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
	
}
