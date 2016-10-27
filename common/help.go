package common

func firstString(vals ...string) string {
  for _, v := range vals {
    if v != "" {
      return v
    }
  }
  return ""
}
