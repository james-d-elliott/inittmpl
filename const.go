package main

const (
	rootLong = `Generate a configuration file from environment variables.

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
complex decoding.`

	rootExample = `  ssit file.yaml
  ssit file -f yaml
  ssit file.toml
  ssit file.json
  ssit file.yaml -z input1.yaml,input2.toml,input3.json
  ssit file.yaml -z file.yaml,input2.toml,input3.json
  ssit file.yaml -x`
)
