package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/repository"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
)

// AuthService 認證服務接口
type AuthService interface {
	Register(ctx context.Context, req *model.UserRegisterRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.UserLoginRequest) (*model.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*model.UserResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	GetUserByID(ctx context.Context, id string) (*model.UserResponse, error)
}

type authService struct {
	userRepo    repository.UserRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewAuthService 創建認證服務實例
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

// Register 用戶註冊
func (s *authService) Register(ctx context.Context, req *model.UserRegisterRequest) (*model.UserResponse, error) {
	// 檢查用戶名是否已存在
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 檢查郵箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 創建新用戶
	user := &model.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 加密密碼
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// 保存用戶
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Login 用戶登錄
func (s *authService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.LoginResponse, error) {
	// 獲取用戶
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 驗證密碼
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	// 生成訪問令牌
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken 驗證令牌
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*model.UserResponse, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// 獲取聲明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// 檢查過期時間
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// 獲取用戶 ID
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	// 獲取用戶信息
	return s.GetUserByID(ctx, userID)
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	// 解析刷新令牌
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	// 獲取用戶
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 生成新的訪問令牌
	newToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新令牌
	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: newToken,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		RefreshToken: newRefreshToken,
	}, nil
}

// GetUserByID 通過ID獲取用戶
func (s *authService) GetUserByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// generateToken 生成訪問令牌
func (s *authService) generateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.tokenExpiry).Unix(),
		"iat": time.Now().Unix(),
		"role": user.Role,
	})

	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken 生成刷新令牌
func (s *authService) generateRefreshToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.tokenExpiry * 24).Unix(), // 刷新令牌有效期更長
		"iat": time.Now().Unix(),
	})

	return token.SignedString(s.jwtSecret)
}
