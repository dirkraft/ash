package common

import (
  "os"
  "errors"
  "fmt"
  "bufio"
  "strings"
  "regexp"
)

type SshConfig map[string]SshConfigHost
type SshConfigHost map[string]string

var REGEXP_WHITESPACE regexp.Regexp

func InitSshConfig() error {
  regexpWhitespace, err := regexp.Compile("\\s+")
  REGEXP_WHITESPACE = *regexpWhitespace
  return err
}

func mergeSshConfig(path string, swallowNotExists bool, sshConfig *SshConfig) error {

  if _, err := os.Stat(path); os.IsNotExist(err) {
    msg := fmt.Sprintf("Can't read %s. Reason: %s", path, err)
    if swallowNotExists {
      dbg(msg)
      return nil
    } else {
      return errors.New(msg)
    }
  }

  inff("Referencing ssh config path: %s", path)
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  scanner := bufio.NewScanner(bufio.NewReader(file))

  var currentHostPattern string
  for scanner.Scan() {
    line := strings.TrimSpace(string(scanner.Bytes()))
    if len(line) == 0 {
      continue
    }

    if strings.HasPrefix(line, "#") {
      continue
    }

    parts := REGEXP_WHITESPACE.Split(line, 2)
    if len(parts) != 2 {
      wrnf("Skipping failed parse of directive in %s: %s", path, line)
      continue
    }

    key := parts[0]
    val := parts[1]
    if key == "Host" {
      currentHostPattern = val
      trcf("Parsed %s:", currentHostPattern)
    } else {
      trcf("    %s %s", key, val)
    }

    if _, exists := (*sshConfig)[currentHostPattern]; !exists {
      (*sshConfig)[currentHostPattern] = SshConfigHost{}
    }
    (*sshConfig)[currentHostPattern][key] = val
  }

  return nil
}

func ParseSshConfig(paths []string, swallowNotExists bool) (*SshConfig, error) {
  sshConfig := &SshConfig{}

  for _, path := range paths {
    err := mergeSshConfig(path, swallowNotExists, sshConfig)
    if err != nil {
      return nil, err
    }
  }

  return sshConfig, nil
}

func (sshConfig *SshConfig) GetConfigValue(targetHost, key string) (string, error) {

  for _, sshConfigHost := range *sshConfig {
    hostRegex := strings.Replace(sshConfigHost["Host"], "*", ".*", -1)
    matches, err := regexp.MatchString(hostRegex, targetHost)
    trcf("ssh config pattern %s matches %s? %t", sshConfigHost["Host"], targetHost, matches)
    if err != nil {
      return "", err
    } else if matches {

      trcf("%s=%s", key, sshConfigHost[key])
      if val, exist := sshConfigHost[key]; exist {
        return val, nil
      }
    }
  }

  return "", nil
}

