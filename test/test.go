package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	"projek-1/controller"
	"projek-1/model/domain"
	"projek-1/repository"
	"projek-1/service"
	"projek-1/middleware"
	"projek-1/helper"
	"projek-1/app"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3309)/projek 1")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	Repository := repository.NewRepository()
	Service := service.NewService(Repository, db, validate)
	Controller := controller.NewController(Service)
	router := app.NewRouter(Controller)
	return middleware.NewAuthMiddleware(&router)
}

func TruncateCategory(db *sql.DB) {
	db.Exec("truncate category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)
	router := setupRouter(db)
	requestBody := strings.NewReader(`{"name": "gadget test"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("content-type", "application.json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "ok", responseBody["status"])
	assert.Equal(t, "gadget test", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)
	router := setupRouter(db)
	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")
	
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "bad request", responseBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category {
		Name: "gadget",
	})
	tx.Commit()

	router := setupRouter(db)
	requestBody := strings.NewReader(`{"name": "gadget update"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "ok", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "gadget update", responseBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db :=  setupTestDB()
	TruncateCategory(db)
	router := setupRouter(db)
	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http:localhost:3000/api/categories", requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "bad request", responseBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category {
		Name: "gadget",
	})
	tx.Commit()

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("content-type", "appliation/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "ok", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)
	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http:localhost:3000/api/categories/404", nil)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "not found", responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category {
		Name: "gadget",
	})
	tx.Commit()

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "ok", responseBody["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)
	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map [string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "not found",responseBody["status"])
}

func TestGetlistCategorySuccess(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)
	
	tx, _ := db.Begin()
	categoryRepository := repository.NewRepository()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category {
		Name: "gadget",
	})
	category2 := categoryRepository.Save(context.Background(), tx, domain.Category {
		Name: "computer",
	})
	tx.Commit()

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "rahasia")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "ok", responseBody["status"])

	var categories = responseBody["data"].([] interface{})
	
	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])

	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	TruncateCategory(db)

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-api-key", "salah")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "unathorized", responseBody["status"])
}