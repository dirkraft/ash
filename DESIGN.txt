

ash
  -u user OR u@h   # infer from AMI, otherwise.
  -i identity      # ~/.ssh/$identity{,.pem} search path
  -k               # read key from KMS
  -h host OR u@h   # derp
  -g asg           # add group constraint
                   #   support * wildcard
  -t name=val      # add tag constraint, multiple can be given
                   #   support * wildcard in value
  -m i-12345678    # instance id
  # Capital options sort of denote Big Deals, changes the basic effect
  -A               # execute on all 
  -T n{-|/|+}[=]   # open tmux to n instances 
                   #   arranged -_r rows, /|c cols, +t tiles 
                   #   = setw synchronize-panes on
  -C {and|or}      # combine restrictions with AND / OR

multiple targets requires a command, or -T

ash base
  # Sets some default options for all ash commands.

asg alias NAME
  # Sets some options with this alias shortcut. The first arg to
  # ash may be interpreted as an alias.

asg ssh_config
  # Recommends some ssh_config rules that might convenience the
  # user if they added to their ~/.ssh/ssh_config
  # e.g. ec2-*.amazonaws.com StrictHostKeyChecking no
