package queue

type Node struct {
	Priority int64
	TimeSlot int64
	Data     []byte
}

type Queue []*Node

func (q *Queue) Push(n *Node) {
	*q = append(*q, n)
}

func (q *Queue) IsEmpty() bool {
	return q.Len() == 0
}

func (q *Queue) Pop() (n *Node) {
	if !q.IsEmpty() {
		n = (*q)[0]
		*q = (*q)[1:]
		return n
	}
	return nil
}

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Last() (n *Node) {
	if !q.IsEmpty() {
		return (*q)[len(*q)-1]
	}
	return nil

}
