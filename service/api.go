package template

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/crypto/config"
	"github.com/dedis/crypto/ed25519"
	"github.com/dedis/crypto/random"
	"github.com/dedis/crypto/sign"
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

// CreateSkipchain .. can also return hash(kp of our blockhain)
// client connects, issues this, gets the hash of the skipchain to find the skipblock
func (c *Client) CreateSkipchain(r *onet.Roster) (*skipchain.SkipBlock, onet.ClientError) {
	dst := r.RandomServerIdentity()
	log.Lvl4("Sending message to", dst)
	reply := &CreateSkipchainResponse{}
	key := &Key{c.keypair.Public, c.keypair.Suite} //putting c.keypair.Public gives an error (no constructor for abstract.point)
	log.Print(key)
	err := c.SendProtobuf(dst, &CreateSkipchainRequest{r, key}, reply) //putting the key here gives error
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}

//CreateNewCertBlock builds a new CertBlock from the supplied parameters
func (c *Client) CreateNewCertBlock(prevMTR *MerkleTreeRoot, newCerts []byte) *CertBlock {
	//TODO: Build a new MT using newCerts, merge with prevMTR and sign the new root with the public key of the client.
	newMTR := &MerkleTreeRoot{random.Bytes(4, random.Stream)}
	signedMTR, err := sign.Schnorr(c.keypair.Suite, c.keypair.Secret, newMTR.MTRoot)
	if err != nil {
		return nil
	}
	return &CertBlock{prevMTR, &MerkleTreeRoot{signedMTR}, &Key{c.keypair.Public, c.keypair.Suite}} //better way of passing the public key?
}

//AddNewTransaction adds a new transaction to the underlying Skipchain service
func (c *Client) AddNewTransaction(sb *skipchain.SkipBlock, cb *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	reply := &AddNewTransactionResponse{}
	err := c.SendProtobuf(dst, &AddNewTransactionRequest{sb, cb}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
