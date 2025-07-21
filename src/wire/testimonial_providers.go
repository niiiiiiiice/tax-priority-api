package wire

import (
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/handlers"
	persistenceRepos "tax-priority-api/src/infrastructure/persistence/repositories"
	presentationHandlers "tax-priority-api/src/presentation/handlers"

	"gorm.io/gorm"
)

// TestimonialProviders предоставляет зависимости для testimonials
func ProvideTestimonialRepository(db *gorm.DB) repositories.TestimonialRepository {
	return persistenceRepos.NewTestimonialRepository(db)
}

func ProvideTestimonialCommandHandlers(repo repositories.TestimonialRepository) *handlers.TestimonialCommandHandlers {
	return handlers.NewTestimonialCommandHandlers(repo)
}

func ProvideTestimonialQueryHandlers(repo repositories.TestimonialRepository) *handlers.TestimonialQueryHandlers {
	return handlers.NewTestimonialQueryHandlers(repo)
}

func ProvideTestimonialHandler(
	commandHandlers *handlers.TestimonialCommandHandlers,
	queryHandlers *handlers.TestimonialQueryHandlers,
) *presentationHandlers.TestimonialHTTPHandler {
	return presentationHandlers.NewTestimonialHandler(commandHandlers, queryHandlers)
}
