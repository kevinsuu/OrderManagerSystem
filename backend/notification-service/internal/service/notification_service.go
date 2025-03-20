package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/repository"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrTemplateNotFound    = errors.New("template not found")
	ErrInvalidTemplate     = errors.New("invalid template")
)

// NotificationService 通知服務接口
type NotificationService interface {
	CreateNotification(ctx context.Context, req *model.CreateNotificationRequest) (*model.Notification, error)
	CreateNotificationFromTemplate(ctx context.Context, req *model.CreateNotificationFromTemplateRequest) (*model.Notification, error)
	GetNotification(ctx context.Context, id string) (*model.NotificationResponse, error)
	GetUserNotifications(ctx context.Context, userID string, page, limit int) (*model.NotificationListResponse, error)
	ListNotifications(ctx context.Context, page, limit int) (*model.NotificationListResponse, error)
	ProcessPendingNotifications(ctx context.Context) error
	CreateTemplate(ctx context.Context, template *model.NotificationTemplate) error
	GetTemplate(ctx context.Context, id string) (*model.NotificationTemplate, error)
	UpdateTemplate(ctx context.Context, template *model.NotificationTemplate) error
	ListTemplates(ctx context.Context, page, limit int) ([]model.NotificationTemplate, int64, error)
}

type notificationService struct {
	repo repository.NotificationRepository
}

// NewNotificationService 創建通知服務實例
func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

// CreateNotification 創建通知
func (s *notificationService) CreateNotification(ctx context.Context, req *model.CreateNotificationRequest) (*model.Notification, error) {
	metadata, err := json.Marshal(req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	notification := &model.Notification{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Type:       req.Type,
		Status:     model.NotificationStatusPending,
		Priority:   req.Priority,
		Title:      req.Title,
		Content:    req.Content,
		Metadata:   string(metadata),
		MaxRetries: req.MaxRetries,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// CreateNotificationFromTemplate 從模板創建通知
func (s *notificationService) CreateNotificationFromTemplate(ctx context.Context, req *model.CreateNotificationFromTemplateRequest) (*model.Notification, error) {
	// 獲取模板
	tmpl, err := s.repo.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	if tmpl == nil {
		return nil, ErrTemplateNotFound
	}

	// 解析並應用模板
	title, err := s.parseTemplate(tmpl.Title, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to parse title template: %w", err)
	}

	content, err := s.parseTemplate(tmpl.Content, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content template: %w", err)
	}

	// 創建通知請求
	createReq := &model.CreateNotificationRequest{
		UserID:     req.UserID,
		Type:       tmpl.Type,
		Priority:   req.Priority,
		Title:      title,
		Content:    content,
		Metadata:   req.Metadata,
		MaxRetries: req.MaxRetries,
	}

	return s.CreateNotification(ctx, createReq)
}

// GetNotification 獲取通知詳情
func (s *notificationService) GetNotification(ctx context.Context, id string) (*model.NotificationResponse, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return &model.NotificationResponse{
		Notification: *notification,
	}, nil
}

// GetUserNotifications 獲取用戶的通知
func (s *notificationService) GetUserNotifications(ctx context.Context, userID string, page, limit int) (*model.NotificationListResponse, error) {
	notifications, total, err := s.repo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	response := &model.NotificationListResponse{
		Notifications: make([]model.NotificationResponse, len(notifications)),
		Total:        total,
		Page:         page,
		Limit:        limit,
	}

	for i, notification := range notifications {
		response.Notifications[i] = model.NotificationResponse{
			Notification: notification,
		}
	}

	return response, nil
}

// ListNotifications 獲取通知列表
func (s *notificationService) ListNotifications(ctx context.Context, page, limit int) (*model.NotificationListResponse, error) {
	notifications, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	response := &model.NotificationListResponse{
		Notifications: make([]model.NotificationResponse, len(notifications)),
		Total:        total,
		Page:         page,
		Limit:        limit,
	}

	for i, notification := range notifications {
		response.Notifications[i] = model.NotificationResponse{
			Notification: notification,
		}
	}

	return response, nil
}

// ProcessPendingNotifications 處理待發送的通知
func (s *notificationService) ProcessPendingNotifications(ctx context.Context) error {
	notifications, err := s.repo.GetPendingNotifications(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %w", err)
	}

	for _, notification := range notifications {
		if err := s.processNotification(ctx, &notification); err != nil {
			notification.RetryCount++
			notification.Status = model.NotificationStatusFailed
			if notification.RetryCount >= notification.MaxRetries {
				notification.Status = model.NotificationStatusCancelled
			}
		} else {
			now := time.Now()
			notification.Status = model.NotificationStatusSent
			notification.SentAt = &now
		}

		notification.UpdatedAt = time.Now()
		if err := s.repo.Update(ctx, &notification); err != nil {
			return fmt.Errorf("failed to update notification status: %w", err)
		}
	}

	return nil
}

// CreateTemplate 創建通知模板
func (s *notificationService) CreateTemplate(ctx context.Context, template *model.NotificationTemplate) error {
	template.ID = uuid.New().String()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	return s.repo.CreateTemplate(ctx, template)
}

// GetTemplate 獲取通知模板
func (s *notificationService) GetTemplate(ctx context.Context, id string) (*model.NotificationTemplate, error) {
	return s.repo.GetTemplateByID(ctx, id)
}

// UpdateTemplate 更新通知模板
func (s *notificationService) UpdateTemplate(ctx context.Context, template *model.NotificationTemplate) error {
	template.UpdatedAt = time.Now()
	return s.repo.UpdateTemplate(ctx, template)
}

// ListTemplates 獲取模板列表
func (s *notificationService) ListTemplates(ctx context.Context, page, limit int) ([]model.NotificationTemplate, int64, error) {
	return s.repo.ListTemplates(ctx, page, limit)
}

// 輔助方法

// processNotification 處理單個通知
func (s *notificationService) processNotification(ctx context.Context, notification *model.Notification) error {
	switch notification.Type {
	case model.NotificationTypeEmail:
		return s.sendEmail(ctx, notification)
	case model.NotificationTypeSMS:
		return s.sendSMS(ctx, notification)
	case model.NotificationTypePush:
		return s.sendPushNotification(ctx, notification)
	case model.NotificationTypeWebhook:
		return s.sendWebhook(ctx, notification)
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}

// parseTemplate 解析模板並替換變數
func (s *notificationService) parseTemplate(content string, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New("notification").Parse(content)
	if err != nil {
		return "", fmt.Errorf("解析模板失敗: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("執行模板失敗: %w", err)
	}

	return buf.String(), nil
}

// sendEmail 發送郵件
func (s *notificationService) sendEmail(ctx context.Context, notification *model.Notification) error {
	// TODO: 實現郵件發送邏輯
	return nil
}

// sendSMS 發送短信
func (s *notificationService) sendSMS(ctx context.Context, notification *model.Notification) error {
	// TODO: 實現短信發送邏輯
	return nil
}

// sendPushNotification 發送推送通知
func (s *notificationService) sendPushNotification(ctx context.Context, notification *model.Notification) error {
	// TODO: 實現推送通知邏輯
	return nil
}

// sendWebhook 發送 Webhook
func (s *notificationService) sendWebhook(ctx context.Context, notification *model.Notification) error {
	// TODO: 實現 Webhook 發送邏輯
	return nil
}
