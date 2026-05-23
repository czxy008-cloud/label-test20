package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"

	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	secret   string
}

func NewAuthService(userRepo *repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		secret:   secret,
	}
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	hashedPassword := hashPassword(req.Password)
	if user.Password != hashedPassword {
		return nil, errors.New("用户名或密码错误")
	}

	token := s.generateToken(user.ID)
	if err := s.userRepo.UpdateToken(ctx, user.ID, token); err != nil {
		return nil, err
	}

	user.Token = token
	return &model.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*model.User, error) {
	user, err := s.userRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("无效的令牌")
	}
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}
	return user, nil
}

func (s *AuthService) generateToken(userID int64) string {
	uuidStr := uuid.New().String()
	token := s.secret + ":" + uuidStr
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
