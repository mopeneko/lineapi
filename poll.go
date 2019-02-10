package lineapi

import (
	"context"
	"log"
	"github.com/comail/colog"
	"github.com/mopeneko/linethrift"
)

type PollingManager struct {
	Talk *linethrift.TalkServiceClient
	Poll *linethrift.TalkServiceClient
	Processors map[linethrift.OpType]interface{}
}

func NewPollingManager(talk *linethrift.TalkServiceClient) (*PollingManager, error) {
	headers := map[string]string{
		"User-Agent": USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access": talk.AuthToken,
		"x-lpqs": "/P4",
		"X-LPV": "1",
		"X-LHM": "POST",
	}
	client, err := NewThriftClient(HOST + POLLING_ENDPOINT, headers)
	poll := linethrift.NewTalkServiceClient(client)
	if err != nil {
		return nil, err
	}
	pollcon := &PollingManager{talk, poll, map[linethrift.OpType]interface{}{}}
	return pollcon, nil
}

func (p *PollingManager) SetOperationProcessor(opType linethrift.OpType, processor interface{}) {
	p.Processors[opType] = processor
}

func (p *PollingManager) ProcessOperations(isLogged bool) {
	var revision int64 = 0
	ctx := context.Background()
	var err error
	revision, err = p.Talk.GetLastOpRevision(ctx)
	if err != nil {
		panic("Your token had been expired.")
	}
	var operations []*linethrift.Operation
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors:	true,
    	Flag:	log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
	log.Printf("info: Started polling")
	for {
		operations, err = p.Poll.FetchOperations(ctx, revision, 100)
		if err != nil {
			newRevision, err := p.Talk.GetLastOpRevision(ctx)
			if err != nil {
				panic("Your token had been expired.")
			}
			revision = newRevision
			log.Printf("error: %v", err)
			continue
		}
		for _, operation := range operations {
			if isLogged {
				log.Printf("debug: [%d]\t%s", int64(operation.Type), operation.Type.String())
			}
			if operation.Type != linethrift.OpType_END_OF_OPERATION {
				if processor, isContain := p.Processors[operation.Type]; isContain {
					processor.(func(*linethrift.Operation))(operation)
				}
				if revision < operation.Revision {
					revision = operation.Revision
				}
			}
		}
	}
}

func (p *PollingManager) StartPolling() {
	p.ProcessOperations(false)
}

func (p *PollingManager) StartPollingWithLogging() {
	p.ProcessOperations(true)
}
