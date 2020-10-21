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

which will run a sub shell with the environment variables from `okctl show credentials` already set.
