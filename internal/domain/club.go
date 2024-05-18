package domain

import "time"

// Club представляет данные о компьютерном клубе
type Club struct {
	Tables    int       // Количество столов в клубе
	StartTime time.Time // Время начала работы клуба
	EndTime   time.Time // Время окончания работы клуба
	Price     int       // Стоимость часа
	Events    []Event   // Список событий
}

// Event представляет событие, происходящее в клубе
type Event struct {
	Time    time.Time // Время события
	ID      int       // Идентификатор события
	Client  string    // Имя клиента
	TableID int       // Номер стола (если применимо)
	Error   string    // Сообщение об ошибке (если есть)
}

// Table представляет данные о столе в клубе
type Table struct {
	ID       int           // Номер стола
	Revenue  int           // Выручка за день
	Duration time.Duration // Время занятости стола
}
