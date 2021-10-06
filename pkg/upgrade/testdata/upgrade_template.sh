#!/usr/bin/env bash

VER=#REPLACE_VER#
OS=#REPLACE_OS#
ARCH=amd64

has_invalid_param() {
    declare -a valid_args
    valid_args=(--debug --confirm --dry-run --dry-run=true --dry-run=false)
    for arg
        do

        # Thanks to https://stackoverflow.com/questions/3685970/check-if-a-bash-array-contains-a-value
        # shellcheck disable=SC2076
        if [[ ! " ${valid_args[*]} " =~ " ${arg} " ]]; then
            echo "Invalid argument: $arg. Valid args: ${valid_args[*]}"
            return 0
        fi
    done

    return 1
}

# Thanks to https://stackoverflow.com/a/56431189/915441
has_param() {
    local term="$1"
    shift
    for arg; do
        if [[ $arg == "$term" ]]; then
            return 0
        fi
    done
    return 1
}

if has_invalid_param "$@"; then
    exit 1
fi

# This is a test upgrade. We output this line so we can verify that this specific upgrade was run.
echo "This is upgrade file for okctl-upgrade_${VER}_${OS}_${ARCH}"

if has_param '--debug' "$@"; then
    echo "--debug flag was provided, so here is some debug output.."
fi

if has_param '--confirm' "$@"; then
    echo "--confirm flag was provided"
fi

if has_param '--dry-run=true' "$@" || has_param '--dry-run' "$@"; then
  echo "--dry-run true flag was provided, so simulating changes."
else
  echo "Doing actual changes."
fi
