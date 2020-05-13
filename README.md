# `okctl` - Opinionated and effortless infrastructure and application management

We will be building up the functionality of this CLI within the coming weeks, and the aim is to release frequently. At this stage we are mostly setting the CI/CD and release processes.

![okctl](logo/okctl-fit.png)

## Installation

To download the latest release, run the command matching your operating system:

```bash
# Linux
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin

# macOS
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```

## Usage

```bash
okctl --help
```
