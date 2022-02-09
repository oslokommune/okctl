How to create a new binary provider
=====
Okctl stores binaries under `~/.okctl/binaries` that are available after running `okctl venv`.
The `pkg/binaries` package contains code to configure, fetch, validate and set up the different binaries we choose to pre-load as a part of okctl.

# Files to add
* pkg/binaries/run/{new-binary}/{new-binary}.go

Create a new folder called `{new-binary}` and a file inside called `{new-binary}.go`. 
The package `pkg/binaries/run/awsiamauthenticator` or `pkg/binaries/run/kubens` are minimum-implementations of the required setup for a new binary, start with the code in one of these and modify accordingly (changing the name from `kubens` to `{new-binary}` for the package and variables). 

# Files to change:
For all these files: follow the pattern for on of the provided binaries (ex: `awsiamauthenticator` or `kubens`)
* pkg/binaries/binaries.go
    * Updated/add: `import`, `Provider interface`, `provider struct`, the `New()` function and add a separate `p.{new-binary}` function to set up the binary
* pkg/commands/venv.go
    * Update: `import`, `OkctlEnvironment struct`, `o.Validate()`, `GetOkctlEnvironment()` and `getPathWithOkctlBinaries()`
* pkg/config/state/binary.go
    * Add a new `{new-binary}KnownBinaries` function and add to the `KnownBinaries()` `binaries` list 
  * Each binary will have a different release-strategy, there is no hard rule on how this is done

Known issues
=====
* The `Host.Arch` in `pkg/config/state/user.go` is set to `amd64` for everyone
  * Following the pattern in `pkg/config/state/binary.go` there is an issue when the tool doesn't release with `amd64` in the file name
  * See workaround in `KubensKnownBinaries` where we hardcode to `x86_64` in the release path and don't use the `Host.Arch` variable to determine download URL