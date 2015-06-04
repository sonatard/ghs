ghs
======

`ghs` - command-line utility for searching Github repositoy.

![](http://f.st-hatena.com/images/fotolife/s/sona-zip/20141029/20141029212146_original.gif?1414585446)

ghs options
===========

```sh
[sona ~]$ ghs --help
Usage:
  ghs [OPTION] "QUERY"

Application Options:
  -s, --sort=       The sort field. 'stars', 'forks', or 'updated'. (best match)
  -o, --order=      The sort order. 'asc' or 'desc'. (desc)
  -l, --language=   searches repositories based on the language they’re written in.
  -u, --user=       limits searches to a specific user name.
  -r, --repo=       limits searches to a specific repository.
  -v, --version     print version infomation and exit.
  -e, --enterprise= search from github enterprise.

Help Options:
  -h, --help        Show this help message

Github search APIv3 QUERY infomation:
   https://developer.github.com/v3/search/
   https://help.github.com/articles/searching-repositories/

Version:
   ghs 0.0.4 (https://github.com/sona-tar/ghs.git)
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
ghs "dotfiles in:name"
```

Limits searches to a specific user.
```zsh
ghs "dotfiles in:name" -u sona-tar
sona-tar/dotfiles                       dotfiles
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
