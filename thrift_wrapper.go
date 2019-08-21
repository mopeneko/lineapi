package lineapi

import (
	"net"
	"net/http"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
)

var HTTPClient = http.DefaultClient

// Generate thrift.TTransport
func NewThriftTransport(url string, headers map[string]string) (thrift.TTransport, *thrift.THttpClient, error) {
	trans, err := thrift.NewTHttpClient(url)
	if err != nil {
		return nil, nil, err
	}
	httpTrans := trans.(*thrift.THttpClient)
	for k, v := range headers {
		httpTrans.SetHeader(k, v)
	}
	return trans, httpTrans, nil
}

// Generate *thrift.TStandardClient for generating Thrift client
func NewThriftClient(url string, headers map[string]string) (*thrift.TStandardClient, *thrift.THttpClient, error) {
	//Generate transport
	transport, thttp, err := NewThriftTransport(url, headers)
	if err != nil {
		return nil, nil, err
	}

	// Generate client
	protocol := make([]*thrift.TCompactProtocol, 2)
	for i := 0; i < 2; i++ {
		protocol[i] = thrift.NewTCompactProtocol(transport)
	}
	client := thrift.NewTStandardClient(protocol[0], protocol[1])

	return client, thttp, nil
}

func NewThriftTransportForLP(url string, headers map[string]string) (thrift.TTransport, *thrift.THttpClient, error) {
	trans, err := thrift.NewTHttpClientWithOptions(url, thrift.THttpClientOptions{&http.Client{Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          5000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}}})
	if err != nil {
		return nil, nil, err
	}
	httpTrans := trans.(*thrift.THttpClient)
	for k, v := range headers {
		httpTrans.SetHeader(k, v)
	}
	return trans, httpTrans, nil
}

func NewThriftClientForLP(url string, headers map[string]string) (*thrift.TStandardClient, *thrift.THttpClient, error) {
	//Generate transport
	transport, thttp, err := NewThriftTransportForLP(url, headers)
	if err != nil {
		return nil, nil, err
	}

	// Generate client
	protocol := make([]*thrift.TCompactProtocol, 2)
	for i := 0; i < 2; i++ {
		protocol[i] = thrift.NewTCompactProtocol(transport)
	}
	client := thrift.NewTStandardClient(protocol[0], protocol[1])

	return client, thttp, nil
}
