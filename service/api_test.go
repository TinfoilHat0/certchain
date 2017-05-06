package template

import (
	"crypto/sha256"
	"testing"

	"github.com/dedis/onet/crypto"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func TestServiceTemplate(t *testing.T) {
	local := onet.NewTCPTest()
	// generate 3 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	//Create a new skipchain
	client := NewClient()
	sb, err := client.CreateSkipchain(roster)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	//Generate CertBlock from bunch of certificates
	newHash := sha256.New
	hash := newHash()
	n := 5
	leaves := make([]crypto.HashID, n)
	for i := range leaves {
		leaves[i] = make([]byte, hash.Size())
		for j := range leaves[i] {
			leaves[i][j] = byte(i)
		}
	}
	cb := client.CreateNewCertBlock(nil, leaves)
	assert.NotNil(t, cb)

	sb, err = client.AddNewTxn(sb, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

}
