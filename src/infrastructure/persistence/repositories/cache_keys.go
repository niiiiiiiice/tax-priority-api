package repositories

const (
	FAQCategoriesKey       = "faq:categories"
	FAQCategoriesWithCount = "faq:categories:with_counts"

	FAQCategoriesPattern = "faq:categories*"
)

func GenerateFAQCategoriesKey(withCounts bool) string {
	if withCounts {
		return FAQCategoriesWithCount
	}
	return FAQCategoriesKey
}
