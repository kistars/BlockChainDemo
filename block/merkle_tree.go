package block

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil { // leaf node
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	// number of leaves must be even
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1]) // duplicate the last one
	}

	// generate the leaf nodes
	var nodes []MerkleNode
	for _, d := range data {
		n := NewMerkleNode(nil, nil, d)
		nodes = append(nodes, *n)
	}

	// create the tree
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			n := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *n)
		}

		nodes = newLevel
	}

	mTree := &MerkleTree{RootNode: &nodes[0]}
	return mTree
}
