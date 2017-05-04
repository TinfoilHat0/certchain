package template

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/crypto/random"
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

//CreateNewCertBlock builds a new CertBlock from the supplied parameters
func (c *Client) CreateNewCertBlock(prevMTR *MerkleTreeRoot, newCerts []byte) (*CertBlock, onet.ClientError) {
	//TODO: Build a new MT using newCerts, merge with prevMTR and sign the new root with the public key of the client.
	newSignedRoot := &MerkleTreeRoot{random.Bytes(4, random.Stream)}
	return &CertBlock{prevMTR, newSignedRoot, nil, false}, nil
}

//AddNewTransaction adds a new transaction to the underlying Skipchain service
func (c *Client) AddNewTransaction(sb *skipchain.SkipBlock, cb *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	//Reply should come from the Skipchain, either true or false
	reply := &AddNewTransactionResponse{}
	err := c.SendProtobuf(dst, &AddNewTransactionRequest{sb, cb}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
