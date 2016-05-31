package ash

import (
  "os"
  "log"
  "github.com/fatih/color"
)

var serr = log.New(os.Stderr, "", 0)
var logLevel = 2

func setRemLevel(level int) {
  logLevel = level
}

func trcf(format string, args ...interface{}) {
  if logLevel <= 0 {
    serr.Println(color.BlueString("[TRC] " + format, args...))
  }
}

func dbg(s string) {
  dbgf("%s", s)
}

func dbgf(format string, args ...interface{}) {
  if logLevel <= 1 {
    serr.Println(color.CyanString("[DBG] " + format, args...))
  }
}

func inff(format string, args ...interface{}) {
  if logLevel <= 2 {
    serr.Printf("[INF] " + format + "\n", args...)
  }
}

func wrnf(format string, args ...interface{}) {
  if logLevel <=3 {
    serr.Println(color.YellowString("[WRN] " + format, args...))
  }
}

func errf(format string, args ...interface{}) {
  if logLevel <= 4 {
    serr.Println(color.RedString("[ERR] " + format, args...))
  }
}

func erro(e error) {
  if logLevel <= 4 && e != nil {
    serr.Println(color.RedString("[ERR] %s", e))
  }
}
