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
	_, roster, _ := local.GenTree(3, true)
	defer local.CloseAll()

	// Send a request to the service
	client := NewClient()
	log.Lvl1("Sending request to service...")

	sb, err := client.CreateSkipchain(roster)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

	cb, err := client.CreateNewCertBlock(nil, nil)
	log.ErrFatal(err, "Can't create new CertBlock")
	log.Print(cb)

	sb, err = client.AddNewTransaction(sb, cb)
	log.ErrFatal(err, "Couldn't send")
	assert.NotNil(t, sb)

}
