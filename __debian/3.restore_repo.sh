#!/usr/bin/env bash

git restore control preinst prerm postinst postrm
rm -rf "vmman"*

apt remove -y gcc g++ pkg-config libvirt-dev
apt autoremove -y
