package log

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type contextKey string

// Prefix is the context key for the log prefix.
const Prefix contextKey = "prefix"

// ContextPrefixFormatter prefixes log messages based on the entry context.
type ContextPrefixFormatter struct {
	*logrus.TextFormatter
}

// Format implements logrus.Formatter. Will prepend the message of entry with a
// prefix present in the log context if exists. For all other formatting the
// default logrus.TextFormatter is used.
func (f *ContextPrefixFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if prefix := f.prefix(entry); prefix != nil {
		entry.Message = fmt.Sprintf("%s %s", prefix, entry.Message)
	}

	return f.TextFormatter.Format(entry)
}

// prefix extracts the log prefix from entry.Context if it exists.
func (f *ContextPrefixFormatter) prefix(entry *logrus.Entry) interface{} {
	if entry.Context == nil {
		return nil
	}

	return entry.Context.Value(Prefix)
}

// ContextWithPrefix creates a context.Context which has given prefix as value.
func ContextWithPrefix(prefix string) context.Context {
	return context.WithValue(context.Background(), Prefix, prefix)
}
