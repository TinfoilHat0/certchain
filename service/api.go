package certchain

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"bytes"
	"crypto/sha256"

	"github.com/TinfoilHat0/certchain/merkle_tree"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"github.com/coniks-sys/coniks-go/merkletree"
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/onet/log"
	"gopkg.in/dedis/crypto.v0/config"
	"gopkg.in/dedis/crypto.v0/random"
	"gopkg.in/dedis/crypto.v0/sign"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
)

// Suite used in signing
var suite = network.Suite

// 32 bytes
var hashSize = sha256.New().Size()

// Client is a structure to communicate with the CoSi
// service
type Client struct {
	*onet.Client
	keyPair     *config.KeyPair
	coniksKey   string
	coniksIndex []byte
}

// NewClient instantiates a new cosi.Client
func NewClient() *Client {
	kp := config.NewKeyPair(suite)
	coniksKey := "my_id@epfl.ch"
	vrfPrivKey, err := vrf.GenerateKey(bytes.NewReader(
		[]byte("deterministic tests need 32 byte")))
	log.ErrFatal(err)
	coniksIndex := vrfPrivKey.Compute([]byte(coniksKey))
	return &Client{onet.NewClient(Name), kp, coniksKey, coniksIndex}
}

// GenerateNewKeyPair generetes a new keypair for the client
func (c *Client) GenerateNewKeyPair() {
	kp := config.NewKeyPair(suite)
	c.keyPair = kp
}

// GenerateCertificates generates n random certificates and returns them in a slice of slice of bytes format
func (c *Client) GenerateCertificates(n int) []crypto.HashID {
	leaves := make([]crypto.HashID, n)
	for i := range leaves {
		leaves[i] = random.Bytes(hashSize, random.Stream)
	}
	return leaves
}

// CreateCertBlock builds a new CertBlock from the supplied certificates
func (c *Client) CreateCertBlock(certifs []crypto.HashID, prevMTR []byte, keyPair *config.KeyPair) *CertBlock {
	certMTR, _ := crypto.ProofTree(sha256.New, certifs)
	leaves := make([]crypto.HashID, 2)
	leaves[0] = prevMTR
	leaves[1] = certMTR
	latestMTR, _ := crypto.ProofTree(sha256.New, leaves)
	latestSignedMTR, err := sign.Schnorr(suite, keyPair.Secret, latestMTR)
	if err != nil {
		return nil
	}
	return &CertBlock{latestSignedMTR, latestMTR, prevMTR, keyPair.Public}
}

// CreateCertBlockCONIKS builds a new CertBlock from the supplied certificates using CONIKs Merkle Tree algorithm
func (c *Client) CreateCertBlockCONIKS(certifs []crypto.HashID, prevMTR []byte, keyPair *config.KeyPair) *CertBlock {
	// create a new merkle tree
	m, err := merkletree.NewMerkleTree()
	if err != nil {
		return nil
	}
	// put the previous root
	if err := m.Set(c.coniksIndex, c.coniksKey, prevMTR); err != nil {
		return nil
	}
	// add new certificates
	for _, cert := range certifs {
		if err := m.Set(c.coniksIndex, c.coniksKey, cert); err != nil {
			return nil
		}
	}
	m.RecomputeHash()
	latestMTR := m.GetRootHash()
	// Or use Get() and then lookupIndex of the returned value?
	/*
		// get the authentication path from the tree
		ap := m.Get(c.coniksIndex)
		if ap.Leaf.Value == nil {
			return nil
		}
	*/
	latestSignedMTR, err := sign.Schnorr(suite, keyPair.Secret, latestMTR)
	if err != nil {
		return nil
	}
	return &CertBlock{latestSignedMTR, latestMTR, prevMTR, keyPair.Public}
}

// CreateSkipchain initializes the Skipchain which is the underlying blockchain service
func (c *Client) CreateSkipchain(r *onet.Roster, genesisCertBlock *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := r.RandomServerIdentity()
	reply := &CreateSkipchainResponse{}
	err := c.SendProtobuf(dst, &CreateSkipchainRequest{r, genesisCertBlock}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}

// AddNewTxn adds a new transaction to the underlying Skipchain service
func (c *Client) AddNewTxn(r *onet.Roster, sb *skipchain.SkipBlock, cb *CertBlock) (*skipchain.SkipBlock, onet.ClientError) {
	dst := sb.Roster.RandomServerIdentity()
	reply := &AddNewTxnResponse{}
	err := c.SendProtobuf(dst, &AddNewTxnRequest{r, sb, cb}, reply)
	if err != nil {
		return nil, err
	}
	return reply.SkipBlock, nil
}
