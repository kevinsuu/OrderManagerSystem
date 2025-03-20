package model

import "time"

// UserPreference 用戶偏好設置
type UserPreference struct {
	UserID            string    `json:"user_id" firestore:"user_id"`
	Language          string    `json:"language" firestore:"language"`
	Currency          string    `json:"currency" firestore:"currency"`
	NotificationEmail bool      `json:"notification_email" firestore:"notification_email"`
	NotificationSMS   bool      `json:"notification_sms" firestore:"notification_sms"`
	Theme             string    `json:"theme" firestore:"theme"`
	CreatedAt         time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" firestore:"updated_at"`
}

// NewDefaultPreference 創建預設的用戶偏好設置
func NewDefaultPreference(userID string) *UserPreference {
	now := time.Now()
	return &UserPreference{
		UserID:            userID,
		Language:          "zh-TW", // 預設繁體中文
		Currency:          "TWD",   // 預設新台幣
		NotificationEmail: true,    // 預設開啟郵件通知
		NotificationSMS:   false,   // 預設關閉簡訊通知
		Theme:             "light", // 預設淺色主題
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// PreferenceRequest 偏好設置請求
type PreferenceRequest struct {
	Language          string `json:"language" binding:"required"`
	Currency          string `json:"currency" binding:"required"`
	NotificationEmail bool   `json:"notification_email"`
	NotificationSMS   bool   `json:"notification_sms"`
	Theme             string `json:"theme" binding:"required"`
}
