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
		CreateSkipchainRequest{},
		AddNewTransactionRequest{},
		MerkleTreeRoot{},
		CertBlock{},
	} {
		network.RegisterMessage(msg)
	}
}

// VerifyMerkleTreeRoot is the ID of the verifier for the Certchain service
var VerifyMerkleTreeRoot = skipchain.VerifierID(uuid.NewV5(uuid.NamespaceURL, "Certchain"))

//CreateSkipchainRequest ...
type CreateSkipchainRequest struct {
	Roster *onet.Roster
}

//CreateSkipchainResponse returns a block from the underlying Skipchain service
type CreateSkipchainResponse struct {
	SkipBlock *skipchain.SkipBlock
}

//AddNewTransactionRequest is the structure for a new transaction request
type AddNewTransactionRequest struct {
	SkipBlock *skipchain.SkipBlock //what would that contain?
	CertBlock *CertBlock
}

//AddNewTransactionResponse is the structure for a transaction response
type AddNewTransactionResponse struct {
	SkipBlock *skipchain.SkipBlock
}

//MerkleTreeRoot is a wrapper for the signed MTR
type MerkleTreeRoot struct {
	SignedRoot []byte
}

//CertBlock stores a transaction of the Certchain
type CertBlock struct {
	PrevMTR   *MerkleTreeRoot
	CurrMTR   *MerkleTreeRoot
	PublicKey []abstract.Point //is this the correct type for a key?
	spent     bool             //Indicates whether PrevMTR has been spent or not
}
