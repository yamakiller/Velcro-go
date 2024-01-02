package network

import (
	"fmt"
	"strings"

	"github.com/yamakiller/velcro-go/debugs/metrics"
	"github.com/yamakiller/velcro-go/extensions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// var extensionId = extensions.NextExtensionID()

type Metrics struct {
	_metrics     *metrics.ProtoMetrics
	_enabled     bool
	_system      *NetworkSystem
	_extensionId extensions.ExtensionID
}

var _ extensions.Extension = &Metrics{}

func (m *Metrics) Enabled() bool {
	return m._enabled
}

func (m *Metrics) ExtensionID() extensions.ExtensionID {
	return m._extensionId
}

func NewMetrics(system *NetworkSystem, provider metric.MeterProvider) *Metrics {
	if provider == nil {
		return &Metrics{}
	}

	return &Metrics{
		_metrics:     metrics.NewProtoMetrics(system.getLogger(), system.Config.meriicsKey),
		_enabled:     true,
		_system:      system,
		_extensionId: extensions.NextExtensionID(),
	}
}

func (m *Metrics) PrepareSendQueueLengthGauge() {
	meter := otel.Meter(metrics.LibName)
	gauge, err := meter.Int64ObservableGauge("protonetwork_client_send_queue_length",
		metric.WithDescription("Client's Send Queue Length"),
		metric.WithUnit("1"))
	if err != nil {
		err = fmt.Errorf("failed to create ActorMailBoxLength instrument, %w", err)
		m._system.Error("[error:%s]", err.Error())
	}
	m._metrics.Instruments().SetClientSendQueueLengthGauge(gauge)
}

func (m *Metrics) CommonLabels(ctx Context) []attribute.KeyValue {
	labels := []attribute.KeyValue{
		attribute.String("address", ctx.NetworkSystem().Address()),
		attribute.String("clienttype", strings.Replace(fmt.Sprintf("%T", ctx.Client()), "*", "", 1)),
	}

	return labels
}
