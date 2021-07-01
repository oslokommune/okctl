#!/usr/bin/env bash

if [[ $* == "-h" || -z "$1" ]]
then
    ME=$(basename $0)
    echo "Creates an okctl upgrade test case to be used in tests."
    echo ""
    echo "USAGE:"
    echo "./$ME <okctl upgrade version>"
    echo ""
    echo "EXAMPLE:"
    echo "./$ME 0.0.45"

    return 0 2> /dev/null || exit 0
fi

VER=$1
ARCH=amd64

mkdir "$VER"
cd "$VER" || exit 1

for OS in {Linux,Darwin} ; do
  UPGRADE_FILE=okctl-upgrade_${VER}
  VERIFICATITON_FILE=okctl-upgrade_${VER}_${OS}_${ARCH}_ran_successfully

  cat <<EOF > "$UPGRADE_FILE"
#!/usr/bin/env sh
# This is a test upgrade. We create a file so we can verify that this upgrade was run.
echo This is upgrade file for okctl-upgrade_${VER}_${OS}_${ARCH}
touch ${VERIFICATITON_FILE}
EOF

  ARCHIVE_FILE=okctl-upgrade_${VER}_${OS}_${ARCH}.tar.gz
  tar czf "$ARCHIVE_FILE" "$UPGRADE_FILE"
done

DIGEST_FILE=okctl-upgrade-checksums.txt
sha256sum *.tar.gz > "$DIGEST_FILE"


rm $UPGRADE_FILE

