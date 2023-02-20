package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/live"
)

func (ds *AMQPDatasource) QueryData(_ context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := ds.Query(req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}



func (ds *AMQPDatasource) Query(pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}
	var qm JsonQueryModel
	response.Error = json.Unmarshal(query.JSON, &qm)

	if response.Error != nil {
		return response
	}

	frame := data.NewFrame("response")

	channel := live.Channel{
		Scope:     live.ScopeDatasource,
		Namespace: pCtx.DataSourceInstanceSettings.UID,
		Path:      fmt.Sprintf("%v\\%v", qm.JsonKeyPath, qm.RegexValue),
	}
	frame.SetMeta(&data.FrameMeta{Channel: channel.String()})

	response.Frames = append(response.Frames, frame)

	return response
}
