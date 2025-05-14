package account

import (
	"context"
	"encoding/json"
	KiteError "kite/internal/errors"
	"kite/internal/models"
	"kite/internal/repositories/accounts"
	"slices"
)

type UserService interface {
	FindByUserId(ctx context.Context, userId uint) (*models.Users, error)
	PushOneLabelForUser(ctx context.Context, user *models.Users, label string) error
	DeleteOneLabelForUser(ctx context.Context, user *models.Users, label string) error
}

type userService struct {
	repo accounts.UserRepository
}

func NewUserService(repo accounts.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) PushOneLabelForUser(ctx context.Context, user *models.Users, label string) error {
	// 反序列化 label
	var labels []string
	if err := json.Unmarshal(user.Label, &labels); err != nil {
		return KiteError.New(KiteError.UnmarshalError, err)
	}
	// 检查标签是否存在, 如果存在，直接返回
	if slices.Contains(labels, label) {
		return nil
	}
	// 序列化后持久化
	labels = append(labels, label)
	newLabels, err := json.Marshal(labels)
	if err != nil {
		return KiteError.New(KiteError.MarshalError, err)
	}
	return s.repo.UpdateLabel(ctx, user.UserId, newLabels)
}

func (s *userService) DeleteOneLabelForUser(ctx context.Context, user *models.Users, label string) error {
	// 反序列化 label
	var labels []string
	if err := json.Unmarshal(user.Label, &labels); err != nil {
		return KiteError.New(KiteError.UnauthorizedError, err)
	}
	// 检查标签是否存在
	index := slices.Index(labels, label)
	if index == -1 {
		return nil
	}
	labels = slices.Delete(labels, index, index+1)
	// 序列化 Label
	newLabel, err := json.Marshal(labels)
	if err != nil {
		return KiteError.New(KiteError.MarshalError, err)
	}
	// 保存 label
	return s.repo.UpdateLabel(ctx, user.UserId, newLabel)
}

func (s *userService) FindByUserId(ctx context.Context, userId uint) (*models.Users, error) {
	return s.repo.FindByUserId(ctx, userId)
}
