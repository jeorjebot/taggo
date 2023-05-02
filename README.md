# Taggo


![GitHub release](https://img.shields.io/github/release/jeorjebot/taggo.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeorjebot/taggo)](https://goreportcard.com/report/github.com/jeorjebot/taggo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<!-- [![GoDoc](https://godoc.org/github.com/jeorje/taggo?status.svg)](https://godoc.org/github.com/jeorjebot/taggo) -->

Easy peasy `git tag` utility for lazy people who don't want to remember git commands.

![taggo gopher](https://www.jeorje.net/images/taggo_v1.png)

**Taggo** handle the creation of lightweight tags and push them to the remote repository.
Tags are created with the format `vX.Y.Z` where `X` is the major version, `Y` is the minor version and `Z` is the patch version.

Future versions will allow to add pre-release and handle annotated tags.

## Table of Contents

- [Installation](#installation)
  - [Go install command](#go-install-command)
  - [From releases](#from-releases)
- [Usage](#usage)
- [Examples](#examples)
- [Git Commands for reference](#git-commands-for-reference)
- [License](#license)
- [Thanks](#thanks)

## Installation
### Go install command
If you have Go installed, you can use the `go install` command to install the binary.

```bash
go install github.com/jeorjebot/taggo
```
The binary will be installed in `$GOPATH/bin` or `$GOBIN` if set.
Make sure you have `$GOPATH/bin` in your path.

### From releases
Download the binary for your OS from the [releases page](https://github.com/jeorjebot/taggo/releases).
Make sure the binary is executable, then move it to your path.

```bash
chmod +x /path/to/taggo
mv /path/to/taggo /usr/local/bin
```


## Usage
- `taggo` ==> show last tag
- `taggo init` ==> create first tag v0.0.0
- `taggo -p` ==> create patch tag. Example: v0.0.1
- `taggo -m` ==> create minor tag. Example: v0.1.0
- `taggo -M` ==> create major tag. Example: v1.0.0
- `taggo -t` ==> create tag specifying version. Example: v1.0.0
- `taggo -d` ==> delete last tag

## Examples
- Show last tag
```bash
$ taggo
[*] Current tag: v1.0.0
```

- Create first tag
```bash
$ taggo init
[*] Initializing git repo
[*] Added tag v0.0.0
```

- Create patch tag
```bash
$ taggo -p
[*] Current tag: v0.0.0
[*] New tag: v0.0.1
[*] Tag pushed successfully
```

- Create a specific tag
```bash
$ taggo -t v1.0.0
[*] Current tag: v0.0.1
[*] New tag: v1.0.0
[*] Tag pushed successfully
```

- Delete last tag
```bash
$ taggo -d
[*] Current tag: v0.0.2
[*] Deleting tag v0.0.2
[*] Tag deleted successfully
```

## Git Commands for reference
- `git tag --sort=committerdate | tail -1` ==> last tag
- `git describe --tags --abbrev=0` ==> last tag 
- `git tag` ==> all tags
- `git tag v1.0.0` ==> create tag
- `git push origin v1.0.0` ==> push tag to remote
- `git tag --delete v1.0.0` ==> delete tag local
- `git push --delete origin v1.0.0` ==> delete tag remote

## License
This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE.md) file for details

## Thanks
- [autotag](https://github.com/pantheon-systems/autotag) for the inspiration
- [gopherize.me](https://gopherize.me/) for the gopher image
