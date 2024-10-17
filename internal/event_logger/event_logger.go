package event_logger

import "gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/event"

type Facade interface {
	Send(event event.Event) error
}
