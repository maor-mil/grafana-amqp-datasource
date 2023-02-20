package plugin

import (
	"context"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func (ds *AMQPDatasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	var qm JsonQueryModel
	chunks := strings.Split(req.Path, "\\")
	qm.JsonKeyPath = chunks[0]
	qm.RegexValue = chunks[1]

	handleMessages := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		log.DefaultLogger.Info("Full Byte Array: %s", message)
		for i := 0; i < len(message.Data); i += 1 {
			log.DefaultLogger.Info("Byte Array At %d: %s", i, message.Data[i])
		}
	}

	ds.Client.Consume(handleMessages)
	defer ds.Client.CloseConsumers()

	for {
		select {
		case <-ctx.Done():
			log.DefaultLogger.Debug("stopped streaming (context canceled)")
			ds.Client.CloseConsumers()
			return nil
		default:
			log.DefaultLogger.Info("Shit is happening")
		}
	}
}

func (ds *AMQPDatasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	return &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusOK,
	}, nil
}

func (ds *AMQPDatasource) PublishStream(_ context.Context, _ *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}
