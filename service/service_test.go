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

func TestServiceTemplate(t *testing.T) {
	local := onet.NewTCPTest()
	// generate 5 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, roster, _ := local.GenTree(5, true)
	defer local.CloseAll()

	// Send a request to the service
	client := NewClient()
	log.Lvl1("Sending request to service...")
	sb, err := client.CreateSkipchain(roster)
	log.ErrFatal(err, "Couldn't send")
	log.Print(sb)
	assert.NotNil(t, sb)
	sb, err = client.AddMerkleTreeRoot(sb, &MerkleTreeRoot{})
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)
}
