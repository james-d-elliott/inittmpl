## inittmpl

Generate a configuration file from environment variables

### Synopsis

Generate a configuration file from environment variables.

```
inittmpl <file> [flags]
```

### Options

```
  -d, --delimiter string        override the environment delimiter (default "__")
  -e, --disable-lowercase       disable lowercase conversion of environment variables
  -c, --disable-opportunistic   disable opportunistic type conversion of environment variables
  -f, --format string           override the output format; yaml, toml, json
  -h, --help                    help for inittmpl
  -x, --overwrite               overwrite existing file
  -p, --prefix string           override the environment prefix (default "INITTMPL")
```