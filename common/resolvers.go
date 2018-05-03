package common

import (
  "github.com/aws/aws-sdk-go/aws"
  "strings"
  "os"
  "path/filepath"
  "errors"
  "net"
  "github.com/aws/aws-sdk-go/service/ec2"
  "os/user"
)


// TODO return lazily initialized *ec2.Instance search. We may not actually need it.
func resolveHost(at, explicitHost, instanceId, group, tag string, private bool) (string, *ec2.Instance, error) {

  if at != "" || explicitHost != "" {
    if at != "" {
      inff("Using explicitly given EC2 host in user@host arg: %s", at)
      explicitHost = strings.Split(at, "@")[1]
    } else {
      inff("Using explicitly given EC2 host: %s", explicitHost)
    }

    dbgf("Is it already an ip address? %s", explicitHost)
    ipAddr := net.ParseIP(explicitHost)
    if ipAddr == nil {
      dbgf("Turn it into ip for public ip address EC2 filter.")
      ipAddrs, err := net.LookupIP(explicitHost)
      if err == nil {
        ipAddr = ipAddrs[0]
      } else {
        dbgf("It may be an alias. Not much else we can do.")
      }
    }

    ipAttr := "ip-address"
    if private {
      ipAttr = "private-ip-address"
    }

    dbgf("Finding EC2 instance by %s=%s", ipAttr, ipAddr)
    ec2_, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String(ipAttr), Values: []*string{aws.String(ipAddr.String())}},
      },
    })
    return explicitHost, ec2_, err
  }

  if instanceId != "" {
    inff("Finding EC2 by instance id: %s", instanceId)
    ec2_, err := findEc2(&ec2.DescribeInstancesInput{
      InstanceIds: []*string{aws.String(instanceId)},
    })
    if err != nil {
      return "", nil, err
    }
    return firstAddress(ec2_, private), ec2_, err
  }

  if group != "" {
    inff("Finding EC2 by auto-scaling group: %s", group)
    ec2_, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String("tag:aws:autoscaling:groupName"), Values: []*string{aws.String(group)}},
      },
    })
    if err != nil {
      return "", nil, err
    }
    return firstAddress(ec2_, private), ec2_, err
  }

  if tag != "" {
    inff("Finding EC2 by tag: %s", tag)
    tagParts := strings.Split(tag, "=")
    ec2_, err := findEc2(&ec2.DescribeInstancesInput{
      Filters: []*ec2.Filter{
        {Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
        {Name: aws.String("tag:" + tagParts[0]), Values: []*string{aws.String(tagParts[1])}},
      },
    })
    if err != nil {
      return "", nil, err
    }
    return firstAddress(ec2_, private), ec2_, err
  }

  return "", nil, errors.New("Unable locate suitable EC2 instance.")
}

func firstAddress(ec2_ *ec2.Instance, private bool) string {
  if private && *ec2_.PrivateDnsName != "" {
    dbgf("Using PrivateDnsName.")
    return *ec2_.PrivateDnsName
  } else if *ec2_.PublicDnsName != "" {
    dbgf("Using PublicDnsName.")
    return *ec2_.PublicDnsName
  // TODO This one is problematic because it is omitted from raw JSON when the instance has none
  //} else if *ec2_.PublicIpAddress != "" {
  //  dbgf("Using PublicIpAddress.")
  //  return *ec2_.PublicIpAddress
  } else {
    dbgf("Using PrivateDnsName.")
    return *ec2_.PrivateDnsName
  }
}

func resolveSshConfig(configFilePath string) (*SshConfig, error) {

  var paths []string
  if configFilePath != "" {
    paths = []string{configFilePath}
  } else {
    var u, err = user.Current()
    if err != nil {
      return nil, err
    }
    // Should work for Windows and Unix
    windowsCompatHomeDir := strings.Replace(u.HomeDir, "\\", "/", -1)
    paths = []string{windowsCompatHomeDir + "/.ssh/config", "/etc/ssh/ssh_config"}
  }

  return ParseSshConfig(paths, true)
}

func resolveUser(at, explicitUser, resolvedHost string, instance *ec2.Instance, sshConfig *SshConfig) (string, error) {

  if at != "" {
    explicitUser = strings.Split(at, "@")[0]
    inff("Authenticating as given user in user@host arg: %s", at)
    return explicitUser, nil
  }

  if explicitUser != "" {
    inff("Authenticating as given user: %s", explicitUser)
    return explicitUser, nil
  }

  resolvedUser, err := sshConfig.GetConfigValue(resolvedHost, "User")
  if err != nil {
    return "", err
  }
  if resolvedUser != "" {
    dbgf("ssh config specifies user %s for host %s. Will not infer user from EC2 metadata.", resolvedUser, resolvedHost)
    return resolvedUser, nil
  }

  dbgf("Reading details of %s", *instance.ImageId)
  ami, err := findAmi(&ec2.DescribeImagesInput{
    ImageIds:[]*string{aws.String(*instance.ImageId)},
  })
  if err != nil {
    return "", err
  }

  switch {
  case strings.HasPrefix(*ami.Name, "amzn-"):
    explicitUser = "ec2-user"
  case strings.HasPrefix(*ami.Name, "ubuntu/"):
    explicitUser = "ubuntu"
  case strings.HasPrefix(*ami.Name, "debian-"):
    explicitUser = "admin"
  }

  if explicitUser != "" {
    inff("Authenticating based on %s as user: %s", *ami.ImageId, explicitUser)
  }
  return explicitUser, nil
}

func resolveIdent(identity string, useKms bool, ec2 *ec2.Instance, sshConfig *SshConfig) (string, error) {
  if identity != "" {
    // Is it a valid path already?
    _, err := os.Stat(identity)
    if err == nil {
      inff("Identifying by given file: %s", identity)
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
      inff("Identifying by file: %s", expandedPath)
      return expandedPath, nil
    }
    if !os.IsNotExist(err) {
      // If some error other than not existing...
      return "", err
    }
    // Otherwise we're going to look some more.

    // Is it the name of a private key in ~/.ssh/ without the pem suffix?
    expandedPath += ".pem"
    _, err = os.Stat(expandedPath)
    if err == nil {
      inff("Identifying by file: %s", expandedPath)
      return expandedPath, nil
    }

    return "", errors.New("Failed to resolve given identity: " + identity)
  }

  if useKms {
    // TODO read key of ec2 instance from kms
    return "", nil
  }

  // TODO If not already covered by ssh config, find key of ec2 instance locally in ~/.ssh/
  return "", nil
}
