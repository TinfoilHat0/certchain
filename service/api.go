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
	log.LLvl4("Sending message to", dst)
	reply := &CreateSkipchainResponse{}
	err := c.SendProtobuf(dst, &CreateSkipchainRequest{r}, reply) //websocket
	if err != nil {
		return nil, err
	}
	return reply.SkıpBlock, nil
}

// AddMerkleTreeRoot blabla
func (c *Client) AddMerkleTreeRoot(sb *skipchain.SkipBlock, mtr *MerkleTreeRoot) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	log.Lvl4("Sending message to", dst)
	reply := &AddMerkleTreeRootResponse{}
	err := c.SendProtobuf(dst, &AddMerkleTreeRootRequest{sb, mtr}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
