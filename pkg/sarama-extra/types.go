package sarama_extra

type Consumer interface {
	Start(topic string) error
}
