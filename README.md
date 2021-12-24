# gh-app-auth

Authenticate with the GitHub App and generate the installation token.

## Installation

```
$ go install github.com/summerwind/gh-app-auth@latest
```

## Usage

If you specify the GitHub App ID, Private Key, and the target account name, This command will generate and display the installation token.

```
$ gh-app-auth --app-id 12345 --private-key /path/to/private-key.pem --account summerwind
ghs_01234...
```

The flags are used as follows.

```
Usage of gh-app-auth:
  -a, --account string       Target account name
  -i, --app-id int           GitHub App ID
  -h, --help                 Print usage and exit
  -f, --private-key string   Path to the private key file of GitHub App
```

## License

MIT license.
