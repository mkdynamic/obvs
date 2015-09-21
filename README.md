# Obvs

CLI for AWS CloudWatch Logs

## Usage

```bash
$ obvs --help
usage: obvs --access-key-id=ACCESS-KEY-ID --secret-access-key=SECRET-ACCESS-KEY [<flags>] <command> [<args> ...]

CLI for AWS CloudWatch Logs.

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --region="us-east-1"
             AWS region.
  --access-key-id=ACCESS-KEY-ID
             AWS access key ID.
  --secret-access-key=SECRET-ACCESS-KEY
             AWS secret access key.
  --version  Show application version.

Commands:
  help [<command>...]
  groups
  streams <group>
  events* [<flags>] <group> [<pattern>]
```
