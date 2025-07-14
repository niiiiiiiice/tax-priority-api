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
func mapProductSortField(field string) string {
	fieldMap := map[string]string{
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
		"name":        "name",
		"description": "description",
		"price":       "price",
		"currency":    "currency",
		"category":    "category",
		"tags":        "tags",
		"inStock":     "in_stock",
		"quantity":    "quantity",
		"images":      "images",
		"id":          "id",
	}

	if dbField, exists := fieldMap[field]; exists {
		return dbField
	}
	return field // возвращаем исходное поле, если маппинг не найден
}

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		db: database.GetDB(),
	}
}

// GetProducts получает список продуктов с фильтрацией и пагинацией
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	// Параметры пагинации
	limit, _ := strconv.Atoi(c.DefaultQuery("_limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("_offset", "0"))

	// Параметры сортировки
	sortBy := mapProductSortField(c.DefaultQuery("_sort", "createdAt"))
	order := c.DefaultQuery("_order", "desc")

	query := h.db.Model(&models.Product{})

	// Применяем фильтры
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if currency := c.Query("currency"); currency != "" {
		query = query.Where("currency = ?", currency)
	}
	if inStock := c.Query("inStock"); inStock != "" {
		query = query.Where("in_stock = ?", inStock == "true")
	}

	// Фильтр по цене (от)
	if priceGte := c.Query("price_gte"); priceGte != "" {
		if price, err := strconv.ParseFloat(priceGte, 64); err == nil {
			query = query.Where("price >= ?", price)
		}
	}
	// Фильтр по цене (до)
	if priceLte := c.Query("price_lte"); priceLte != "" {
		if price, err := strconv.ParseFloat(priceLte, 64); err == nil {
			query = query.Where("price <= ?", price)
		}
	}

	// Фильтр по количеству (от)
	if quantityGte := c.Query("quantity_gte"); quantityGte != "" {
		if quantity, err := strconv.Atoi(quantityGte); err == nil {
			query = query.Where("quantity >= ?", quantity)
		}
	}
	// Фильтр по количеству (до)
	if quantityLte := c.Query("quantity_lte"); quantityLte != "" {
		if quantity, err := strconv.Atoi(quantityLte); err == nil {
			query = query.Where("quantity <= ?", quantity)
		}
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
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   products,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetProductCount получает количество продуктов
func (h *ProductHandler) GetProductCount(c *gin.Context) {
	var count int64

	query := h.db.Model(&models.Product{})

	// Применяем те же фильтры что и в GetProducts
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if currency := c.Query("currency"); currency != "" {
		query = query.Where("currency = ?", currency)
	}
	if inStock := c.Query("inStock"); inStock != "" {
		query = query.Where("in_stock = ?", inStock == "true")
	}

	if err := query.Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetProduct получает продукт по ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.db.First(&product, "id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct создает новый продукт
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if product.Currency == "" {
		product.Currency = "USD"
	}

	// Валидация
	if !product.IsValidCurrency() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency"})
		return
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct обновляет продукт
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.db.First(&product, "id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация обновляемых полей
	if updateData.Currency != "" && !updateData.IsValidCurrency() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency"})
		return
	}

	if err := h.db.Model(&product).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct удаляет продукт
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.db.Delete(&models.Product{}, "id = ?", productID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetProductsBatch получает продукты по списку ID
func (h *ProductHandler) GetProductsBatch(c *gin.Context) {
	var request struct {
		IDs []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var products []models.Product
	if err := h.db.Where("id IN ?", request.IDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": products})
}

// BulkDeleteProducts удаляет продукты по списку ID
func (h *ProductHandler) BulkDeleteProducts(c *gin.Context) {
	var request struct {
		IDs []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Where("id IN ?", request.IDs).Delete(&models.Product{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Products deleted successfully"})
}

// GetProductsByCategory получает продукты по категории
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")

	var products []models.Product
	if err := h.db.Where("category = ?", category).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProductStock обновляет количество продукта
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	id := c.Param("id")

	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity cannot be negative"})
		return
	}

	// Обновляем количество, inStock обновится автоматически через hook
	if err := h.db.Model(&models.Product{}).Where("id = ?", productID).Update("quantity", request.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
		return
	}

	// Получаем обновленный продукт
	var updatedProduct models.Product
	h.db.First(&updatedProduct, "id = ?", productID)

	c.JSON(http.StatusOK, updatedProduct)
}

// SearchProducts поиск продуктов
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	var products []models.Product
	searchQuery := "%" + query + "%"

	if err := h.db.Where("name ILIKE ? OR description ILIKE ? OR category ILIKE ?",
		searchQuery, searchQuery, searchQuery).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
