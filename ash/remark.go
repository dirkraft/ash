package ash

import (
  "os"
  "log"
  "github.com/fatih/color"
)

var serr = log.New(os.Stderr, "", 0)

func dbg(format string, args ...interface{}) {
  serr.Println(color.BlueString("[DBG] " + format, args...))
}

func rem(format string, args ...interface{}) {
  serr.Println(color.YellowString("[INF] " + format, args...))
}
