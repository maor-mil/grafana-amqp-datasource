package plugin

import (
	"encoding/json"

	"github.com/maor2475/amqp/pkg/amqpstreamclient"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces- only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*AMQPDatasource)(nil)
	_ backend.CheckHealthHandler    = (*AMQPDatasource)(nil)
	_ backend.StreamHandler         = (*AMQPDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*AMQPDatasource)(nil)
)

func NewAMQPInstance(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	client, err := getDatasourceSettings(s)
	if err != nil {
		return nil, err
	}

	_, err = client.Connect()
	if err != nil {
		return nil, err
	}

	return NewAMQPDatasource(client), nil
}

type AMQPDatasource struct {
	Client amqpstreamclient.Client
}

func NewAMQPDatasource(client amqpstreamclient.Client) *AMQPDatasource {
	return &AMQPDatasource{
		Client: client,
	}
}

// Dispose here tells plugin SDK that plugin wants to clean up resources
// when a new instance created. As soon as datasource settings change detected
// by SDK old datasource instance will be disposed and a new one will be created
// using NewMQTTDatasource factory function.
func (ds *AMQPDatasource) Dispose() {
	ds.Client.Dispose()
}

func getDatasourceSettings(s backend.DataSourceInstanceSettings) (*amqpstreamclient.AmqpStreamClient, error) {
	client := amqpstreamclient.NewAmqpStreamClient()

	if err := json.Unmarshal(s.JSONData, client); err != nil {
		return nil, err
	}

	if password, exists := s.DecryptedSecureJSONData["password"]; exists {
		client.AmqpOptions.Password = password
	}

	return client, nil
}
