package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

//Create a new CertBlock
func TestGenerateCertBlock(t *testing.T) {
	client := NewClient()
	certifs := client.GenerateCertificates(5)
	assert.NotNil(t, certifs)

	prevSignedMTR := make([]byte, 32)
	cb := client.CreateCertBlock(prevSignedMTR, certifs, client.keypair)
	assert.NotNil(t, cb)
}

//Create a new CertBlock and store it in the skipchain
func TestCreateSkipChain(t *testing.T) {
	client := NewClient()
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	cb := client.CreateCertBlock(make([]byte, 32), client.GenerateCertificates(5), client.keypair)
	sb, err := client.CreateSkipchain(roster, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

}

/*
//Add a new txn to the skipchain
func TestAddNewTxn(t *testing.T) {
	client := NewClient()
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	cb := client.CreateCertBlock(make([]byte, 32), client.GenerateCertificates(5), client.keypair)
	sb, err := client.CreateSkipchain(roster, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	sb, err = client.AddNewTxn(roster, sb, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

}
*/
