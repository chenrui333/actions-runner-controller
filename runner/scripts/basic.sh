#!/bin/bash
################################################################################
##  File:  basic.sh
##  Desc:  Installs basic command line utilities and dev packages
################################################################################

set -e

echo "Install dnsutils"
apt-get install -y --no-install-recommends dnsutils

echo "Install file"
apt-get install -y --no-install-recommends file

echo "Install ftp"
apt-get install -y --no-install-recommends ftp

echo "Install iproute2"
apt-get install -y --no-install-recommends iproute2

echo "Install iputils-ping"
apt-get install -y --no-install-recommends iputils-ping

echo "Install jq"
apt-get install -y --no-install-recommends jq

echo "Install libcurl3"
apt-get install -y --no-install-recommends libcurl3

echo "Install libunwind8"
apt-get install -y --no-install-recommends libunwind8

echo "Install locales"
apt-get install -y --no-install-recommends locales

echo "Install netcat"
apt-get install -y --no-install-recommends netcat

echo "Install openssh-client"
apt-get install -y --no-install-recommends openssh-client

echo "Install rsync"
apt-get install -y --no-install-recommends rsync

echo "Install shellcheck"
apt-get install -y --no-install-recommends shellcheck

echo "Install sudo"
apt-get install -y --no-install-recommends sudo

echo "Install telnet"
apt-get install -y --no-install-recommends telnet

echo "Install time"
apt-get install -y --no-install-recommends time

echo "Install tzdata"
apt-get install -y --no-install-recommends tzdata

echo "Install unzip"
apt-get install -y --no-install-recommends unzip

echo "Install upx"
apt-get install -y --no-install-recommends upx

echo "Install wget"
apt-get install -y --no-install-recommends wget

echo "Install zip"
apt-get install -y --no-install-recommends zip

echo "Install zstd"
apt-get install -y --no-install-recommends zstd

echo "Install libxkbfile"
apt-get install -y --no-install-recommends libxkbfile-dev

echo "Install pkg-config"
apt-get install -y --no-install-recommends pkg-config

echo "Install libsecret-1-dev"
apt-get install -y --no-install-recommends libsecret-1-dev

echo "Install libxss1"
apt-get install -y --no-install-recommends libxss1

echo "Install libgconf-2-4"
apt-get install -y --no-install-recommends libgconf-2-4

echo "Install dbus"
apt-get install -y --no-install-recommends dbus

echo "Install xvfb"
apt-get install -y --no-install-recommends xvfb

echo "Install libgtk"
apt-get install -y --no-install-recommends libgtk-3-0

echo "Install tk"
apt install -y tk

echo "Install fakeroot"
apt-get install -y --no-install-recommends fakeroot

echo "Install dpkg"
apt-get install -y --no-install-recommends dpkg

echo "Install rpm"
apt-get install -y --no-install-recommends rpm

echo "Install xz-utils"
apt-get install -y --no-install-recommends xz-utils

echo "Install xorriso"
apt-get install -y --no-install-recommends xorriso

echo "Install zsync"
apt-get install -y --no-install-recommends zsync

echo "Install curl"
apt-get install -y --no-install-recommends curl

echo "Install parallel"
apt-get install -y --no-install-recommends parallel

# Run tests to determine that the software installed as expected
echo "Testing to make sure that script performed as expected, and basic scenarios work"
for cmd in curl file ftp jq netcat ssh parallel rsync shellcheck sudo telnet time unzip wget zip; do
    if ! command -v $cmd; then
        echo "$cmd was not installed"
        exit 1
    fi
done

# Document what was added to the image
echo "Lastly, documenting what we added to the metadata file"
