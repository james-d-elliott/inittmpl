# inittmpl

Stupid simple init templating tool.

## Limitations

The only current way to define lists via environment variables is via the `json::` value prefix, and this is the only
current way to define a list of any type. 

## TODO

- dependency management:
  - [ ] add renovate
  - [ ] add dependabot
- continuous integration:
  - [ ] add tags / releases
  - [ ] add tests and test steps
  - [ ] add build steps
  - [ ] add deploy steps
- [ ] add ghcr.io container 
- [ ] add step security

## Installation

Like everything else with this tool, installing is stupid simple:

```shell
go install github.com/james-d-elliott/inittmpl@2cb3e92997477c991323365c7ab007b2b3dc7daf
```

## Building

Like everything else with this tool, building is stupid simple:

```shell
git clone https://github.com/james-d-elliott/inittmpl.git
cd inittmpl
go mod download
go build
```

## Docker

Use any of the images detailed below. The command used is the binary, just add the arguments and
environment variables for your desired output.

| Repository  |               Image                |
|:-----------:|:----------------------------------:|
| `docker.io` | `docker.io/jamesdelliott/inittmpl` |
|  `ghcr.io`  |                N/A                 |

### Kubernetes Init Containers

The following is an initContainers example for Kubernetes assuming a volume is mounted at /config:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: 'example'
spec:
  containers:
    - name: 'example'
      volumeMounts:
        - mountPath: '/config'
          name: 'config-vol'
  initContainers:
    - name: 'inittmpl'
      image: 'jamesdelliott/inittmpl:latest'
      args: ['/config/config.yaml', '-xec']
      volumeMounts:
        - mountPath: '/config'
          name: 'config-vol'
      env:
        - name: 'INITTMPL__example_integer'
          value: 'int::123'
        - name: 'INITTMPL__example_string'
          value: 'string::123'
        - name: 'INITTMPL__example_boolean'
          value: 'bool::true'
        - name: 'INITTMPL__example__multilevel_string'
          value: 'string::true'
        - name: 'INITTMPL__example__multilevel_list'
          value: 'json::["abc","123"]'
        - name: 'INITTMPL__example__multilevel_object'
          value: 'json::{"abc":123,"xyz":456,"boolean":true,"string":"value"}'
  volumes:
    - name: 'config-vol'
      persistentVolumeClaim:
        claimName: 'config-pvc'
```

Output:

```yaml
example:
  multilevel_list:
    - abc
    - '123'
  multilevel_object:
    abc: 123.0
    boolean: true
    string: value
    xyz: 456.0
  multilevel_string: 'true'
example_boolean: true
example_integer: 123
example_string: '123'
```

## Behaviour

Because environment variables have certain characteristics and there may be some nuanced use cases some behaviour flags 
and defaults are necessary. See the full [usage documentation](USAGE.md) for usage information.

### Format

This tool automatically determines the appropriate output format using the file path. This however doesn't fit all use
cases and this can be changed using the flag `--format` or `-f` for short.

Formats:

|   Extensions    | Name | Value  |
|:---------------:|:----:|:------:|
| `.yaml`, `.yml` | YAML | `yaml` |
|     `.toml`     | TOML | `toml` |
|     `.json`     | JSON | `json` |

### Overwrite

By default this tool does not overwrite files so that this can be used as an init-only tool. You can however enable 
overwrite behaviour casing by using
`--overwrite` or `-x` for short.

### Lowercase

By default this tool converts all uppercase keys to lowercase. This is because env vars are often expressed all 
uppercase but it's rare that configuration files have uppercase keys. You can however enable explicit casing by using
`--disable-lowercase` or `-e` for short.

### Type Conversion

By default this tool opportunistically performs type conversions. You can disable this behaviour using 
`--disable-opportunistic` or `-o` for short.

You can define explicit type conversions regardless of this setting using the special prefixes `string::`, `int::`, 
`bool::`, and `float::`. Examples of various behaviour are illustrated in the table below.

|  Type  | Opportunistic |                    Example                    |     Output (YAML)     |
|:------:|:-------------:|:---------------------------------------------:|:---------------------:|
| string |      Yes      |        `INITTMPL__EXAMPLE=string::123`        |        `'123'`        |
| string |      No       |            `INITTMPL__EXAMPLE=123`            |        `'123'`        |
| string |      No       |             `INITTMPL__EXAMPLE=1`             |         `'1'`         |
| string |      No       |           `INITTMPL__EXAMPLE=true`            |       `'true'`        |
|  int   |      Yes      |         `INITTMPL__EXAMPLE=int::123`          |         `123`         |
|  int   |      Yes      |            `INITTMPL__EXAMPLE=123`            |         `123`         |
|  bool  |      Yes      |          `INITTMPL__EXAMPLE=bool::1`          |        `true`         |
|  bool  |      Yes      |        `INITTMPL__EXAMPLE=bool::false`        |        `false`        |
|  bool  |      Yes      |           `INITTMPL__EXAMPLE=false`           |        `false`        |
| float  |      Yes      |        `INITTMPL__EXAMPLE=float::123`         |        `123.0`        |
|  json  |      No       |        `INITTMPL__EXAMPLE=json::[123]`        |       `[123.0]`       |
|  json  |      No       | `INITTMPL__EXAMPLE=json::["123","456","abc"]` | `['123','456','abc']` |
