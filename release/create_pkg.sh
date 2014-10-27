#!/bin/sh
#
# @(#) create_pkg.sh create release packages.
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
# └── work
#     └── 0.0.1
#         ├── ghs-0.0.1-darwin_386
#         │   ├── CHANGES
#         │   ├── ghs
#         │   └── README.md
#         ├── ghs-0.0.1-darwin_amd64
#         │   ├── CHANGES
#         │   ├── ghs
#         │   └── README.md
#         ├── ghs-0.0.1-linux_386
#         │   ├── CHANGES
#         │   ├── ghs
#         │   └── README.md
#         ├── ghs-0.0.1-linux_amd64
#         │   ├── CHANGES
#         │   ├── ghs
#         │   └── README.md
#         ├── ghs-0.0.1-windows_386
#         │   ├── CHANGES
#         │   ├── ghs.exe
#         │   └── README.md
#         └── ghs-0.0.1-windows_amd64
#             ├── CHANGES
#             ├── ghs.exe
#             └── README.md

repo_root=$(pwd)

XC_VERSION=$1
[ -z "${XC_VERSION}" ] && echo "usage : create_pkg.sh <version>" && exit 1

XC_ARCH=${XC_ARCH:-386 amd64}
XC_OS=${XC_OS:-linux darwin windows}

work_dir=${repo_root}/pkg/work/${XC_VERSION}
rm -rf pkg/
gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -output "${work_dir}/{{.Dir}}-${XC_VERSION}-{{.OS}}_{{.Arch}}/{{.Dir}}"



archive_dir=${repo_root}/pkg/archive/${XC_VERSION}
mkdir -p ${work_dir} ${archive_dir}

cd ${work_dir}
for target in *;
do
    cp ${repo_root}/README.md ${target}
    cp ${repo_root}/CHANGES   ${target}
    if [ $(echo $target | grep linux) ]; then
        tar zcvf ${archive_dir}/${target}.tar.gz ${target}
    else
        zip -r ${archive_dir}/${target}.zip ${target}
    fi
done
