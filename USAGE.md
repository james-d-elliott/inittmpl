## kissit

Generate a configuration file from Keep It Simple Stupid Markup Language / KISSML (environment variables or random files you have laying around)

### Synopsis

Generate a configuration file from Keep It Simple Stupid Markup Language / KISSML (environment variables or random files you have laying around).

Environment variables take the highest precedence overwriting any other input.

By default this command will not overwrite the output file unless the one of
the input files is also the same absolute path as it, or if you use a flag
to instruct the overwriting should take place.

Some types automatically and opportunistically are cast before marshalling the
output file; however you can specify the type in the environment variables using 
the 'string::',' int::', 'uint::', 'bool::', and 'float::' prefixes before the value.

In addition to the standard types a 'json::' prefix treats the following text as
JSON and performs a json unmarshal step allowing complex types including lists
of objects. 

The decoder and encoder formats are automatically detected via the file extension
though you can specify the format to use.

Because environment keys are often expressed entirely uppercase and configuration
keys are usually expressed entirely lowercase we automatically perform lowercase
correction. This behaviour can be turned off however all of your key names must
match the target case after the prefix.

Input files can be utilized to effectively modify existing configs however both
the key order and the comments may be sacrificed. Alternatively you may just opt
not to overwrite and use this solely for generating first time configs.

Levels within the data structure are separated by default by a double underscore.
This allows easily defining multi-level datastructures without the need for
complex decoding.

```
kissit <file> [flags]
```

### Examples

```
  kissit file.yaml
  kissit file -f yaml
  kissit file.toml
  kissit file.json
  kissit file.yaml -z input1.yaml,input2.toml,input3.json
  kissit file.yaml -z file.yaml,input2.toml,input3.json
  kissit file.yaml -x
```

### Options

```
  -d, --delimiter string        override the environment delimiter (default "__")
  -e, --disable-lowercase       disable lowercase conversion of environment variables
  -c, --disable-opportunistic   disable opportunistic type conversion of environment variables
  -z, --files strings           uses the specified files as additional input, if any of these are the same file as the output assumes overwrite
  -f, --format string           override the output format; yaml, toml, json
  -h, --help                    help for kissit
  -x, --overwrite               overwrite existing file
  -p, --prefix string           override the environment prefix (default "kissit")
```

