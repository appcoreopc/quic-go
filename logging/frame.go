package logging

import "github.com/lucas-clemente/quic-go/internal/wire"

type Frame interface{}

type AckRange = wire.AckRange

type (
	AckFrame                = wire.AckFrame
	ConnectionCloseFrame    = wire.ConnectionCloseFrame
	DataBlockedFrame        = wire.DataBlockedFrame
	HandshakeDoneFrame      = wire.HandshakeDoneFrame
	MaxDataFrame            = wire.MaxDataFrame
	MaxStreamDataFrame      = wire.MaxStreamDataFrame
	MaxStreamsFrame         = wire.MaxStreamsFrame
	NewConnectionIDFrame    = wire.NewConnectionIDFrame
	NewTokenFrame           = wire.NewTokenFrame
	PathChallengeFrame      = wire.PathChallengeFrame
	PathResponseFrame       = wire.PathResponseFrame
	PingFrame               = wire.PingFrame
	ResetStreamFrame        = wire.ResetStreamFrame
	RetireConnectionIDFrame = wire.RetireConnectionIDFrame
	StopSendingFrame        = wire.StopSendingFrame
	StreamsBlockedFrame     = wire.StreamsBlockedFrame
	StreamDataBlockedFrame  = wire.StreamDataBlockedFrame
)

type CryptoFrame struct {
	Offset ByteCount
	Length ByteCount
}

type StreamFrame struct {
	StreamID StreamID
	Offset   ByteCount
	Length   ByteCount
	FinBit   bool
}
