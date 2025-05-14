package account

import (
	"context"
	"encoding/json"
	KiteError "kite/internal/errors"
	"kite/internal/models"
	"kite/internal/repositories/accounts"
)

type UserService interface {
	FindByUserId(ctx context.Context, userId uint) (*models.Users, error)
	DeleteOneLabelForUser(ctx context.Context, user *models.Users, label string) error
}

type userService struct {
	repo accounts.UserRepository
}

func NewUserService(repo accounts.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) DeleteOneLabelForUser(ctx context.Context, user *models.Users, label string) error {
	// 反序列化 label
	var labels []string
	if err := json.Unmarshal([]byte(label), &labels); err != nil {
		return KiteError.New(KiteError.DataError, err)
	}
	// 检查标签是否存在
	exists := false
	newLabels := make([]string, 0, len(labels))
	for _, l := range labels {
		if l != label {
			newLabels = append(newLabels, l)
		} else {
			exists = true
		}
	}
	// 不存在直接返回
	if !exists {
		return nil
	}
	// 序列化 Label
	newLabelJson, err := json.Marshal(newLabels)
	if err != nil {
		return KiteError.New(KiteError.MarshalError, err)
	}
	// 保存 label
	return s.repo.UpdateLabel(ctx, user.UserId, newLabelJson)
}

func (s *userService) FindByUserId(ctx context.Context, userId uint) (*models.Users, error) {
	return s.repo.FindByUserId(ctx, userId)
}
