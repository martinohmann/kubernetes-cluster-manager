package revision

import (
	"fmt"
	"strings"

	"github.com/kr/text"
)

func indent(s string) string {
	return text.Indent(s, "  * ")
}

func errorFormatFunc(es []error) string {
	if len(es) == 1 {
		return es[0].Error()
	}

	messages := make([]string, len(es))
	for i, err := range es {
		messages[i] = err.Error()
	}

	return fmt.Sprintf("%d errors occurred:\n%s", len(es), indent(strings.Join(messages, "\n")))
}
