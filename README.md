ash
===

ash is an AWS EC2 ssh command-line tool focused on convenience.

The intention is that the options illustrated by the examples below are implemented, but there may be some incomplete stubs here and there.

<a href="https://travis-ci.org/dirkraft/ash">
<img src="https://travis-ci.org/dirkraft/ash.svg?branch=master">
</a>



### Get It ###

Dev builds: https://github.com/dirkraft/ash/releases/tag/dev

For OSX, that might look like this:

    curl -o /usr/local/bin/ash -L https://github.com/dirkraft/ash/releases/download/dev/ash.darwin.amd64
    chmod +x /usr/local/bin/ash

Only a few platforms are covered right now but I can add more supported
build targets. I need to know which GOOS+GOARCH:
https://golang.org/doc/install/source#environment



### Examples ###

SSH to a running instance in auto-scaling **-g**roup `prod/webapp1`
as **-u**ser `ubuntu` using all available keys matching
`~/.ssh/{id_rsa,*.pem}` with the ssh**-A**gent.

    ash -g prod/webapp1 -u ubuntu -A

SSH by some arbitrary **-t**ag, **-i**dentified by some private key.
No user has been explicitly given, so it will be guessed based on
the EC2 instance's AMI, unless the user's ssh_config defines a specific
user to connect with.

    ash -t Name=jenkins -i ~/.ssh/devkey.pem
    
SSH by EC2 **-m**achine's instance id, **-i**dentified by a private key
in ~/.ssh/devkey.pem using a shortcut. 

    ash -m i-12345678 -i devkey

Other options to come.



### AWS Config ###

ash relies on [aws-sdk-go](https://github.com/aws/aws-sdk-go) which
reads auth credentials from certain places (i.e. `~/.aws/credentials`)
and requires `AWS_REGION` to be exported. ash otherwise defaults to
`AWS_REGION=us-east-1`.



### SSH Config ###

ash combined with ssh config settings can make for an optimally
convenient EC2 SSH experience. Say you use the same master key for
all EC2 instances. In your `~/.ssh/config` this will use that
IdentityFile for all EC2 connections, unless you specify one explicitly
with the `-i/--identity` option.

    Host ec2-*.compute-1.amazonaws.com
        IdentityFile path/to/privatekey.pem

Ash uses DNS names instead of IP addresses which enables ssh config host
pattern options such as this. SSH config is preferred over *inferenced
parameters* (like guessing the user based on AMI) whenever the local
ssh config has relevant configuration. 



### Development ###

Get any dependencies via a make target.

    make develop

To run commands in development, replace `ash` with `go run ...`, e.g. 

    go run ash/main.go ash_args...

To build binaries for all targeted platforms

    make clean all

For me to publish new dev releases to
https://github.com/dirkraft/ash/releases/tag/dev

    export GITHUB_TOKEN
    make clean publish-dev



### CLI Documentation ###

```
$ ash --help
NAME:
   ash - AWS EC2 ssh tool

USAGE:
   ash [global options] [arguments...]

GLOBAL OPTIONS:
   --host value, -h value                       ssh by hostname [$ASH_HOST]
   --instance value, --machine value, -m value  ssh by EC2 instance id [$ASH_MACHINE]
   --group value, -g value                      ssh by auto-scaling group name [$ASH_GROUP]
   --tag value, -t value                        ssh by EC2 tag [$ASH_TAG]
   --private, -p                                When resolving host, use AWS private DNS name. [$ASH_PRIVATE_IP]
   --user value, -u value                       ssh as this username [$ASH_USER]
   --identity value, -i value                   ssh identified by this private key file [$ASH_IDENTITY]
   --kms, -k                                    NOT IMPLEMENTED ssh identified by a private key from KMS [$ASH_IDENTITY]
   --Agent, -A                                  ssh identified by any private key in ~/.ssh/{id_rsa,*.pem} via ssh-agent [$ASH_AGENT]
   --config-file value, -F value                ssh -F option: use a config file other than the default (usually ~/.ssh/{ssh_,}config [$ASH_CONFIG_FILE]
   --verbosity value, -v value                  ash verbosity: 0 - TRACE, 1 - DEBUG, 2 - INFO (default level), 3 - WARN, 4 - ERROR) (default: 2) [$ASH_VERBOSITY]
   --version                                    print the version
   --help                                       show help
 ```



### License ###

```
The MIT License (MIT)
Copyright (c) 2016 Jason Dunkelberger (a.k.a. "dirkraft")

Permission is hereby granted, free of charge, to any person obtaining a 
copy of this software and associated documentation files (the 
"Software"), to deal in the Software without restriction, including 
without limitation the rights to use, copy, modify, merge, publish, 
distribute, sublicense, and/or sell copies of the Software, and to 
permit persons to whom the Software is furnished to do so, subject to 
the following conditions:

The above copyright notice and this permission notice shall be included 
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS 
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF 
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. 
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY 
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, 
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE 
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
