package template

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"bytes"

	"github.com/dedis/cothority/messaging"
	"github.com/dedis/cothority/skipchain"
	"gopkg.in/dedis/crypto.v0/sign"
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
	perr := s.startPropagation(cs.Roster, cs.CertBlock.LatestMTR, sb.Hash)
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
	perr := s.startPropagation(txn.Roster, txn.CertBlock.LatestMTR, sb.Latest.Hash)
	log.ErrFatal(perr)
	return &AddNewTxnResponse{sb.Latest}, nil
}

// VerifyTxn verifies a Certchain txn as follows:
// 1. Get the public key from the previous block
// 2. Verify the signature on the blocks latestMTR, if previousMTR is all 0 (this is the genesis certblock) this is the only verification
// 3. Check whether the block is in unspentTxn map. If it is, remove block from the map and return true. Otherwise, return false
func (s *Service) VerifyTxn(newID []byte, newSB *skipchain.SkipBlock) bool {
	client := skipchain.NewClient()
	previousSB, cerr := client.GetSingleBlock(newSB.Roster, newSB.BackLinkIDs[0])
	if cerr != nil {
		return false
	}
	// Get the public key from the previous block as verification has to be done using that key
	_, cbPrev, err := network.Unmarshal(previousSB.Data)
	log.ErrFatal(err)
	publicKey := cbPrev.(*CertBlock).PublicKey
	// Verify the signature
	_, cb, _ := network.Unmarshal(newSB.Data)
	signErr := sign.VerifySchnorr(suite, publicKey, cb.(*CertBlock).LatestMTR, cb.(*CertBlock).LatestSignedMTR)
	if signErr != nil {
		return false
	}
	// If block is the genesis block, verification only consists of checking the signature
	if bytes.Equal(cb.(*CertBlock).PrevMTR, make([]byte, 32)) {
		return true
	}
	// Check if the block is unspent. If it is spent, i.e. it can't be found in the map, return false
	if _, exists := s.unspentTxnMap[string(cb.(*CertBlock).PrevMTR)]; !exists {
		return false
	}
	// Spend the txn by removing it from the map
	delete(s.unspentTxnMap, string(cb.(*CertBlock).PrevMTR))
	return true
}

// notify other services about new/updated unspentTxnMap
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

// PropagateTxnMap propagates TxnMap across nodes of the roster
func (s *Service) propagateTxnMap(msg network.Message) {
	txnInfo, ok := msg.(*PropagateTxnInfo)
	if !ok {
		log.Error("Couldn't convert to PropagateTxnInfo")
		return
	}
	s.unspentTxnMap[string(txnInfo.BlockMTR)] = txnInfo.BlockMTR
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
