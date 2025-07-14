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
	CreateFAQGenericRepository,
	infraRepos.NewFAQRepository,

	// Application handlers aggregators
	appHandlers.NewFAQCommandHandlers,
	appHandlers.NewFAQQueryHandlers,

	// HTTP handler
	httpHandlers.NewFAQHTTPHandler,
)

// InitializeFAQHTTPHandler инициализирует HTTP обработчик FAQ
func InitializeFAQHTTPHandler(db *gorm.DB) *httpHandlers.FAQHTTPHandler {
	wire.Build(FAQProviderSet)
	return &httpHandlers.FAQHTTPHandler{}
}
