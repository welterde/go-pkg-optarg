package optarg

import "strings"
import "bytes"
import "regexp"

const (
	ALIGN_LEFT = iota
	ALIGN_CENTER
	ALIGN_RIGHT
	ALIGN_JUSTIFY
)

var reg_multilinewrap = regexp.MustCompile("[^a-zA-Z0-9,.]")

func multilineWrap(text string, linesize, leftmargin, rightmargin, alignment int) []string {
	lines := make([]string, 0)
	pad := ""

	for n := 0; n < leftmargin; n++ {
		pad += " "
	}

	if linesize < 1 {
		linesize = 80
	}

	wordboundary := 0
	size := linesize - leftmargin - rightmargin

	if len(text) <= size {
		lines = []string{align(text, pad, linesize, size, alignment)}
		return lines
	}

	for n := 0; n < len(text); n++ {
		if reg_multilinewrap.MatchString(text[n : n+1]) {
			wordboundary = n
		}

		if n > size {
			listappend(
				align(
					strings.TrimSpace(text[0:wordboundary]),
					pad, linesize, size, alignment,
				),
				&lines,
			)
			text = text[wordboundary:len(text)]
			n = 0
		}
	}

	listappend(align(strings.TrimSpace(text), pad, linesize, size, alignment), &lines)
	return lines
}

func align(v, pad string, linesize, size, alignment int) (line string) {
	switch alignment {
	case ALIGN_LEFT:
		line = pad + v

	case ALIGN_RIGHT:
		diff := linesize - len(v) - len(pad)
		for n := 0; n < diff; n++ {
			line += " "
		}
		line += v

	case ALIGN_CENTER:
		diff := (size - len(v)) / 2
		line = pad
		for n := 0; n < diff; n++ {
			line += " "
		}
		line += v

	case ALIGN_JUSTIFY:
		if strings.Index(v, " ") == -1 {
			line = pad + v
			return
		}

		diff := size - len(v)
		if diff == 0 {
			line = pad + v
			break
		}

		spread := "  "
		for {
			v = replaceAll(v, spread[0:len(spread)-1], spread)
			if len(v) > size {
				break
			}
			spread += " "
		}

		for {
			if strings.Index(v, spread) == -1 {
				spread = spread[0 : len(spread)-1]
			}

			v = replaceSingle(v, spread, spread[0:len(spread)-1])
			if len(v) <= size {
				break
			}
		}

		line = pad + v
	}
	return
}

// Resizes @list by 1 and assigns @item to the new slot.
func listappend(item string, list *[]string) {
	c := make([]string, len(*list)+1)
	copy(c, *list)
	c[len(c)-1] = item
	*list = c
}

/*
	ReplaceSingle() replaces first occurance of @n in @h with @r and exits.
*/
func replaceSingle(h, n, r string) (ret string) {
	return string(bytes.Join(bytes.Split([]byte(h), []byte(n), 1), []byte(r)))
}

/*
	ReplaceAll() replaces all non-overlapping occurrences of @n in @h with @r 
	and returns the resulting string.
*/
func replaceAll(h, n, r string) (ret string) {
	return string(bytes.Join(bytes.Split([]byte(h), []byte(n), 0), []byte(r)))
}
