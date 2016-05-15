package ash

import (
  "github.com/codegangsta/cli"
  "os"
  "os/exec"
  "fmt"
  "strings"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "errors"
)

var ec2Svc = ec2.New(session.New(), &aws.Config{Region:aws.String("us-east-1")})

func findEc2(params *ec2.DescribeInstancesInput) (*ec2.Instance, error) {
  resp, err := ec2Svc.DescribeInstances(params)
  if err != nil {
    return nil, err
  }
  ec2 := resp.Reservations[0].Instances[0]
  fmt.Printf("Located EC2: id %s, pub %s, priv %s, launch %s\n",
    *ec2.InstanceId, *ec2.PublicIpAddress, *ec2.PrivateIpAddress, *ec2.LaunchTime)
  return ec2, nil
}

func resolveHost(instanceId, group, tag string) (string, *ec2.Instance, error) {

  if instanceId != "" {
    fmt.Println("Finding EC2 by instance id: %s", instanceId)
    ec2, err := findEc2(&ec2.DescribeInstancesInput{
      InstanceIds: []*string{aws.String(instanceId)},
    })
    return *ec2.PublicIpAddress, ec2, err
  }

  if group != "" {
    fmt.Println("Finding EC2 by auto-scaling group: %s", group)
    ec2, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String("tag:aws:autoscaling:groupName"), Values: []*string{aws.String(group)}},
      },
    })
    return *ec2.PublicIpAddress, ec2, err
  }

  if tag != "" {
    fmt.Println("Finding EC2 by tag: %s", tag)
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
    fmt.Println("Authenticating as given user: %s", user)
    return user, nil
  }
  // TODO guess user based on AMI: debian>admin, ubuntu>ubuntu, amzn>ec2-user
  return "", nil
}

func resolveIdent(identity string, useKms bool, ec2 *ec2.Instance) (string, error) {
  if identity != "" {
    // TODO resolve to identity file
    return "", nil
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
    },
    cli.StringFlag{
      Name: "instance, machine, m",
      Value: "",
    },
    cli.StringFlag{
      Name: "group, g",
      Value: "",
    },
    cli.StringFlag{
      Name: "tag, t",
      Value: "",
    },
    cli.StringFlag{
      Name: "user, u",
      Value: "",
    },
    cli.StringFlag{
      Name: "identity, i",
      Value: "",
    },
    cli.BoolFlag{
      Name: "kms, k",
    },
  }
  app.Action = func(c *cli.Context) error {
    parts := make([]string, 0, 20)

    args := c.Args()
    //at := ""
    //for i, arg := range args {
    //  if strings.Contains(arg, "@") {
    //    at = arg
    //    // Cut it out
    //    args = append(args[:i], args[i + 1:]...)
    //    break
    //  }
    //}

    host, ec2, err := resolveHost(c.String("instance"), c.String("group"), c.String("tag"))
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
      parts = append(parts, "-i", ident)
    }

    if user != "" {
      parts = append(parts, user + "@" + host)
    } else if host != "" {
      parts = append(parts, host)
    }

    fmt.Printf("Command: %s\n", parts)
    fmt.Printf("Remaining: %s\n", args)

    cmd := exec.Command("ssh", parts...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    return cmd.Run()
  }

  err := app.Run(os.Args)
  if err != nil {
    fmt.Printf("Error: %s\n", err)
  }
}
