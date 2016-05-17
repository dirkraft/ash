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

or by some arbitrary **-t**ag, **-i**dentified by some private key

    ash -t appName=jenkins -i ~/.ssh/devkey.pem
    
or by EC2 **-m**achine's instance id, **-i**dentified by a private key
in ~/.ssh/devkey.pem using a shortcut

    ash -m i-12345678 -i devkey

Other options to come.



### Development ###

Get any dependencies.

    go get 

To run commands in development, just replace `ash` with `go run main.go` 

    go run main.go ash_args...

To build binaries for all targeted platforms

    make clean all

For me to publish new dev releases to
https://github.com/dirkraft/ash/releases/tag/dev

    export GITHUB_TOKEN
    make clean publish-dev