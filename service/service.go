package template

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"bytes"

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
	path string
	//A map for the unspent transactions. Key is the string of prevSignedMTR and value is the hash of the skipblock
	unspentTxns map[string]skipchain.SkipBlockID
}

//CreateSkipchain creates a new skipchain
func (s *Service) CreateSkipchain(cs *CreateSkipchainRequest) (*CreateSkipchainResponse, onet.ClientError) {
	client := skipchain.NewClient()
	log.Print("From service")
	sb, err := client.CreateGenesis(cs.Roster, 1, 1, []skipchain.VerifierID{VerifyTxn}, cs.CertBlock, nil)
	if err != nil {
		return nil, err
	}
	s.unspentTxns[string(cs.CertBlock.PrevSignedMTR)] = sb.Hash //add block to the map, keyed by PrevSignedMTR Q:Should I just use the hash of certblock?
	return &CreateSkipchainResponse{sb}, nil
}

//AddNewTxn stores a new transaction in the underlying Skipchain
func (s *Service) AddNewTxn(txn *AddNewTxnRequest) (*AddNewTxnResponse, onet.ClientError) {
	client := skipchain.NewClient()                                     //Shouldn't I use a single client?
	sb, err := client.StoreSkipBlock(txn.SkipBlock, nil, txn.CertBlock) //txn.CertBlock is passed as Data right ?
	if err != nil {
		return nil, err
	}
	s.unspentTxns[string(txn.CertBlock.PrevSignedMTR)] = sb.Latest.Hash //add block to the map, keyed by PrevSignedMTR
	return &AddNewTxnResponse{sb.Latest}, nil
}

//VerifyTxn verifies a CertChain txn
//Verification is done as follows:
//1. Get the public key from the previous block
//2. Verify the signature on the blocks latestMTR, if previousMTR is all 0(this is the genesis certblock) return true
//3. Check whether the block is in unspentTxn map. If it is, remove block from the map and return true. Otherwise, returnfalse
func (s *Service) VerifyTxn(newID []byte, newSB *skipchain.SkipBlock) bool {
	//Do I need to call verifications for skipchain block as well or is it handled by itself?
	client := skipchain.NewClient()
	parentSB, getBlockErr := client.GetSingleBlock(newSB.Roster, newSB.ParentBlockID) //Does this return genesis when called by the genesis ?
	if getBlockErr != nil {
		return false
	}
	//Get the public key from the previous block
	_, cbPrev, _ := network.Unmarshal(parentSB.Data)
	publicKey := cbPrev.(*CertBlock).PublicKey //Verification has to be done using the public key of the previous block
	suite := cbPrev.(*CertBlock).Suite

	//Verify the signature
	_, cb, _ := network.Unmarshal(newSB.Data)
	signErr := sign.VerifySchnorr(suite, publicKey, cb.(*CertBlock).LatestMTR, cb.(*CertBlock).LatestSignedMTR)
	if signErr != nil {
		return false
	}
	//If block doesn't have any previous MTR, verification only consists of checking the signature
	if bytes.Equal(cb.(*CertBlock).PrevSignedMTR, make([]byte, 32)) {
		return true
	}
	//Check if the block is unspent. If it is spent, i.e. it can't be found in the map, return false
	if _, exists := s.unspentTxns[string(cb.(*CertBlock).PrevSignedMTR)]; !exists {
		return false

	}
	//Spend the txn by removing it from map
	delete(s.unspentTxns, string(cb.(*CertBlock).PrevSignedMTR))
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
	log.ErrFatal(skipchain.RegisterVerification(c, VerifyTxn, s.VerifyTxn))
	return s
}
