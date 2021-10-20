# Contributing

This project depends on contributions from its users to succeed, thank you so much for considering becoming a part of our small community ❤️

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change. 

Please note we have a [code of conduct](CODE_OF_CONDUCT.md), please follow it in all your interactions with the project. Also, by contributing to this project you agree to the [Developer Certificate of Origin (DCO)](DCO), which simply states that you have the legal right to make the contribution.

Parts of this contributing guide has been borrowed from the excellent [eksctl](https://github.com/weaveworks/eksctl/blob/master/CONTRIBUTING.md) project.

## Geting started

- Fork the repository on GitHub
- Read the [README](README.md) for getting started as a user and learn how/where to ask for help
- If you want to contribute as a developer, continue reading this document for further instructions
- Play with the project, submit bugs, submit pull requests!

### Contribution workflow

#### 1. Set up your Go environment

This project is written in Go. To be able to contribute you will need:

1. A working Go installation of Go >= 1.14. You can check the
[official installation guide](https://golang.org/doc/install).

2. Make sure that `$(go env GOPATH)/bin` is in your shell's `PATH`. You can do so by
   running `export PATH="$(go env GOPATH)/bin:$PATH"`

#### 2. Fork and clone the repo

Make a fork of this repository and clone it by running:

```bash
git clone git@github.com:<yourusername>/okctl.git
```

It is not recommended to clone under your `GOPATH` (if you define one). Otherwise, you will need to set
`GO111MODULE=on` explicitly.

#### 3. Run the tests and build okctl

Make sure you can run the tests and build the binary.

```bash
make test
make release-local
```

After running the initialization process that will automatically start when running a command, please edit your
`~/.okctl/conf.yml` and add the following:

```yaml
...
metrics:
  userAgent: okctldev
...
```

This will allow us to differentiate between actual production usage and development usage. Thank you!

#### 4. Write your feature

- Find an [issue](https://github.com/oslokommune/okctl/issues) to work on or create your own. If you are a new
  contributor take a look at issues marked
  with [good first issue](https://github.com/oslokommune/okctl/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
  .

- Then create a topic branch from where you want to base your work (usually branched from master):

    ```bash
    git checkout -b <feature-name>
    ```

- Write your feature. Make commits of logical units and make sure your
  commit messages are in the [proper format](#format-of-the-commit-message).

- Add automated tests to cover your changes.

- If needed, update the documentation, either in the [README](README.md) or in the [docs](docs/) folder.

- Make sure the tests are running successfully.

#### 5. Submit a pull request

Push your changes to your fork and submit a pull request to the original repository. If your PR is a work in progress
then make sure you prefix the title with `WIP: `. This lets everyone know that this is still being worked on. Once its
ready remove the `WIP: ` title prefix and where possible squash your commits. Alternatively, you can use `Draft PR`
feature of Github as mentioned [here](https://github.blog/2019-02-14-introducing-draft-pull-requests/)

```bash
git push <username> <feature-name>
```

Our github actions will run the automated tests and give you feedback in the review section. We will review your
changes and give you feedback as soon as possible.

## Acceptance policy

These things will make a PR more likely to be accepted:

- a well-described requirement
- tests for new code
- tests for old code!
- new code and tests follow the conventions in old code and tests
- a good commit message (see below)

In general, we will merge a PR once a maintainer has reviewed and approved it.
Trivial changes (e.g., corrections to spelling) may get waved through.
For substantial changes, more people may become involved, and you might get asked to resubmit the PR or divide the
changes into more than one PR.

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```text
Added AWS Profile Support

Changes to ensure that AWS profiles are supported. This involved making
sure that the AWS config file is loaded (SharedConfigEnabled) and
also making sure we have a TokenProvider set.

Added an explicit --profile flag that can be used to explicitly specify
which AWS profile you would like to use. This will override any profile
that you have specified via AWS_PROFILE.

If endpoints are being overridden then the credentials from the initial
session creation are shared with any subsequent session creation to
ensure that the tokens are shared (otherwise you may get multiple MFA
prompts).

Issue #57
```

Or, in case of shorter messages:

```text
#57: Added AWS Profile Support
```

The format can be described more formally as follows:

```text
<short title for what changed>
<BLANK LINE>
<why this change was made and what changed>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various git tools.