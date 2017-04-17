package template

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/satori/go.uuid"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
)

// We need to register all messages so the network knows how to handle them.
func init() {
	for _, msg := range []interface{}{
		CreateSkipchainRequest{},
		MerkleTreeRoot{},
		AddMerkleTreeRootRequest{},
	} {
		network.RegisterMessage(msg)
	}
}

// VerifyCert id for the verification function
var VerifyCert = skipchain.VerifierID(uuid.NewV5(uuid.NamespaceURL, "Certchain"))

// CreateSkipchain will run the tepmlate-protocol on the roster and return
// the time spent doing so.
type CreateSkipchainRequest struct {
	Roster *onet.Roster
}

type CreateSkipchainResponse struct {
	SkipBlock *skipchain.SkipBlock
}

type MerkleTreeRoot struct {
}

type AddMerkleTreeRootRequest struct {
	SkipBlock *skipchain.SkipBlock
	TreeRoot  *MerkleTreeRoot
}

type AddMerkleTreeRootResponse struct {
	SkipBlock *skipchain.SkipBlock
}
