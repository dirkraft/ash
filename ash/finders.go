package ash

import (
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
)

var ec2Svc = ec2.New(session.New(), &aws.Config{Region:aws.String("us-east-1")})

func findEc2(params *ec2.DescribeInstancesInput) (*ec2.Instance, error) {
  resp, err := ec2Svc.DescribeInstances(params)
  if err != nil {
    return nil, err
  }
  ec2 := resp.Reservations[0].Instances[0]
  inff("Located EC2: id %s, pub %s, launch %s", *ec2.InstanceId, *ec2.PublicDnsName, *ec2.LaunchTime)
  return ec2, nil
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

