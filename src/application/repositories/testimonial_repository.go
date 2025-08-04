package repositories

import (
	"tax-priority-api/src/domain/entities"
)

type TestimonialRepository interface {
	GenericRepository[*entities.Testimonial, string]
}
