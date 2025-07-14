package entities

import (
	"errors"
	"strings"
	"time"
)

// FAQ представляет сущность часто задаваемых вопросов
type FAQ struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Category  string    `json:"category"`
	IsActive  bool      `json:"isActive"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Реализация интерфейса Entity
func (f *FAQ) GetID() string {
	return f.ID
}

func (f *FAQ) SetID(id string) {
	f.ID = id
}

func (f *FAQ) GetCreatedAt() time.Time {
	return f.CreatedAt
}

func (f *FAQ) SetCreatedAt(t time.Time) {
	f.CreatedAt = t
}

func (f *FAQ) GetUpdatedAt() time.Time {
	return f.UpdatedAt
}

func (f *FAQ) SetUpdatedAt(t time.Time) {
	f.UpdatedAt = t
}

// Бизнес-логика

// NewFAQ создает новую FAQ сущность
func NewFAQ(question, answer, category string) (*FAQ, error) {
	faq := &FAQ{
		Question:  strings.TrimSpace(question),
		Answer:    strings.TrimSpace(answer),
		Category:  strings.TrimSpace(category),
		IsActive:  true,
		Priority:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := faq.Validate(); err != nil {
		return nil, err
	}

	return faq, nil
}

// Validate проверяет валидность FAQ
func (f *FAQ) Validate() error {
	if f.Question == "" {
		return errors.New("question cannot be empty")
	}

	if len(f.Question) < 10 {
		return errors.New("question must be at least 10 characters long")
	}

	if len(f.Question) > 500 {
		return errors.New("question cannot exceed 500 characters")
	}

	if f.Answer == "" {
		return errors.New("answer cannot be empty")
	}

	if len(f.Answer) < 10 {
		return errors.New("answer must be at least 10 characters long")
	}

	if len(f.Answer) > 2000 {
		return errors.New("answer cannot exceed 2000 characters")
	}

	if f.Category == "" {
		return errors.New("category cannot be empty")
	}

	if len(f.Category) > 100 {
		return errors.New("category cannot exceed 100 characters")
	}

	if f.Priority < 0 || f.Priority > 100 {
		return errors.New("priority must be between 0 and 100")
	}

	return nil
}

// UpdateQuestion обновляет вопрос
func (f *FAQ) UpdateQuestion(question string) error {
	f.Question = strings.TrimSpace(question)
	f.UpdatedAt = time.Now()
	return f.Validate()
}

// UpdateAnswer обновляет ответ
func (f *FAQ) UpdateAnswer(answer string) error {
	f.Answer = strings.TrimSpace(answer)
	f.UpdatedAt = time.Now()
	return f.Validate()
}

// UpdateCategory обновляет категорию
func (f *FAQ) UpdateCategory(category string) error {
	f.Category = strings.TrimSpace(category)
	f.UpdatedAt = time.Now()
	return f.Validate()
}

// SetPriority устанавливает приоритет
func (f *FAQ) SetPriority(priority int) error {
	if priority < 0 || priority > 100 {
		return errors.New("priority must be between 0 and 100")
	}
	f.Priority = priority
	f.UpdatedAt = time.Now()
	return nil
}

// Activate активирует FAQ
func (f *FAQ) Activate() {
	f.IsActive = true
	f.UpdatedAt = time.Now()
}

// Deactivate деактивирует FAQ
func (f *FAQ) Deactivate() {
	f.IsActive = false
	f.UpdatedAt = time.Now()
}

// IsValidForPublishing проверяет готовность к публикации
func (f *FAQ) IsValidForPublishing() bool {
	return f.IsActive && f.Question != "" && f.Answer != "" && f.Category != ""
}

// GetSearchableText возвращает текст для поиска
func (f *FAQ) GetSearchableText() string {
	return strings.ToLower(f.Question + " " + f.Answer + " " + f.Category)
}
