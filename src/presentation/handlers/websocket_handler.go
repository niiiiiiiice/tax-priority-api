package handlers

import (
	"net/http"
	"time"

	"tax-priority-api/src/application/events"
	"tax-priority-api/src/infrastructure/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebSocketHandler HTTP обработчик для WebSocket
type WebSocketHandler struct {
	hub                 *websocket.Hub
	notificationService events.NotificationService
}

// NewWebSocketHandler создает новый WebSocket обработчик
func NewWebSocketHandler(hub *websocket.Hub, notificationService events.NotificationService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:                 hub,
		notificationService: notificationService,
	}
}

// HandleWebSocket обрабатывает WebSocket подключения
// @Summary WebSocket подключение
// @Description Устанавливает WebSocket соединение для получения уведомлений в реальном времени
// @Tags WebSocket
// @Param clientId query string false "ID клиента (если не указан, будет сгенерирован)"
// @Router /ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	clientID := c.Query("clientId")
	if clientID == "" {
		clientID = uuid.New().String()
	}

	h.hub.ServeWS(c.Writer, c.Request, clientID)
}

// GetWebSocketStats возвращает статистику WebSocket подключений
// @Summary Статистика WebSocket
// @Description Возвращает статистику WebSocket подключений и подписок
// @Tags WebSocket
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ws/stats [get]
func (h *WebSocketHandler) GetWebSocketStats(c *gin.Context) {
	stats := h.notificationService.GetStats()
	c.JSON(http.StatusOK, stats)
}

// SendTestNotification отправляет тестовое уведомление
// @Summary Тестовое уведомление
// @Description Отправляет тестовое уведомление всем подключенным клиентам
// @Tags WebSocket
// @Produce json
// @Success 200 {object} gin.H
// @Router /ws/test [post]
func (h *WebSocketHandler) SendTestNotification(c *gin.Context) {
	h.notificationService.NotifySystemEvent(c.Request.Context(), "test", map[string]interface{}{
		"message":   "This is a test notification",
		"timestamp": time.Now(),
		"from":      "API",
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test notification sent",
	})
}

// BroadcastMessage отправляет сообщение всем подключенным клиентам
// @Summary Широковещательное сообщение
// @Description Отправляет сообщение всем подключенным WebSocket клиентам
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param message body BroadcastMessageRequest true "Сообщение для отправки"
// @Success 200 {object} gin.H
// @Router /ws/broadcast [post]
func (h *WebSocketHandler) BroadcastMessage(c *gin.Context) {
	var req BroadcastMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := websocket.Message{
		Type:      "broadcast",
		Event:     req.Event,
		Data:      req.Data,
		Timestamp: time.Now(),
	}

	h.hub.BroadcastToAll(message)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Broadcast message sent",
		"clients": h.hub.GetClientCount(),
	})
}

// GetConnectionInfo возвращает информацию о подключениях
// @Summary Информация о подключениях
// @Description Возвращает подробную информацию о WebSocket подключениях
// @Tags WebSocket
// @Produce json
// @Success 200 {object} ConnectionInfoResponse
// @Router /ws/info [get]
func (h *WebSocketHandler) GetConnectionInfo(c *gin.Context) {
	stats := h.notificationService.GetStats()

	response := ConnectionInfoResponse{
		Enabled:        stats["enabled"].(bool),
		TotalClients:   stats["connections"].(int),
		FAQSubscribers: stats["faq_subscribers"].(int),
		ServerTime:     time.Now(),
		AvailableEvents: []string{
			"faq.created",
			"faq.updated",
			"faq.deleted",
			"faq.activated",
			"faq.deactivated",
			"faq.priority_changed",
			"faq.category_changed",
			"faq.batch_created",
			"faq.batch_deleted",
		},
		SubscriptionTypes: []string{
			"faq",    // Все FAQ события
			"faq:ID", // События конкретного FAQ
			"system", // Системные события
		},
	}

	c.JSON(http.StatusOK, response)
}

// BroadcastMessageRequest запрос для отправки широковещательного сообщения
type BroadcastMessageRequest struct {
	Event string      `json:"event" binding:"required"`
	Data  interface{} `json:"data"`
}

// ConnectionInfoResponse ответ с информацией о подключениях
type ConnectionInfoResponse struct {
	Enabled           bool      `json:"enabled"`
	TotalClients      int       `json:"totalClients"`
	FAQSubscribers    int       `json:"faqSubscribers"`
	ServerTime        time.Time `json:"serverTime"`
	AvailableEvents   []string  `json:"availableEvents"`
	SubscriptionTypes []string  `json:"subscriptionTypes"`
}

// GetHub возвращает WebSocket хаб
func (h *WebSocketHandler) GetHub() *websocket.Hub {
	return h.hub
}
