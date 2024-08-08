package sarama_ex

type Consumer interface {
	Start(topic string) error
}
