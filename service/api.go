package certchain

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"bytes"

	coniks_crypto "github.com/coniks-sys/coniks-go/crypto"
	coniks_sign "github.com/coniks-sys/coniks-go/crypto/sign"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"github.com/coniks-sys/coniks-go/merkletree"
	"github.com/dedis/cothority/skipchain"
	"github.com/dedis/onet/crypto"
	"github.com/dedis/onet/log"
	"gopkg.in/dedis/crypto.v0/random"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
)

// Suite used in signing
var suite = network.Suite

// Digest size of hash functions
var hashSize = 32

// CONIKS setup
var vrfPrivKey, _ = vrf.GenerateKey(bytes.NewReader(
	[]byte("deterministic tests need 32 byte")))

// Client is a structure to communicate with the CoSi
// service
type Client struct {
	*onet.Client
	// merkle tree structure of coniks
	pad *merkletree.PAD
	// secretKey used to sign merkle tree roots, also used to derive public key
	secretKey coniks_sign.PrivateKey
	// number of certificates issued by this client
	certCtr uint64
	// email of the client
	email string
}

// NewClient instantiates a new cosi.Client
func NewClient(email string) *Client {
	secretKey, err := coniks_sign.GenerateKey(nil)
	if err != nil {
		return nil
	}
	pad, err := merkletree.NewPAD(PadAd{""}, secretKey, vrfPrivKey, 10)
	if err != nil {
		return nil
	}
	return &Client{onet.NewClient(Name), pad, secretKey, 0, email}
}

// GenerateNewSecretKey generetes a new secret key for the client
func (c *Client) GenerateNewSecretKey() {
	secretKey, err := coniks_sign.GenerateKey(nil)
	if err != nil {
		return
	}
	c.secretKey = secretKey
}

// GenerateCertificates generates n random certificates and returns them in a slice of slice of bytes format
func (c *Client) GenerateCertificates(n int) []crypto.HashID {
	leaves := make([]crypto.HashID, n)
	for i := range leaves {
		leaves[i] = random.Bytes(hashSize, random.Stream)
	}
	return leaves
}

// CreateCertBlock builds a new CertBlock from the supplied certificates using CONIKs Merkle Tree algorithm
func (c *Client) CreateCertBlock(certifs []crypto.HashID) *CertBlock {
	for _, cert := range certifs {
		if err := c.pad.Set(c.email, cert); err != nil {
			return nil
		}
	}
	c.pad.Update(nil)
	str := c.pad.LatestSTR()
	latestSignedMTR := str.Signature
	latestSignedMTRHash := coniks_crypto.Digest(latestSignedMTR)
	prevSignedMTRHash := str.PreviousSTRHash
	latestMTR := str.Serialize()
	publicKey, ok := c.secretKey.Public()
	if !ok {
		return nil
	}
	ap, _ := c.pad.Lookup(c.email)
	return &CertBlock{latestSignedMTR, latestSignedMTRHash, prevSignedMTRHash, latestMTR, publicKey, ap}
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

// GetConiksAuth returns an inclusion proof from the pad of a client given an email
func GetConiksAuth(email string, pad *merkletree.PAD) (*merkletree.AuthenticationPath, error) {
	ap, err := pad.Lookup(email)
	if ap.Leaf.Value == nil {
		log.Error("Cannot find key:", email)
		return nil, err
	}
	if !bytes.Equal(ap.Leaf.Index, ap.LookupIndex) {
		log.Error("Leaf index and LookupIndex values don't match!")
		return nil, err
	}
	return ap, nil
}

// VerifyConiksAuth returns the latest SkipBlock associated with the given AuthenticationPath as a proof of inclusion
func VerifyConiksAuth(ap *merkletree.AuthenticationPath, latestSB *skipchain.SkipBlock) *skipchain.SkipBlock {
	_, cb, err := network.Unmarshal(latestSB.Data)
	log.ErrFatal(err)
	if bytes.Equal(ap.VrfProof, cb.(*CertBlock).ConiksAuthPath.VrfProof) {
		return latestSB
	}
	client, ctr := skipchain.NewClient(), 0
	for {
		prevSB, cerr := client.GetSingleBlock(latestSB.Roster, latestSB.BackLinkIDs[ctr])
		if cerr != nil {
			return nil
		}
		_, cb, err := network.Unmarshal(prevSB.Data)
		log.ErrFatal(err)
		if bytes.Equal(ap.VrfProof, cb.(*CertBlock).ConiksAuthPath.VrfProof) {
			return latestSB
		}
		ctr++
	}

}
