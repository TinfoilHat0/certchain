package template

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"crypto/sha256"
	"log"

	"github.com/TinfoilHat0/certchain/merkle_tree"
	"github.com/dedis/cothority/skipchain"
	"gopkg.in/dedis/crypto.v0/config"
	"gopkg.in/dedis/crypto.v0/ed25519"
	"gopkg.in/dedis/crypto.v0/random"
	"gopkg.in/dedis/crypto.v0/sign"
	"gopkg.in/dedis/onet.v1"
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

//GenerateKeyPair generetes a new keypair for the client
func (c *Client) GenerateKeyPair() {
	suite := ed25519.NewAES128SHA256Ed25519(false)
	kp := config.NewKeyPair(suite)
	c.keypair = kp
}

//GenerateCertificates generates n random certificates and returns them in a slice of slice of bytes format
func (c *Client) GenerateCertificates(n int) []crypto.HashID {
	newHash := sha256.New
	hash := newHash()
	leaves := make([]crypto.HashID, n)
	for i := range leaves {
		leaves[i] = random.Bytes(hash.Size(), random.Stream)
	}
	return leaves
}

//CreateCertBlock builds a new CertBlock from the supplied certificates
func (c *Client) CreateCertBlock(prevSignedMTR []byte, certifs []crypto.HashID, keyPair *config.KeyPair) *CertBlock {
	//Should signing done with the previous' blocks secret key?
	certRoot, _ := crypto.ProofTree(sha256.New, certifs) //Create a MTR from the supplied certificates
	leaves := make([]crypto.HashID, 2)
	leaves[0] = certRoot
	leaves[1] = prevSignedMTR
	certRoot, _ = crypto.ProofTree(sha256.New, leaves)
	latestSignedMTR, err := sign.Schnorr(keyPair.Suite, keyPair.Secret, certRoot)
	if err != nil {
		return nil
	}
	return &CertBlock{prevSignedMTR, latestSignedMTR, certRoot, keyPair.Public, keyPair.Suite}
}

// CreateSkipchain initializes the Skipchain which is the underlying blockchain of the service
func (c *Client) CreateSkipchain(r *onet.Roster, genesisCertBlock *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := r.RandomServerIdentity()
	reply := &CreateSkipchainResponse{}
	log.Print("From api")
	err := c.SendProtobuf(dst, &CreateSkipchainRequest{r, genesisCertBlock}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}

//AddNewTxn adds a new transaction to the underlying Skipchain service
func (c *Client) AddNewTxn(r *onet.Roster, sb *skipchain.SkipBlock, cb *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	reply := &AddNewTxnResponse{}
	err := c.SendProtobuf(dst, &AddNewTxnRequest{r, sb, cb}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
