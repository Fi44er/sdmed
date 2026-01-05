package chat_entity

import "time"

type MessageType string

const (
	// REGULAR — Основной тип сообщения.
	// Содержит текст (Payload) и/или любые файлы (Attachments: картинки, документы).
	// Именно этот тип используется для 99% общения.
	MsgTypeRegular MessageType = "regular"

	// SYSTEM — Технические уведомления.
	// "Оператор Иван подключился", "Чат закрыт", "Чат переведен на отдел доставки".
	// Фронтенд рендерит это без "пузыря", просто текстом по центру.
	MsgTypeSystem MessageType = "system"

	// ORDER_CARD — Тот самый тип для ваших заказов.
	// В Payload или Metadata кладется JSON заказа.
	// Фронтенд рисует красивую карточку заказа с кнопками и статусом.
	MsgTypeOrderCard MessageType = "order_card"

	// FEEDBACK — Карточка оценки качества.
	// Отправляется автоматически при закрытии чата.
	// Содержит кнопки или звезды (1-5).
	MsgTypeFeedback MessageType = "feedback"

	// ACTION — Интерактивные кнопки/меню.
	// Если у вас будет бот или "быстрые ответы" (например, кнопки "Где мой заказ?", "Проблема с оплатой").
	MsgTypeAction MessageType = "action"
)

type Message struct {
	ID       string
	ChatID   string // В какой чат летит (обязательно для индексов в БД)
	SenderID string // Кто отправил (User.ID)

	Type    MessageType
	Payload string // Текст сообщения или JSON с метаданными

	// Вложения (Attachments)
	Attachments []Attachment

	// Для функций мессенджера
	ReplyToID string // ID сообщения, на которое отвечаем
	IsEdited  bool

	// Статусы доставки
	CreatedAt time.Time
	UpdatedAt time.Time
	ReadAt    *time.Time // nil если не прочитано
}

type Attachment struct {
	ID       string
	FileName string
	FileSize int64
	MimeType string
	URL      string // Временная ссылка для фронтенда
}
