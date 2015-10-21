#!/bin/sh
#
# @(#) release_pkg.sh create release packages.
#
# pkg
# ├── archive
# │   └── 0.0.1
# │       ├── ghs-0.0.1-darwin_386.zip
# │       ├── ghs-0.0.1-darwin_amd64.zip
# │       ├── ghs-0.0.1-linux_386.tar.gz
# │       ├── ghs-0.0.1-linux_amd64.tar.gz
# │       ├── ghs-0.0.1-windows_386.zip
# │       └── ghs-0.0.1-windows_amd64.zip
#
# Release Frow
# $ git tag -a ${VERSION}
# $ git push --tags
# $ ./release/create_pkg.sh ${VERSION}
# $ ./release/release_pkg.sh ${VERSION}

XC_VERSION=$1
[ -z "${XC_VERSION}" ] && echo "usage : release_pkg.sh <version>" && exit 1

ghr $2 ${XC_VERSION} pkg/archive/${XC_VERSION}/

openssl sha1 pkg/archive/${XC_VERSION}/*
