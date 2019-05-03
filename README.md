# gitPD - Git Police Department


## About
gitPD is a tool meant to download and scan github organizations' repositories, team members' repositories and gists or a custom given list of them.

Scanning is done by open source tool:
* [truffleGopher](https://github.com/michaelgrifalconi/trugglegopher) - high performance commit scanner for user provided regular expressions which is based on the same concepts of [truffleHog](https://github.com/dxa4481/truffleHog).

## Architecture

* *sergeant* is a reworked fork of git-all-secrets, with similar functionality. It is responsible to find download and scan repositories/gists in a organization, team or organization members. Seargent can be called as standalone service without involving other components.
* *sheriff* package is the brain of the project, meant to be a server who orchestrate various components like the web interface, sergeant's scan schedules and the database of scanned items and findings.

## Getting started
The only currently supported way to run gitpd is using its Docker image.
Standalone support may be added in future based on user's feedback.

* TODO: document way to start scan and mount results directory as well
* TODO: deADLOCK timeout, avoid stuck on huge repos

## Currently available features and wishlist

* [ ] Full-org scan: find, download and scan all repos of an organization, its users' repos and gists
* [ ] Download only flag: to scan later on
* [ ] Fech instead of re-download if possible: to improve performance
* [ ] Blacklist repos: skip to download/scan a given repo based on user-provided list
* [ ] Stateful runs, keep track of truffleGopher scanned items
* [ ] Evaluate refactoring of github repo enumeration and cloning to drop git binary and go-github packages in favor of libgit2
* [ ] Web interface to manage target repositories, browse scan results
* [ ] Update dependencies
* [ ] Better tests
* [ ] Refactor to follow https://github.com/golang-standards/project-layout
* [ ] support for custom user-provided tool

### Is this a fork of git-all-secrets?
* package `sergeant` was born as a fork of [git-all-secrets](https://github.com/anshumanbh/git-all-secrets), but is evolving to a different form:
  * different name 

TODO: https://github.com/TrueFurby/go-callvis graph