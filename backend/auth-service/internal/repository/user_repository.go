package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/model"
)

// IUserRepository 用戶存儲庫接口
type IUserRepository interface {
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

// UserRepository Realtime Database 實現
type UserRepository struct {
	client *db.Client
}

// NewUserRepository 創建用戶存儲實例
func NewUserRepository(client *db.Client) IUserRepository {
	return &UserRepository{
		client: client,
	}
}

// Create 創建用戶
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// 先檢查 email 是否已存在
	log.Printf("Checking if email exists: %s", user.Email)
	existingUser, err := r.GetByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if existingUser != nil {
		log.Printf("Email already exists: %s", user.Email)
		return fmt.Errorf("email already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	log.Printf("Creating user with ID: %s, Username: %s, Email: %s", user.ID, user.Username, user.Email)
	log.Printf("Password hash length: %d", len(user.Password))

	// 創建一個臨時的 map 來確保所有字段都被正確序列化
	userData := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.Password,
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt.Format(time.RFC3339),
		"updated_at": user.UpdatedAt.Format(time.RFC3339),
	}

	ref := r.client.NewRef("users/" + user.ID)
	if err := ref.Set(ctx, userData); err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 驗證數據是否正確保存
	var savedUser map[string]interface{}
	if err := ref.Get(ctx, &savedUser); err != nil {
		log.Printf("Error verifying saved user: %v", err)
		return nil
	}

	log.Printf("Saved user data: %+v", savedUser)
	return nil
}

// GetByID 通過ID獲取用戶
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.client.NewRef("users/"+id).Get(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 通過用戶名獲取用戶
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	log.Printf("Searching for user with username: %s", username)

	var users map[string]interface{}
	ref := r.client.NewRef("users").OrderByChild("username").EqualTo(username)
	if err := ref.Get(ctx, &users); err != nil {
		log.Printf("Error getting user by username: %v", err)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	log.Printf("Found users data: %+v", users)

	if len(users) == 0 {
		log.Printf("No user found with username: %s", username)
		return nil, nil
	}

	// 將 map 轉換為 User 結構體
	for _, userData := range users {
		userMap, ok := userData.(map[string]interface{})
		if !ok {
			log.Printf("Error: user data is not a map: %+v", userData)
			continue
		}

		log.Printf("Processing user data: %+v", userMap)

		user := &model.User{
			ID:       userMap["id"].(string),
			Username: userMap["username"].(string),
			Email:    userMap["email"].(string),
			Password: userMap["password"].(string),
			Role:     userMap["role"].(string),
			Status:   userMap["status"].(string),
		}

		return user, nil
	}

	return nil, nil
}

// GetByEmail 通過郵箱獲取用戶
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	log.Printf("Searching for user with email: %s", email)

	var users map[string]interface{}
	ref := r.client.NewRef("users").OrderByChild("email").EqualTo(email)
	if err := ref.Get(ctx, &users); err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	log.Printf("Found users data: %+v", users)

	if len(users) == 0 {
		log.Printf("No user found with email: %s", email)
		return nil, nil
	}

	// 將 map 轉換為 User 結構體
	for _, userData := range users {
		userMap, ok := userData.(map[string]interface{})
		if !ok {
			log.Printf("Error: user data is not a map: %+v", userData)
			continue
		}

		log.Printf("Processing user data: %+v", userMap)

		user := &model.User{
			ID:       userMap["id"].(string),
			Username: userMap["username"].(string),
			Email:    userMap["email"].(string),
			Password: userMap["password"].(string),
			Role:     userMap["role"].(string),
			Status:   userMap["status"].(string),
		}

		// 檢查密碼是否存在
		log.Printf("Retrieved password hash length: %d", len(user.Password))

		return user, nil
	}

	return nil, nil
}

// Update 更新用戶
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	return r.client.NewRef("users/"+user.ID).Set(ctx, user)
}

// Delete 刪除用戶
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.client.NewRef("users/" + id).Delete(ctx)
}

// List 獲取用戶列表
func (r *UserRepository) List(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	var users map[string]model.User
	err := r.client.NewRef("users").Get(ctx, &users)
	if err != nil {
		return nil, 0, err
	}

	userList := make([]model.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	return userList, int64(len(userList)), nil
}

// CreateAddress 創建地址
func (r *UserRepository) CreateAddress(ctx context.Context, address *model.Address) error {
	log.Printf("Creating address for user ID: %s", address.UserID)

	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()

	// 創建一個臨時的 map 來確保所有字段都被正確序列化
	addressData := map[string]interface{}{
		"id":          address.ID,
		"user_id":     address.UserID,
		"name":        address.Name,
		"phone":       address.Phone,
		"street":      address.Street,
		"city":        address.City,
		"district":    address.District,
		"postal_code": address.PostalCode,
		"is_default":  address.IsDefault,
		"created_at":  address.CreatedAt.Format(time.RFC3339),
		"updated_at":  address.UpdatedAt.Format(time.RFC3339),
	}

	ref := r.client.NewRef("addresses/" + address.ID)
	if err := ref.Set(ctx, addressData); err != nil {
		log.Printf("Error creating address: %v", err)
		return fmt.Errorf("failed to create address: %w", err)
	}

	// 驗證數據是否正確保存
	var savedAddress map[string]interface{}
	if err := ref.Get(ctx, &savedAddress); err != nil {
		log.Printf("Error verifying saved address: %v", err)
		return nil
	}

	log.Printf("Saved address data: %+v", savedAddress)
	return nil
}

// GetAddresses 獲取用戶的所有地址
func (r *UserRepository) GetAddresses(ctx context.Context, userID string) ([]model.Address, error) {
	log.Printf("Fetching addresses for user ID: %s", userID)

	var addresses map[string]interface{}
	ref := r.client.NewRef("addresses").OrderByChild("user_id").EqualTo(userID)
	if err := ref.Get(ctx, &addresses); err != nil {
		log.Printf("Error getting addresses: %v", err)
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	log.Printf("Found addresses data: %+v", addresses)

	addressList := make([]model.Address, 0)
	for _, addrData := range addresses {
		addrMap, ok := addrData.(map[string]interface{})
		if !ok {
			log.Printf("Error: address data is not a map: %+v", addrData)
			continue
		}

		log.Printf("Processing address data: %+v", addrMap)

		// 檢查所需的字段是否存在
		if addrMap["id"] == nil || addrMap["user_id"] == nil ||
			addrMap["name"] == nil || addrMap["phone"] == nil ||
			addrMap["street"] == nil || addrMap["city"] == nil ||
			addrMap["district"] == nil || addrMap["postal_code"] == nil {
			log.Printf("Error: missing required fields in address data")
			continue
		}

		address := model.Address{
			ID:         addrMap["id"].(string),
			UserID:     addrMap["user_id"].(string),
			Name:       addrMap["name"].(string),
			Phone:      addrMap["phone"].(string),
			Street:     addrMap["street"].(string),
			City:       addrMap["city"].(string),
			District:   addrMap["district"].(string),
			PostalCode: addrMap["postal_code"].(string),
			IsDefault:  addrMap["is_default"].(bool),
		}

		if createdAt, ok := addrMap["created_at"].(string); ok {
			t, err := time.Parse(time.RFC3339, createdAt)
			if err == nil {
				address.CreatedAt = t
			}
		}

		if updatedAt, ok := addrMap["updated_at"].(string); ok {
			t, err := time.Parse(time.RFC3339, updatedAt)
			if err == nil {
				address.UpdatedAt = t
			}
		}

		addressList = append(addressList, address)
	}

	log.Printf("Returning %d addresses", len(addressList))
	return addressList, nil
}

// GetAddressByID 通過ID獲取地址
func (r *UserRepository) GetAddressByID(ctx context.Context, id string) (*model.Address, error) {
	var address model.Address
	err := r.client.NewRef("addresses/"+id).Get(ctx, &address)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// UpdateAddress 更新地址
func (r *UserRepository) UpdateAddress(ctx context.Context, address *model.Address) error {
	address.UpdatedAt = time.Now()
	return r.client.NewRef("addresses/"+address.ID).Set(ctx, address)
}

// DeleteAddress 刪除地址
func (r *UserRepository) DeleteAddress(ctx context.Context, id string) error {
	return r.client.NewRef("addresses/" + id).Delete(ctx)
}

// GetPreference 獲取用戶偏好設置
func (r *UserRepository) GetPreference(ctx context.Context, userID string) (*model.UserPreference, error) {
	log.Printf("Getting preference for user ID: %s", userID)

	var pref model.UserPreference
	ref := r.client.NewRef("preferences/" + userID)
	if err := ref.Get(ctx, &pref); err != nil {
		log.Printf("Error getting preference: %v", err)
		// 如果找不到偏好設置，返回預設值
		if err.Error() == "http error status: 404; reason: path preferences/"+userID+" not found" {
			defaultPref := model.NewDefaultPreference(userID)
			// 保存預設偏好設置
			if err := r.UpdatePreference(ctx, defaultPref); err != nil {
				log.Printf("Error saving default preference: %v", err)
				return nil, fmt.Errorf("failed to save default preference: %w", err)
			}
			return defaultPref, nil
		}
		return nil, fmt.Errorf("failed to get preference: %w", err)
	}

	// 如果所有字段都是空值，返回預設值
	if pref.Language == "" && pref.Currency == "" && pref.Theme == "" {
		defaultPref := model.NewDefaultPreference(userID)
		// 保存預設偏好設置
		if err := r.UpdatePreference(ctx, defaultPref); err != nil {
			log.Printf("Error saving default preference: %v", err)
			return nil, fmt.Errorf("failed to save default preference: %w", err)
		}
		return defaultPref, nil
	}

	return &pref, nil
}

// UpdatePreference 更新用戶偏好設置
func (r *UserRepository) UpdatePreference(ctx context.Context, pref *model.UserPreference) error {
	log.Printf("Updating preference for user ID: %s", pref.UserID)

	pref.UpdatedAt = time.Now()

	// 創建一個臨時的 map 來確保所有字段都被正確序列化
	prefData := map[string]interface{}{
		"user_id":            pref.UserID,
		"language":           pref.Language,
		"currency":           pref.Currency,
		"notification_email": pref.NotificationEmail,
		"notification_sms":   pref.NotificationSMS,
		"theme":              pref.Theme,
		"created_at":         pref.CreatedAt.Format(time.RFC3339),
		"updated_at":         pref.UpdatedAt.Format(time.RFC3339),
	}

	ref := r.client.NewRef("preferences/" + pref.UserID)
	if err := ref.Set(ctx, prefData); err != nil {
		log.Printf("Error updating preference: %v", err)
		return fmt.Errorf("failed to update preference: %w", err)
	}

	// 驗證數據是否正確保存
	var savedPref map[string]interface{}
	if err := ref.Get(ctx, &savedPref); err != nil {
		log.Printf("Error verifying saved preference: %v", err)
		return nil
	}

	log.Printf("Saved preference data: %+v", savedPref)
	return nil
}
