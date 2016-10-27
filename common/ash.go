package common

import (
  "github.com/codegangsta/cli"
  "os"
  "os/exec"
  "strings"
  "errors"
  "fmt"
)

func cliRun(c *cli.Context) error {
  if c.IsSet("version") {
    fmt.Println(c.App.Version)
    return nil
  }

  if c.IsSet("help") {
    cli.ShowAppHelp(c)
    return nil
  }

  setRemLevel(c.Int("verbosity"))

  ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  // Build the ssh command

  sshParts := make([]string, 1, 30)
  sshParts[0] = "ssh"

  args := c.Args()
  at := ""
  for i, arg := range args {
    if strings.Contains(arg, "@") {
      at = arg
      // Cut it out
      args = append(args[:i], args[i + 1:]...)
      break
    }
  }

  if !(at != "" || c.IsSet("host") || c.IsSet("instance") || c.IsSet("group") || c.IsSet("tag")) {
    return errors.New("To whom do I connect? Specify one of: user@host, --host/-h, --instance/-m, --group/-g, --tag/-t")
  }

  dbgf("Resolving ssh_config. Some inferences are sensitive to ssh_config.")
  sshConfigFile := c.String("config-file")
  parsedSshConfig, err := resolveSshConfig(sshConfigFile)
  if err != nil {
    return err
  }

  dbgf("Resolving hosts first. We may need EC2 info to infer other ssh params.")
  host, ec2, err := resolveHost(at, c.String("host"), c.String("instance"), c.String("group"), c.String("tag"))
  if err != nil {
    return err
  }

  dbgf("Resolving user.")
  user, err := resolveUser(at, c.String("user"), host, ec2, parsedSshConfig)
  if err != nil {
    return err
  }

  dbgf("Resolving identity.")
  ident, err := resolveIdent(c.String("identity"), c.Bool("kms"), ec2, parsedSshConfig)
  if err != nil {
    return err
  }

  dbgf("Assembling ssh args.")
  if sshConfigFile != "" {
    sshParts = append(sshParts, "-F", sshConfigFile)
  }

  if ident != "" {
    sshParts = append(sshParts, "-i", ident)
  }

  if user != "" {
    sshParts = append(sshParts, user + "@" + host)
  } else if host != "" {
    sshParts = append(sshParts, host)
  }

  inff("SSH command: %s", sshParts)
  if c.Bool("Agent") {
    inff("  with ssh-agent")
  }
  inff("Remote command: %s", args)

  sshCmd := strings.Join(sshParts, " ") + " " + strings.Join(args, " ")

  ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  // Build the script

  scriptLines := make([]string, 1, 10)

  if c.Bool("Agent") {
    scriptLines = append(scriptLines, "eval $(ssh-agent) >> /dev/null", "ssh-add ~/.ssh/{id_rsa,*.pem} 2> /dev/null")
  }

  scriptLines = append(scriptLines, sshCmd)

  if c.Bool("Agent") {
    scriptLines = append(scriptLines, "kill $SSH_AGENT_PID")
  }

  ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  // Execution

  cmd := exec.Command("bash", "-c", strings.Join(scriptLines, "\n"))
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  return cmd.Run()
  return err
}

func Run() {
  err := InitSshConfig()
  if err != nil {
    erro(err)
    return
  }

  app := cli.NewApp()
  app.Name = "ash"
  app.Version = "0.1.0"
  app.Usage = "AWS EC2 ssh tool"
  app.HideHelp = true  // Help conflicts with that of host flag. Disable it.
  app.HideVersion = true // -v conflicts with verbosity flag. Disabled it.
  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "host, h",
      EnvVar: "ASH_HOST",
      Value: "",
      Usage: "ssh by hostname",
    },
    cli.StringFlag{
      Name: "instance, machine, m",
      EnvVar: "ASH_MACHINE",
      Value: "",
      Usage: "ssh by EC2 instance id",
    },
    cli.StringFlag{
      Name: "group, g",
      EnvVar: "ASH_GROUP",
      Value: "",
      Usage: "ssh by auto-scaling group name",
    },
    cli.StringFlag{
      Name: "tag, t",
      EnvVar: "ASH_TAG",
      Value: "",
      Usage: "ssh by EC2 tag",
    },
    cli.StringFlag{
      Name: "user, u",
      EnvVar: "ASH_USER",
      Value: "",
      Usage: "ssh as this username",
    },
    cli.StringFlag{
      Name: "identity, i",
      EnvVar: "ASH_IDENTITY",
      Value: "",
      Usage: "ssh identified by this private key file",
    },
    cli.BoolFlag{
      Name: "kms, k",
      EnvVar: "ASH_IDENTITY",
      Usage: "ssh identified by a private key from KMS",
    },
    cli.BoolFlag{
      Name: "Agent, A",
      EnvVar: "ASH_AGENT",
      Usage: "ssh identified by any private key in ~/.ssh/{id_rsa,*.pem} via ssh-agent",
    },
    cli.StringFlag{
      Name: "config-file, F",
      EnvVar: "ASH_CONFIG_FILE",
      Value: "",
      Usage: "ssh -F option: use a config file other than the default (usually ~/.ssh/{ssh_,}config",
    },
    cli.IntFlag{
      Name: "verbosity, v",
      EnvVar: "ASH_VERBOSITY",
      Value: 2,
      Usage: "ash verbosity: 0 - TRACE, 1 - DEBUG, 2 - INFO (default level), 3 - WARN, 4 - ERROR)",
    },
    // Re-instate long version flag.
    cli.BoolFlag{
      Name: "version",
      Usage: "print the version",
    },
    // Re-instate long help flag.
    cli.BoolFlag {
      Name: "help",
      Usage: "show help",
    },
  }
  app.Action = cliRun

  err = app.Run(os.Args)
  if err != nil {
    erro(err)
    return
  }
}
