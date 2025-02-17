package forwarder

import (
	"bytes"
	"context"
	"fmt"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/utils"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const taskStartedEvent = `{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "my-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.task.started",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`
const taskFinishedEvent = `{
				"data": "",
				"id": "5de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "my-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.task.fnished",
				"shkeptncontext": "c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`

func Test_ForwardEventsToNATS(t *testing.T) {
	expectedReceivedMessageCount := 0

	svr, shutdownNats := runNATSServer()
	defer shutdownNats()

	cfg := config.EnvConfig{}
	envconfig.Process("", &cfg)
	cfg.PubSubURL = svr.Addr().String()

	natsClient, err := nats.Connect(svr.Addr().String())
	if err != nil {
		t.Errorf("could not initialize nats client: %s", err.Error())
	}
	defer natsClient.Close()
	_, _ = natsClient.Subscribe("sh.keptn.events.task.*", func(m *nats.Msg) {
		expectedReceivedMessageCount++
	})

	apiset, _ := keptnapi.New(config.DefaultShipyardControllerBaseURL)
	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		keptnEventAPI:     apiset.APIV1(),
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
		env:               cfg,
	}

	ctx, cancel := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go f.Start(executionContext)

	time.Sleep(2 * time.Second)
	numEvents := 1000
	for i := 0; i < numEvents; i++ {
		eventFromService(taskFinishedEvent)
	}

	assert.Eventually(t, func() bool {
		return expectedReceivedMessageCount == numEvents
	}, time.Second*time.Duration(10), time.Second)

	cancel()
	executionContext.Wg.Wait()
}

func Test_ForwardEventsToKeptnAPI(t *testing.T) {

	receivedMessageCount := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) { receivedMessageCount++ }))

	cfg := config.EnvConfig{}
	envconfig.Process("", &cfg)
	cfg.KeptnAPIEndpoint = ts.URL
	apiset, _ := keptnapi.New(ts.URL)

	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		keptnEventAPI:     apiset.APIV1(),
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
		env:               cfg,
	}
	ctx, cancel := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go f.Start(executionContext)

	//TODO: remove waiting
	time.Sleep(2 * time.Second)
	eventFromService(taskStartedEvent)
	eventFromService(taskFinishedEvent)

	assert.Eventually(t, func() bool {
		return receivedMessageCount == 2
	}, time.Second*time.Duration(10), time.Second)
	cancel()
	executionContext.Wg.Wait()
}

func Test_APIProxy(t *testing.T) {
	proxyEndpointCalled := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			proxyEndpointCalled++
		}))

	cfg := config.EnvConfig{}
	envconfig.Process("", &cfg)
	cfg.KeptnAPIEndpoint = ""
	config.InClusterAPIProxyMappings["/testpath"] = strings.TrimPrefix(ts.URL, "http://")

	apiset, _ := keptnapi.New(ts.URL)

	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		keptnEventAPI:     apiset.APIV1(),
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
		env:               cfg,
	}
	ctx, cancel := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go f.Start(executionContext)

	//TODO: remove wait
	time.Sleep(2 * time.Second)
	apiCallFromService()

	assert.Eventually(t, func() bool {
		return proxyEndpointCalled == 1
	}, time.Second*time.Duration(10), time.Second)

	cancel()
	executionContext.Wg.Wait()
}

func apiCallFromService() {
	http.Get(fmt.Sprintf("http://127.0.0.1:%d/testpath", 8081))

}

func eventFromService(event string) {
	payload := bytes.NewBuffer([]byte(event))
	http.Post(fmt.Sprintf("http://127.0.0.1:%d/event", 8081), "application/cloudevents+json", payload)
}

func runNATSServer() (*server.Server, func()) {
	svr := natsserver.RunRandClientPortServer()
	return svr, func() { svr.Shutdown() }
}
