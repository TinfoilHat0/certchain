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
}

//CreateSkipchain creates a new skipchain
func (s *Service) CreateSkipchain(cs *CreateSkipchainRequest) (*CreateSkipchainResponse, onet.ClientError) {
	client := skipchain.NewClient()
	sb, err := client.CreateGenesis(cs.Roster, 1, 1, []skipchain.VerifierID{VerifyMerkleTreeRoot}, nil, nil) //create genesis&store skip block calls the verification function
	if err != nil {
		return nil, err
	}
	//log.Print(sb)
	return &CreateSkipchainResponse{sb}, nil
}

//AddMerkleTreeRoot stores a merkle tree root in the blockhain
func (s *Service) AddMerkleTreeRoot(mtr *AddMerkleTreeRootRequest) (*AddMerkleTreeRootResponse, onet.ClientError) {
	//Call VerifyMerkleTreeRoot here, add to SkipChain only if it returns true
	client := skipchain.NewClient()
	sb, err := client.StoreSkipBlock(mtr.SkipBlock, nil, mtr.TreeRoot) //nil will be replaced by storeskipblock dat
	if err != nil {
		return nil, err
	}
	return &AddMerkleTreeRootResponse{sb.Latest}, nil
}

//VerifyMerkleTreeRoot verifies a signed Merkle tree root
func (s *Service) VerifyMerkleTreeRoot(newID []byte, newSB *skipchain.SkipBlock) bool {
	//Input: Signed Merkle tree root from the client and its K_p
	//Run verification algorithm here, depending on it ret true or false
	log.Print("Verify is called!")
	//log.Print(s.ServerIdentity())
	return false
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
	log.ErrFatal(skipchain.RegisterVerification(c, VerifyMerkleTreeRoot, s.VerifyMerkleTreeRoot))
	return s
}
