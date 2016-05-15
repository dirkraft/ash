

ash \
  -u user OR user@ # infer from AMI, otherwise.
  -i identity      # ~/.ssh/$identity{,.pem} search path
  -k               # read key from KMS
  -h host OR @host # derp
  -g asg           # add group constraint
  -t name=val      # add tag constraint, multiple can be given
  -m i-12345678    # instance id
  # Capital options sort of denote Big Deals, changes the basic effect
  -A               # requires command, execute on all 
  -T n{-|/|+}[=]   # forbids command open tmux to n instances 
                   #   arranged -_r rows, /|c cols, +t tiles 
                   #   = setw synchronize-panes on
  -C {and|or}      # combine restrictions with AND / OR

ash base
  # TODO

asg alias NAME
  # TODO
