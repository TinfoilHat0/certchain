package template

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/crypto/random"
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
	prevMTR := &MerkleTreeRoot{make([]byte, 5)}
	latestMTR := &MerkleTreeRoot{random.Bytes(4, random.Stream)}
	genesisData := &CertBlock{prevMTR, latestMTR, cs.PublicKey}
	sb, err := client.CreateGenesis(cs.Roster, 1, 1, []skipchain.VerifierID{VerifyMerkleTreeRoot}, genesisData, nil) //create genesis&store skip block calls the verification function
	if err != nil {
		return nil, err
	}
	return &CreateSkipchainResponse{sb}, nil
}

//AddNewTxn stores a new transaction in the underlying Skipchain
func (s *Service) AddNewTxn(txnRequest *AddNewTxnRequest) (*AddNewTxnResponse, onet.ClientError) {
	client := skipchain.NewClient()
	sb, err := client.StoreSkipBlock(txnRequest.SkipBlock, nil, txnRequest.CertBlock) //Is this the proper way to call that ?
	if err != nil {
		return nil, err
	}
	return &AddNewTxnResponse{sb.Latest}, nil
}

//VerifyMerkleTreeRoot verifies a signed Merkle tree root
func (s *Service) VerifyMerkleTreeRoot(newID []byte, newSB *skipchain.SkipBlock) bool {
	//What does newID contain? How can I access the tree roots inside this function ?
	//Run verification algorithm here, depending on it ret true or false
	log.Print("Verify is called!")
	return true
}

// newService receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *onet.Context) onet.Service {
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	if err := s.RegisterHandlers(s.CreateSkipchain, s.AddNewTxn); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	log.ErrFatal(skipchain.RegisterVerification(c, VerifyMerkleTreeRoot, s.VerifyMerkleTreeRoot))
	return s
}
