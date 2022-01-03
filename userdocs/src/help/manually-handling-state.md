Okctl uses [BoltDB](https://dbdb.io/db/boltdb) to manage necessary state. When running a command in okctl, we download
the state database file from an S3 bucket, make the required changes, then upload it again. This ensures synchronization
between multiple users. To prevent users from operating the state at the same time, we use DynamoDB as a locking
mechanism.

If something were to go wrong, the hidden command `okctl maintenance` contains subcommands useful when manually handling
state.

⚠ Be careful! Manually modifying state can render your okctl environment useless.

## Useful commands

```shell
okctl maintenance state-acquire-lock
```

will acquire a lock preventing other users from modifying state.

```shell
okctl maintenance state-release-lock
```

will release a lock allowing other users to modify the state.

```shell
okctl maintenance state-download
```

will download the state database file to a local location.

```shell
okctl maintenance state-upload <path>
```

will upload the state database file to S3 allowing local modifications to propagate.

## Useful tools

We use [BoltBrowser](https://github.com/br0xen/boltbrowser) to view and edit a BoltDB database.

## Lock implementation

Before each command that potentially modifies the state, `okctl` will look for a DynamoDB table called
`okctl-<cluster name>-state-lock` and verify that it doesn't contain any items.

If an item exists (meaning someone else is running a command that can potentially modify the state), `okctl` will stop
the current operation ensuring only the first user can run their command.

If no items exist, `okctl` will insert a string which it will remove after it's done running the requested command.