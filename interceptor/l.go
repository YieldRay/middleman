package interceptor

import (
	"github.com/mborders/logmatic"
	"github.com/yieldray/middleman/cmd/flags"
)

var l = logmatic.NewLogger()

func init() {
	l.SetLevel(logmatic.LogLevel(flags.LogLevel))
}