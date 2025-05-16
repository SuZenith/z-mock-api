package repositories

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"kite/internal/api/payloads"
	"kite/internal/database"
	KiteError "kite/internal/errors"
	"kite/internal/models"
)

type ApiRepository interface {
	CreateApi(ctx context.Context, payload payloads.MockApiPayload, uuid string) error
	QueryApiWithUidAndPathAndMethod(ctx context.Context, uid string, path string, method string) (*models.Api, error)
}

type apiRepository struct {
	db *gorm.DB
}

func NewApiRepository(connection *database.MySQLConnection) ApiRepository {
	return &apiRepository{db: connection.GetDB()}
}

func (r *apiRepository) CreateApi(ctx context.Context, payload payloads.MockApiPayload, uuid string) error {
	var headers json.RawMessage
	headers, err := payload.GetHeadersJSON()
	if err != nil {
		return KiteError.New(KiteError.MarshalError, err)
	}
	api := &models.Api{
		Uuid:         uuid,
		UserId:       payload.UserId,
		Path:         payload.Path,
		Method:       payload.Method,
		StatusCode:   payload.StatusCode,
		ContentType:  payload.ContentType,
		Headers:      headers,
		ResponseBody: payload.ResponseBody,
	}
	result := r.db.WithContext(ctx).Create(api)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *apiRepository) QueryApiWithUidAndPathAndMethod(ctx context.Context, uid string, path string, method string) (*models.Api, error) {
	var api *models.Api
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND path = ? AND method = ?", uid, path, method).
		First(&api)
	if result.Error != nil {
		return nil, result.Error
	}
	return api, nil
}
