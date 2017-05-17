package certchain

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/satori/go.uuid"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
)

// We need to register all messages so the network knows how to handle them. //move them to messsages?
func init() {
	for _, msg := range []interface{}{
		&CreateSkipchainRequest{}, //What to put , what not to put ? (Things that will be marshalled ?)
		&CreateSkipchainResponse{},
		&AddNewTxnRequest{},
		&AddNewTxnResponse{},
		&PropagateTxnInfo{},
		&CertBlock{},
		&Service{},
	} {
		network.RegisterMessage(msg)
	}
}

// VerifyTxn is the ID of the verifier for the Certchain service
var VerifyTxn = skipchain.VerifierID(uuid.NewV5(uuid.NamespaceURL, "Certchain"))

// How many msec to wait before a timeout is generated in the propagation.
const propagateTimeout = 10000

// CreateSkipchainRequest is the structure for a new skipchain addition request
type CreateSkipchainRequest struct {
	Roster    *onet.Roster
	CertBlock *CertBlock
}

// CreateSkipchainResponse is the structure for a skipchain addition response
type CreateSkipchainResponse struct {
	SkipBlock *skipchain.SkipBlock
}

// AddNewTxnRequest is the structure for a txn addition request
type AddNewTxnRequest struct {
	Roster    *onet.Roster
	SkipBlock *skipchain.SkipBlock
	CertBlock *CertBlock
}

// AddNewTxnResponse is the structure for a txn addition request response
type AddNewTxnResponse struct {
	SkipBlock *skipchain.SkipBlock
}

// PropagateTxnInfo is a wrapper to propagate a txn info across nodes
type PropagateTxnInfo struct {
	BlockMTR  []byte
	BlockHash skipchain.SkipBlockID
}

// CertBlock stores a transaction of the Certchain (this is stored in data field of a Skipblock)
type CertBlock struct {
	LatestSignedMTR []byte
	LatestMTR       []byte
	PrevMTR         []byte
	PublicKey       abstract.Point
}
