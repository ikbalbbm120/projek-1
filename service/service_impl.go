package service

import (
	"projek-1/repository"
	"projek-1/model/web"
	"projek-1/model/domain"
	"projek-1/helper"
	"projek-1/exception"
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
)

type ServiceImpl struct {
	Repository repository.Repository
	DB		*sql.DB
	Validate	*validator.Validate
}

func NewService(CategoryRepository repository.Repository, DB *sql.DB, validate *validator.Validate) Service {
	return &ServiceImpl {
		Repository: CategoryRepository,
		DB: 				DB,
		Validate: 	  validate,
	}
}

func (service *ServiceImpl) Create(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	category := domain.Category {

		
	}

	category = service.Repository.Save(ctx, tx, category)

	return helper.ToCategoryResponse(category)
}

func (service *ServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	category, err := service.Repository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	category.Name = request.Name

	category = service.Repository.Update(ctx, tx, category)

	return helper.ToCategoryResponse(category)
}

func (service *ServiceImpl) Delete(ctx context.Context, categoryid int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	category, err := service.Repository.FindById(ctx, tx, categoryid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	service.Repository.Delete(ctx, tx, category)
}

func (service *ServiceImpl)FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	category, err := service.Repository.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCategoryResponse(category)
}

func (service *ServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	categories := service.Repository.FindAll(ctx, tx)

	return helper.ToCategoryResponses(categories)
}