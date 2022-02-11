#!/usr/bin/env bash

# This script makes everything ready for release. It doesn't push any changes, but makes everything ready and outputs a command
# for you to push changes to the repository.
#
# USAGE:
# cd okctl
# ./release.sh
#
# For nicer git diff, install https://github.com/dandavison/delta

function get_next_version() {
    VERSION="$1"
    VERSION="${VERSION#[vV]}"
    VERSION_MAJOR="${VERSION%%\.*}"
    VERSION_MINOR="${VERSION#*.}"
    VERSION_MINOR="${VERSION_MINOR%.*}"
    VERSION_PATCH="${VERSION##*.}"

    (( NEXT_PATCH=VERSION_PATCH+1 ))

    echo "$VERSION_MAJOR.$VERSION_MINOR.$NEXT_PATCH"
}

if [[ $(pwd) != */okctl* ]]; then
    echo ERROR: You must be in an okctl repository
    exit 1
fi

if [[ -n $(git status --porcelain) ]]; then
    echo Dirty git status. Clean up before retrying.
    exit 1
fi

if ! command -v bat &> /dev/null
then
    echo 'bat' could not be found. Install before retrying.
    exit
fi

git checkout master
git pull --rebase

CURRENT_VERSION_RAW=$(git tag | sort -V | tail -1 | cut -c2-30)
RELEASE_VERSION_RAW=$(get_next_version "$CURRENT_VERSION_RAW")
RELEASE_VERSION="v$RELEASE_VERSION_RAW"
NEXT_VERSION_RAW=$(get_next_version "$RELEASE_VERSION_RAW")

echo
echo CURRENT_VERSION_RAW: $CURRENT_VERSION_RAW
echo RELEASE_VERSION_RAW: $RELEASE_VERSION_RAW
echo NEXT_VERSION_RAW: $NEXT_VERSION_RAW
echo

cat <<EOF > docs/release_notes/$NEXT_VERSION_RAW.md
# Release $NEXT_VERSION_RAW

## Features

## Bugfixes

## Changes

## Other

EOF

(bat docs/release_notes/$RELEASE_VERSION_RAW.md || cat docs/release_notes/$RELEASE_VERSION_RAW.md)

git add docs/release_notes/$NEXT_VERSION_RAW.md
git diff --cached
git commit -m "Add changelog file for $NEXT_VERSION_RAW"

git tag $RELEASE_VERSION
if [[ $? != 0 ]]; then
    echo "Git tag failed"
    exit 1
fi

echo
echo ---------------------------------------------------------------------
echo "* Verify that docs/release_notes/$RELEASE_VERSION_RAW.md looks ok"
echo "* Verify that git log looks OK:"
echo
git log --graph --pretty=format:'%Cred%h%Creset |%Cblue%>(12)%cr%Creset | %s%C(auto)%d%Creset %C(dim italic)%an%Creset' --abbrev-commit -2
echo
echo ""
echo "Running dry run of git push"
echo "git push --dry-run --atomic origin master $RELEASE_VERSION"
git push --dry-run --atomic origin master $RELEASE_VERSION
echo
echo "To release, run"
echo "git push --atomic origin master $RELEASE_VERSION"
echo
echo "Or abort with:"
echo "git tag -d $RELEASE_VERSION && git reset --hard origin/master"
echo ---------------------------------------------------------------------
