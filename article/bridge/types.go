package bridge

import "github.com/IBM/sarama"

type LargeMessagesProducer sarama.SyncProducer

type SmallMessagesProducer sarama.SyncProducer
