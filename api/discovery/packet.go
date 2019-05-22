package discovery

import (
	"errors"
	"io"
)

type Header [8]byte

func parseHeader(buf []byte) Header {
	h := Header{}
	copy(h[:], buf)
	return h
}

func (this Header) Len() int {
	//blockSize := this[6-1] //TODO need follow the protocol
	blockSize := 64
	blockCnt := ((int(this[8-1]) & 0xFF) << 8) + (int(this[8-1]) & 0xFF)
	return blockSize * blockCnt
}

type PacketType int

const (
	TypeRequestConnClose PacketType = iota + 1
	TypeAcquireSessionReq
	TypeAcquireSessionResp
	TypeKeepAliveReq
	TypeKeepAliveResp
	TypePublishServiceReq
	TypePublishServiceResp

	UNKNOWN PacketType = 0xFFFFFFFF
)

func (this Header) Type() PacketType {
	switch int(this[2-1]) {
	case 1:
		return TypeRequestConnClose
	case 2:
		return TypeAcquireSessionReq
	case 3:
		return TypeAcquireSessionResp
	case 4:
		return TypeKeepAliveReq
	case 5:
		return TypeKeepAliveResp
	case 6:
		return TypePublishServiceReq
	case 7:
		return TypePublishServiceResp
	default:
		return UNKNOWN
	}
}

func (PacketType) ParseFrom(prevBuf []byte, r io.Reader) (interface{}, error) {
	return nil, errors.New("unknown packet")
}

type RequestConnClose struct {
}

type AcquireSessionReq struct {
}

type AcquireSessionResp struct {
}

type KeepAliveReq struct {
}

type KeepAliveResp struct {
}

type PublishServiceReq struct {
}

type PublishServiceResp struct {
}

// packet format:
// header: 8bytes
//      1: 0001 0000 - high 4 bits protocol version, low 4 bits options are 0 and reserved
//      2: packet type, self defined
//      3: 0 reserved
//      4: 0 reserved
//      5: streaming mode - high 6 bits streaming id, low 2 bits streaming block type
//              000000 XX - non streaming packet
//              YYYYYY 00 - start of streaming id, could have packet body
//              YYYYYY 11 - end of streaming id, no packet body
//              YYYYYY ZZ - streaming block packet
//      6: 0000 0000 - high 4 bits reserved option, low 4 bits block size
//                     block size mapping:
//                     0000 - 64 bytes
//                     0001 - 256 bytes
//                     0010 - 1024 bytes
//                     0011 - 4096 bytes
//                     other - reserved
//    7-8: block count of request, minimum is 1. 0 is illegal. includes header (8 bytes).
func ParseFrom(reader io.Reader) (Header, interface{}, error) {
	headerBuf := make([]byte, 64)
	_, err := io.ReadFull(reader, headerBuf)
	if err != nil {
		return Header{}, nil, err
	}
	h := parseHeader(headerBuf[:8])

	i, err := h.Type().ParseFrom(headerBuf[8:], reader)
	return h, i, err
}
