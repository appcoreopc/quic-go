// Package logging defines a logging interface for quic-go.
// This package should not be considered stable
package logging

import (
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/internal/congestion"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

type (
	ByteCount           = protocol.ByteCount
	ConnectionID        = protocol.ConnectionID
	EncryptionLevel     = protocol.EncryptionLevel
	KeyPhase            = protocol.KeyPhase
	PacketNumber        = protocol.PacketNumber
	Perspective         = protocol.Perspective
	StreamID            = protocol.StreamID
	StreamNum           = protocol.StreamNum
	StreamType          = protocol.StreamType
	VersionNumber       = protocol.VersionNumber
	Header              = wire.Header
	ExtendedHeader      = wire.ExtendedHeader
	TransportParameters = wire.TransportParameters

	RTTStats = congestion.RTTStats
)

const (
	// PerspectiveServer is used for a QUIC server
	PerspectiveServer Perspective = protocol.PerspectiveServer
	// PerspectiveClient is used for a QUIC client
	PerspectiveClient Perspective = protocol.PerspectiveClient
)

const (
	// StreamTypeUni is a unidirectional stream
	StreamTypeUni = protocol.StreamTypeUni
	// StreamTypeBidi is a bidirectional stream
	StreamTypeBidi = protocol.StreamTypeBidi
)

type Tracer interface {
	TracerForServer(odcid ConnectionID) ConnectionTracer
	TracerForClient(odcid ConnectionID) ConnectionTracer
}

// A ConnectionTracer records events.
type ConnectionTracer interface {
	StartedConnection(local, remote net.Addr, version VersionNumber, srcConnID, destConnID ConnectionID)
	ClosedConnection(CloseReason)
	SentTransportParameters(*TransportParameters)
	ReceivedTransportParameters(*TransportParameters)
	SentPacket(hdr *ExtendedHeader, packetSize ByteCount, ack *wire.AckFrame, frames []wire.Frame)
	ReceivedVersionNegotiationPacket(*Header)
	ReceivedRetry(*Header)
	ReceivedPacket(hdr *ExtendedHeader, packetSize ByteCount, frames []wire.Frame)
	ReceivedStatelessReset(token *[16]byte)
	BufferedPacket(PacketType)
	DroppedPacket(PacketType, ByteCount, PacketDropReason)
	UpdatedMetrics(rttStats *RTTStats, cwnd ByteCount, bytesInFLight ByteCount, packetsInFlight int)
	LostPacket(EncryptionLevel, PacketNumber, PacketLossReason)
	UpdatedPTOCount(value uint32)
	UpdatedKeyFromTLS(EncryptionLevel, Perspective)
	UpdatedKey(generation KeyPhase, remote bool)
	DroppedEncryptionLevel(EncryptionLevel)
	SetLossTimer(TimerType, EncryptionLevel, time.Time)
	LossTimerExpired(TimerType, EncryptionLevel)
	LossTimerCanceled()
	// Close is called when the connection is closed.
	Close()
}
