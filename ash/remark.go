package ash

import (
  "os"
  "log"
  "github.com/fatih/color"
)

var serr = log.New(os.Stderr, "", 0)

func rem(format string, args ...interface{}) {
  serr.Println(color.YellowString(format, args...))
}
