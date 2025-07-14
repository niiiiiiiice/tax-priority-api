package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"tax-priority-api/src/database"
	"tax-priority-api/src/models"
)

// mapSortField маппит поля сортировки от AdminJS к именам колонок БД
func mapSortField(field string) string {
	fieldMap := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"email":     "email",
		"name":      "name",
		"role":      "role",
		"status":    "status",
		"avatar":    "avatar",
		"id":        "id",
	}

	if dbField, exists := fieldMap[field]; exists {
		return dbField
	}
	return field // возвращаем исходное поле, если маппинг не найден
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		db: database.GetDB(),
	}
}

// GetUsers получает список пользователей с фильтрацией и пагинацией
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	var total int64

	// Параметры пагинации
	limit, _ := strconv.Atoi(c.DefaultQuery("_limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("_offset", "0"))

	// Параметры сортировки
	sortBy := mapSortField(c.DefaultQuery("_sort", "createdAt"))
	order := c.DefaultQuery("_order", "desc")

	query := h.db.Model(&models.User{})

	// Применяем фильтры
	if email := c.Query("email"); email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Фильтр по дате создания (от)
	if createdAtGte := c.Query("createdAt_gte"); createdAtGte != "" {
		query = query.Where("created_at >= ?", createdAtGte)
	}
	// Фильтр по дате создания (до)
	if createdAtLte := c.Query("createdAt_lte"); createdAtLte != "" {
		query = query.Where("created_at <= ?", createdAtLte)
	}

	// Подсчет общего количества
	query.Count(&total)

	// Применяем сортировку
	orderClause := sortBy + " " + strings.ToUpper(order)
	query = query.Order(orderClause)

	// Применяем пагинацию
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUserCount получает количество пользователей
func (h *UserHandler) GetUserCount(c *gin.Context) {
	var count int64

	query := h.db.Model(&models.User{})

	// Применяем те же фильтры что и в GetUsers
	if email := c.Query("email"); email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetUser получает пользователя по ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser создает нового пользователя
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Status == "" {
		user.Status = "active"
	}

	// Валидация
	if !user.IsValidRole() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	if !user.IsValidStatus() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser обновляет пользователя
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация обновляемых полей
	if updateData.Role != "" && !updateData.IsValidRole() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	if updateData.Status != "" && !updateData.IsValidStatus() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.db.Model(&user).Updates(updateData).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser удаляет пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.db.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUsersBatch получает пользователей по списку ID
func (h *UserHandler) GetUsersBatch(c *gin.Context) {
	var request struct {
		IDs []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users []models.User
	if err := h.db.Where("id IN ?", request.IDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": users})
}

// BulkDeleteUsers удаляет пользователей по списку ID
func (h *UserHandler) BulkDeleteUsers(c *gin.Context) {
	var request struct {
		IDs []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Where("id IN ?", request.IDs).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users deleted successfully"})
}

// ChangeUserStatus изменяет статус пользователя
func (h *UserHandler) ChangeUserStatus(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация статуса
	user := models.User{Status: request.Status}
	if !user.IsValidStatus() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.db.Model(&models.User{}).Where("id = ?", userID).Update("status", request.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	// Получаем обновленного пользователя
	var updatedUser models.User
	h.db.First(&updatedUser, "id = ?", userID)

	c.JSON(http.StatusOK, updatedUser)
}

// GetUsersByRole получает пользователей по роли
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
	role := c.Param("role")

	var users []models.User
	if err := h.db.Where("role = ?", role).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// SearchUsers поиск пользователей
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	var users []models.User
	searchQuery := "%" + query + "%"

	if err := h.db.Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
