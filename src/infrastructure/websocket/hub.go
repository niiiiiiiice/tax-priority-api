package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client представляет WebSocket клиента
type Client struct {
	ID            string
	Conn          *websocket.Conn
	Hub           *Hub
	Send          chan []byte
	Subscriptions map[string]bool // Подписки клиента
	mu            sync.RWMutex
}

// Hub управляет WebSocket подключениями
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

// Message представляет сообщение WebSocket
type Message struct {
	Type      string      `json:"type"`
	Event     string      `json:"event,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SubscriptionMessage сообщение для подписки/отписки
type SubscriptionMessage struct {
	Type          string   `json:"type"` // "subscribe" или "unsubscribe"
	Subscriptions []string `json:"subscriptions"`
}

// EventMessage сообщение о событии
type EventMessage struct {
	Type      string      `json:"type"`
	Event     string      `json:"event"`
	Entity    string      `json:"entity"`
	EntityID  string      `json:"entityId,omitempty"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// В продакшене здесь должна быть проверка origin
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewHub создает новый WebSocket хаб
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run запускает хаб
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client %s connected. Total clients: %d", client.ID, len(h.clients))

			// Отправляем приветственное сообщение
			welcome := Message{
				Type:      "system",
				Event:     "connected",
				Data:      map[string]string{"clientId": client.ID},
				Timestamp: time.Now(),
			}
			client.SendMessage(welcome)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
				log.Printf("Client %s disconnected. Total clients: %d", client.ID, len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
			h.mu.RUnlock()

		case <-ctx.Done():
			log.Println("WebSocket hub shutting down...")
			return
		}
	}
}

// ServeWS обрабатывает WebSocket подключения
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, clientID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:            clientID,
		Conn:          conn,
		Hub:           h,
		Send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}

	h.register <- client

	// Запускаем горутины для чтения и записи
	go client.writePump()
	go client.readPump()
}

// BroadcastToAll отправляет сообщение всем подключенным клиентам
func (h *Hub) BroadcastToAll(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.broadcast <- data
}

// BroadcastToSubscribers отправляет сообщение только подписчикам
func (h *Hub) BroadcastToSubscribers(subscription string, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		client.mu.RLock()
		if client.Subscriptions[subscription] {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client.ID)
			}
		}
		client.mu.RUnlock()
	}
}

// BroadcastEvent отправляет событие подписчикам
func (h *Hub) BroadcastEvent(entity, action, entityID string, data interface{}) {
	subscription := entity // Подписка по типу сущности
	eventMessage := EventMessage{
		Type:      "event",
		Event:     entity + "." + action,
		Entity:    entity,
		EntityID:  entityID,
		Action:    action,
		Data:      data,
		Timestamp: time.Now(),
	}

	if entityID != "" {
		// Также отправляем подписчикам конкретной сущности
		specificSubscription := entity + ":" + entityID
		h.BroadcastEventMessage(specificSubscription, eventMessage)
	}

	// Отправляем общим подписчикам типа сущности
	h.BroadcastEventMessage(subscription, eventMessage)
}

// BroadcastEventMessage отправляет событие подписчикам
func (h *Hub) BroadcastEventMessage(subscription string, eventMessage EventMessage) {
	data, err := json.Marshal(eventMessage)
	if err != nil {
		log.Printf("Error marshaling event message: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		client.mu.RLock()
		if client.Subscriptions[subscription] {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client.ID)
			}
		}
		client.mu.RUnlock()
	}
}

// GetClientCount возвращает количество подключенных клиентов
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetSubscriptionCount возвращает количество подписок для конкретной подписки
func (h *Hub) GetSubscriptionCount(subscription string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, client := range h.clients {
		client.mu.RLock()
		if client.Subscriptions[subscription] {
			count++
		}
		client.mu.RUnlock()
	}
	return count
}

// writePump отправляет сообщения клиенту
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Добавляем дополнительные сообщения из очереди
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump читает сообщения от клиента
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Обрабатываем сообщение от клиента
		c.handleMessage(message)
	}
}

// handleMessage обрабатывает сообщения от клиента
func (c *Client) handleMessage(message []byte) {
	var msg SubscriptionMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling client message: %v", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	switch msg.Type {
	case "subscribe":
		for _, subscription := range msg.Subscriptions {
			c.Subscriptions[subscription] = true
			log.Printf("Client %s subscribed to %s", c.ID, subscription)
		}

		// Отправляем подтверждение подписки
		response := Message{
			Type:      "subscription",
			Event:     "subscribed",
			Data:      map[string]interface{}{"subscriptions": msg.Subscriptions},
			Timestamp: time.Now(),
		}
		c.SendMessage(response)

	case "unsubscribe":
		for _, subscription := range msg.Subscriptions {
			delete(c.Subscriptions, subscription)
			log.Printf("Client %s unsubscribed from %s", c.ID, subscription)
		}

		// Отправляем подтверждение отписки
		response := Message{
			Type:      "subscription",
			Event:     "unsubscribed",
			Data:      map[string]interface{}{"subscriptions": msg.Subscriptions},
			Timestamp: time.Now(),
		}
		c.SendMessage(response)

	default:
		log.Printf("Unknown message type from client %s: %s", c.ID, msg.Type)
	}
}

// SendMessage отправляет сообщение клиенту
func (c *Client) SendMessage(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for client %s: %v", c.ID, err)
		return
	}

	select {
	case c.Send <- data:
	default:
		close(c.Send)
		delete(c.Hub.clients, c.ID)
	}
}

// GetSubscriptions возвращает список подписок клиента
func (c *Client) GetSubscriptions() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	subscriptions := make([]string, 0, len(c.Subscriptions))
	for subscription := range c.Subscriptions {
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions
}
