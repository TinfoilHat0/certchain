package template

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

// Client is a structure to communicate with the CoSi
// service
type Client struct {
	*onet.Client
}

// NewClient instantiates a new cosi.Client
func NewClient() *Client {
	return &Client{Client: onet.NewClient(Name)}
}

// CreateSkipchain .. can also return hash(kp of our blockhain)
// client connects, issues this, gets the hash of the skipchain to find the skipblock
func (c *Client) CreateSkipchain(r *onet.Roster) (*skipchain.SkipBlock, onet.ClientError) {
	dst := r.RandomServerIdentity()
	log.Lvl4("Sending message to", dst)
	reply := &CreateSkipchainResponse{}
	err := c.SendProtobuf(dst, &CreateSkipchainRequest{r}, reply) //websocket
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}

// AddMerkleTreeRoot builds a merkle tree of certificates, signs it and sends the root to the certchain for verification
func (c *Client) AddMerkleTreeRoot(sb *skipchain.SkipBlock, mtr *MerkleTreeRoot) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	//Build the merkle tree root here, sign in with K_s and creates a transaction
	//A transaction must contain the following: Previous signed tree root, current signed tree root, public key to verify the sign
	//We must have a signing algorithm in cothority, how to use it properly ? Also, where's the secret/public key of client
	log.Lvl4("Sending message to", dst)
	//Reply should come from the Skipchain, either true or false
	reply := &AddMerkleTreeRootResponse{}
	err := c.SendProtobuf(dst, &AddMerkleTreeRootRequest{sb, mtr}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
