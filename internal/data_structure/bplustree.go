package data_structure

type Item struct {
	Score  float64
	Member string
}

func (i *Item) CompateTo(other *Item) int {
	if i.Score > other.Score {
		return 1
	}

	if i.Score < other.Score {
		return -1
	}

	// In case that, the score of both item is equal then we will compare their string based on lexicographical order
	if i.Member < other.Member {
		return -1
	}

	if i.Member > other.Member {
		return 1
	}
	return 0
}

type Node struct {
	Items    *[]Item
	Children *[]Node // Pointer(s) to child nodes
	IsLeaf   bool    // True if it's a leaf node
	Parent   *Node   // Pointer back to its parent
	Leaf     *Node   // For leaf nodes, a pointer to the next leaf in the sequence
}

type BPlusTree struct {
	Root   *Node
	Degree int // The maximum number of children that a node can have
}

func NewBPlusTree(degree int) *BPlusTree {
	return &BPlusTree{
		Root:   &Node{IsLeaf: true},
		Degree: degree,
	}
}

func (t *BPlusTree) Score(member string) (float64, bool) {
	
}
