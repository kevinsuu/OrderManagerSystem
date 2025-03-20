package model

import (
	"time"
)

// NotificationType 通知類型
type NotificationType string

const (
	NotificationTypeEmail    NotificationType = "email"
	NotificationTypeSMS      NotificationType = "sms"
	NotificationTypePush     NotificationType = "push"
	NotificationTypeWebhook  NotificationType = "webhook"
)

// NotificationStatus 通知狀態
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusCancelled NotificationStatus = "cancelled"
)

// NotificationPriority 通知優先級
type NotificationPriority string

const (
	NotificationPriorityLow     NotificationPriority = "low"
	NotificationPriorityNormal  NotificationPriority = "normal"
	NotificationPriorityHigh    NotificationPriority = "high"
	NotificationPriorityUrgent  NotificationPriority = "urgent"
)

// Notification 通知模型
type Notification struct {
	ID          string              `json:"id" gorm:"primaryKey"`
	UserID      string              `json:"userId" gorm:"index"`
	Type        NotificationType    `json:"type"`
	Status      NotificationStatus  `json:"status"`
	Priority    NotificationPriority `json:"priority"`
	Title       string              `json:"title"`
	Content     string              `json:"content"`
	Metadata    string              `json:"metadata,omitempty"` // JSON 字符串，存儲額外信息
	RetryCount  int                 `json:"retryCount"`
	MaxRetries  int                 `json:"maxRetries"`
	SentAt      *time.Time          `json:"sentAt,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	DeletedAt   *time.Time          `json:"deletedAt,omitempty" gorm:"index"`
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex"`
	Type      NotificationType `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Variables []string  `json:"variables" gorm:"-"` // 不存儲在數據庫中
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateNotificationRequest 創建通知請求
type CreateNotificationRequest struct {
	UserID      string              `json:"userId" binding:"required"`
	Type        NotificationType    `json:"type" binding:"required"`
	Priority    NotificationPriority `json:"priority" binding:"required"`
	Title       string              `json:"title" binding:"required"`
	Content     string              `json:"content" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	MaxRetries  int                 `json:"maxRetries,omitempty"`
}

// CreateNotificationFromTemplateRequest 從模板創建通知請求
type CreateNotificationFromTemplateRequest struct {
	UserID       string                 `json:"userId" binding:"required"`
	TemplateID   string                 `json:"templateId" binding:"required"`
	Priority     NotificationPriority    `json:"priority" binding:"required"`
	Variables    map[string]interface{} `json:"variables" binding:"required"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	MaxRetries   int                    `json:"maxRetries,omitempty"`
}

// NotificationResponse 通知響應
type NotificationResponse struct {
	Notification
}

// NotificationListResponse 通知列表響應
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                 `json:"total"`
	Page          int                   `json:"page"`
	Limit         int                   `json:"limit"`
}

// EmailConfig 郵件配置
type EmailConfig struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	HTML     bool     `json:"html"`
}

// SMSConfig 短信配置
type SMSConfig struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// PushConfig 推送配置
type PushConfig struct {
	DeviceToken string                 `json:"deviceToken"`
	Title       string                 `json:"title"`
	Body        string                 `json:"body"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// WebhookConfig Webhook配置
type WebhookConfig struct {
	URL     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    map[string]interface{} `json:"body"`
}
