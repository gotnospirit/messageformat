package messageformat

import (
	"fmt"
)

// A parseError is used to embed an error message occurring while the processing an input.
type parseError struct {
	msg string // description of the error
	pos int    // offset read
}

func (x parseError) Error() string {
	return fmt.Sprintf("ParseError: `%s` at %d", x.msg, x.pos)
}
