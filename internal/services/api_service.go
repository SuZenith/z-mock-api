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
