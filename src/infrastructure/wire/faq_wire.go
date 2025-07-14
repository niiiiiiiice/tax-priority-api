//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	appHandlers "tax-priority-api/src/application/faq/handlers"
	infraRepos "tax-priority-api/src/infrastructure/persistence/repositories"
	httpHandlers "tax-priority-api/src/presentation/handlers"
)

// FAQProviderSet набор провайдеров для FAQ
var FAQProviderSet = wire.NewSet(
	// Repository
	infraRepos.NewFAQRepository,

	// Application handlers
	appHandlers.NewFAQCommandHandler,
	appHandlers.NewFAQQueryHandler,

	// HTTP handler
	httpHandlers.NewFAQHTTPHandler,
)

// InitializeFAQHTTPHandler инициализирует HTTP обработчик FAQ
func InitializeFAQHTTPHandler(db *gorm.DB) *httpHandlers.FAQHTTPHandler {
	wire.Build(FAQProviderSet)
	return &httpHandlers.FAQHTTPHandler{}
}

// InitializeFAQCommandHandler инициализирует обработчик команд FAQ
func InitializeFAQCommandHandler(db *gorm.DB) *appHandlers.FAQCommandHandler {
	wire.Build(
		infraRepos.NewFAQRepository,
		appHandlers.NewFAQCommandHandler,
	)
	return &appHandlers.FAQCommandHandler{}
}

// InitializeFAQQueryHandler инициализирует обработчик запросов FAQ
func InitializeFAQQueryHandler(db *gorm.DB) *appHandlers.FAQQueryHandler {
	wire.Build(
		infraRepos.NewFAQRepository,
		appHandlers.NewFAQQueryHandler,
	)
	return &appHandlers.FAQQueryHandler{}
}
