package lineapi

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/mopeneko/linethrift"
)

// Generate thrift.TTransport
func NewThriftTransport(url string, headers map[string]string) (thrift.TTransport, error) {
	trans, err := thrift.NewTHttpClient(url)
	if err != nil {
		return nil, err
	}
	httpTrans := trans.(*thrift.THttpClient)
	for k, v := range headers {
		httpTrans.SetHeader(k, v)
	}
	return trans, nil
}

// Generate *thrift.TStandardClient for generating Thrift client
func NewThriftClient(url string, headers map[string]string) (*thrift.TStandardClient, error) {
	//Generate transport
	transport, err := NewThriftTransport(url, headers)
	if err != nil {
		return nil, err
	}

	// Generate client
	protocol := make([]*thrift.TCompactProtocol, 2)
	for i := 0; i < 2; i++ {
		protocol[i] = thrift.NewTCompactProtocol(transport)
	}
	client := thrift.NewTStandardClient(protocol[0], protocol[1])

	return client, nil
}

// Generate LINE TalkService client
func NewLineClient(authToken string) (*linethrift.TalkServiceClient, error) {
	headers := map[string]string{
		"User-Agent": USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access": authToken,
		"x-lpqs": "/S4",
		"X-LPV": "1",
		"X-LHM": "POST",
	}
	client, err := NewThriftClient(HOST + TALKSERVICE_ENDPOINT, headers)
	talk := linethrift.NewTalkServiceClient(client)
	if err != nil {
		return nil, err
	}
	return talk, nil
}
