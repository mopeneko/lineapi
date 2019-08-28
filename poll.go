package lineapi

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/mopeneko/linethrift"
)

type PollingManager struct {
	Talk       *linethrift.TalkServiceClient
	Poll       *linethrift.TalkServiceClient
	Processors map[linethrift.OpType]func(*linethrift.Operation)
}

func NewPollingManager(talk *linethrift.TalkServiceClient, options thrift.THttpClientOptions) (*PollingManager, error) {
	headers := map[string]string{
		"User-Agent":         USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access":      talk.AuthToken,
	}
	client, err := NewThriftClient(HOST+POLLING_ENDPOINT, headers, options)
	poll := linethrift.NewTalkServiceClient(client)
	if err != nil {
		return nil, err
	}
	pollcon := &PollingManager{talk, poll, map[linethrift.OpType]func(*linethrift.Operation){}}
	return pollcon, nil
}

func (p *PollingManager) SetOperationProcessor(opType linethrift.OpType, processor func(*linethrift.Operation)) {
	p.Processors[opType] = processor
}

func (p *PollingManager) StartPolling() {
	var revision int64 = 0
	ctx := context.Background()
	var err error
	revision, err = p.Talk.GetLastOpRevision(ctx)
	if err != nil {
		log.Fatal(err.Error())
		panic("Your token had been expired.")
	}
	var ops []*linethrift.Operation
	log.Printf("info: Started polling")
	for {
		ops, err = p.Poll.FetchOperations(ctx, revision, 100)
		if err != nil {
			fmt.Println(err)
			revision, _ = p.Talk.GetLastOpRevision(ctx)
			continue
		}
		for _, op := range ops {
			if op.Type == linethrift.OpType_END_OF_OPERATION {
				continue
			}
			if receiver, ok := p.Processors[op.Type]; ok {
				receiver(op)
			}
			if revision < op.Revision {
				revision = op.Revision
			}
		}
	}
}
