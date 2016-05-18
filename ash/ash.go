package ash

import (
  "github.com/codegangsta/cli"
  "os"
  "os/exec"
  "strings"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/go-errors/errors"
  "path/filepath"
  "net"
)

var ec2Svc = ec2.New(session.New(), &aws.Config{Region:aws.String("us-east-1")})

func findEc2(params *ec2.DescribeInstancesInput) (*ec2.Instance, error) {
  resp, err := ec2Svc.DescribeInstances(params)
  if err != nil {
    return nil, err
  }
  ec2 := resp.Reservations[0].Instances[0]
  rem("Located EC2: id %s, pub %s, priv %s, launch %s",
    *ec2.InstanceId, *ec2.PublicIpAddress, *ec2.PrivateIpAddress, *ec2.LaunchTime)
  return ec2, nil
}

func resolveHost(at, host, instanceId, group, tag string) (string, *ec2.Instance, error) {

  if at != "" || host != "" {
    if at != "" {
      rem("Using explicitly given EC2 host in user@host arg: %s", at)
      host = strings.Split(at, "@")[1]
    } else {
      rem("Using explicitly given EC2 host: %s", host)
    }

    // Is it already an ip address?
    ipAddr := net.ParseIP(host)
    if ipAddr == nil {
      // Turn it into ip for public ip address EC2 filter.
      ipAddrs, err := net.LookupIP(host)
      if err != nil {
        ipAddr = ipAddrs[0]
      }
      // If it's an alias, not much else we can do.
    }

    ec2, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String("ip-address"), Values: []*string{aws.String(ipAddr.String())}},
      },
    })
    return host, ec2, err
  }

  if instanceId != "" {
    rem("Finding EC2 by instance id: %s", instanceId)
    ec2, err := findEc2(&ec2.DescribeInstancesInput{
      InstanceIds: []*string{aws.String(instanceId)},
    })
    return *ec2.PublicIpAddress, ec2, err
  }

  if group != "" {
    rem("Finding EC2 by auto-scaling group: %s", group)
    ec2, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String("tag:aws:autoscaling:groupName"), Values: []*string{aws.String(group)}},
      },
    })
    return *ec2.PublicIpAddress, ec2, err
  }

  if tag != "" {
    rem("Finding EC2 by tag: %s", tag)
    tagParts := strings.Split(tag, "=")
    ec2, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String("tag:" + tagParts[0]), Values: []*string{aws.String(tagParts[1])}},
      },
    })
    return *ec2.PublicIpAddress, ec2, err
  }

  return "", nil, errors.New("Unable locate suitable EC2 instance.")
}

func resolveUser(user string, ec2 *ec2.Instance) (string, error) {
  if user != "" {
    rem("Authenticating as given user: %s", user)
    return user, nil
  }
  // TODO guess user based on AMI: debian>admin, ubuntu>ubuntu, amzn>ec2-user
  return "", nil
}

func resolveIdent(identity string, useKms bool, ec2 *ec2.Instance) (string, error) {
  if identity != "" {
    // Is it a valid path already?
    _, err := os.Stat(identity)
    if err == nil {
      rem("Identifying by given file: %s", identity)
      return identity, nil
    }
    if !os.IsNotExist(err) {
      // If some error other than not existing...
      return "", err
    }
    // Otherwise we're going to look some more.

    // Is it the name of a private key in ~/.ssh/ ?
    expandedPath := filepath.Join(os.Getenv("HOME"), ".ssh", identity)
    _, err = os.Stat(expandedPath)
    if err == nil {
      rem("Identifying by file: %s", expandedPath)
      return expandedPath, nil
    }
    if !os.IsNotExist(err) {
      // If some error other than not existing...
      return "", err
    }
    // Otherwise we're going to look some more.

    // Is it the name of a private key in ~/.ssh/ without the pem suffix.
    expandedPath += ".pem"
    _, err = os.Stat(expandedPath)
    if err == nil {
      rem("Identifying by file: %s", expandedPath)
      return expandedPath, nil
    }

    return "", errors.New("Failed to resolve given identity: " + identity)
  }

  if useKms {
    // TODO read key of ec2 instance from kms
    return "", nil
  }

  // TODO find key of ec2 instance locally in ~/.ssh/
  return "", nil
}

func Run() {

  app := cli.NewApp()
  app.Name = "ash"
  app.HideHelp = true  // Help conflicts with that of host flag. Disable it.
  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "host, h",
      Value: "",
      Usage: "ssh to instance by hostname",
    },
    cli.StringFlag{
      Name: "instance, machine, m",
      Value: "",
      Usage: "ssh to instance by EC2 instance id",
    },
    cli.StringFlag{
      Name: "group, g",
      Value: "",
      Usage: "ssh to instance(s) by auto-scaling group name",
    },
    cli.StringFlag{
      Name: "tag, t",
      Value: "",
      Usage: "ssh to instance(s) by EC2 tag",
    },
    cli.StringFlag{
      Name: "user, u",
      Value: "",
      Usage: "ssh to instance(s) as this username",
    },
    cli.StringFlag{
      Name: "identity, i",
      Value: "",
      Usage: "ssh to instance(s) identified by this private key file",
    },
    cli.BoolFlag{
      Name: "kms, k",
      Usage: "ssh to instance(s) identified by a private key from KMS",
    },
    cli.BoolFlag{
      Name: "Agent, A",
      Usage: "ssh to instance(s) identified by any private key in ~/.ssh/{id_rsa,*.pem} via ssh-agent",
    },
  }
  app.Action = func(c *cli.Context) error {

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
      return errors.New("To whom do I connect? Specific one of: user@host, --host, --instance, --group, --tag")
    }

    // Resolve hosts first. We may need EC2 information for resolution of other ssh parts.
    host, ec2, err := resolveHost(at, c.String("host"), c.String("instance"), c.String("group"), c.String("tag"))
    if err != nil {
      return err
    }

    user, err := resolveUser(c.String("user"), ec2)
    if err != nil {
      return err
    }

    ident, err := resolveIdent(c.String("identity"), c.Bool("kms"), ec2)
    if err != nil {
      return err
    }
    if ident != "" {
      sshParts = append(sshParts, "-i", ident)
    }

    if user != "" {
      sshParts = append(sshParts, user + "@" + host)
    } else if host != "" {
      sshParts = append(sshParts, host)
    }

    rem("SSH command: %s", sshParts)
    if c.Bool("Agent") {
      rem("  with ssh-agent")
    }
    rem("Remote command: %s", args)

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
  }

  err := app.Run(os.Args)
  if err != nil {
    rem("Error: %s", err)
    // TODO with verbose only
    rem(err.(*errors.Error).ErrorStack())
  }
}
