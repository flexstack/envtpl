# envtpl

A CLI tool for generating .env files from .env.template files via prompts and random values.

## Installation

Using `curl`:
  
 ```sh
 curl -sfL https://raw.githubusercontent.com/flexstack/envtpl/main/install.sh | sh
 ```

Using `go install`:

```sh
go install github.com/flexstack/envtpl
```

## Usage

```sh
envtpl [template-file] -o [output-file]
```

### Supported types

```sh
# .env.template
UUID=<uuid>
PASSWORD=<alpha:24>
ENCRYPTION_KEY=<ascii85:32>
ENCRYPTION_IV=<base64:16>
ENCRYPTION_SALT=<hex:16>
PORT=<int:1-65535>
ENUM=<enum:development,production>
PROMPT=<text:Enter your name>
SECRET_PROMPT=<password:Enter your password>
PROMPT_2= # Empty values elicit a prompt
```

| Type | Description | Example |
| --- | --- | --- |
| alpha | Random alphanumeric characters | `PASSWORD=<alpha:24>` |
| ascii85 | Random ASCII characters | `ENCRYPTION_KEY=<ascii85:32>` |
| base64 | Random base64 characters | `ENCRYPTION_IV=<base64:16>` |
| hex | Random hexadecimal characters | `ENCRYPTION_SALT=<hex:16>` |
| uuid | Random UUIDv4 | `UUID=<uuid>` |
| enum | Enumerated list of values | `ENUM=<enum:development,production>` |
| int | Random integer | `PORT=<int:1-65535>` |
| text | Prompt user for text input | `PROMPT1=<text:Enter your name>` |
| password | Prompt user for secret input | `SECRET_PROMPT=<password:Enter your API key>` |
