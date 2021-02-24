To download the latest release, make sure you have the required pre-requisites. Then run the command matching your
operating system.

## Prerequisites

## Linux

```bash
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin
```

## macOS

### Prerequisites

### Installation

```bash
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```

## Shell autocompletion

To get autocompletion in your shell working run `okctl completion -h` and follow the steps for your preferred shell.
