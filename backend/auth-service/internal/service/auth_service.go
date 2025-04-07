package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// IAuthService 定義認證服務接口
type IAuthService interface {
	Register(ctx context.Context, req *model.UserRegisterRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.UserLoginRequest) (*model.LoginResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*model.UserResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	GetUserByID(ctx context.Context, id string) (*model.UserResponse, error)
	CreateAddress(ctx context.Context, userID string, req *model.AddressRequest) (*model.Address, error)
	GetAddresses(ctx context.Context, userID string) ([]model.Address, error)
	UpdateAddress(ctx context.Context, userID string, addressID string, req *model.AddressRequest) (*model.Address, error)
	DeleteAddress(ctx context.Context, userID string, addressID string) error
	GetPreference(ctx context.Context, userID string) (*model.UserPreference, error)
	UpdatePreference(ctx context.Context, userID string, req *model.PreferenceRequest) (*model.UserPreference, error)
	GetAddressByID(ctx context.Context, addressID string) (*model.Address, error)
	ResetPassword(ctx context.Context, tokenString, newPassword string) error
	ForgetPassword(ctx context.Context, emailString, newPassword string) error
	SetDefaultAddress(ctx context.Context, userID string, addressID string) (*model.Address, error)
}

// authService 實現 IAuthService 接口
type authService struct {
	userRepo    repository.IUserRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewAuthService 創建新的認證服務實例
func NewAuthService(repo repository.IUserRepository, jwtSecret string, tokenExpiry time.Duration) IAuthService {
	return &authService{
		userRepo:    repo,
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
		return nil, ErrUsernameTaken
	}

	// 檢查郵箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
		Status:   "active",
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// Login 用戶登錄
func (s *authService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.LoginResponse, error) {
	log.Printf("Attempting login for email: %s", req.Email)

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		log.Printf("No user found with email: %s", req.Email)
		return nil, ErrInvalidCredentials
	}

	log.Printf("Found user with email: %s, checking password", req.Email)
	if !user.CheckPassword(req.Password) {
		log.Printf("Invalid password for user: %s", req.Email)
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	response := user.ToResponse()
	log.Printf("Login successful for user: %s", req.Email)
	return &model.LoginResponse{
		User:         response,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken 驗證令牌
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*model.UserResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
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

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
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

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	newToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &model.LoginResponse{
		User:         response,
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetUserByID 通過ID獲取用戶
func (s *authService) GetUserByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// generateToken 生成訪問令牌
func (s *authService) generateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(s.tokenExpiry).Unix(),
		"role": user.Role,
	})

	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken 生成刷新令牌
func (s *authService) generateRefreshToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.tokenExpiry * 24).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

// CreateAddress 創建地址
func (s *authService) CreateAddress(ctx context.Context, userID string, req *model.AddressRequest) (*model.Address, error) {
	// 如果新地址設為預設，將該用戶的所有其他地址設為非預設
	if req.IsDefault {
		addresses, err := s.userRepo.GetAddresses(ctx, userID)
		if err != nil {
			return nil, err
		}

		// 將所有地址的預設狀態設為false
		for _, addr := range addresses {
			if addr.IsDefault {
				addr.IsDefault = false
				if err := s.userRepo.UpdateAddress(ctx, &addr); err != nil {
					return nil, err
				}
			}
		}
	}

	address := &model.Address{
		ID:         uuid.New().String(),
		UserID:     userID,
		Name:       req.Name,
		Phone:      req.Phone,
		Street:     req.Street,
		City:       req.City,
		District:   req.District,
		PostalCode: req.PostalCode,
		IsDefault:  req.IsDefault,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.userRepo.CreateAddress(ctx, address); err != nil {
		return nil, err
	}

	return address, nil
}

// GetAddresses 獲取地址列表（預設地址排在最前面）
func (s *authService) GetAddresses(ctx context.Context, userID string) ([]model.Address, error) {
	addresses, err := s.userRepo.GetAddresses(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 排序：將預設地址放在最前面
	if len(addresses) > 1 {
		// 尋找預設地址並移至第一位
		var defaultIndex int = -1
		for i, addr := range addresses {
			if addr.IsDefault {
				defaultIndex = i
				break
			}
		}

		// 如果找到預設地址且不是第一個，則調整順序
		if defaultIndex > 0 {
			defaultAddr := addresses[defaultIndex]
			// 將預設地址前的所有地址往後移一位
			for i := defaultIndex; i > 0; i-- {
				addresses[i] = addresses[i-1]
			}
			// 將預設地址放到第一位
			addresses[0] = defaultAddr
		}
	}

	return addresses, nil
}

// UpdateAddress 更新地址
func (s *authService) UpdateAddress(ctx context.Context, userID string, addressID string, req *model.AddressRequest) (*model.Address, error) {
	// 先獲取現有地址資訊
	existingAddress, err := s.userRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// 檢查地址是否存在
	if existingAddress == nil {
		return nil, fmt.Errorf("address not found")
	}

	// 檢查地址是否屬於該用戶
	if existingAddress.UserID != userID {
		return nil, fmt.Errorf("address does not belong to the user")
	}

	// 如果要將地址設為預設，且當前不是預設
	if req.IsDefault && !existingAddress.IsDefault {
		// 獲取該用戶的所有地址
		addresses, err := s.userRepo.GetAddresses(ctx, userID)
		if err != nil {
			return nil, err
		}

		// 將所有其他地址設為非預設
		for _, addr := range addresses {
			if addr.IsDefault && addr.ID != addressID {
				addr.IsDefault = false
				if err := s.userRepo.UpdateAddress(ctx, &addr); err != nil {
					return nil, err
				}
			}
		}
	}

	// 更新地址信息
	address := &model.Address{
		ID:         addressID,
		UserID:     userID,
		Name:       req.Name,
		Phone:      req.Phone,
		Street:     req.Street,
		City:       req.City,
		District:   req.District,
		PostalCode: req.PostalCode,
		IsDefault:  req.IsDefault,
		UpdatedAt:  time.Now(),
	}

	// 保留現有的創建時間
	if existingAddress.CreatedAt.Unix() > 0 {
		address.CreatedAt = existingAddress.CreatedAt
	} else {
		address.CreatedAt = time.Now()
	}

	if err := s.userRepo.UpdateAddress(ctx, address); err != nil {
		return nil, err
	}

	return address, nil
}

// DeleteAddress 刪除地址
func (s *authService) DeleteAddress(ctx context.Context, userID string, addressID string) error {
	return s.userRepo.DeleteAddress(ctx, addressID)
}

// GetPreference 獲取用戶偏好
func (s *authService) GetPreference(ctx context.Context, userID string) (*model.UserPreference, error) {
	return s.userRepo.GetPreference(ctx, userID)
}

// UpdatePreference 更新用戶偏好
func (s *authService) UpdatePreference(ctx context.Context, userID string, req *model.PreferenceRequest) (*model.UserPreference, error) {
	pref := &model.UserPreference{
		UserID:            userID,
		Language:          req.Language,
		Currency:          req.Currency,
		NotificationEmail: req.NotificationEmail,
		NotificationSMS:   req.NotificationSMS,
		Theme:             req.Theme,
		UpdatedAt:         time.Now(),
	}

	if err := s.userRepo.UpdatePreference(ctx, pref); err != nil {
		return nil, err
	}

	return pref, nil
}

// GetAddressByID 獲取地址
func (s *authService) GetAddressByID(ctx context.Context, addressID string) (*model.Address, error) {
	address, err := s.userRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	return address, nil
}

// ForgetPassword 忘記密碼
func (s *authService) ForgetPassword(ctx context.Context, emailString, newPassword string) error {
	// 透過email找到userID
	user, err := s.userRepo.GetByEmail(ctx, emailString)
	if err != nil {
		return err
	}

	// 如果找不到用戶
	if user == nil {
		return ErrUserNotFound
	}

	// 更新密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword 重置密碼
func (s *authService) ResetPassword(ctx context.Context, tokenString, newPassword string) error {
	// 1. 驗證 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	// 從 token 中獲取用戶 ID
	userID, ok := claims["sub"].(string)
	if !ok {
		return ErrInvalidToken
	}

	// 2. 更新密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

// SetDefaultAddress 設置預設地址
func (s *authService) SetDefaultAddress(ctx context.Context, userID string, addressID string) (*model.Address, error) {
	// 檢查地址是否存在
	address, err := s.userRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, fmt.Errorf("address not found")
	}

	// 檢查地址是否屬於該用戶
	if address.UserID != userID {
		return nil, fmt.Errorf("address does not belong to the user")
	}

	// 如果已經是預設地址，則無需操作
	if address.IsDefault {
		return address, nil
	}

	// 獲取該用戶的所有地址
	addresses, err := s.userRepo.GetAddresses(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 將所有其他地址設為非預設
	for _, addr := range addresses {
		if addr.IsDefault {
			addr.IsDefault = false
			if err := s.userRepo.UpdateAddress(ctx, &addr); err != nil {
				return nil, err
			}
		}
	}

	// 將指定地址設為預設
	address.IsDefault = true
	address.UpdatedAt = time.Now()
	if err := s.userRepo.UpdateAddress(ctx, address); err != nil {
		return nil, err
	}

	return address, nil
}
