package metrics

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/lucas-clemente/quic-go/logging"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Measures
var (
	connections = stats.Int64("connections", "number of QUIC connections", stats.UnitDimensionless)
)

// Tags
var (
	keyPerspective, _ = tag.NewKey("perspective")
	keyIPVersion, _   = tag.NewKey("ip_version")
)

// Views
var (
	ConnectionsView = &view.View{
		Measure:     connections,
		TagKeys:     []tag.Key{keyPerspective, keyIPVersion},
		Aggregation: view.Count(),
	}
)

func init() {
	if err := view.Register(connectionsView); err != nil {
		log.Fatalf("Failed to register view: %s", err)
	}

	go func() {
		pe, err := prometheus.NewExporter(prometheus.Options{
			Namespace: "quic",
		})
		if err != nil {
			log.Fatal(err)
		}
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(":8888", mux); err != nil {
			log.Fatalf("Failed to run Prometheus /metrics endpoint: %v", err)
		}
	}()
}

type tracer struct{}

var _ logging.Tracer = &tracer{}

// NewTracer creates a new metrics tracer.
func NewTracer() logging.Tracer { return &tracer{} }

func (t *tracer) TracerForServer(logging.ConnectionID) logging.ConnectionTracer {
	return newConnTracer(t, logging.PerspectiveServer)
}

func (t *tracer) TracerForClient(logging.ConnectionID) logging.ConnectionTracer {
	return newConnTracer(t, logging.PerspectiveClient)
}

type connTracer struct {
	perspective logging.Perspective
	tracer      logging.Tracer
}

func newConnTracer(tracer logging.Tracer, perspective logging.Perspective) logging.ConnectionTracer {
	return &connTracer{
		perspective: perspective,
		tracer:      tracer,
	}
}

var _ logging.ConnectionTracer = &connTracer{}

func (t *connTracer) StartedConnection(local, _ net.Addr, _ logging.VersionNumber, _, _ logging.ConnectionID) {
	var perspectiveTag tag.Mutator
	switch t.perspective {
	case logging.PerspectiveClient:
		perspectiveTag = tag.Upsert(keyPerspective, "client")
	case logging.PerspectiveServer:
		perspectiveTag = tag.Upsert(keyPerspective, "server")
	}

	var ipVersionTag tag.Mutator
	if udpAddr, ok := local.(*net.UDPAddr); ok {
		// If ip is not an IPv4 address, To4 returns nil.
		// Note that there might be some corner cases, where this is not correct.
		// See https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6.
		if udpAddr.IP.To4() == nil {
			ipVersionTag = tag.Upsert(keyIPVersion, "IPv6")
		} else {
			ipVersionTag = tag.Upsert(keyIPVersion, "IPv4")
		}
	} else {
		ipVersionTag = tag.Upsert(keyIPVersion, "unknown")
	}

	stats.RecordWithTags(
		context.Background(),
		[]tag.Mutator{perspectiveTag, ipVersionTag},
		connections.M(1),
	)
}

func (t *connTracer) ClosedConnection(logging.CloseReason)                     {}
func (t *connTracer) SentTransportParameters(*logging.TransportParameters)     {}
func (t *connTracer) ReceivedTransportParameters(*logging.TransportParameters) {}
func (t *connTracer) SentPacket(*logging.ExtendedHeader, logging.ByteCount, *logging.AckFrame, []logging.Frame) {
}
func (t *connTracer) ReceivedVersionNegotiationPacket(*logging.Header) {}
func (t *connTracer) ReceivedRetry(*logging.Header)                    {}
func (t *connTracer) ReceivedPacket(*logging.ExtendedHeader, logging.ByteCount, []logging.Frame) {
}
func (t *connTracer) ReceivedStatelessReset(token *[16]byte)                                        {}
func (t *connTracer) BufferedPacket(logging.PacketType)                                             {}
func (t *connTracer) DroppedPacket(logging.PacketType, logging.ByteCount, logging.PacketDropReason) {}
func (t *connTracer) UpdatedMetrics(*logging.RTTStats, logging.ByteCount, logging.ByteCount, int)   {}
func (t *connTracer) LostPacket(logging.EncryptionLevel, logging.PacketNumber, logging.PacketLossReason) {
}
func (t *connTracer) UpdatedPTOCount(value uint32)                                       {}
func (t *connTracer) UpdatedKeyFromTLS(logging.EncryptionLevel, logging.Perspective)     {}
func (t *connTracer) UpdatedKey(logging.KeyPhase, bool)                                  {}
func (t *connTracer) DroppedEncryptionLevel(logging.EncryptionLevel)                     {}
func (t *connTracer) SetLossTimer(logging.TimerType, logging.EncryptionLevel, time.Time) {}
func (t *connTracer) LossTimerExpired(logging.TimerType, logging.EncryptionLevel)        {}
func (t *connTracer) LossTimerCanceled()                                                 {}
func (t *connTracer) Close()                                                             {}
