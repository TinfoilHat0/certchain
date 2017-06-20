package certchain

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"github.com/dedis/cothority/messaging"
	"github.com/dedis/cothority/skipchain"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
)

// Name is the name to refer to the CertChain service from another package.
const Name = "CertChain"

func init() {
	onet.RegisterNewService(Name, newService)
}

// Service is our CertChain-service
type Service struct {
	*onet.ServiceProcessor
	propagate messaging.PropagationFunc
	// A map for the unspent transactions. Key is the string of latestMTR and value is the hash of the skipblock
	unspentTxnMap map[string]skipchain.SkipBlockID
}

// CreateSkipchain creates a new skipchain
func (s *Service) CreateSkipchain(cs *CreateSkipchainRequest) (*CreateSkipchainResponse, onet.ClientError) {
	client := skipchain.NewClient()
	sb, err := client.CreateGenesis(cs.Roster, 1, 1, []skipchain.VerifierID{VerifyTxn}, cs.CertBlock, nil)
	if err != nil {
		return nil, err
	}
	perr := s.startPropagation(cs.Roster, cs.CertBlock.LatestSignedMTRHash, sb.Hash)
	log.ErrFatal(perr)
	return &CreateSkipchainResponse{sb}, nil
}

// AddNewTxn stores a new transaction in the underlying Skipchain service
func (s *Service) AddNewTxn(txn *AddNewTxnRequest) (*AddNewTxnResponse, onet.ClientError) { //Where do I use roster here ?
	client := skipchain.NewClient()
	sb, err := client.StoreSkipBlock(txn.SkipBlock, nil, txn.CertBlock)
	if err != nil {
		return nil, err
	}
	perr := s.startPropagation(txn.Roster, txn.CertBlock.LatestSignedMTRHash, sb.Latest.Hash)
	log.ErrFatal(perr)
	return &AddNewTxnResponse{sb.Latest}, nil
}

// VerifyTxn verifies a txn as follows:
// 1. Get the public key from block
// 2. Verify the signature on the blocks latestMTRW
// 3. Check whether the block is in unspentTxnMap. If it is, remove block from the map and return true. Otherwise, return false
func (s *Service) VerifyTxn(newID []byte, newSB *skipchain.SkipBlock) bool {
	//  Get the public key from the block and verify the signature
	_, cb, _ := network.Unmarshal(newSB.Data)
	publicKey := cb.(*CertBlock).PublicKey
	if !publicKey.Verify(cb.(*CertBlock).LatestMTR, cb.(*CertBlock).LatestSignedMTR) {
		return false
	}
	// Check if the block references to an unspent txn
	if _, exists := s.unspentTxnMap[string(cb.(*CertBlock).PrevSignedMTRHash)]; !exists {
		return false
	}
	// Spend the referred txn by removing it from the map
	delete(s.unspentTxnMap, string(cb.(*CertBlock).PrevSignedMTRHash))
	return true
}

// StartPropagation is a convenience function to call propagate so that we don't duplicate code
func (s *Service) startPropagation(roster *onet.Roster, blockMTR []byte, blockHash skipchain.SkipBlockID) error {
	log.Lvl3("Starting to propagate for service", s.ServerIdentity())
	replies, err := s.propagate(roster, &PropagateTxnInfo{blockMTR, blockHash}, propagateTimeout)
	if err != nil {
		return err
	}
	if replies != len(roster.List) {
		log.Warn("Did only get", replies, "out of", len(roster.List))
	}
	return nil
}

// PropagateTxnMap is the propagation function of this service which is
// used to share data across nodes of a roster
func (s *Service) propagateTxnMap(msg network.Message) {
	txnInfo, ok := msg.(*PropagateTxnInfo)
	if !ok {
		log.Error("Couldn't convert to PropagateTxnInfo")
		return
	}
	s.unspentTxnMap[string(txnInfo.BlockMTR)] = txnInfo.BlockHash
}

// newService receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *onet.Context) onet.Service {
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
		unspentTxnMap:    make(map[string]skipchain.SkipBlockID),
	}
	if err := s.RegisterHandlers(s.CreateSkipchain, s.AddNewTxn); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	var err error
	s.propagate, err = messaging.NewPropagationFunc(c, "TxnMapPropagate", s.propagateTxnMap)
	log.ErrFatal(err)
	log.ErrFatal(skipchain.RegisterVerification(c, VerifyTxn, s.VerifyTxn))
	return s
}
