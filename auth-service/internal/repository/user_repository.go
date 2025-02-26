package repository

import (
	"context"
	"errors"
	"time"

	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用戶存儲庫接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]model.User, int64, error)
	CreateAddress(ctx context.Context, address *model.Address) error
	GetAddresses(ctx context.Context, userID string) ([]model.Address, error)
	GetAddressByID(ctx context.Context, id string) (*model.Address, error)
	UpdateAddress(ctx context.Context, address *model.Address) error
	DeleteAddress(ctx context.Context, id string) error
	GetPreference(ctx context.Context, userID string) (*model.UserPreference, error)
	UpdatePreference(ctx context.Context, pref *model.UserPreference) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 創建用戶存儲實例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 創建用戶
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 通過ID獲取用戶
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 通過用戶名獲取用戶
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 通過郵箱獲取用戶
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用戶
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 刪除用戶
func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// List 獲取用戶列表
func (r *userRepository) List(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CreateAddress 創建地址
func (r *userRepository) CreateAddress(ctx context.Context, address *model.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

// GetAddresses 獲取用戶的所有地址
func (r *userRepository) GetAddresses(ctx context.Context, userID string) ([]model.Address, error) {
	var addresses []model.Address
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetAddressByID 通過ID獲取地址
func (r *userRepository) GetAddressByID(ctx context.Context, id string) (*model.Address, error) {
	var address model.Address
	if err := r.db.WithContext(ctx).First(&address, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

// DeleteAddress 刪除地址
func (r *userRepository) DeleteAddress(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Address{}, "id = ?", id).Error
}

// GetPreference 獲取用戶偏好設置
func (r *userRepository) GetPreference(ctx context.Context, userID string) (*model.UserPreference, error) {
	var preference model.UserPreference
	if err := r.db.WithContext(ctx).First(&preference, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果沒有找到偏好設置，返回默認值
			return &model.UserPreference{
				UserID:            userID,
				Language:          "zh-TW",
				Currency:          "TWD",
				NotificationEmail: true,
				NotificationSMS:   false,
				Theme:             "light",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}, nil
		}
		return nil, err
	}
	return &preference, nil
}

// UpdatePreference 更新用戶偏好設置
func (r *userRepository) UpdatePreference(ctx context.Context, pref *model.UserPreference) error {
	return r.db.WithContext(ctx).Save(pref).Error
}

// UpdateAddress 更新地址
func (r *userRepository) UpdateAddress(ctx context.Context, address *model.Address) error {
	// 檢查地址是否存在
	var existingAddress model.Address
	if err := r.db.WithContext(ctx).First(&existingAddress, "id = ?", address.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	// 更新地址
	return r.db.WithContext(ctx).Save(address).Error
}
