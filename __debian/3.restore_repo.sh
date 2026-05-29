#!/usr/bin/env bash

git restore control preinst prerm postinst postrm
rm -rf "vmman"*

sudo apt remove -y gcc pkg-config libvirt-dev
sudo apt autoremove -y
