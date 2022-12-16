package harness

import (
	"sync"

	. "github.com/ipfs/kubo/test/cli/testutils"
	"github.com/multiformats/go-multiaddr"
)

// Nodes is a collection of Kubo nodes along with operations on groups of nodes.
type Nodes []*Node

func (n Nodes) Init(args ...string) Nodes {
	ForEachPar(n, func(node *Node) { node.Init(args...) })
	return n
}

func (n Nodes) Connect() Nodes {
	wg := sync.WaitGroup{}
	for i, node := range n {
		for j, otherNode := range n {
			if i == j {
				continue
			}
			node := node
			otherNode := otherNode
			wg.Add(1)
			go func() {
				defer wg.Done()
				node.Connect(otherNode)
			}()
		}
	}
	wg.Wait()
	for _, node := range n {
		firstPeer := node.Peers()[0]
		if _, err := firstPeer.ValueForProtocol(multiaddr.P_P2P); err != nil {
			log.Panicf("unexpected state for node %d with peer ID %s: %s", node.ID, node.PeerID(), err)
		}
	}
	return n
}

func (n Nodes) StartDaemons() Nodes {
	ForEachPar(n, func(node *Node) { node.StartDaemon() })
	return n
}

func (n Nodes) StopDaemons() Nodes {
	ForEachPar(n, func(node *Node) { node.StopDaemon() })
	return n
}
