package services

import (
	uuid2 "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kite/internal/api/payloads"
	KiteError "kite/internal/errors"
	"kite/internal/repositories"
)

type ApiService interface {
	Create(ctx echo.Context, payload payloads.MockApiPayload) error
	Request(ctx echo.Context, uid string, path string, method string) (contextType *string, responseBody *string, err error)
}

type apiService struct {
	repo repositories.ApiRepository
}

func NewApiService(repo repositories.ApiRepository) ApiService {
	return &apiService{repo}
}

func (s *apiService) Create(ctx echo.Context, payload payloads.MockApiPayload) error {
	uuid := uuid2.NewString()
	err := s.repo.CreateApi(ctx.Request().Context(), payload, uuid)
	if err != nil {
		return KiteError.New(KiteError.ApiCreateError, err)
	}
	return nil
}

func (s *apiService) Request(ctx echo.Context, uid string, path string, method string) (contextType *string, responseBody *string, err error) {
	api, err := s.repo.QueryApiWithUidAndPathAndMethod(ctx.Request().Context(), uid, path, method)
	if err != nil {
		return nil, nil, KiteError.New(KiteError.InternalServerError, err)
	}
	return &api.ContentType, &api.ResponseBody, nil
}
