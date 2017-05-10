package template

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

//Create a new CertBlock from 5 random certificates
func TestGenerateCertBlock(t *testing.T) {
	client := NewClient()
	certifs := client.GenerateCertificates(5)
	assert.NotNil(t, certifs)

	prevMTR := make([]byte, hashSize)
	cb := client.CreateCertBlock(certifs, prevMTR, client.keyPair)
	assert.NotNil(t, cb)
	assert.Equal(t, cb.PrevMTR, prevMTR)
	assert.Equal(t, client.keyPair.Public, cb.PublicKey)
}

//Initialize a new SkipChain and store a CertBlock in it
func TestCreateSkipChain(t *testing.T) {
	client := NewClient()
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	cb := client.CreateCertBlock(client.GenerateCertificates(5), make([]byte, hashSize), client.keyPair)
	sb, err := client.CreateSkipchain(roster, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	_, sbRawData, merr := network.Unmarshal(sb.Data)
	log.ErrFatal(merr)
	assert.NotNil(t, sbRawData)
	assert.Equal(t, cb.LatestMTR, sbRawData.(*CertBlock).LatestMTR)
	assert.Equal(t, cb.LatestSignedMTR, sbRawData.(*CertBlock).LatestSignedMTR)
	assert.Equal(t, cb.PrevMTR, sbRawData.(*CertBlock).PrevMTR)
	assert.True(t, cb.PublicKey.Equal(sbRawData.(*CertBlock).PublicKey))

}

//Add a new txn to the skipchain
func TestAddNewTxn(t *testing.T) {
	client := NewClient()
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	cb := client.CreateCertBlock(client.GenerateCertificates(5), make([]byte, hashSize), client.keyPair)
	sb, err := client.CreateSkipchain(roster, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	cbNew := client.CreateCertBlock(client.GenerateCertificates(5), cb.LatestMTR, client.keyPair)
	assert.NotNil(t, cbNew)
	sb, err = client.AddNewTxn(roster, sb, cbNew)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

}
