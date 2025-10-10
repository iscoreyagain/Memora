package data_structure

type ListNode struct {
	lp   *ListPack
	prev *ListNode
	next *ListNode
}

type QuickList struct {
	head  *ListNode
	tail  *ListNode
	count int // the number of elems in ALL OF THE NODES belongs to that quicklist - Each node points to a listpack
}

type List struct {
	key   string
	qlist *QuickList
}

func NewNode(prev, next *ListNode, lp *ListPack) *ListNode {
	return &ListNode{
		lp:   lp,
		prev: prev,
		next: next,
	}
}

func NewQuickList(head, tail *ListNode, num int) *QuickList {
	return &QuickList{
		head:  head,
		tail:  tail,
		count: num,
	}
}

func NewList(key string) *List {
	lp := NewListPack(0)
	node := NewNode(nil, nil, lp)
	ql := NewQuickList(node, node, 0)
	list := &List{
		key:   key,
		qlist: ql,
	}

	return list
}

func (l *List) RPush(member ...interface{}) int {
	if l.qlist.tail == nil {
		return 0
	}

	added := l.qlist.tail.lp.PushRight(member...)
	l.qlist.count += added
	return added
}

func (l *List) LPush(member ...interface{}) int {
	if l.qlist.head == nil {
		return 0
	}

	added := l.qlist.head.lp.PushLeft(member...)
	l.qlist.count += added
	return added
}

func (l *List) Len() int {
	return l.qlist.count
}
