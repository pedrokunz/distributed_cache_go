package test

import (
	"github.com/pedrokunz/distributed_cache_go/internal/data_node"
	"github.com/pedrokunz/distributed_cache_go/pkg/consistent_hash"
	"testing"
)

func TestHashRing_AddAndGetNode(t *testing.T) {
	hashRing := consistent_hash.New(3)
	node := data_node.New()

	hashRing.AddNode(node)

	retrievedNode := hashRing.GetNode("some_key")
	if retrievedNode == nil {
		t.Fatal("expected a node, got nil")
	}

	if retrievedNode.ID() != node.ID() {
		t.Errorf("expected node ID %s, got %s", node.ID(), retrievedNode.ID())
	}
}

func TestHashRing_RemoveNode(t *testing.T) {
	hashRing := consistent_hash.New(3)
	node := data_node.New()

	hashRing.AddNode(node)
	hashRing.RemoveNode(node)

	retrievedNode := hashRing.GetNode("some_key")
	if retrievedNode != nil {
		t.Errorf("expected no node, got node ID %s", retrievedNode.ID())
	}
}

func TestHashRing_ConsistentHashing(t *testing.T) {
	hashRing := consistent_hash.New(3)
	node1 := data_node.New()
	node2 := data_node.New()

	hashRing.AddNode(node1)
	hashRing.AddNode(node2)

	key := "some_key"
	node := hashRing.GetNode(key)
	if node == nil {
		t.Fatal("expected a node, got nil")
	}

	// Remove the node and check if the key is reassigned
	hashRing.RemoveNode(node)
	newNode := hashRing.GetNode(key)
	if newNode == nil {
		t.Fatal("expected a node, got nil")
	}

	if newNode.ID() == node.ID() {
		t.Errorf("expected a different node, got the same node ID %s", newNode.ID())
	}
}
