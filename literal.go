package messageformat

import (
	"bytes"
)

func formatLiteral(expr Expression, ptr_output *bytes.Buffer, _ *map[string]interface{}, _ *MessageFormat, pound string) error {
	content := expr.([]string)

	for _, c := range content {
		if c != "" {
			ptr_output.WriteString(c)
		} else if pound != "" {
			ptr_output.WriteString(pound)
		} else {
			ptr_output.WriteRune(PoundChar)
		}
	}
	return nil
}

func parseLiteral(start, end int, ptr_input *[]rune) []string {
	var items []int

	input := *ptr_input
	escaped := false

	s, e := start, start
	gap := 0
	for i := start; i < end; i++ {
		c := input[i]

		if c == EscapeChar {
			gap++
			e++
			escaped = true
		} else {
			switch c {
			default:
				e++

			case OpenChar, CloseChar, PoundChar:
				if escaped {
					if i-s > gap {
						if gap > 1 {
							items = append(items, s, i)
						} else {
							items = append(items, s, i-1)
						}
					}
					s = i
				} else {
					if s != e {
						items = append(items, s, e, i, i)
					} else if s != i {
						items = append(items, s, i, i, i)
					} else {
						items = append(items, i, i)
					}
					s = i + 1
				}
				e = s
			}

			escaped = false
			gap = 0
		}
	}

	if s < end {
		items = append(items, s, end)
	}

	n := len(items)
	result := make([]string, n/2)
	for i := 0; i < n; i += 2 {
		result[i/2] = string(input[items[i]:items[i+1]])
	}
	return result
}
