#!/usr/bin/env bash

CMD='\033[1;32m'
BOLD='\033[1m'
CX='\033[0m'
USER_AGENT=$(grep -Po "userAgent: \K[a-zA-Z0-9]+$" ~/.okctl/conf.yml 2>/dev/null || echo okctl)
METRICS_URL="https://metrics.kjoremiljo.oslo.systems/v1/metrics/events"

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
    echo "❗ Installation NOT complete. You need to remove the old version of okctl before continuing. Run the following command:"
    echo -e "${CMD}sudo rm /usr/local/bin/okctl${CX}"
    echo
    echo "and then re-run this installation."
    echo
    echo -e "Reason: This installation has detected that you have ${BOLD}~/.local/bin${CX} added to your PATH, so the okctl"\
      "binary will be put there instead. By doing this, future installations should be completely automatic, without the need"\
      "for running the ${CMD}sudo rm${CX} command above."
    exit 1
  fi

  fetch_binary $1

  mv /tmp/okctl $HOME/.local/bin
  echo "Successfully installed okctl. You can test it by running:"
  echo -e "${CMD}okctl version${CX}"
}

function install_usr() {
  fetch_binary $1

  echo "Done."
  echo
  echo "-------------------------------------------------------------------------------------------------------------------------"
  echo "To complete installation, run:"
  printf $CMD
  echo "sudo mv /tmp/okctl /usr/local/bin"
  printf $CX
  echo "-------------------------------------------------------------------------------------------------------------------------"
  echo
  echo -e "💡 Tip: To avoid having to run the ${CMD}sudo mv${CX} command above for future installations, you can do the following:"
  echo -e "- Create the directory ${BOLD}~/.local/bin${CX}"
  echo -e "- Add ${BOLD}~/.local/bin${CX} to your PATH."
  echo "- Re-run this installation. This installation will detect the directory, and put okctl there."
}

function publish_event() {
  ACTION=$1

  curl $METRICS_URL \
      -X POST \
      -H "User-Agent: $USER_AGENT" \
      -H "Content-Type: application/json" \
      -d "{\"category\": \"install\", \"action\": \"$ACTION\" }"
}

function publish_start_stop() {
  PHASE_KEY=$1

  curl $METRICS_URL \
    -X POST \
    -H "User-Agent: $USER_AGENT" \
    -H "Content-Type: application/json" \
    -d "{\"category\": \"install\", \"action\": \"okctl\", \"labels\": { \"phase\": \"$PHASE_KEY\" } }"
}

function publish_start() {
  PHASE_KEY="start"
  publish_start_stop $PHASE_KEY
}

function publish_stop() {
  PHASE_KEY="stop"
  publish_start_stop $PHASE_KEY
}

# publish_start

if [[ -z $1 ]]; then
  VERSION=latest
else
  VERSION=$1
fi

# Check if brew exists
if command -v brew &> /dev/null; then
  # Check if okctl exists
  if brew list okctl &> /dev/null; then
    # publish_event brew_uninstall

    echo Uninstalling okctl from brew
    brew uninstall okctl &> /tmp/okctl_brew_uninstall.txt
    brew untap oslokommune/tap &>> /tmp/okctl_brew_uninstall.txt
  fi
fi

if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]] && [[ -d $HOME/.local/bin ]]
then
  install_local_bin $VERSION
else
  install_usr $VERSION
fi

# publish_stop
