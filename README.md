# Taggo - git tag utility

![taggo gopher](https://www.jeorje.net/images/taggo.png)
Easy peasy git tag utility for lazy people who don't want to remember git commands.


Taggo handle the creation of lightweight tags and push them to the remote repository.
Tags are created with the format `vX.Y.Z` where `X` is the major version, `Y` is the minor version and `Z` is the patch version.

Future versions will allow to add pre-release and handle annotated tags.

## Installation
**Work in progress**


## Usage
- `taggo` ==> show last tag
- `taggo init` ==> create first tag v0.0.0
- `taggo -p` ==> create patch tag. Example: v0.0.1
- `taggo -m` ==> create minor tag. Example: v0.1.0
- `taggo -M` ==> create major tag. Example: v1.0.0
- `taggo -d` ==> delete last tag



## Git Commands for reference
- `git tag --sort=committerdate | tail -1` ==> last tag
- `git tag` ==> all tags
- `git tag v1.0.0` ==> create tag
- `git push origin v1.0.0` ==> push tag to remote
- `git tag --delete v1.0.0` ==> delete tag local
- `git push --delete origin v1.0.0` ==> delete tag remote
