package api

type MessageMarshal interface {
	MarshalPublishRequest(p *PublishRequest) ([]byte, error)
	UnmarshalPublishRequest(data []byte) (*PublishRequest, error)
	MarshalPublishResponse(p *PublishResponse) ([]byte, error)
	UnmarshalPublishResponse(data []byte) (*PublishResponse, error)

	MarshalAcquireRequest(a *AcquireRequest) ([]byte, error)
	UnmarshalAcquireRequest(data []byte) (*AcquireRequest, error)
	MarshalAcquireResponse(a *AcquireResponse) ([]byte, error)
	UnmarshalAcquireResponse(data []byte) (*AcquireResponse, error)
}

var (
	marshallingProvider = make(map[string]MessageMarshal)
)

func RegisterMarshallingProvider(name string, m MessageMarshal) {
	marshallingProvider[name] = m
}

func GetmarshallingProvider(name string) MessageMarshal {
	m, ok := marshallingProvider[name]
	if !ok {
		return nil
	} else {
		return m
	}
}

type PublishRequest struct {
	Topic       string
	MessageList []*struct {
		MsgId    string
		MsgOutId string
		MsgBody  []byte
	}
}

type PublishResponse struct {
	SuccessIdList []*struct {
		MsgId    string
		MsgOutId string
	}
	FailIdList []*struct {
		MsgId    string
		MsgOutId string
		Code     string
	}
}

type AcquireRequest struct {
	Queue string
}

type AcquireResponse struct {
	MessageList []*struct {
		MsgId   string
		MsgBody []byte
	}
}
