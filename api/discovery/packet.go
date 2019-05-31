package discovery

import (
	"errors"
	"io"
)

type Header [8]byte

func NewHeader(packetType int, blockSize, blockCount int) Header {
	return Header{1 << 4, byte(packetType), 0, 0, 0, 0, byte(blockCount >> 8), byte(blockCount)}
}

func NewHeaderWithRestBodySize(packetType int, blockSize int, restSize int) (Header, int) {
	h := NewHeader(packetType, blockSize, 1+restSize/64+1) //TODO need follow the protocol
	padding := (restSize/64+1)*64 - restSize
	return h, padding
}

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

func (this Header) BlockSize() int {
	return 64 //TODO need follow the protocol
}

type PacketType interface {
	ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error)
	WriteTo(w io.Writer) error
}

func (this Header) Type() PacketType {
	switch int(this[2-1]) {
	case 1:
		return new(RequestConnClose)
	case 2:
		return new(AcquireSessionReq)
	case 3:
		return new(AcquireSessionResp)
	case 4:
		return new(KeepAliveReq)
	case 5:
		return new(KeepAliveResp)
	case 6:
		return new(PublishServiceReq)
	case 7:
		return new(PublishServiceResp)
	default:
		return UnknownReq{}
	}
}

type UnknownReq struct {
}

func (UnknownReq) WriteTo(w io.Writer) error {
	return errors.New("unknown packet")
}

func (UnknownReq) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	return nil, errors.New("unknown packet")
}

type RequestConnClose struct {
}

func (this *RequestConnClose) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	return this, nil
}

func (this *RequestConnClose) WriteTo(w io.Writer) error {
	h := NewHeader(1, 0, 1)
	b := make([]byte, 64)
	copy(b[:8], h[:])
	_, err := w.Write(b)
	return err
}

type AcquireSessionReq struct {
}

func (this *AcquireSessionReq) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	return this, nil
}

func (this *AcquireSessionReq) WriteTo(w io.Writer) error {
	h := NewHeader(2, 0, 1)
	b := make([]byte, 64)
	copy(b[:8], h[:])
	_, err := w.Write(b)
	return err
}

type AcquireSessionResp struct {
	SessionId int64
	Timeout   int32
}

func (this *AcquireSessionResp) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	this.SessionId = bytesToInt64(prevBuf[:8])
	this.Timeout = bytesToInt32(prevBuf[8:12])
	return this, nil
}

func (this *AcquireSessionResp) WriteTo(w io.Writer) error {
	h := NewHeader(3, 0, 1)
	b := make([]byte, 64)
	copy(b[:8], h[:])
	copy(b[8:16], int64toBytes(this.SessionId))
	copy(b[16:20], int32toBytes(this.Timeout))
	_, err := w.Write(b)
	return err
}

type KeepAliveReq struct {
	SessionId int64
}

func (this *KeepAliveReq) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	this.SessionId = bytesToInt64(prevBuf[:8])
	return this, nil
}

func (this *KeepAliveReq) WriteTo(w io.Writer) error {
	h := NewHeader(4, 0, 1)
	b := make([]byte, 64)
	copy(b[:8], h[:])
	copy(b[8:16], int64toBytes(this.SessionId))
	_, err := w.Write(b)
	return err
}

type KeepAliveResp struct {
	Result byte // 0 - success, 1 - session invalid
}

func (this *KeepAliveResp) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	this.Result = prevBuf[0]
	return this, nil
}

func (this *KeepAliveResp) WriteTo(w io.Writer) error {
	h := NewHeader(5, 0, 1)
	b := make([]byte, 64)
	copy(b[:8], h[:])
	b[8] = this.Result
	_, err := w.Write(b)
	return err
}

type PublishServiceReq struct {
	SessionId int64
	Data      []byte

	Detail struct {
		Service  string
		Version  string
		Tags     []string
		NodeData []byte
	}
}

func (this *PublishServiceReq) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	this.SessionId = bytesToInt64(prevBuf[0:8])
	dataLen := bytesToInt32(prevBuf[8:12])
	restLen := h.Len() - h.BlockSize()
	if restLen < int(dataLen) {
		return nil, errors.New("length field value is not <= rest data length")
	}
	buf := make([]byte, h.Len()-h.BlockSize())
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	this.Data = buf[:dataLen]
	return this, nil
}

func (this *PublishServiceReq) WriteTo(w io.Writer) error {
	dataLen := len(this.Data)

	h, padding := NewHeaderWithRestBodySize(6, 0, len(this.Data))
	b := make([]byte, 64)
	copy(b[:8], h[:])
	copy(b[8:16], int64toBytes(this.SessionId))
	copy(b[16:20], int32toBytes(int32(dataLen)))

	_, err := w.Write(b)
	if err != nil {
		return err
	}
	_, err = w.Write(this.Data)
	if err != nil {
		return err
	}
	if padding > 0 {
		_, err = w.Write(make([]byte, padding))
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *PublishServiceReq) ParseData() error {
	//TODO
}

func (this *PublishServiceReq) GenerateData() error {
	//TODO
}

type PublishServiceResp struct {
	Result byte  // 0 - success, 1 - error
	Code   int32 // error code
}

func (this *PublishServiceResp) ParseFrom(h Header, prevBuf []byte, r io.Reader) (interface{}, error) {
	this.Result = prevBuf[0]
	this.Code = bytesToInt32(prevBuf[1:5])
	return this, nil
}

func (this *PublishServiceResp) WriteTo(w io.Writer) error {
	h := NewHeader(7, 0, 1)
	b := make([]byte, 64)
	copy(b[:8], h[:])
	b[8] = this.Result
	copy(b[9:13], int32toBytes(this.Code))
	_, err := w.Write(b)
	return err
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

	i, err := h.Type().ParseFrom(h, headerBuf[8:], reader)
	return h, i, err
}
