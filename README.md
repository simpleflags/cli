## CLI commands

```shell
add flag test1 -n="Test flag" -o=flags -p=default --on-value=true --off-value=false --rule="target.identifier == 'enver'":true --rule="target.identifier == 'hamza'":false
```

remove resource from simple flags:
```shell
remove flag test1 -o=flags -p=default
```