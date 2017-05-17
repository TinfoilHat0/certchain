package certchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

// Generate a CertBlock
func TestGenerateCertBlock(t *testing.T) {
	client := NewClient()
	certifs := client.GenerateCertificates(5)
	assert.NotNil(t, certifs)

	cb := client.CreateCertBlock(certifs)
	assert.NotNil(t, cb)
}

// Initialize a new SkipChain and store a CertBlock in it
func TestCreateSkipChain(t *testing.T) {
	client := NewClient()
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	cb := client.CreateCertBlock(client.GenerateCertificates(1))
	sb, err := client.CreateSkipchain(roster, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	_, sbRawData, merr := network.Unmarshal(sb.Data)
	log.ErrFatal(merr)
	assert.NotNil(t, sbRawData)

	assert.Equal(t, cb.LatestMTR, sbRawData.(*CertBlock).LatestMTR)
	assert.Equal(t, cb.LatestSignedMTR, sbRawData.(*CertBlock).LatestSignedMTR)
	assert.Equal(t, cb.PrevSignedMTRHash, sbRawData.(*CertBlock).PrevSignedMTRHash)
	assert.Equal(t, cb.PublicKey, sbRawData.(*CertBlock).PublicKey)

}

// Add a new txn to the SkipChain by running the verification function using the CONIKS' Merkle Tree Algorithm
func TestAddNewTxn(t *testing.T) {
	client := NewClient()
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	cb := client.CreateCertBlock(client.GenerateCertificates(1))
	sb, err := client.CreateSkipchain(roster, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	// Adding more than one certificate at a time causes error
	cb = client.CreateCertBlock(client.GenerateCertificates(1))
	assert.NotNil(t, cb)
	sb, err = client.AddNewTxn(roster, sb, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	_, sbRawData, merr := network.Unmarshal(sb.Data)
	log.ErrFatal(merr)

	assert.NotNil(t, sbRawData)

	assert.Equal(t, cb.LatestMTR, sbRawData.(*CertBlock).LatestMTR)
	assert.Equal(t, cb.LatestSignedMTR, sbRawData.(*CertBlock).LatestSignedMTR)
	assert.Equal(t, cb.PrevSignedMTRHash, sbRawData.(*CertBlock).PrevSignedMTRHash)
	assert.Equal(t, cb.PublicKey, sbRawData.(*CertBlock).PublicKey)

}
