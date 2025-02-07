package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/service"
)

type Handler struct {
	notificationService service.NotificationService
}

func NewHandler(notificationService service.NotificationService) *Handler {
	return &Handler{
		notificationService: notificationService,
	}
}

// HealthCheck 健康檢查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateNotification 創建通知
func (h *Handler) CreateNotification(c *gin.Context) {
	var req model.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationService.CreateNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// CreateNotificationFromTemplate 從模板創建通知
func (h *Handler) CreateNotificationFromTemplate(c *gin.Context) {
	var req model.CreateNotificationFromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationService.CreateNotificationFromTemplate(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetNotification 獲取通知詳情
func (h *Handler) GetNotification(c *gin.Context) {
	id := c.Param("id")
	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// GetUserNotifications 獲取用戶通知
func (h *Handler) GetUserNotifications(c *gin.Context) {
	userID := c.Param("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notifications, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// ListNotifications 獲取通知列表
func (h *Handler) ListNotifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notifications, err := h.notificationService.ListNotifications(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// CreateTemplate 創建通知模板
func (h *Handler) CreateTemplate(c *gin.Context) {
	var template model.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.CreateTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// GetTemplate 獲取通知模板
func (h *Handler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.notificationService.GetTemplate(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate 更新通知模板
func (h *Handler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var template model.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	if err := h.notificationService.UpdateTemplate(c.Request.Context(), &template); err != nil {
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// ListTemplates 獲取模板列表
func (h *Handler) ListTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	templates, total, err := h.notificationService.ListTemplates(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}
