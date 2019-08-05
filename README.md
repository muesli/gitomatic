git-o-matic
===========

A tool to monitor git repositories and automatically pull & push changes

## Installation

### Packages & Binaries

- Arch Linux: [gitomatic](https://aur.archlinux.org/packages/gitomatic/)
- [Binaries](https://github.com/muesli/gitomatic/releases) for Linux, macOS & Windows

### From Source

Make sure you have a working Go environment (Go 1.11 or higher is required).
See the [install instructions](http://golang.org/doc/install.html).

Compiling gitomatic is easy, simply run:

    git clone https://github.com/muesli/gitomatic.git
    cd gitomatic
    go build

## Usage

Monitor a repository for changes and automatically pull & push changes:

```
gitomatic <path>

2019/08/03 00:16:48 Checking repository: /tmp/gitomatic-test/
2019/08/03 00:16:48 Pulling changes...
2019/08/03 00:16:49 New file detected: hello_world.txt
2019/08/03 00:16:49 Adding file to work-tree: hello_world.txt
2019/08/03 00:16:49 Creating commit: Add hello_world.txt.
2019/08/03 00:16:49 Pushing changes...
2019/08/03 00:16:53 Sleeping until next check in 10s...
2019/08/03 00:17:03 Checking repository: /tmp/gitomatic-test/
2019/08/03 00:17:03 Pulling changes...
2019/08/03 00:17:07 Deleted file detected: hello_world.txt
2019/08/03 00:17:07 Removing file from work-tree: hello_world.txt
2019/08/03 00:17:07 Creating commit: Remove hello_world.txt.
2019/08/03 00:17:07 Pushing changes...
```

Auth methods:

```
gitomatic -privkey ~/.ssh/id_rsa <path>
gitomatic -username "someone" -password "mypass" <path>
```

If you want to pull new changes but don't create commits (or vice versa):

```
gitomatic -pull=true -push=false <path>
```

You can control how often gitomatic checks for changes:

```
gitomatic -interval 30m <path>
```

Change the commit author's name and email:

```
gitomatic -author "John Doe" -email "some@mail.tld" <path>
```
