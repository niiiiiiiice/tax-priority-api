//go:build wireinject
// +build wireinject

package wire

import (
	"log"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	appCache "tax-priority-api/src/application/cache"
	appEvents "tax-priority-api/src/application/events"
	appFaqHandlers "tax-priority-api/src/application/faq/handlers"
	appTestimonialHandlers "tax-priority-api/src/application/testimonial/handlers"
	infraCache "tax-priority-api/src/infrastructure/cache"
	infraEvents "tax-priority-api/src/infrastructure/events"
	infraPersistence "tax-priority-api/src/infrastructure/persistence"
	infraRepos "tax-priority-api/src/infrastructure/persistence/repositories"
	infraWebSocket "tax-priority-api/src/infrastructure/websocket"
	httpHandlers "tax-priority-api/src/presentation/handlers"
)

// DependencyContainer содержит все основные зависимости
type DependencyContainer struct {
	DB                  *gorm.DB
	RedisClient         *redis.Client
	Hub                 *infraWebSocket.Hub
	NotificationService appEvents.NotificationService
	Cache               appCache.Cache
}

// NewDependencyContainer создает контейнер зависимостей
func NewDependencyContainer(
	db *gorm.DB,
	redisClient *redis.Client,
	hub *infraWebSocket.Hub,
	notificationService appEvents.NotificationService,
	cache appCache.Cache,
) *DependencyContainer {
	return &DependencyContainer{
		DB:                  db,
		RedisClient:         redisClient,
		Hub:                 hub,
		NotificationService: notificationService,
		Cache:               cache,
	}
}

// CreateRedisClient создает Redis клиента для Wire (без ошибки)
func CreateRedisClient(config *infraPersistence.RedisConfig) *redis.Client {
	client, err := infraPersistence.ConnectRedis(config)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return client
}

// BaseProviderSet базовый набор провайдеров для всех модулей
var BaseProviderSet = wire.NewSet(
	// WebSocket
	infraWebSocket.NewHub,

	// Redis
	infraPersistence.NewRedisConfig,
	CreateRedisClient,

	// Cache
	appCache.NewCacheConfig,
	infraCache.NewRedisCache,

	// Events
	infraEvents.NewNotificationService,

	// Container
	NewDependencyContainer,
)

// FAQProviderSet набор провайдеров для FAQ
var FAQProviderSet = wire.NewSet(
	BaseProviderSet,

	// Repository
	CreateFAQGenericRepository,
	infraRepos.NewFAQRepository,
	infraRepos.NewCachedFAQRepository,

	// Application handlers aggregators
	appFaqHandlers.NewFAQCommandHandlers,
	appFaqHandlers.NewFAQQueryHandlers,

	// HTTP handler
	httpHandlers.NewFAQHTTPHandler,
)

// TestimonialProviderSet набор провайдеров для Testimonials
var TestimonialProviderSet = wire.NewSet(
	BaseProviderSet,

	// Repository
	CreateTestimonialGenericRepository,

	infraRepos.NewTestimonialRepository,
	infraRepos.NewCachedFAQRepository,

	// Application handlers
	appTestimonialHandlers.NewTestimonialCommandHandlers,
	appTestimonialHandlers.NewTestimonialQueryHandlers,

	// HTTP handler
	httpHandlers.NewTestimonialHTTPHandler,
)

// InitializeFAQHTTPHandler инициализирует HTTP обработчик FAQ
func InitializeFAQHTTPHandler(db *gorm.DB) *httpHandlers.FAQHTTPHandler {
	wire.Build(FAQProviderSet)
	return &httpHandlers.FAQHTTPHandler{}
}

// InitializeTestimonialHandler инициализирует HTTP обработчик Testimonials
func InitializeTestimonialHandler(db *gorm.DB) *httpHandlers.TestimonialHandler {
	wire.Build(TestimonialProviderSet)
	return &httpHandlers.TestimonialHandler{}
}

// HandlerFactory фабрика для создания обработчиков
type HandlerFactory struct {
	container *DependencyContainer
}

// NewHandlerFactory создает новую фабрику обработчиков
func NewHandlerFactory(container *DependencyContainer) *HandlerFactory {
	return &HandlerFactory{
		container: container,
	}
}

// CreateWebSocketHandler создает WebSocket обработчик
func (f *HandlerFactory) CreateWebSocketHandler() *httpHandlers.WebSocketHandler {
	return httpHandlers.NewWebSocketHandler(f.container.Hub, f.container.NotificationService)
}

// CreateFAQHandler создает FAQ обработчик
func (f *HandlerFactory) CreateFAQHandler() *httpHandlers.FAQHTTPHandler {
	return InitializeFAQHTTPHandler(f.container.DB)
}

// CreateTestimonialHandler создает Testimonial обработчик
func (f *HandlerFactory) CreateTestimonialHandler() *httpHandlers.TestimonialHandler {
	return InitializeTestimonialHandler(f.container.DB)
}

// InitializeHandlerFactory инициализирует фабрику обработчиков
func InitializeHandlerFactory(db *gorm.DB) *HandlerFactory {
	wire.Build(BaseProviderSet, NewHandlerFactory)
	return &HandlerFactory{}
}
