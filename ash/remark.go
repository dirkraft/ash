package ash

import (
  "os"
  "log"
  "github.com/fatih/color"
)

var serr = log.New(os.Stderr, "", 0)

func dbgf(format string, args ...interface{}) {
  serr.Println(color.BlueString("[DBG] " + format, args...))
}

func inff(format string, args ...interface{}) {
  serr.Printf("[INF] " + format + "\n", args...)
}

func wrnf(format string, args ...interface{}) {
  serr.Println(color.YellowString("[WRN] " + format, args...))
}

func errf(format string, args ...interface{}) {
  serr.Println(color.RedString("[ERR] " + format, args...))
}

func erro(e error) {
  if e != nil {
    serr.Println(color.RedString("[ERR] %s", e))
  }
}
