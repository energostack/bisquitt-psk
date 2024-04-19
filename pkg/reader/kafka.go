package reader

import (
	"context"

	"bisquitt-psk/pkg/clientmap"
	"bisquitt-psk/pkg/config"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

// KafkaReader is a reader that reads clients from a Kafka topic.
type KafkaReader struct {
	updatesCh   chan bool
	KafkaClient sarama.ConsumerGroup
	consumer    *Consumer
}

// NewKafkaReader creates a new KafkaReader with the specified configuration.
func NewKafkaReader(cfg *config.Config, ctx context.Context) (*KafkaReader, error) {
	reader := KafkaReader{
		updatesCh: make(chan bool),
	}

	reader.consumer = &Consumer{
		clientMap: clientmap.New(),
		updatesCh: reader.updatesCh,
	}

	kafkaClient, err := reader.createConsumerGroup(cfg)
	if err != nil {
		log.Err(err).Msg("Failed to create Kafka consumer group")
		return nil, err
	}

	reader.KafkaClient = kafkaClient
	reader.consume(cfg, ctx)

	return &reader, nil
}

// Read reads clients from a Kafka topic.
func (kr *KafkaReader) Read() (*clientmap.Map, error) {
	clientMap := clientmap.New()
	clientMap.Set(kr.consumer.clientMap.Get())

	return clientMap, nil
}

// Updates returns a channel that receives updates when the Kafka topic changes.
func (kr *KafkaReader) Updates() <-chan bool {
	return kr.updatesCh
}

func (kr *KafkaReader) consume(cfg *config.Config, ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				kr.KafkaClient.Close()
				return
			default:
				err := kr.KafkaClient.Consume(ctx, []string{cfg.KafkaTopic}, kr.consumer)
				if err != nil {
					log.Err(err).Msg("Failed to consume messages from Kafka")
					return
				}
			}
		}
	}()
}

func (kr *KafkaReader) createConsumerGroup(cfg *config.Config) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroup, config)
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	clientMap *clientmap.Map
	updatesCh chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	session.ResetOffset(claim.Topic(), claim.Partition(), sarama.OffsetOldest, "")
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Info().Msg("Message channel closed")
				return nil
			}

			go func() {
				consumer.clientMap.Store(string(message.Key), message.Value)
				consumer.updatesCh <- true
			}()
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
