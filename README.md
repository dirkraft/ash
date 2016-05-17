ash
===

ash is an AWS EC2 ssh tool focused on convenience.


### Get It ###



### Examples ###

SSH to a running instance in auto-scaling `-g`roup **prod/webapp1**
as user **ubuntu** using all available keys with the ssh`-A`gent
matching `~/.ssh/{id_rsa,*.pem}`.

    ash -g prod/webapp1 -u ubuntu -A

or by some arbitrary `-t`ag 

    ash -t appName=jenkins
    
or by EC2 `-m`achine's instance id

    ash -m i-12345678

Other options to come.