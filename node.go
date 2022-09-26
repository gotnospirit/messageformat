package messageformat

type Expression interface{}

type ParseTree struct {
	Children []*Node `json:"children"`
}

type Node struct {
	Type string     `json:"type"`
	Expr Expression `json:"expr"`
}

func (x *ParseTree) add(ctype string, child Expression) {
	x.Children = append(x.Children, &Node{ctype, child})
}

func (x *ParseTree) ForEach(fn func(n *Node) error) error {
	var err error
	for _, child := range x.Children {
		err = fn(child)
		if err != nil {
			return err
		}
	}
	return nil
}
