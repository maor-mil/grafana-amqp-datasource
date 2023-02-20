package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (ds *AMQPDatasource) CheckHealth(_ context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if !ds.Client.IsConnected() {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "AMQP Disconnected",
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "AMQP Connected",
	}, nil
}
