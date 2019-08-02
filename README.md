git-o-matic
===========

A tool to monitor a git repository and automatically pull & push changes

## Installation

Make sure you have a working Go environment (Go 1.11 or higher is required).
See the [install instructions](http://golang.org/doc/install.html).

To install gitomatic, simply run:

    go get github.com/muesli/gitomatic

## Usage

Monitor a repository for changes and automatically pull & push changes:

```
gitomatic <path>
```

Available parameters:

```
gitomatic -interval 30m
gitomatic -privkey ~/.ssh/id_rsa
gitomatic -author "John Doe"
gitomatic -email "some@mail.tld"
```
