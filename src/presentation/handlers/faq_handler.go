package handlers

import (
	"net/http"
	"strconv"

	"tax-priority-api/src/application/faq/commands"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/faq/handlers"
	"tax-priority-api/src/application/faq/queries"
	"tax-priority-api/src/presentation/models"

	"github.com/gin-gonic/gin"
)

// FAQHTTPHandler HTTP обработчик для FAQ
type FAQHTTPHandler struct {
	commandHandlers *handlers.FAQCommandHandlers
	queryHandlers   *handlers.FAQQueryHandlers
}

// NewFAQHTTPHandler создает новый HTTP обработчик FAQ
func NewFAQHTTPHandler(commandHandlers *handlers.FAQCommandHandlers, queryHandlers *handlers.FAQQueryHandlers) *FAQHTTPHandler {
	return &FAQHTTPHandler{
		commandHandlers: commandHandlers,
		queryHandlers:   queryHandlers,
	}
}

// GetFAQ получает FAQ по ID
// @Summary Получить FAQ по ID
// @Description Возвращает FAQ по указанному ID
// @Tags FAQ
// @Produce json
// @Param id path string true "ID FAQ"
// @Success 200 {object} models.FAQResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/{id} [get]
func (h *FAQHTTPHandler) GetFAQ(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	query := queries.GetFAQByIDQuery{ID: id}
	result, err := h.queryHandlers.GetByID.HandleGetFAQByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, dtos.ToFAQResponse(result.FAQ))
}

// GetFAQs получает список FAQ
// @Summary Получить список FAQ
// @Description Возвращает список FAQ с пагинацией и фильтрацией
// @Tags FAQ
// @Produce json
// @Param _limit query int false "Лимит записей" default(10)
// @Param _offset query int false "Смещение" default(0)
// @Param _sort query string false "Поле сортировки" default(createdAt)
// @Param _order query string false "Порядок сортировки" Enums(asc,desc) default(desc)
// @Param category query string false "Фильтр по категории"
// @Param isActive query bool false "Фильтр по активности"
// @Success 200 {object} models.PaginatedFAQResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq [get]
func (h *FAQHTTPHandler) GetFAQs(c *gin.Context) {
	// Парсим параметры запроса
	limit, _ := strconv.Atoi(c.DefaultQuery("_limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("_offset", "0"))
	sortBy := c.DefaultQuery("_sort", "createdAt")
	sortOrder := c.DefaultQuery("_order", "desc")
	category := c.Query("category")
	isActive, _ := strconv.ParseBool(c.Query("isActive"))

	// Создаем HTTP-модель запроса
	req := models.GetFAQsQuery{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Category:  category,
		IsActive:  isActive,
	}

	query := req.ToGetFAQsQuery()
	result, err := h.queryHandlers.GetMany.HandleGetFAQs(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, dtos.ToPaginatedFAQResponse(result.Paginated))
}

// UpdateFAQ обновляет FAQ
// @Summary Обновить FAQ
// @Description Обновляет существующую FAQ
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path string true "ID FAQ"
// @Param faq body models.UpdateFAQRequest true "Данные для обновления"
// @Success 200 {object} models.CommandResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/{id} [put]
func (h *FAQHTTPHandler) UpdateFAQ(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req models.UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := req.ToUpdateFAQCommand(id)
	result, err := h.commandHandlers.Update.HandleUpdateFAQ(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteFAQ удаляет FAQ
// @Summary Удалить FAQ
// @Description Удаляет FAQ по ID
// @Tags FAQ
// @Produce json
// @Param id path string true "ID FAQ"
// @Success 200 {object} models.CommandResult
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/{id} [delete]
func (h *FAQHTTPHandler) DeleteFAQ(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	cmd := commands.DeleteFAQCommand{ID: id}
	result, err := h.commandHandlers.Delete.HandleDeleteFAQ(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetFAQCount получает количество FAQ
// @Summary Получить количество FAQ
// @Description Возвращает общее количество FAQ
// @Tags FAQ
// @Produce json
// @Param category query string false "Фильтр по категории"
// @Param isActive query bool false "Фильтр по активности"
// @Success 200 {object} models.CountResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/count [get]
func (h *FAQHTTPHandler) GetFAQCount(c *gin.Context) {
	category := c.Query("category")
	isActive, _ := strconv.ParseBool(c.Query("isActive"))

	req := models.GetFAQCountQuery{
		Category: category,
		IsActive: isActive,
	}

	query := req.ToGetFAQCountQuery()
	result, err := h.queryHandlers.GetCount.HandleGetFAQCount(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, models.CountResponse{Count: result.Count})
}

// GetFAQsByIDs получает FAQ по списку ID
// @Summary Получить FAQ по списку ID
// @Description Возвращает FAQ по списку ID (batch запрос)
// @Tags FAQ
// @Accept json
// @Produce json
// @Param ids body models.GetFAQsByIDsRequest true "Список ID"
// @Success 200 {array} models.FAQResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/batch [post]
func (h *FAQHTTPHandler) GetFAQsByIDs(c *gin.Context) {
	var req models.GetFAQsByIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := req.ToGetFAQsByIDsQuery()
	result, err := h.queryHandlers.GetByIDs.HandleGetFAQsByIDs(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, dtos.ToFAQResponses(result.FAQs))
}

// BulkDeleteFAQs массовое удаление FAQ
// @Summary Массовое удаление FAQ
// @Description Удаляет несколько FAQ по списку ID
// @Tags FAQ
// @Accept json
// @Produce json
// @Param ids body models.BulkDeleteFAQRequest true "Список ID для удаления"
// @Success 200 {object} models.BatchCommandResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/bulk-delete [delete]
func (h *FAQHTTPHandler) BulkDeleteFAQs(c *gin.Context) {
	var req models.BulkDeleteFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := req.ToBulkDeleteFAQCommand()
	result, err := h.commandHandlers.BulkDelete.HandleBulkDeleteFAQ(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ActivateFAQ активирует FAQ
// @Summary Активировать FAQ
// @Description Активирует FAQ по ID
// @Tags FAQ
// @Produce json
// @Param id path string true "ID FAQ"
// @Success 200 {object} models.CommandResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/{id}/activate [patch]
func (h *FAQHTTPHandler) ActivateFAQ(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	cmd := commands.ActivateFAQCommand{ID: id}
	result, err := h.commandHandlers.Activate.HandleActivateFAQ(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeactivateFAQ деактивирует FAQ
// @Summary Деактивировать FAQ
// @Description Деактивирует FAQ по ID
// @Tags FAQ
// @Produce json
// @Param id path string true "ID FAQ"
// @Success 200 {object} models.CommandResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/{id}/deactivate [patch]
func (h *FAQHTTPHandler) DeactivateFAQ(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	cmd := commands.DeactivateFAQCommand{ID: id}
	result, err := h.commandHandlers.Deactivate.HandleDeactivateFAQ(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateFAQPriority обновляет приоритет FAQ
// @Summary Обновить приоритет FAQ
// @Description Обновляет приоритет FAQ по ID
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path string true "ID FAQ"
// @Param priority body models.UpdateFAQPriorityRequest true "Новый приоритет"
// @Success 200 {object} models.CommandResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq/{id}/priority [patch]
func (h *FAQHTTPHandler) UpdateFAQPriority(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req models.UpdateFAQPriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := req.ToUpdateFAQPriorityCommand(id)
	result, err := h.commandHandlers.UpdatePriority.HandleUpdateFAQPriority(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RegisterFAQRoutes регистрирует маршруты для FAQ
func RegisterFAQRoutes(r *gin.Engine, handler *FAQHTTPHandler) {
	api := r.Group("/api")
	faq := api.Group("/faq")
	{
		// CRUD операции
		faq.GET("/:id", handler.GetFAQ)
		faq.PUT("/:id", handler.UpdateFAQ)
		faq.DELETE("/:id", handler.DeleteFAQ)

		// Получение списков
		faq.GET("", handler.GetFAQs)
		faq.GET("/count", handler.GetFAQCount)

		// Batch операции
		faq.POST("/batch", handler.GetFAQsByIDs)
		faq.DELETE("/bulk-delete", handler.BulkDeleteFAQs)

		// Управление состоянием
		faq.PATCH("/:id/activate", handler.ActivateFAQ)
		faq.PATCH("/:id/deactivate", handler.DeactivateFAQ)
		faq.PATCH("/:id/priority", handler.UpdateFAQPriority)
	}
}
