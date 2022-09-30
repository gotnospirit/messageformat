package messageformat

import "bytes"

type (
	Expression interface{}

	node struct {
		children []*nodeExpr
	}

	nodeExpr struct {
		ctype string
		expr  Expression
	}
)

func (x *node) add(ctype string, child Expression) {
	x.children = append(x.children, &nodeExpr{ctype, child})
}

func (x *node) format(ptr_output *bytes.Buffer, data *map[string]interface{}, ptr_mf *MessageFormat, pound string) error {
	for _, child := range x.children {
		ctype := child.ctype

		fn, err := ptr_mf.getFormatter(ctype)
		if err != nil {
			return err
		}

		err = fn(child.expr, ptr_output, data, ptr_mf, pound)
		if err != nil {
			return err
		}
	}
	return nil
}
