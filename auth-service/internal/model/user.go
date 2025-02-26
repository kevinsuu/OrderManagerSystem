package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用戶模型
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"not null"` // 密碼不會在 JSON 中返回
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserLoginRequest 用戶登錄請求
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserRegisterRequest 用戶註冊請求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserResponse 用戶響應
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginResponse 登錄響應
type LoginResponse struct {
	Token        string       `json:"token"`
	User         UserResponse `json:"user"`
	RefreshToken string       `json:"refreshToken"`
}

// HashPassword 加密密碼
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 檢查密碼
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Address 用戶地址
type Address struct {
    ID         string    `json:"id" gorm:"primaryKey"`
    UserID     string    `json:"userId" gorm:"index"`
    Name       string    `json:"name"`
    Recipient  string    `json:"recipient"`
    Phone      string    `json:"phone"`
    PostalCode string    `json:"postal_code"`
    City       string    `json:"city"`
    District   string    `json:"district"`
    Street     string    `json:"street"`
    IsDefault  bool      `json:"is_default"`
    CreatedAt  time.Time `json:"createdAt"`
    UpdatedAt  time.Time `json:"updatedAt"`
}

// UserPreference 用戶偏好
type UserPreference struct {
	UserID            string    `json:"userId" gorm:"primaryKey"`
	Language          string    `json:"language" gorm:"default:'zh-TW'"`
	Currency          string    `json:"currency" gorm:"default:'TWD'"`
	NotificationEmail bool      `json:"notificationEmail" gorm:"default:true"`
	NotificationSMS   bool      `json:"notificationSMS" gorm:"default:false"`
	Theme             string    `json:"theme" gorm:"default:'light'"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// AddressRequest 地址請求
type AddressRequest struct {
    Name       string `json:"name" binding:"required"`
    Recipient  string `json:"recipient" binding:"required"`
    Phone      string `json:"phone" binding:"required"`
    PostalCode string `json:"postal_code" binding:"required"`
    City       string `json:"city" binding:"required"`
    District   string `json:"district" binding:"required"`
    Street     string `json:"street" binding:"required"`
    IsDefault  bool   `json:"is_default"`
}

// PreferenceRequest 偏好設置請求
type PreferenceRequest struct {
	Language          string `json:"language"`
	Currency          string `json:"currency"`
	NotificationEmail bool   `json:"notificationEmail"`
	NotificationSMS   bool   `json:"notificationSMS"`
	Theme             string `json:"theme"`
}
