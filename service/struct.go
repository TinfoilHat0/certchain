package template

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/crypto/abstract"
	"github.com/satori/go.uuid"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
)

// We need to register all messages so the network knows how to handle them.
func init() {
	for _, msg := range []interface{}{
		CreateSkipchainRequest{}, //What to put , what not to put
		AddNewTransactionRequest{},
		MerkleTreeRoot{},
		CertBlock{},
		Key{},
	} {
		network.RegisterMessage(msg)
	}
}

// VerifyMerkleTreeRoot is the ID of the verifier for the Certchain service
var VerifyMerkleTreeRoot = skipchain.VerifierID(uuid.NewV5(uuid.NamespaceURL, "Certchain"))

//CreateSkipchainRequest is the structure for a new skipchain request
type CreateSkipchainRequest struct {
	Roster    *onet.Roster
	PublicKey *Key
}

//CreateSkipchainResponse returns a block from the underlying Skipchain service
type CreateSkipchainResponse struct {
	SkipBlock *skipchain.SkipBlock
}

//AddNewTransactionRequest is the structure for a new transaction request
type AddNewTransactionRequest struct {
	SkipBlock *skipchain.SkipBlock
	CertBlock *CertBlock
}

//AddNewTransactionResponse is the structure for a transaction response
type AddNewTransactionResponse struct {
	SkipBlock *skipchain.SkipBlock
}

//MerkleTreeRoot is a wrapper for the signed MTR
type MerkleTreeRoot struct {
	MTRoot []byte
}

//CertBlock stores a transaction of the Certchain(this is stored as Data in Skipchain)
type CertBlock struct {
	PrevMTR   *MerkleTreeRoot
	LatestMTR *MerkleTreeRoot
	PublicKey *Key
}

//Key is a wrapper structure for the key used in Schnorr Signature
type Key struct {
	PublicKey abstract.Point
	suite     abstract.Suite
}
