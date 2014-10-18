ghs
======

`ghs` - command-line utility for searching Github repositoy.

ghs usage
===========
```sh
[sona ~]$ ghs --help
Usage:
  ghs [OPTION] "QUERY"(The search keywords, as well as any qualifiers.)

Application Options:
  -s, --sort=  The sort field. 'stars', 'forks', or 'updated'. (best match)
  -o, --order= The sort order. 'asc' or 'desc'. (desc)

Help Options:
  -h, --help   Show this help message

Github search APIv3 QUERY infomation:
   https://developer.github.com/v3/search/
   https://help.github.com/articles/searching-repositories/
```

Exapmle
===========
```sh
[sona ~]$ ghs github
michael/github                          A higher-level wrapper around the Github API. Intended for the browser.
peter-murach/github                     Ruby interface to github API v3
jwiegley/github                         The github API for Haskell
isaacs/github                           Just a place to track issues and feature requests that I have for github
gulinghao1847/github
chscodecamp/github                      GitHub Seminar
opauth/github                           GitHub authentication strategy for Opauth
Kdyby/Github                            Github API client with authorization for Nette Framework
JeroenDeDauw/GitHub
```

With [motemen/ghq](https://github.com/motemen/ghq)
===========
```sh
ghs  | peco | awk '{print $1}' | ghq import

```


Zsh function
===========
```zsh
function gpi () {
  [ "$#" -eq 0 ] && echo "Usage : gpi QUERY" && return 1
  ghs "$@" | peco | awk '{print $1}' | ghq import
}
```

gpi usage
===========
1. exec gpi QUERY
![1](http://f.st-hatena.com/images/fotolife/s/sona-zip/20141018/20141018194948_original.png?1413630026)
2. filtering by peco
![ghsWithghq](http://f.st-hatena.com/images/fotolife/s/sona-zip/20141018/20141018194949_original.gif?1413630039)
3. clone repository by ghq
![ghqimport](http://f.st-hatena.com/images/fotolife/s/sona-zip/20141018/20141018194950_original.png)

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

