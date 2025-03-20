package repository

import (
	"context"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/model"
)

// NotificationRepository 通知存儲接口
type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, id string) (*model.Notification, error)
	Update(ctx context.Context, notification *model.Notification) error
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Notification, int64, error)
	List(ctx context.Context, page, limit int) ([]model.Notification, int64, error)
	GetPendingNotifications(ctx context.Context, limit int) ([]model.Notification, error)
	CreateTemplate(ctx context.Context, template *model.NotificationTemplate) error
	GetTemplateByID(ctx context.Context, id string) (*model.NotificationTemplate, error)
	GetTemplateByName(ctx context.Context, name string) (*model.NotificationTemplate, error)
	UpdateTemplate(ctx context.Context, template *model.NotificationTemplate) error
	ListTemplates(ctx context.Context, page, limit int) ([]model.NotificationTemplate, int64, error)
}

type notificationRepository struct {
	db *db.Client
}

// NewNotificationRepository 創建通知存儲實例
func NewNotificationRepository(db *db.Client) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

// Create 創建通知
func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	ref := r.db.NewRef("notifications")
	return ref.Child(notification.ID).Set(ctx, notification)
}

// GetByID 根據ID獲取通知
func (r *notificationRepository) GetByID(ctx context.Context, id string) (*model.Notification, error) {
	var notification model.Notification
	ref := r.db.NewRef("notifications").Child(id)
	if err := ref.Get(ctx, &notification); err != nil {
		return nil, err
	}
	if notification.ID == "" {
		return nil, nil
	}
	return &notification, nil
}

// Update 更新通知
func (r *notificationRepository) Update(ctx context.Context, notification *model.Notification) error {
	ref := r.db.NewRef("notifications").Child(notification.ID)
	notification.UpdatedAt = time.Now()
	return ref.Set(ctx, notification)
}

// GetByUserID 獲取用戶的通知
func (r *notificationRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Notification, int64, error) {
	var result map[string]model.Notification
	ref := r.db.NewRef("notifications")
	if err := ref.OrderByChild("userId").EqualTo(userID).Get(ctx, &result); err != nil {
		return nil, 0, err
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if start >= int(total) {
		return []model.Notification{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	notifications := make([]model.Notification, 0, len(result))
	for _, notification := range result {
		notifications = append(notifications, notification)
	}

	return notifications[start:end], total, nil
}

// List 獲取通知列表
func (r *notificationRepository) List(ctx context.Context, page, limit int) ([]model.Notification, int64, error) {
	var result map[string]model.Notification
	ref := r.db.NewRef("notifications")
	if err := ref.Get(ctx, &result); err != nil {
		return nil, 0, err
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if start >= int(total) {
		return []model.Notification{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	notifications := make([]model.Notification, 0, len(result))
	for _, notification := range result {
		notifications = append(notifications, notification)
	}

	return notifications[start:end], total, nil
}

// GetPendingNotifications 獲取待處理的通知
func (r *notificationRepository) GetPendingNotifications(ctx context.Context, limit int) ([]model.Notification, error) {
	var result map[string]model.Notification
	ref := r.db.NewRef("notifications")
	if err := ref.OrderByChild("status").EqualTo(string(model.NotificationStatusPending)).Get(ctx, &result); err != nil {
		return nil, err
	}

	notifications := make([]model.Notification, 0, len(result))
	for _, notification := range result {
		notifications = append(notifications, notification)
	}

	if len(notifications) > limit {
		notifications = notifications[:limit]
	}

	return notifications, nil
}

// CreateTemplate 創建通知模板
func (r *notificationRepository) CreateTemplate(ctx context.Context, template *model.NotificationTemplate) error {
	ref := r.db.NewRef("templates")
	return ref.Child(template.ID).Set(ctx, template)
}

// GetTemplateByID 根據ID獲取模板
func (r *notificationRepository) GetTemplateByID(ctx context.Context, id string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	ref := r.db.NewRef("templates").Child(id)
	if err := ref.Get(ctx, &template); err != nil {
		return nil, err
	}
	if template.ID == "" {
		return nil, nil
	}
	return &template, nil
}

// GetTemplateByName 根據名稱獲取模板
func (r *notificationRepository) GetTemplateByName(ctx context.Context, name string) (*model.NotificationTemplate, error) {
	var result map[string]model.NotificationTemplate
	ref := r.db.NewRef("templates")
	if err := ref.OrderByChild("name").EqualTo(name).Get(ctx, &result); err != nil {
		return nil, err
	}
	for _, template := range result {
		return &template, nil
	}
	return nil, nil
}

// UpdateTemplate 更新模板
func (r *notificationRepository) UpdateTemplate(ctx context.Context, template *model.NotificationTemplate) error {
	ref := r.db.NewRef("templates").Child(template.ID)
	template.UpdatedAt = time.Now()
	return ref.Set(ctx, template)
}

// ListTemplates 獲取模板列表
func (r *notificationRepository) ListTemplates(ctx context.Context, page, limit int) ([]model.NotificationTemplate, int64, error) {
	var result map[string]model.NotificationTemplate
	ref := r.db.NewRef("templates")
	if err := ref.Get(ctx, &result); err != nil {
		return nil, 0, err
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if start >= int(total) {
		return []model.NotificationTemplate{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	templates := make([]model.NotificationTemplate, 0, len(result))
	for _, template := range result {
		templates = append(templates, template)
	}

	return templates[start:end], total, nil
}
