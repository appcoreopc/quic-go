package logutils

import (
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/logging"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CRYPTO frame", func() {
	It("converts CRYPTO frames", func() {
		f := ConvertFrame(&wire.CryptoFrame{
			Offset: 1234,
			Data:   []byte("foobar"),
		})
		Expect(f).To(BeAssignableToTypeOf(&logging.CryptoFrame{}))
		cf := f.(*logging.CryptoFrame)
		Expect(cf.Offset).To(Equal(logging.ByteCount(1234)))
		Expect(cf.Length).To(Equal(logging.ByteCount(6)))
	})

	It("converts CRYPTO frames", func() {
		f := ConvertFrame(&wire.StreamFrame{
			StreamID: 42,
			Offset:   1234,
			Data:     []byte("foo"),
			FinBit:   true,
		})
		Expect(f).To(BeAssignableToTypeOf(&logging.StreamFrame{}))
		sf := f.(*logging.StreamFrame)
		Expect(sf.StreamID).To(Equal(logging.StreamID(42)))
		Expect(sf.Offset).To(Equal(logging.ByteCount(1234)))
		Expect(sf.Length).To(Equal(logging.ByteCount(3)))
		Expect(sf.FinBit).To(BeTrue())
	})

	It("converts other frames", func() {
		f := ConvertFrame(&wire.MaxDataFrame{ByteOffset: 1234})
		Expect(f).To(BeAssignableToTypeOf(&logging.MaxDataFrame{}))
		Expect(f).ToNot(BeAssignableToTypeOf(&logging.MaxStreamDataFrame{}))
		mdf := f.(*logging.MaxDataFrame)
		Expect(mdf.ByteOffset).To(Equal(logging.ByteCount(1234)))
	})
})
