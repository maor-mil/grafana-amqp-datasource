package amqpstreamclient

import (
	"fmt"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type Stream interface {
	CreateStream(env *stream.Environment) error
	Consume(*stream.Environment, stream.MessagesHandler) (*stream.Consumer, error)
	CloseConsumers() error
}

type StreamOptions struct {
	StreamName          string               `json:"streamName"`
	MaxAge              time.Duration        `json:"maxAge"`
	MaxLengthBytes      *stream.ByteCapacity `json:"maxLengthBytes"`
	MaxSegmentSizeBytes *stream.ByteCapacity `json:"maxSegmentSizeBytes"`
	Consumers           []*stream.Consumer
}

func NewStreamOptions() *StreamOptions {
	return &StreamOptions{}
}

func (streamOptions *StreamOptions) CreateStream(env *stream.Environment) error {
	err := env.DeclareStream(streamOptions.StreamName,
		stream.NewStreamOptions().
			SetMaxAge(streamOptions.MaxAge).
			SetMaxLengthBytes(streamOptions.MaxLengthBytes).
			SetMaxSegmentSizeBytes(streamOptions.MaxSegmentSizeBytes))

	return err
}

func (streamOptions *StreamOptions) Consume(env *stream.Environment, messagesHandler stream.MessagesHandler) (*stream.Consumer, error) {
	consumer, err := env.NewConsumer(
		streamOptions.StreamName,
		messagesHandler,
		stream.NewConsumerOptions().
			SetConsumerName("my_consumer").                  // set a consumer name
			SetOffset(stream.OffsetSpecification{}.First()). // start consuming from the beginning
			SetCRCCheck(false))                              // Disable crc control, increase the performances
	if err != nil {
		return nil, err
	}
	streamOptions.Consumers = append(streamOptions.Consumers, consumer)
	defer consumerClose(consumer.NotifyClose())
	return consumer, nil
}

func consumerClose(channelClose stream.ChannelClose) {
	event := <-channelClose
	fmt.Printf("Consumer: %s closed on the stream: %s, reason: %s \n", event.Name, event.StreamName, event.Reason)
}

func (streamOptions *StreamOptions) CloseConsumers() error {
	for consumerIndex := 0; consumerIndex < len(streamOptions.Consumers); consumerIndex += 1 {
		if err := streamOptions.Consumers[consumerIndex].Close(); err != nil {
			return err
		}
	}
	return nil
}
