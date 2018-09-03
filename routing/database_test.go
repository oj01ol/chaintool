package routing

import (
	"fmt"
	"io/ioutil"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tmpDBName = "nodeDB_test"
)

func Test_newNodeDB(t *testing.T) {
	path, err := ioutil.TempDir("", tmpDBName)
	assert.Nil(t, err, "make tempdir err")
	self := Hash{}
	db, err := newNodeDB(path, self)
	assert.Nil(t, err, "new node db err")
	db.close()

}

func Test_Node(t *testing.T) {
	path, err := ioutil.TempDir("", tmpDBName)
	assert.Nil(t, err, "make tempdir err")
	self := Hash{}
	db, err := newNodeDB(path, self)
	assert.Nil(t, err, "new node db err")
	defer db.close()
	//ti := time.Now()
	node := &Node{Addr: "na", ID: Hash{45}, Time: time.Now().Unix()}
	err = db.updateNode(node)
	require.Nil(t, err, "update node err")
	key := makeKey(node.ID, nodeDBDiscoverRoot)
	nget := db.getNode(node.ID)

	if !node.Equal(nget) {
		t.Error(fmt.Sprintf("node get from db wrong,key:%v\nori:%v\nget:%v", key, node, nget))
	}

	err = db.deleteNode(node.ID)
	require.Nil(t, err, "delete node err")

	nget = db.getNode(node.ID)
	//db.ensureExpirer()
	if nget != nil {
		t.Error("delete node fail,get:%v", nget)
	}

}

func Test_querySeeds(t *testing.T) {
	path, err := ioutil.TempDir("", tmpDBName)
	require.Nil(t, err, "make tempdir err")
	self := Hash{7, 63, 74}
	db, err := newNodeDB(path, self)
	require.Nil(t, err, "new node db err")
	defer db.close()
	var node *Node

	have := make(map[Hash]struct{})
	want := make(map[Hash]struct{})
	node = &Node{Addr: "na", ID: Hash{7, 63, 74}, Time: time.Now().Unix()}
	err = db.updateNode(node)
	if err != nil {
		t.Error(err)
	}
	//want[node.ID] = struct{}{}
	node = &Node{Addr: "nb", ID: Hash{24, 26, 84}, Time: time.Now().Unix()}
	err = db.updateNode(node)
	if err != nil {
		t.Error(err)
	}
	want[node.ID] = struct{}{}
	node = &Node{Addr: "nc", ID: Hash{60, 14, 24}, Time: time.Now().Unix()}
	err = db.updateNode(node)
	if err != nil {
		t.Error(err)
	}
	want[node.ID] = struct{}{}
	err = db.updateLastPongReceived(Hash{7, 63, 74}, time.Now())
	err = db.updateLastPongReceived(Hash{24, 26, 84}, time.Now())
	err = db.updateLastPongReceived(Hash{60, 14, 24}, time.Now())
	nodes := db.querySeeds(5, time.Hour*12)
	for _, node = range nodes {
		have[node.ID] = struct{}{}
	}
	if len(have) != len(want) {
		t.Error("quert count mistake", "have:", len(have), "want:", len(want))
	}

	for id := range have {
		if _, ok := want[id]; !ok {
			t.Error("extra missed : ", id)
		}
	}

}
