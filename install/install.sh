#!/usr/bin/env bash

C='\033[0;32m'
CX='\033[0m'

function fetch_binary() {
  VERSION=$1
  if [[ $VERSION == "latest" ]]; then
    echo "Downloading latest version of okctl..."
    curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp okctl
  else
    echo "Downloading version '$VERSION' of okctl..."
    curl --silent --location "https://github.com/oslokommune/okctl/releases/download/v$VERSION/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp okctl
  fi
}

function install_local_bin() {
  if [[ -f /usr/local/bin/okctl ]]; then
    echo "You need to remove the old version of okctl before continuing. Run the following command"
    echo -e "${C}sudo rm /usr/local/bin/okctl${CX}"
    echo
    echo "and then re-run this installation."
    exit 1
  fi

  fetch_binary $1

  mv /tmp/okctl $HOME/.local/bin
  echo "Successfully installed okctl. Test it by running:"
  echo -e "${C}okctl version${CX}"
}

function install_usr() {
  fetch_binary $1

  echo "okctl downloaded. To complete installation, run"
  printf $C
  echo sudo mv /tmp/okctl /usr/local/bin
  printf $CX
  echo
  echo "Then test it by running"
  echo -e "${C}okctl version${CX}"
  echo
  echo -e "💡 Tip: To avoid having to run the ${C}sudo mv${CX} command above in the future, create the directory"\
    "\033[1m~/.local/bin${CX} and add it to your PATH."
}

if [[ -z $1 ]]; then
  VERSION=latest
else
  VERSION=$1
fi

# Check if brew exists
if command -v brew &> /dev/null; then
  # Check if okctl exists
  if brew list okctl &> /dev/null; then
    echo Uninstalling okctl from brew
    brew uninstall okctl
    brew untap oslokommune/tap
    echo
  fi
fi

if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]] && [[ -d $HOME/.local/bin ]]
then
  install_local_bin $VERSION
else
  install_usr $VERSION
fi

