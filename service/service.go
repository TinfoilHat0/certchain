package template

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

// Name is the name to refer to the Template service from another
// package.
const Name = "CertChain"

func init() {
	onet.RegisterNewService(Name, newService)
}

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor
	path string
	// Count holds the number of calls to 'ClockRequest'
	Count int
}

// ClockRequest starts a template-protocol and returns the run-time.
func (s *Service) CreateSkipchain(cs *CreateSkipchainRequest) (*CreateSkipchainResponse, onet.ClientError) {
	log.Print("create skipchain")
	client := skipchain.NewClient()
	sb, err := client.CreateGenesis(cs.Roster, 1, 1, []skipchain.VerifierID{VerifyCert}, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CreateSkipchainResponse{sb}, nil
}

// CountRequest returns the number of instantiations of the protocol.
func (s *Service) AddMerkleTreeRoot(mtr *AddMerkleTreeRootRequest) (*AddMerkleTreeRootResponse, onet.ClientError) {
	log.Print("create mtr")
	client := skipchain.NewClient()
	sb, err := client.StoreSkipBlock(mtr.SkipBlock, nil, mtr.TreeRoot)
	if err != nil {
		return nil, err
	}
	return &AddMerkleTreeRootResponse{sb.Latest}, nil
}

func (s *Service) verifyCert(newID []byte, newSB *skipchain.SkipBlock) bool {
	log.Print(s.ServerIdentity())
	return true
}

// newTemplate receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *onet.Context) onet.Service {
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	if err := s.RegisterHandlers(s.CreateSkipchain, s.AddMerkleTreeRoot); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	// call s.verifyCert when called
	log.ErrFatal(skipchain.RegisterVerification(c, VerifyCert, s.verifyCert))
	return s
}
