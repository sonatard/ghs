ghs
======

[![Go Report Card](https://goreportcard.com/badge/github.com/sona-tar/ghs)](https://goreportcard.com/report/github.com/sona-tar/ghs)

`ghs` - command-line utility for searching Github repositoy.

![](http://f.st-hatena.com/images/fotolife/s/sona-zip/20141029/20141029212146_original.gif?1414585446)

Trial
===========
[ghs trial page](https://codepicnic.com/consoles/ghs/embed?sidebar=closed&hide=save,show_result,download,options,run,second_terminal,readme)


ghs options
===========

```sh
[sona ~]$ ghs --help
Usage:
  ghs [OPTION] "QUERY"

Application Options:
  -f, --fields=     limits what fields are searched. 'name', 'description', or 'readme'. (default: name,description)
  -s, --sort=       The sort field. 'stars', 'forks', or 'updated'. (default: best match)
  -o, --order=      The sort order. 'asc' or 'desc'. (default: desc)
  -l, --language=   searches repositories based on the language they’re written in.
  -u, --user=       limits searches to a specific user name.
  -r, --repo=       limits searches to a specific repository.
  -m, --max=        limits number of result. range 1-1000 (default: 100)
  -v, --version     print version infomation and exit.
  -e, --enterprise= search from github enterprise.
  -t, --token=      Github API token to avoid Github API rate limit

Help Options:
  -h, --help        Show this help message

Github search APIv3 QUERY infomation:
   https://developer.github.com/v3/search/
   https://help.github.com/articles/searching-repositories/

Version:
   ghs 0.0.7 (https://github.com/sona-tar/ghs.git)
```

Install
===========

[homebrew](http://brew.sh/index_ja.html), [linuxbrew](http://brew.sh/linuxbrew/)

```zsh
brew install sona-tar/tools/ghs
```

for Windows
[Releases sona-tar/ghs](https://github.com/sona-tar/ghs/releases)


Usage
===========

basic usage.
default search target.(name, description and readme)
```zsh
ghs "dotfiles"
```

You can restrict the search to just the repository name.
```zsh
ghs -f name "dotfiles"
```

Limits searches to a specific user.
```zsh
ghs -f name -u sona-tar "dotfiles"
sona-tar/dotfiles                       dotfiles
```

Github Authentication to avoid Github API rate limit
===========

Priority of authentication token

1. Exec `ghs` with `-t` or `--token` option

```bash
$ ghs -t "....."
```

2. `GITHUB_TOKEN` environmental variable
```bash
$ export GITHUB_TOKEN="....."
```

3. github.token in gitconfig

```bash
$ git config --global github.token "....."
```


With [motemen/ghq](https://github.com/motemen/ghq) and [peco/peco](https://github.com/peco/peco)
===========

```sh
ghs QUERY | peco | awk '{print $1}' | ghq import
```

crate zsh function

```zsh
function gpi () {
  [ "$#" -eq 0 ] && echo "Usage : gpi QUERY" && return 1
  ghs "$@" | peco | awk '{print $1}' | ghq import
}
```

gpi usage

1. exec gpi QUERY
2. filtering by peco
3. clone repository by ghq

![](http://f.st-hatena.com/images/fotolife/s/sona-zip/20141029/20141029210952_original.gif?1414584687)

more
===========

```zsh
function gpr () {
  ghq list --full-path | peco | xargs rm -r
}
```

```sh
gpr
```


Contributors
===========

[kou-m](https://github.com/kou-m)


Author
===========

[sona-tar](https://github.com/sona-tar)
