package kafka

type Producer interface {
	Send(message Message) error
}
