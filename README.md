# <img src="./images/vmman_banner.png" alt="vmman logo" height="384" width="1024" />
___

A tool to manage a libvirtd-based KVM/QEMU hypervisor so you can handle the virtual machines it hosts. With it you can :
- list / create / stop / start / remove VMs
- list / create / remove VM snapshots
- edit the VM's resources (CPU, Memory, storage)

**TABLE OF CONTENTS**<br>

[Basic concepts](#concepts)

[Build / Install requirements](#build-install-requirements)

[Using the tool](#using-the-tool)

[Command summary](#command-summary)
- [Connection operations](#conn-ops)
- [VM operations](#vm-ops)
- [Snnapshot operations](#snap-ops)
- [Resource operations](#resource-ops)

<a id="concepts"></a>
# Basic Virtualisation Concepts
TO BE CONTINUED


<a id="build-install-requirements"></a>
# Build/install requirements

You have three alternatives:
- [Install from source](#install-from-source)
- [Install from a binary package](#install-from-a-binary-package)
- [Build your own APK, DEB, RPM, ARCH packages, then manually install those packages](#build-your-own-package)

Installing from source requires a bit more work in the sense that GO has to be installed on your system

<a id="install-from-source"></a>
## Install from source
1. Clone/fork the repo : either `git clone https://github.com/jeanfrancoisgratton/nxtools` or `git clone https://git.famillegratton.net:3000/devops/nxtools`
2. Ensure that you have the proper GO version, as stated in the `go.version` file in the root of the repo. Your GO version should be equal or higher than the one in that file. To ensure, run `go version`
3. cd to `src`, and then run: `./updateBuildDeps.sh`, to ensure that all build dependencies are up to date; this might be overkill, but I always run it nonetheless
4. run `./build.sh`. By default, the binary will be created in /opt (check the dir's permission ahead of running it). Examine that script, you can taylor the output as you see fit

<a id="install-from-a-binary-package"></a>
## Install from a binary package
The simplest way : just go in the RELEASES tab of the repo, select your format, download it, and then install through you package manager

<a id="build-your-own-package"></a>
## Build your own package
The scripts and files (__alpine/, __debian/, __redhat, __archlinux/) are there for my own ease of work; I usually build my tools using "builder containers" for each format: `apkbuilder`, `debbuilder`, `rpmbuilder`
I'll leave you with homeworks, and will show you how to roughly reproduce my environment

**Whatever the format, the build is a CGO build**: `libvirt.org/go/libvirt` binds to the C library through `pkg-config`, and it needs four `.pc` files -- `libvirt`, `libvirt-admin`, `libvirt-qemu` and `libvirt-lxc` -- even though vmman4 only ever calls the plain libvirt API. All four ship in a single package:

| Distro | Development package | Runtime package |
| --- | --- | --- |
| Alpine | `libvirt-dev` | `libvirt-libs` (resolved automatically by abuild) |
| Debian / Ubuntu | `libvirt-dev` | `libvirt0` |
| RedHat / Fedora / Rocky | `libvirt-devel` (in the `crb` repo on RHEL clones) | resolved automatically by rpm |
| Archlinux | `libvirt` (no split -devel package) | `libvirt` |

Note that the Fedora/RHEL package literally named `libvirt-admin` is *not* what you want: that one only ships the `virt-admin` CLI. The `libvirt-admin.pc` file the build needs comes from `libvirt-devel`.

Build with `CGO_ENABLED=1`; `CGO_ENABLED=0` cannot work here.

### APKBUILDER : Alpine Linux
1. In an Alpine container or VM, you need the following packages: `abuild-doc pax-utils git alpine-sdk libvirt-dev`. Some other packages might be needed, depending on the config in __alpine/APKBUILD
2. From the `__alpine`, run: `abuild -r`

This should give you an Alpine package

### DEBBUILDER : Debian-based distros (Debian, Ubuntu, Mint, etc)
1. cd to `__debian`
2. Run `./install-build-deps.sh`; it pulls in `build-essential`, `pkg-config` and `libvirt-dev`, then installs the GO version named in `../go.version`
3. Run `make build` to produce the .deb, `make release` to build and upload it to your nexus repository, and `make clean` to drop the staging tree

### RPMBUILDER : RedHat-based distros (RedHat, CentOS, Fedora, RockyLinux)
1. cd to `__redhat`. Everything is run from there
2. run `./rpmbuild-deps.sh` to ensure that everything needed to build; you might also need to install `rpm-build` and `rpmdevtools`
3. run `make` :
   - `make rpm`: build the package 
   - `make rpmcl`: the above, with the latest changes integrated in the specfile's Changelog section
   - `make commitcl` : to commit the changelog changes to the git repo
   - `make upload` : upload the binary package to your nexus repository (assuming you have one, and that `nxtools` is installed)
4. If you do not have a `nexus repository manager` server, or `nxtools`, it is your own responsibility to manage the binary package once built.

### ARCHBUILDER : Archlinux-based distros
TO BE CONTINUED

<a id="using-the-tool"></a>
# USING THE TOOL

*FIX ME*

ADD DOCUMENTATION 