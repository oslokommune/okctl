## Installation

To download the latest release, make sure you have the required pre-requisites. Then run the command matching your
operating system.

### Prerequisites

* Your own GPG key
    * List existing keys with `gpg --list-keys`
    * If you don't have any keys, create a new one with `gpg --gen-key`

### Linux

```bash
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin
```

### macOS

#### Prerequisites

* `pass`
    * If you do not have it, run `brew install pass`
* `pass` must be initalized
    * Run `pass ls` to see if you have a password store. If it returns `Error: password store is empty. Try "pass init".`
      then run `pass init <GPG ID>` where GPG ID is the ID of your GPG key - use `gpg --list-keys` to see your keys (both
      e-mail and the ID are valid values, but if you have several keys for the same e-mail, use the ID).

#### Installation

```bash
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```