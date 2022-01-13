# Working on nightly

## How to use the dummify patches

### Workflow

First, create a branch for your task.

Then apply the patches using the following commands:

```shell
git am < .github/workflows/nightly/0001-DROP-dummify-nightly.patch
git am < .github/workflows/nightly/0002-DROP-dummyfying-apply-cluster-and-delete-cluster.patch
```

Write some nightly based code.

Clean up by rebasing on master, and then dropping the commits prefixed with DROP.

### Dummify nightly

On your branch, in the repository root directory, run the following:

`git am < .github/workflows/nightly/0001-DROP-dummify-nightly.patch`

This will add a commit to your branch prefixed with "DROP" that will:

1. Change the nightly trigger to `on.push` instead of `on.schedule` which will make the action trigger when you push
   commits to `origin/<branch>`
2. Remove the test step.
3. Add a cache step to the nightly action so it doesn't have to build the binary every run. N.B.: caching will only
   start after the first successful nightly run.
4. Remove the Slack notify step, preventing #kjøremiljø to be spammed when developing.

### Dummify apply and delete cluster

On your branch, in the repository root directory, run the following:

`git am < .github/workflows/nightly/0002-DROP-dummyfying-apply-cluster-and-delete-cluster.patch`

This will add a commit to your branch prefixed with "DROP" that will:

1. Remove all reconcilers, command metrics and state uploads/downloads from `apply cluster`
2. Remove all reconcilers, command metrics and state uploads/downloads from `delete cluster`
