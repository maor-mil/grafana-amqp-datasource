package amqpstreamclient

import "github.com/grafana/grafana-plugin-sdk-go/backend/log"

func failOnError(err error, msg string) {
	if err != nil {
		log.DefaultLogger.Error("%s: %s\n", msg, err)
	}
}
