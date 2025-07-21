package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tax-priority-api/src/application/testimonial/dtos"
	"tax-priority-api/src/application/testimonial/handlers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestimonialHTTPHandler struct {
	commandHandlers *handlers.TestimonialCommandHandlers
	queryHandlers   *handlers.TestimonialQueryHandlers
}

func NewTestimonialHandler(
	commandHandlers *handlers.TestimonialCommandHandlers,
	queryHandlers *handlers.TestimonialQueryHandlers,
) *TestimonialHTTPHandler {
	return &TestimonialHTTPHandler{
		commandHandlers: commandHandlers,
		queryHandlers:   queryHandlers,
	}
}

// CreateTestimonial создает новый отзыв
// @Summary Создать отзыв
// @Description Создает новый отзыв с возможностью загрузки файла
// @Tags testimonials
// @Accept multipart/form-data
// @Produce json
// @Param content formData string true "Содержание отзыва"
// @Param author formData string true "Автор отзыва"
// @Param authorEmail formData string true "Email автора"
// @Param rating formData int true "Рейтинг (1-5)"
// @Param company formData string false "Компания"
// @Param position formData string false "Должность"
// @Param file formData file false "Файл (PDF или изображение)"
// @Success 201 {object} dtos.CommandResult
// @Failure 400 {object} dtos.CommandResult
// @Failure 500 {object} dtos.CommandResult
// @Router /testimonials [post]
func (h *TestimonialHTTPHandler) CreateTestimonial(c *gin.Context) {
	var cmd dtos.CreateTestimonialCommand

	// Получаем данные из формы
	cmd.Content = c.PostForm("content")
	cmd.Author = c.PostForm("author")
	cmd.AuthorEmail = c.PostForm("authorEmail")
	cmd.Company = c.PostForm("company")
	cmd.Position = c.PostForm("position")

	// Парсим рейтинг
	ratingStr := c.PostForm("rating")
	if ratingStr != "" {
		rating, err := strconv.Atoi(ratingStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dtos.CommandResult{
				Success:   false,
				Error:     "Invalid rating format",
				Timestamp: time.Now(),
			})
			return
		}
		cmd.Rating = rating
	}

	// Создаем отзыв
	result, err := h.commandHandlers.CreateTestimonial(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	// Обрабатываем загрузку файла, если он есть
	file, fileHeader, err := c.Request.FormFile("file")
	if err == nil && file != nil {
		defer file.Close()

		// Проверяем тип файла
		allowedTypes := map[string]bool{
			"application/pdf": true,
			"image/jpeg":      true,
			"image/png":       true,
			"image/gif":       true,
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, dtos.CommandResult{
				Success:   false,
				Error:     "Unsupported file type. Only PDF and images are allowed",
				Timestamp: time.Now(),
			})
			return
		}

		// Создаем папку для файлов, если её нет
		uploadDir := "uploads/testimonials"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, dtos.CommandResult{
				Success:   false,
				Error:     "Failed to create upload directory",
				Timestamp: time.Now(),
			})
			return
		}

		// Генерируем уникальное имя файла
		ext := filepath.Ext(fileHeader.Filename)
		fileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, fileName)

		// Сохраняем файл
		dst, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dtos.CommandResult{
				Success:   false,
				Error:     "Failed to save file",
				Timestamp: time.Now(),
			})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			c.JSON(http.StatusInternalServerError, dtos.CommandResult{
				Success:   false,
				Error:     "Failed to save file",
				Timestamp: time.Now(),
			})
			return
		}

		// Обновляем отзыв с информацией о файле
		if testimonial, ok := result.Data.(*dtos.CommandResult); ok {
			if testimonialData, ok := testimonial.Data.(map[string]interface{}); ok {
				if id, ok := testimonialData["id"].(string); ok {
					fileCmd := dtos.UploadTestimonialFileCommand{
						ID:       id,
						FilePath: filePath,
						FileName: fileHeader.Filename,
						FileType: contentType,
						FileSize: fileHeader.Size,
					}
					// Здесь можно добавить обработку команды загрузки файла
					_ = fileCmd
				}
			}
		}
	}

	c.JSON(http.StatusCreated, result)
}

// GetTestimonials получает список отзывов
// @Summary Получить список отзывов
// @Description Получает список отзывов с пагинацией и фильтрацией
// @Tags testimonials
// @Accept json
// @Produce json
// @Param limit query int false "Лимит записей" default(10)
// @Param offset query int false "Смещение" default(0)
// @Param sortBy query string false "Поле для сортировки" default("createdAt")
// @Param sortOrder query string false "Порядок сортировки" Enums(asc, desc) default("desc")
// @Param approved query bool false "Фильтр по статусу одобрения"
// @Param rating query int false "Фильтр по рейтингу"
// @Success 200 {object} dtos.QueryResult
// @Failure 400 {object} dtos.QueryResult
// @Failure 500 {object} dtos.QueryResult
// @Router /testimonials [get]
func (h *TestimonialHTTPHandler) GetTestimonials(c *gin.Context) {
	var query dtos.GetTestimonialsQuery

	// Парсим параметры запроса
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			query.Offset = offset
		}
	}

	query.SortBy = c.Query("sortBy")
	query.SortOrder = c.Query("sortOrder")

	// Парсим фильтры
	filters := make(map[string]interface{})

	if approvedStr := c.Query("approved"); approvedStr != "" {
		if approved, err := strconv.ParseBool(approvedStr); err == nil {
			filters["is_approved"] = approved
		}
	}

	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil {
			filters["rating"] = rating
		}
	}

	if author := c.Query("author"); author != "" {
		filters["author ILIKE"] = fmt.Sprintf("%%%s%%", author)
	}

	query.Filters = filters

	result, err := h.queryHandlers.GetTestimonials(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTestimonialByID получает отзыв по ID
// @Summary Получить отзыв по ID
// @Description Получает отзыв по указанному ID
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path string true "ID отзыва"
// @Success 200 {object} dtos.QueryResult
// @Failure 404 {object} dtos.QueryResult
// @Failure 500 {object} dtos.QueryResult
// @Router /testimonials/{id} [get]
func (h *TestimonialHTTPHandler) GetTestimonialByID(c *gin.Context) {
	id := c.Param("id")

	query := dtos.GetTestimonialByIDQuery{ID: id}
	result, err := h.queryHandlers.GetTestimonialByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetApprovedTestimonials получает одобренные отзывы
// @Summary Получить одобренные отзывы
// @Description Получает список одобренных и активных отзывов
// @Tags testimonials
// @Accept json
// @Produce json
// @Param limit query int false "Лимит записей" default(10)
// @Param offset query int false "Смещение" default(0)
// @Param sortBy query string false "Поле для сортировки" default("createdAt")
// @Param sortOrder query string false "Порядок сортировки" Enums(asc, desc) default("desc")
// @Success 200 {object} dtos.QueryResult
// @Failure 500 {object} dtos.QueryResult
// @Router /testimonials/approved [get]
func (h *TestimonialHTTPHandler) GetApprovedTestimonials(c *gin.Context) {
	var query dtos.GetApprovedTestimonialsQuery

	// Парсим параметры запроса
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			query.Offset = offset
		}
	}

	query.SortBy = c.Query("sortBy")
	query.SortOrder = c.Query("sortOrder")

	result, err := h.queryHandlers.GetApprovedTestimonials(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTestimonialStats получает статистику отзывов
// @Summary Получить статистику отзывов
// @Description Получает статистику по отзывам (количество, рейтинги и т.д.)
// @Tags testimonials
// @Accept json
// @Produce json
// @Success 200 {object} dtos.QueryResult
// @Failure 500 {object} dtos.QueryResult
// @Router /testimonials/stats [get]
func (h *TestimonialHTTPHandler) GetTestimonialStats(c *gin.Context) {
	query := dtos.GetTestimonialStatsQuery{}
	result, err := h.queryHandlers.GetTestimonialStats(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateTestimonial обновляет отзыв
// @Summary Обновить отзыв
// @Description Обновляет существующий отзыв
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path string true "ID отзыва"
// @Param testimonial body dtos.UpdateTestimonialCommand true "Данные для обновления"
// @Success 200 {object} dtos.CommandResult
// @Failure 400 {object} dtos.CommandResult
// @Failure 404 {object} dtos.CommandResult
// @Failure 500 {object} dtos.CommandResult
// @Router /testimonials/{id} [put]
func (h *TestimonialHTTPHandler) UpdateTestimonial(c *gin.Context) {
	id := c.Param("id")

	var cmd dtos.UpdateTestimonialCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CommandResult{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	cmd.ID = id
	result, err := h.commandHandlers.UpdateTestimonial(c.Request.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, result)
		} else {
			c.JSON(http.StatusInternalServerError, result)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// ApproveTestimonial одобряет отзыв
// @Summary Одобрить отзыв
// @Description Одобряет отзыв для публикации
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path string true "ID отзыва"
// @Param approveData body dtos.ApproveTestimonialCommand true "Данные для одобрения"
// @Success 200 {object} dtos.CommandResult
// @Failure 400 {object} dtos.CommandResult
// @Failure 404 {object} dtos.CommandResult
// @Failure 500 {object} dtos.CommandResult
// @Router /testimonials/{id}/approve [patch]
func (h *TestimonialHTTPHandler) ApproveTestimonial(c *gin.Context) {
	id := c.Param("id")

	var cmd dtos.ApproveTestimonialCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CommandResult{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	cmd.ID = id
	result, err := h.commandHandlers.ApproveTestimonial(c.Request.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, result)
		} else {
			c.JSON(http.StatusInternalServerError, result)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteTestimonial удаляет отзыв
// @Summary Удалить отзыв
// @Description Удаляет отзыв по ID
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path string true "ID отзыва"
// @Success 200 {object} dtos.CommandResult
// @Failure 404 {object} dtos.CommandResult
// @Failure 500 {object} dtos.CommandResult
// @Router /testimonials/{id} [delete]
func (h *TestimonialHTTPHandler) DeleteTestimonial(c *gin.Context) {
	id := c.Param("id")

	cmd := dtos.DeleteTestimonialCommand{ID: id}
	result, err := h.commandHandlers.DeleteTestimonial(c.Request.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, result)
		} else {
			c.JSON(http.StatusInternalServerError, result)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// BulkApproveTestimonials массово одобряет отзывы
// @Summary Массово одобрить отзывы
// @Description Одобряет несколько отзывов одновременно
// @Tags testimonials
// @Accept json
// @Produce json
// @Param bulkApprove body dtos.BulkApproveTestimonialsCommand true "Данные для массового одобрения"
// @Success 200 {object} dtos.CommandResult
// @Failure 400 {object} dtos.CommandResult
// @Failure 500 {object} dtos.CommandResult
// @Router /testimonials/bulk/approve [patch]
func (h *TestimonialHTTPHandler) BulkApproveTestimonials(c *gin.Context) {
	var cmd dtos.BulkApproveTestimonialsCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CommandResult{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	result, err := h.commandHandlers.BulkApproveTestimonials(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// BulkDeleteTestimonials массово удаляет отзывы
// @Summary Массово удалить отзывы
// @Description Удаляет несколько отзывов одновременно
// @Tags testimonials
// @Accept json
// @Produce json
// @Param bulkDelete body dtos.BulkDeleteTestimonialsCommand true "Данные для массового удаления"
// @Success 200 {object} dtos.CommandResult
// @Failure 400 {object} dtos.CommandResult
// @Failure 500 {object} dtos.CommandResult
// @Router /testimonials/bulk/delete [delete]
func (h *TestimonialHTTPHandler) BulkDeleteTestimonials(c *gin.Context) {
	var cmd dtos.BulkDeleteTestimonialsCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CommandResult{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	result, err := h.commandHandlers.BulkDeleteTestimonials(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterTestimonialRoutes(router *gin.Engine, handler *TestimonialHTTPHandler) {
	testimonialGroup := router.Group("/testimonials")
	{
		// Публичные маршруты
		testimonialGroup.GET("/approved", handler.GetApprovedTestimonials)
		testimonialGroup.GET("/stats", handler.GetTestimonialStats)
		testimonialGroup.POST("", handler.CreateTestimonial)

		// Маршруты для управления отзывами
		testimonialGroup.GET("", handler.GetTestimonials)
		testimonialGroup.GET("/:id", handler.GetTestimonialByID)
		testimonialGroup.PUT("/:id", handler.UpdateTestimonial)
		testimonialGroup.DELETE("/:id", handler.DeleteTestimonial)
		testimonialGroup.PATCH("/:id/approve", handler.ApproveTestimonial)

		// Массовые операции
		bulkGroup := testimonialGroup.Group("/bulk")
		{
			bulkGroup.PATCH("/approve", handler.BulkApproveTestimonials)
			bulkGroup.DELETE("/delete", handler.BulkDeleteTestimonials)
		}
	}
}
