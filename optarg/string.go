/*
Copyright (c) 2009-2010 Jim Teeuwen.

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

    1. The origin of this software must not be misrepresented; you must not
    claim that you wrote the original software. If you use this software
    in a product, an acknowledgment in the product documentation would be
    appreciated but is not required.

    2. Altered source versions must be plainly marked as such, and must not be
    misrepresented as being the original software.

    3. This notice may not be removed or altered from any source distribution.

*/
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

	linesize--

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

func align(v, pad string, linesize, size, alignment int) string {
	var data []byte
	buf := bytes.NewBuffer(data)

	switch alignment {
	case ALIGN_LEFT:
		buf.WriteString(pad)
		buf.WriteString(v)

	case ALIGN_RIGHT:
		diff := linesize - len(v) - len(pad)
		for n := 0; n < diff; n++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(v)

	case ALIGN_CENTER:
		diff := (size - len(v)) / 2
		buf.WriteString(pad)
		for n := 0; n < diff; n++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(v)

	case ALIGN_JUSTIFY:
		if strings.Index(v, " ") == -1 {
			buf.WriteString(pad)
			buf.WriteString(v)
			return buf.String()
		}

		diff := size - len(v)
		if diff == 0 {
			buf.WriteString(pad)
			buf.WriteString(v)
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

		buf.WriteString(pad)
		buf.WriteString(v)
	}

	return buf.String()
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
