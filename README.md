# SafeIP

SafeIP is a command-line tool written in Go that helps you mask public IPv4 addresses and DNS-like entries from your text input. It's useful for redacting sensitive information from logs or other textual data to avoid people eyeballing your data 🕵️

## Demo

![SafeIP Demo](output.gif)

## Build

### Installation

To use SafeIP, you'll need to compile it from source. Make sure you have Go installed on your system.

```bash
git clone https://github.com/darklight147/safeip.git
cd safeip

go mod download

go build -o safeip

mv safeip /usr/local/bin

# or

mv safeip /usr/bin

safeip --help

```

## Get a pre-built binary

You can download a pre-built binary for your system from the [releases page](https://github.com/darklight147/safeip/releases).

### Usage

**Basic Usage:**

This will read from `input.txt`, mask public IPv4 addresses, and print the modified text to the console.

```bash
cat input.txt | safeip
```

Masking public IPv4 addresses is the default behavior of SafeIP. The command above will read from `input.txt`, mask public IPv4 addresses, and print the modified text to the console.

```bash
kubectl describe nodes | safeip
```

Masking IPs from logs

```bash
kubectl logs pod-name | safeip
```

**Masking DNS-like Entries:**

```bash
cat input.txt | safeip --mask-dns
```

When the `--mask-dns` flag is set to `true`, SafeIP will also mask DNS-like entries, providing an extra layer of privacy.

### Example

**Before Using SafeIP:**

```
Addresses:
  InternalIP:   172.31.16.128
  ExternalIP:   54.87.142.142
  InternalDNS:  ip-172-31-16-128.ec2.internal
  Hostname:     ip-172-31-16-128.ec2.internal
  ExternalDNS:  ec2-54-87-142-142.compute-1.amazonaws.com
```

**After Using SafeIP:**

```
Addresses:
  InternalIP:   172.31.16.128
  ExternalIP:   XXX.XXX.XXX.XXX
  InternalDNS:  ip-172-31-16-128.ec2.internal
  Hostname:     ip-172-31-16-128.ec2.internal
  ExternalDNS:  XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

In the example above, SafeIP masked the external IP address and DNS-like entry, providing an added layer of security.

### Flags

| Flag          | Description                   | Default                                          |
| ------------- | ----------------------------- | ------------------------------------------------ |
| `--mask-dns`  | Mask DNS-like entries         | `false`                                          |
| `--mask`      | Custom String to replace with | `"XXXXXX"`                                       |
| `--dns-regex` | Custom String to replace with | `"(\\b(?:[a-zA-Z0-9-]+\\.){2,}[a-zA-Z]{2,}\\b)"` |
| `--help`      | Help for SafeIP               |                                                  |

## Completion

### Bash

```bash
source <(safeip completion bash)
```

### Zsh

```bash
source <(safeip completion zsh)
```

### Fish

```bash
safeip completion fish | source
```

## License

This project is licensed under the [MIT License](LICENSE).
