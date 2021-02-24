# Venv

When running `okctl show credentials`, you'll get a bunch of environment variables that need to be set for you to
access the cluster. To be able to set these quickly in a new terminal, you could do something like

```bash
okctl show credentials myenv > ~/.okctl/binaries/venv-myenv
echo "alias activate-myenv='source ~/.okctl/binaries/venv-myenv' >> ~/.bash_aliases"
```

However, a simpler solution is to just run

```bash
okctl venv myenv
```

which will
* run a sub shell with the environment variables from `okctl show credentials` already set.
* set a command prompt like this (example):

```bash
~/somepath (myenv:mynamespace) $
```

where

* `myenv` will be replaced with the name of your okctl environment
* `mynamespace` will be replaced with the name of the Kubernetes namespace of your current kube context

This command prompt can be turned off or configured, see below.

## Configuration

| Environment variable      | Default   | Description                                                            |
| ------------------        | --------- | ---------------------------------------------------------------------- |
| OKCTL_NO_PS1              | true      | If true, disables overriding the PS1. This will keep the existing PS1. |
| OKCTL_PS1                 | *not set* | If set, `venv` will use this as the PS1 in the executed subshell.      |
| OKCTL_SHELL               | true      | Override which shell to run. For instance `/bin/sh`.                   |

Any occurence `%env` in `OKCTL_PS1` will be replaced by the okctl environment. This makes it possible to get the okctl
environment in your custom OKCTL_PS1. A use case for this can be when combining with the `venv_ps1` built-in:

```bash
export OKCTL_PS1="\w \$(venv_ps1 %env) $"
okctl venv myenv
```

The command prompt will then be like this

```bash
/tmp myenv:mynamespace $
```
