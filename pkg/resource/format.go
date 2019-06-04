package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/kr/text"
)

// hintPrefixMap contains a mapping of hints to prefix symbols for the output.
var hintPrefixMap = map[Hint]string{
	NoChange: "*",
	Addition: "+",
	Update:   "~",
	Removal:  "-",
}

// hintColorFuncMap contains a mapping of hints to color printing functions.
var hintColorFuncMap = map[Hint]func(string, ...interface{}) string{
	Addition: color.GreenString,
	Update:   color.YellowString,
	Removal:  color.RedString,
}

// Format formats the resource as string.
func Format(r *Resource) string {
	return text.Indent(format(r), "  ")
}

// format formats the resource. It will enrich the output based on the resource
// hints. If the Updated hint is set on the resource, and it also received a
// contentHint via WithContentHint, a diff will be added to the formatted
// output only if the diff itself is not empty.
func format(r *Resource) string {
	colorFunc := hintColorFunc(r.hint)
	prefix := hintPrefix(r.hint)
	s := r.String()

	switch r.hint {
	case Update:
		prefix = colorFunc(prefix)
		s = colorFunc(s)

		diff := r.diff()

		if diff != "" {
			return fmt.Sprintf("%s %s\n\n%s", prefix, s, strings.TrimSpace(diff))
		}
	}

	return colorFunc("%s %s", prefix, s)
}

// FormatSlice formats a slice of resources as string. If s is nil or empty the
// formatted string will also be empty. The formatted output will be prepended
// with a summary of the counts of different resource hints (e.g. updates).
func FormatSlice(s Slice) string {
	if len(s) == 0 {
		return ""
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "%d resources (%s)\n\n", len(s), summarize(s))

	for _, r := range s {
		sb.WriteString(Format(r))
		sb.WriteString("\n\n")
	}

	return sb.String()
}

// hintPrefix returns the prefix symbol for h or ? if no mapping exists.
func hintPrefix(h Hint) string {
	prefix, ok := hintPrefixMap[h]
	if ok {
		return prefix
	}

	return "?"
}

// hintColorFunc will return the color func for given hint. Will fall back to
// fmt.Sprintf if h does not exist in the hintColorFuncMap.
func hintColorFunc(h Hint) func(string, ...interface{}) string {
	colorFunc := hintColorFuncMap[h]
	if colorFunc == nil {
		return fmt.Sprintf
	}

	return colorFunc
}

// summarize walks s and counts all the different hints it finds on the
// resources. It will then compile a summary of these and return it as a
// string.
func summarize(s Slice) string {
	buckets := make(map[Hint]int)

	for _, r := range s {
		buckets[r.hint]++
	}

	keys := make([]int, 0, len(buckets))

	for k := range buckets {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	summary := make([]string, len(keys))

	for i, k := range keys {
		h := Hint(k)
		colorFunc := hintColorFunc(h)
		prefix := hintPrefix(h)

		summary[i] = fmt.Sprintf("%s %s: %d", colorFunc(prefix), h, buckets[h])
	}

	return strings.Join(summary, ", ")
}
