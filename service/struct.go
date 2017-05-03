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
		SkipBlockData{},
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

//MerkleTreeRoot ..
type MerkleTreeRoot struct {
}

//AddMerkleTreeRootRequest ..
type AddMerkleTreeRootRequest struct {
	SkipBlock *skipchain.SkipBlock
	TreeRoot  *MerkleTreeRoot
	//Should previous Merkle Tree root be added here along with secret/public keys of the client ?
}

//AddMerkleTreeRootResponse ..
type AddMerkleTreeRootResponse struct {
	SkipBlock *skipchain.SkipBlock
}

//SkipBlockData ..
type SkipBlockData struct {
	//What I want in my transaction goes here
}
