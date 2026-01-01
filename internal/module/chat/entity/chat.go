package chat_entity

import "time"

type ChatStatus string

const (
	StatusNew      ChatStatus = "new"      // Только создан, оператор еще не зашел
	StatusOpen     ChatStatus = "open"     // Оператор взял в работу
	StatusPending  ChatStatus = "pending"  // Ждем ответа от пользователя
	StatusClosed   ChatStatus = "closed"   // Вопрос решен
	StatusArchived ChatStatus = "archived" // Перенесен в архив
)

type Chat struct {
	ID         string
	ClientID   string // Ссылка на User.ID (клиент)
	OperatorID string // Текущий ответственный оператор (может меняться)

	Subject  string // Тема обращения (опционально)
	Status   ChatStatus
	Priority int // 1-Низкий, 2-Средний, 3-Высокий

	Tags []string // Например: ["billing", "bug", "ios"]

	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time // Ссылка, так как может быть nil

	// Метаданные для аналитики
	FirstResponseTime *time.Duration // Время до первого ответа оператора
	Rating            int            // Оценка пользователя после закрытия (1-5)
}
