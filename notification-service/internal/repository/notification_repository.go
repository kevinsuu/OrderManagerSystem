package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/model"
	"gorm.io/gorm"
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
	db    *gorm.DB
	redis *redis.Client
}

// NewNotificationRepository 創建通知存儲實例
func NewNotificationRepository(db *gorm.DB, redis *redis.Client) NotificationRepository {
	return &notificationRepository{
		db:    db,
		redis: redis,
	}
}

// Create 創建通知
func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return err
	}

	// 如果是待發送狀態，加入待處理隊列
	if notification.Status == model.NotificationStatusPending {
		return r.addToPendingQueue(ctx, notification)
	}

	return nil
}

// GetByID 根據ID獲取通知
func (r *notificationRepository) GetByID(ctx context.Context, id string) (*model.Notification, error) {
	// 嘗試從快取獲取
	notification, err := r.getFromCache(ctx, id)
	if err == nil {
		return notification, nil
	}

	// 從數據庫獲取
	notification = &model.Notification{}
	if err := r.db.WithContext(ctx).First(notification, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// 設置快取
	r.setCache(ctx, notification)

	return notification, nil
}

// Update 更新通知
func (r *notificationRepository) Update(ctx context.Context, notification *model.Notification) error {
	if err := r.db.WithContext(ctx).Save(notification).Error; err != nil {
		return err
	}

	// 更新快取
	r.setCache(ctx, notification)

	return nil
}

// GetByUserID 獲取用戶的通知
func (r *notificationRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// List 獲取通知列表
func (r *notificationRepository) List(ctx context.Context, page, limit int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetPendingNotifications 獲取待處理的通知
func (r *notificationRepository) GetPendingNotifications(ctx context.Context, limit int) ([]model.Notification, error) {
	var notifications []model.Notification

	if err := r.db.WithContext(ctx).
		Where("status = ?", model.NotificationStatusPending).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// CreateTemplate 創建通知模板
func (r *notificationRepository) CreateTemplate(ctx context.Context, template *model.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetTemplateByID 根據ID獲取模板
func (r *notificationRepository) GetTemplateByID(ctx context.Context, id string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	if err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// GetTemplateByName 根據名稱獲取模板
func (r *notificationRepository) GetTemplateByName(ctx context.Context, name string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	if err := r.db.WithContext(ctx).First(&template, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// UpdateTemplate 更新模板
func (r *notificationRepository) UpdateTemplate(ctx context.Context, template *model.NotificationTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// ListTemplates 獲取模板列表
func (r *notificationRepository) ListTemplates(ctx context.Context, page, limit int) ([]model.NotificationTemplate, int64, error) {
	var templates []model.NotificationTemplate
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// 輔助方法

// addToPendingQueue 將通知加入待處理隊列
func (r *notificationRepository) addToPendingQueue(ctx context.Context, notification *model.Notification) error {
	key := "pending_notifications"
	score := float64(time.Now().Unix())
	member := notification.ID

	return r.redis.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// getFromCache 從快取獲取通知
func (r *notificationRepository) getFromCache(ctx context.Context, id string) (*model.Notification, error) {
	key := "notification:" + id
	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var notification model.Notification
	if err := json.Unmarshal(data, &notification); err != nil {
		return nil, err
	}

	return &notification, nil
}

// setCache 設置通知快取
func (r *notificationRepository) setCache(ctx context.Context, notification *model.Notification) error {
	key := "notification:" + notification.ID
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, key, data, 1*time.Hour).Err()
}
