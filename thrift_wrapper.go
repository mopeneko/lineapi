package lineapi

import (
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"
)

var HTTPClient = http.DefaultClient

// Generate thrift.TTransport
func NewThriftTransport(url string, headers map[string]string, options thrift.THttpClientOptions) (thrift.TTransport, error) {
	trans, err := thrift.NewTHttpClientWithOptions(url, options)
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
func NewThriftClient(url string, headers map[string]string, options thrift.THttpClientOptions) (*thrift.TStandardClient, error) {
	//Generate transport
	transport, err := NewThriftTransport(url, headers, options)
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
