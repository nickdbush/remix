package remix

func (n *Node) findNearestLoad() *loadOp {
	current := n
	for {
		if op, ok := current.Operation.(loadOp); ok {
			return &op
		}

		switch len(current.parents()) {
		case 0:
			return nil
		case 1:
			current = current.parents()[0]
		default:
			panic("multiple parents")
		}
	}
}
