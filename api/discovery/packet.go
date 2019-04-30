package discovery

import "io"

type Header uint64

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
//      1: 0001 0000 - high 4 bits protocol version, low 4 bits option reserved
//      2: packet type, self defined
//      3:
//      4:
//      5: streaming mode - high 3 bits streaming id, low 5 bits streaming block offset
//              000 XXXXX - non streaming packet
//              YYY 11111 - end of streaming id, no matter how many blocks
//              YYY ZZZZZ - streaming block packet
//      6: 0000 0000 - high 4 bits reserved option, low 4 bits block size
//                     block size mapping:
//                     0000 - 64 bytes
//                     0001 - 256 bytes
//                     0010 - 1024 bytes
//                     0011 - 4096 bytes
//                     other - reserved
//    7-8: block count of request, minimum is 1. 0 is illegal
func ParseFrom(reader io.Reader) (Header, interface{}, error) {
	return 0, nil, nil
}
