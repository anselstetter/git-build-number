# git-build-number
![Coverage](https://img.shields.io/badge/Coverage-86.8%25-brightgreen)

[Installation](#installation) | [Docs](#docs) | [Notes](#notes)

`git-build-number` lets you manage build numbers within a Git repository. 

It provides consistent build numbering across CI/CD systems, enables environment-specific[^1] build numbers, and replaces `${{ github.run_number }}` or Jenkins `${BUILD_NUMBER}`.

```
Manage build numbers within a Git repository

Usage:
  git-build-number [flags]
  git-build-number [command]

Available Commands:
  fetch         Fetch build number(s)
  get           Get the latest build number
  hash          Show the hash for a specific build number
  help          Help about any command
  inc           Increment the build number
  namespace     Manage namespaces
  push          Push build number(s)
  set           Set the build number
  version       Print the version

Flags:
  -h, --help   help for git-build-number

Use "git-build-number [command] --help" for more information about a command.
```

## How it works:

`git-build-number` uses Git’s object model to store and track build numbers.

1. A file describing the build number and its corresponding commit hash is generated:

```
{build-number} {hash}
```

2. This file is saved as a Git blob and added to a tree under the filename `build-number`.
3. A new commit referencing this tree is created.  
The commit message records the build number and associated commit:

```
key: {build-number}
value: {hash}
```

This structure makes it possible to look up build numbers and the commits they refer to by traversing the history of these build-number commits, without needing to inspect the stored blob contents directly.

Each build-number sequence is stored under a dedicated ref:

```
refs/build-number/{namespace}
```

If no namespace is specified, the default reference is used:

```
refs/build-number/default
```

For background on how blobs, trees, and commits work in Git, see [Git Internals — Git Objects](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)


## Setup

Here is an example GitHub Actions workflow step to increment and use a build number:

```yaml
steps:
  - name: Checkout
    uses: actions/checkout@v4

  - name: Incrememt build number
    run: |
      echo "BUILD_NUMBER=$(git build-number inc)" >> ${GITHUB_ENV}
      git build-number push

  - name: Build stuff
    run: |
      build app --version ${BUILD_NUMBER}
```

After running this step, the new build number will be available in `BUILD_NUMBER` environment variable.

## Installation

### Go

```
go install github.com/anselstetter/git-build-number/cmd/git-build-number@latest
```

### Manually

Download the latest binaries from [here](https://github.com/anselstetter/git-build-number/releases) and copy them to your desired location.

### Other

There are no releases to other package repositories yet.

## Docs:

Info about all available commands are [here](./docs/git-build-number.md).

## Notes:

### macOS:

The releases for macOS are not signed, so macOS will deny running the binary.

To get around this issue you can delete the quarantine extended attribute with:

```
xattr -d com.apple.quarantine <binary>
```

## Windows:

This binary is cross compiled for Windows `AMD64` and `ARM64`, but has never been tested.

Your mileage may vary on these platforms.

[^1]: Environments are defined by `namespaces`, which are simply names. For example, a production branch might use the `production` namespace, while a development branch uses the `dev` namespace - each maintaining its own build numbers.
