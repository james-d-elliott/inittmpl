# inittmpl

Stupid simple init templating tool.

## Installation

`go install github.com/james-d-elliott/inittmpl@447fc491b612393a3595e072cc41fd233c896ef4`

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

|  Type  | Opportunistic |             Example             | Output (YAML) |
|:------:|:-------------:|:-------------------------------:|:-------------:|
| string |      Yes      | `INITTMPL__EXAMPLE=string::123` |    `'123'`    |
| string |      No       |     `INITTMPL__EXAMPLE=123`     |    `'123'`    |
| string |      No       |      `INITTMPL__EXAMPLE=1`      |     `'1'`     |
| string |      No       |    `INITTMPL__EXAMPLE=true`     |   `'true'`    |
|  int   |      Yes      |  `INITTMPL__EXAMPLE=int::123`   |     `123`     |
|  int   |      Yes      |     `INITTMPL__EXAMPLE=123`     |     `123`     |
|  bool  |      Yes      |   `INITTMPL__EXAMPLE=bool::1`   |    `true`     |
|  bool  |      Yes      | `INITTMPL__EXAMPLE=bool::false` |    `false`    |
|  bool  |      Yes      |    `INITTMPL__EXAMPLE=false`    |    `false`    |
| float  |      Yes      | `INITTMPL__EXAMPLE=float::123`  |    `123.0`    |
