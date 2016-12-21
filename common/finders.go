package common

import (
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "os"
  "errors"
)

var ec2Svc = ec2.New(session.New(), &aws.Config{
  Region:aws.String(firstString(os.Getenv("AWS_REGION"), "us-east-1")),
})

func findEc2(params *ec2.DescribeInstancesInput) (*ec2.Instance, error) {
  resp, err := ec2Svc.DescribeInstances(params)
  if err != nil {
    return nil, err
  }
  if len(resp.Reservations) == 0 {
    return nil, errors.New("Could not find any matching EC2 instances.")
  }
  if len(resp.Reservations[0].Instances) == 0 {
    return nil, errors.New("Could not find any matching EC2 instances.")
  }
  ec2_ := resp.Reservations[0].Instances[0]
  inff("Located EC2: id %s, dns %s, launch %s", *ec2_.InstanceId,
    firstString(*ec2_.PublicDnsName, /* *ec2_.PublicIpAddress,*/ *ec2_.PrivateDnsName), *ec2_.LaunchTime)
  return ec2_, nil
}

func findAmi(params *ec2.DescribeImagesInput) (*ec2.Image, error) {
  resp, err := ec2Svc.DescribeImages(params)
  if err != nil {
    return nil, err
  }

  ami := resp.Images[0]
  dbgf("Located AMI: id %s, name %s", *ami.ImageId, *ami.Name)
  return ami, nil
}

