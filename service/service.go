package template

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/crypto/random"
	"github.com/dedis/onet/log"
	"gopkg.in/dedis/onet.v1"
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
	//publicKey := cs.Roster.Publics gives error
	prevMTR := &MerkleTreeRoot{make([]byte, 4)} //for the very first transaction, prevMTR is just a slice of 0 bytes
	curMTR := &MerkleTreeRoot{random.Bytes(4, random.Stream)}
	genesisTxn := &CertBlock{prevMTR, curMTR, nil, false}                                                           //how to access the public key of the service?
	sb, err := client.CreateGenesis(cs.Roster, 1, 1, []skipchain.VerifierID{VerifyMerkleTreeRoot}, genesisTxn, nil) //create genesis&store skip block calls the verification function
	if err != nil {
		return nil, err
	}
	//log.Print(network.Unmarshal(sb.Data))
	return &CreateSkipchainResponse{sb}, nil //What does sb.Data contain at this point? Marshalled genesisTxn ?
}

//AddNewTransaction stores a new transaction in the underlying Skipchain
func (s *Service) AddNewTransaction(txnRequest *AddNewTransactionRequest) (*AddNewTransactionResponse, onet.ClientError) {
	client := skipchain.NewClient()
	//StoreSkipBlock already calls its verifier function
	sb, err := client.StoreSkipBlock(txnRequest.SkipBlock, nil, txnRequest.CertBlock.CurrMTR) //No idea how this should be called.
	if err != nil {
		return nil, err
	}
	return &AddNewTransactionResponse{sb.Latest}, nil
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
	if err := s.RegisterHandlers(s.CreateSkipchain, s.AddNewTransaction); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	log.ErrFatal(skipchain.RegisterVerification(c, VerifyMerkleTreeRoot, s.VerifyMerkleTreeRoot))
	return s
}
