#!/usr/bin/env bash

set -e

if [ $# -ne 3 ]; then
    echo "Usage: $0 <os> <arch> <version>"
    echo "eg: $0 Linux x86_64 0.13.0"
    exit 1
fi

OS="$1"
ARCH="$2"
VERSION="$3"

cwd=$(pwd)

temp_dir=$(mktemp -d)
if [ ! -e ${temp_dir} ]; then
    echo "Failed to create temporary directory."
    exit 1
fi

cd $temp_dir

curl -sSLO "https://github.com/google/yamlfmt/releases/download/v${VERSION}/yamlfmt_${VERSION}_${OS}_${ARCH}.tar.gz"
curl -sSLO "https://github.com/google/yamlfmt/releases/download/v${VERSION}/checksums.txt"

sha256sum --ignore-missing -c checksums.txt

tar -xzf "yamlfmt_${VERSION}_${OS}_${ARCH}.tar.gz" -C ${temp_dir}/
cd $cwd

cp "${temp_dir}/yamlfmt" .

rm -r ${temp_dir}
