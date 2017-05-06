package template

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"crypto/sha256"

	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/crypto/config"
	"github.com/dedis/crypto/ed25519"
	"github.com/dedis/crypto/sign"
	"github.com/dedis/onet/crypto"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

// Client is a structure to communicate with the CoSi
// service
type Client struct {
	*onet.Client
	keypair *config.KeyPair
}

// NewClient instantiates a new cosi.Client
func NewClient() *Client {
	suite := ed25519.NewAES128SHA256Ed25519(false)
	kp := config.NewKeyPair(suite)
	return &Client{onet.NewClient(Name), kp} //by reference or by value?
}

// CreateSkipchain initializes the skipchain which is the underlying blockchain
func (c *Client) CreateSkipchain(r *onet.Roster) (*skipchain.SkipBlock, onet.ClientError) {
	dst := r.RandomServerIdentity()
	log.Lvl4("Sending message to", dst)
	reply := &CreateSkipchainResponse{}
	key := &Key{c.keypair.Public, c.keypair.Suite}
	err := c.SendProtobuf(dst, &CreateSkipchainRequest{r, key}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}

//CreateNewCertBlock generates a signed MTR from certifs and wraps the parameters in a CertBlock
func (c *Client) CreateNewCertBlock(prevMTR *MerkleTreeRoot, certifs []crypto.HashID) *CertBlock {
	root, _ := crypto.ProofTree(sha256.New, certifs) //Why can't I just use a [][]byte ?
	signedMTR, err := sign.Schnorr(c.keypair.Suite, c.keypair.Secret, root)
	if err != nil {
		return nil
	}
	return &CertBlock{prevMTR, &MerkleTreeRoot{signedMTR}, &Key{c.keypair.Public, c.keypair.Suite}}
}

//AddNewTxn adds a new transaction to the underlying Skipchain service
func (c *Client) AddNewTxn(sb *skipchain.SkipBlock, cb *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	reply := &AddNewTxnResponse{}
	err := c.SendProtobuf(dst, &AddNewTxnRequest{sb, cb}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
