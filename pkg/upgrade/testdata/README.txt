Example of how to create an archive to be used for testing:
- Create a runnable bash script
- Run: tar czf okctl-upgrade_0.0.61_Linux_amd64.tar.gz okctl-upgrade_0.0.61
- Run: `sha256sum okctl-upgrade_0.0.61_Linux_amd64.tar.gz okctl-upgrade_0.0.61`
- Put the resulting sha256 hash into okctl-upgrade-checksum.txt
