package amqpstreamclient

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type Client interface {
	IsConnected() bool
	Connect() (Client, error)
	Reconnect() Client
	Consume(stream.MessagesHandler) (interface{}, error)
	CloseConsumers() error
	Dispose() error
}

type AmqpStreamOptions struct {
	Host                  string             `json:"host"`
	StreamPort            int                `json:"streamPort"`
	AmqpPort              int                `json:"amqpPort"`
	VHost                 string             `json:"vhost"`
	User                  string             `json:"user"`
	Password              string             `json:"password"`
	MaxProducersPerClient int                `json:"maxProducersPerClient"`
	MaxConsumersPerClient int                `json:"maxConsumersPerClient"`
	IsTLS                 bool               `json:"isTLS"`
	TLSConfig             bool               `json:"TLSConfig"`
	RequestedHeartbeat    time.Duration      `json:"requestedHeartbeat"`
	RequestedMaxFrameSize int                `json:"requestedMaxFrameSize"`
	WriteBuffer           int                `json:"writeBuffer"`
	ReadBuffer            int                `json:"readBuffer"`
	NoDelay               bool               `json:"noDelay"`
	StreamOptions         *StreamOptions     `json:"streamOptions"`
	ExchangesOptions      []*ExchangeOptions `json:"exchanges"`
	BindingsOptions       []*BindingOptions  `json:"bindings"`
}

type AmqpStreamClient struct {
	AmqpOptions *AmqpStreamOptions `json:"amqpStreamSettings"`
	Env         *stream.Environment
	Stream      Stream
	Exchanges   []Exchange
	Bindings    []Binding
}

const timeToReconnect time.Duration = 2000 * time.Millisecond

func NewAmqpStreamClient() *AmqpStreamClient {
	return &AmqpStreamClient{}
}

func (client *AmqpStreamClient) SetEnv() (*AmqpStreamClient, error) {
	// Connect to the broker
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(client.AmqpOptions.Host).
			SetPort(client.AmqpOptions.StreamPort).
			SetVHost(client.AmqpOptions.VHost).
			SetUser(client.AmqpOptions.User).
			SetPassword(client.AmqpOptions.Password).
			SetMaxProducersPerClient(client.AmqpOptions.MaxProducersPerClient).
			SetMaxConsumersPerClient(client.AmqpOptions.MaxConsumersPerClient).
			IsTLS(client.AmqpOptions.IsTLS).
			//SetTLSConfig(&tls.Config{}).
			SetRequestedHeartbeat(client.AmqpOptions.RequestedHeartbeat * time.Second).
			SetRequestedMaxFrameSize(client.AmqpOptions.RequestedMaxFrameSize).
			SetWriteBuffer(client.AmqpOptions.WriteBuffer).
			SetReadBuffer(client.AmqpOptions.ReadBuffer).
			SetNoDelay(client.AmqpOptions.NoDelay),
	)

	client.Env = env

	return client, err
}

func (client *AmqpStreamClient) SetStream() *AmqpStreamClient {
	client.Stream = client.AmqpOptions.StreamOptions
	return client
}

func (client *AmqpStreamClient) SetExchanges() *AmqpStreamClient {
	for exchangeIndex := 0; exchangeIndex < len(client.AmqpOptions.ExchangesOptions); exchangeIndex += 1 {
		client.Exchanges = append(client.Exchanges, client.AmqpOptions.ExchangesOptions[exchangeIndex])
	}
	return client
}

func (client *AmqpStreamClient) SetBindings() *AmqpStreamClient {
	for bindingIndex := 0; bindingIndex < len(client.AmqpOptions.BindingsOptions); bindingIndex += 1 {
		client.Bindings = append(client.Bindings, client.AmqpOptions.BindingsOptions[bindingIndex])
	}
	return client
}

func (client *AmqpStreamClient) CreateStream() (*AmqpStreamClient, error) {
	return client, client.Stream.CreateStream(client.Env)
}

func (client *AmqpStreamClient) CreateExchanges() (*AmqpStreamClient, error) {
	for exchangeIndex := 0; exchangeIndex < len(client.AmqpOptions.ExchangesOptions); exchangeIndex += 1 {
		if err := client.Exchanges[exchangeIndex].CreateExchange(client.AmqpOptions); err != nil {
			return client, err
		}
	}
	return client, nil
}

func (client *AmqpStreamClient) CreateBindings() (*AmqpStreamClient, error) {
	for bindingIndex := 0; bindingIndex < len(client.AmqpOptions.BindingsOptions); bindingIndex += 1 {
		if err := client.Bindings[bindingIndex].CreateBinding(client.AmqpOptions); err != nil {
			return client, err
		}
	}
	return client, nil
}

func (client *AmqpStreamClient) IsConnected() bool {
	return !client.Env.IsClosed()
}

func (client *AmqpStreamClient) Connect() (Client, error) {
	_, err := client.SetEnv()
	if err != nil {
		return client, err
	}

	client.SetStream()
	client.SetExchanges()
	client.SetBindings()

	_, err = client.CreateStream()
	if err != nil {
		return client, err
	}
	_, err = client.CreateExchanges()
	if err != nil {
		return client, err
	}
	_, err = client.CreateBindings()
	if err != nil {
		return client, err
	}

	return client, nil
}

func (client *AmqpStreamClient) CloseConnection() error {
	if err := client.CloseConsumers(); err != nil {
		return err
	}
	if err := client.Env.DeleteStream(client.AmqpOptions.StreamOptions.StreamName); err != nil {
		return err
	}
	if err := client.Env.Close(); err != nil {
		return err
	}
	return nil
}

func (client *AmqpStreamClient) Reconnect() Client {
	for {
		time.Sleep(timeToReconnect)
		log.DefaultLogger.Info(
			"Trying to reconnect to AMQP Stream: {AMQP Host: %v ; Stream Name: %v}",
			client.AmqpOptions.Host,
			client.AmqpOptions.StreamOptions.StreamName,
		)
		err := client.CloseConnection()
		if err != nil {
			continue
		}
		_, err = client.Connect()
		if err != nil {
			continue
		}
		break
	}
	return client
}

func (client *AmqpStreamClient) Consume(messageHandler stream.MessagesHandler) (interface{}, error) {
	return client.Stream.Consume(client.Env, messageHandler)
}

func (client *AmqpStreamClient) CloseConsumers() error {
	return client.Stream.CloseConsumers()
}

func (client *AmqpStreamClient) Dispose() error {
	log.DefaultLogger.Info(
		"Disposing AMQP Stream: {AMQP Host: %v ; Stream Name: %v}",
		client.AmqpOptions.Host,
		client.AmqpOptions.StreamOptions.StreamName,
	)
	return client.CloseConnection()
}

/*
Default stream env options:
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetVHost("/").
			SetUser("guest").
			SetPassword("guest").
			SetMaxProducersPerClient(1).
			SetMaxConsumersPerClient(1).
			IsTLS(false).
			// SetTLSConfig(&tls.Config{}).
			SetRequestedHeartbeat(60 * time.Second).
			SetRequestedMaxFrameSize(1048576).
			SetWriteBuffer(8192).
			SetReadBuffer(65536).
			SetNoDelay(false),
*/
