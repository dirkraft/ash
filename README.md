ash
===

ash is an AWS EC2 ssh tool focused on convenience.

This is a work in progress, but the intention is that the options
illustrated by the examples below are implemented.



### Get It ###

Dev builds: https://github.com/dirkraft/ash/releases/tag/dev

Only a few platforms are covered right now but I can add more supported
build targets. I need to know which GOOS+GOARCH:
https://golang.org/doc/install/source#environment



### Examples ###

SSH to a running instance in auto-scaling **-g**roup `prod/webapp1`
as **-u**ser `ubuntu` using all available keys matching
`~/.ssh/{id_rsa,*.pem}` with the ssh**-A**gent.

    ash -g prod/webapp1 -u ubuntu -A

or by some arbitrary **-t**ag 

    ash -t appName=jenkins
    
or by EC2 **-m**achine's instance id

    ash -m i-12345678

Other options to come.