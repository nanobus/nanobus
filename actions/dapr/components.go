package dapr

import (
	"context"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/actors"
	"github.com/dapr/dapr/pkg/channel"
	"github.com/dapr/dapr/pkg/config"
	"github.com/dapr/dapr/pkg/messaging"
	invokev1 "github.com/dapr/dapr/pkg/messaging/v1"
	"github.com/dapr/dapr/pkg/runtime"
	"github.com/dapr/dapr/pkg/runtime/embedded"
)

type DaprComponents struct {
	Entities []string

	Actors          actors.Actors
	DirectMessaging messaging.DirectMessaging
	StateStores     map[string]state.Store
	InputBindings   map[string]bindings.InputBinding
	OutputBindings  map[string]bindings.OutputBinding
	SecretStores    map[string]secretstores.SecretStore
	PubSubs         map[string]pubsub.PubSub

	invokeHandler       InvokeHandler
	inputBindingHandler InputBindingHandler
	pubsubHandler       PubSubHandler
}

type InvokeHandler func(ctx context.Context, method, contentType string, data []byte, metadata map[string][]string) ([]byte, string, error)
type InputBindingHandler func(ctx context.Context, event *embedded.BindingEvent) ([]byte, error)
type PubSubHandler func(ctx context.Context, event *embedded.TopicEvent) (embedded.EventResponseStatus, error)

func (c *DaprComponents) RegisterComponents(reg runtime.ComponentRegistry) error {
	c.Actors = reg.Actors
	c.DirectMessaging = reg.DirectMessaging
	c.StateStores = reg.StateStores
	c.InputBindings = reg.InputBindings
	c.OutputBindings = reg.OutputBindings
	c.SecretStores = reg.SecretStores
	c.PubSubs = reg.PubSubs
	return nil
}

func (c *DaprComponents) CreateLocalChannel(port, maxConcurrency int, spec config.TracingSpec, sslEnabled bool, maxRequestBodySize int, readBufferSize int) (channel.AppChannel, error) {
	return c, nil
}

func (c *DaprComponents) GetBaseAddress() string {
	return "http://localhost:32321"
}

func (c *DaprComponents) GetAppConfig() (*config.ApplicationConfig, error) {
	return &config.ApplicationConfig{
		Entities: c.Entities,
	}, nil
}

func (c *DaprComponents) InvokeMethod(ctx context.Context, req *invokev1.InvokeMethodRequest) (*invokev1.InvokeMethodResponse, error) {
	if c.invokeHandler != nil {
		msg := req.Message()
		md := req.Metadata()
		metadata := make(map[string][]string, len(md))
		for k, v := range md {
			metadata[k] = v.Values
		}
		resp, contentType, err := c.invokeHandler(ctx, msg.Method, msg.ContentType, msg.Data.Value, metadata)
		if err != nil {
			return nil, err
		}

		response := invokev1.NewInvokeMethodResponse(200, "OK", nil)
		response.WithRawData(resp, contentType)
		return response, nil
	}
	return nil, nil
}

func (c *DaprComponents) OnBindingEvent(ctx context.Context, event *embedded.BindingEvent) ([]byte, error) {
	if c.inputBindingHandler != nil {
		return c.inputBindingHandler(ctx, event)
	}
	return nil, nil
}

func (c *DaprComponents) OnTopicEvent(ctx context.Context, event *embedded.TopicEvent) (embedded.EventResponseStatus, error) {
	if c.pubsubHandler != nil {
		return c.pubsubHandler(ctx, event)
	}
	return embedded.EventResponseStatusSuccess, nil
}

func (c *DaprComponents) InvokeHandler(invokeHandler InvokeHandler) {
	c.invokeHandler = invokeHandler
}

func (c *DaprComponents) InputBindingHandler(inputBindingHandler InputBindingHandler) {
	c.inputBindingHandler = inputBindingHandler
}

func (c *DaprComponents) PubSubHandler(pubsubHandler PubSubHandler) {
	c.pubsubHandler = pubsubHandler
}
