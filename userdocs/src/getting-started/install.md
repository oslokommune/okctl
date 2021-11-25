## Automatic installation

To automatically install the latest version of okctl, run

```shell
curl -fsSL https://raw.githubusercontent.com/oslokommune/okctl/HEAD/install/install.sh | /bin/bash
sudo mv /tmp/okctl /usr/local/bin
```

To get a specific version of okctl, run (replace `0.0.73` with the version you want)

```shell
curl -fsSL https://raw.githubusercontent.com/oslokommune/okctl/HEAD/install/install.sh | /bin/bash -s 0.0.73
sudo mv /tmp/okctl /usr/local/bin
```

## Manual installation

### macOS

If you previously installed okctl with `brew`, uninstall it first, by running

```bash
brew uninstall oslokommune/tap/okctl
brew untap oslokommune/tap
```

Then install with the instructions below.

### macOS and Linux

```bash
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin # or somewhere else on your $PATH
```

To get a specific version, see the release assets on the [releases](https://github.com/oslokommune/okctl/releases) page.

## Enable shell autocompletion

To get autocompletion in your shell working run `okctl completion -h` and follow the steps for your preferred shell.
