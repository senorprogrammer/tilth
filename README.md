<p align="center"><img src="till_header.png" width="916" height="306" alt="til" title="til: jot it down" /></p>

`til` is a fast, simple, command line-driven, mini-static site builder. Three commands, that's it (two if you're not too picky about your commit messages).

# tl;rd

```bash
❯ til New title here
❯ til -save
```

And you're done.

# Contents

* [Installation](#installation)
* [Configuration](#configuration)
    * [Example](#config-example)
* [Execution](#execution)
* [Usage](#usage)
    * [Creating a new page](#creating-a-new-page)
    * [Building static pages](#building-static-pages)
    * [Building, saving, committing, and pushing](#building-saving-committing-and-pushing)

## Installation

To use this yourself, the simplest approach is probably to fork this repo and delete the contents of the `/docs` directory.

You can also:

* create your own empty git repo
* copy `til.go` into it
* create the `docs` directory: `mkdir docs`
* push that up to GitHub

Now run `go run til.go --help` to initialize everything and make sure it's working.

## Configuration

When you first ran `go run til.go --help` it either exploded with an error message (open an issue here with the message), or it displayed the help info. If you saw help info, it also will have created a configuration file that you'll need to edit.

The config file is in `~/.config/til/config.yml` (if you're an XDG kind of person, it will be wherever you've set that to).

Open `~/.config/til/config.yml`, change the following, and save:

    * committerEmail
    * committerName
    * editor
    
`committerEmail` and `committerName` are the values it will use to commit with when running `til -build`. 

`editor` is the text editor it will open for writing in when running `til [some title here]`.

### Config Example

```
---
commitMessage: "build, save, push"
committerEmail: test@example.com
committerName: "TIL Autobot"
editor: "mvim"
```

## Execution

I use this `zsh` alias to execute it from whichever directory I'm in:

```shell
alias til='cd ~/Documents/til && go run ./til.go'
```

I don't bother to compile/install it, it's fast enough as-is. However, `go install` works just fine. 

## Usage

`til` only has three usage options: `til`, `til -build`, and `til -save`.

### Creating a new page

```bash
❯ til Testing title
2020-04-20T14-52-57-testing-title.md
```

And then that page will open in [MacVim](https://macvim-dev.github.io/macvim/) or whichever editor you've defined in your config.

### Building static pages

```bash
❯ til -build
```

Builds the index and tag pages, and leaves them uncommitted.

### Building, saving, committing, and pushing

```bash
❯ til -save [optional commit message]
```

Builds the index and tag pages, commits everything to the git repo with the commit message you've defined in your config, and pushes it all up to the remote repo.

`-save` takes an optional commit message. If that message is supplied, it will be used as the commit message. If that message is not supplied, the `commitMessage` value in the config file will be used. If that value is not supplied, an error will be raised.

As an example: [https://senorprogrammer.github.io/til/](https://senorprogrammer.github.io/til/).
