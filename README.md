# <img src="./images/vmman_banner.png" alt="vmman logo" height="384" width="1024" />
___

A CLI tool to manage a libvirtd-based KVM/QEMU hypervisor (local or remote, over SSH) and the virtual machines it hosts.

Today, from the command line, you can:
- list running/defined VMs, and start / stop them (one at a time or all at once)
- open a console session on a VM
- rename a VM
- list storage pools, and inspect a VM's storage (disks/volumes)
- manage named connections to local or remote hypervisors

Some capabilities exist in the codebase but are not yet exposed as commands — see [Roadmap](#roadmap).

**TABLE OF CONTENTS**<br>

[Basic concepts](#concepts)

[Build / Install requirements](#build-install-requirements)

[Using the tool](#using-the-tool)

[Command summary](#command-summary)
- [Global flags](#global-flags)
- [Connection operations](#conn-ops)
- [VM operations](#vm-ops)
- [Storage/pool operations](#storage-ops)
- [Shell completion](#completion-ops)

[Roadmap](#roadmap)

<a id="concepts"></a>
# Basic Virtualisation Concepts
TO BE CONTINUED

<a id="build-install-requirements"></a>
# Build/install requirements

You have three alternatives:
- [Install from source](#install-from-source)
- [Install from a binary package](#install-from-a-binary-package)
- [Build your own APK, DEB, RPM, or Arch package, then manually install it](#build-your-own-package)

Installing from source requires a bit more work in the sense that Go has to be installed on your system.

<a id="install-from-source"></a>
## Install from source
1. Clone/fork the repo: `git clone https://git.famillegratton.net:3000/devops/vmman4` (mirrored at `https://github.com/jeanfrancoisgratton/vmman4`)
2. Ensure that you have the proper Go version, as stated in the `go.version` file at the root of the repo. Your Go version should be equal to or higher than the one in that file. To check, run `go version`
3. cd to `src`, and then run `./updateBuildDeps.sh` to ensure all build dependencies are up to date; this might be overkill, but it's good hygiene
4. Run `./build.sh`. By default, the binary is created in `/opt/bin` (check that directory's permissions ahead of running it). Examine that script if you want to tailor the output — e.g. `./build.sh ~/bin`. Building off a branch other than `main`/`develop` produces a `vmman4-<branch>` binary instead, so a feature-branch build never clobbers your main one.

<a id="install-from-a-binary-package"></a>
## Install from a binary package
The simplest way: go to the RELEASES tab of the repo, pick your format, download it, and install it through your package manager.

<a id="build-your-own-package"></a>
## Build your own package
The scripts and files (`__alpine/`, `__debian/`, `__redhat/`, `__archlinux/`, `__macos/`) are there for my own ease of work; I usually build my packages using "builder containers" for each format: `apkbuilder`, `debbuilder`, `rpmbuilder`, `archbuilder`. I'll leave you with homework, and will show you how to roughly reproduce my environment.

**Whatever the format, the build is a CGO build**: `libvirt.org/go/libvirt` binds to the C library through `pkg-config`, and it needs four `.pc` files — `libvirt`, `libvirt-admin`, `libvirt-qemu` and `libvirt-lxc` — even though vmman4 only ever calls the plain libvirt API. All four ship in a single package:

| Distro | Development package | Runtime package |
| --- | --- | --- |
| Alpine | `libvirt-dev` | `libvirt-libs` (resolved automatically by abuild) |
| Debian / Ubuntu | `libvirt-dev` | `libvirt0` |
| RedHat / Fedora / Rocky | `libvirt-devel` (in the `crb`/`powertools` repo on RHEL clones) | resolved automatically by rpm |
| Archlinux | `libvirt` (no split -devel package) | `libvirt` |
| macOS | `libvirt` + `pkg-config` (via Homebrew) | n/a — no package is produced, see below |

Note that the Fedora/RHEL package literally named `libvirt-admin` is *not* what you want: that one only ships the `virt-admin` CLI. The `libvirt-admin.pc` file the build needs comes from `libvirt-devel`.

Build with `CGO_ENABLED=1`; `CGO_ENABLED=0` cannot work here.

### APKBUILDER : Alpine Linux
1. In an Alpine container or VM, install the build dependencies: `abuild-doc pax-utils git alpine-sdk libvirt-dev` (some others might be needed, depending on the config in `__alpine/APKBUILD`)
2. From `__alpine`, run `abuild -r`, or use the Makefile: `make build` (build the `.apk`), `make release` (build, then upload to the `apkLocal` Nexus repo), `make clean`, `make info` (print name/version/arch extracted from `APKBUILD`)

### DEBBUILDER : Debian-based distros (Debian, Ubuntu, Mint, etc)
1. cd to `__debian`
2. Run `./install-build-deps.sh`; it pulls in `gcc fakeroot devscripts build-essential pkg-config libvirt-dev`, then installs the Go version named in `../go.version`
3. Run `make build` to produce the `.deb`, `make release` to build and upload it to the `aptLocal` Nexus repository, `make clean` to remove the staging tree and the built `.deb`, and `make info` to print name/version/arch extracted from the `control` file

### RPMBUILDER : RedHat-based distros (RedHat, CentOS, Fedora, RockyLinux)
1. cd to `__redhat`. Everything is run from there
2. Run `./rpmbuild-deps.sh` to install the specfile's `BuildRequires` (this also enables the `crb`/`powertools` repo where needed, and installs the Go toolchain from `../go.version`); you might also need `rpm-build` and `rpmdevtools` beforehand
3. Run `make` with one of:
   - `make rpm`: build the package
   - `make rpmcl`: the above, with the latest changes folded into the specfile's `%changelog` section
   - `make commitcl`: push that changelog update to the `develop` branch (`CL_BRANCH`)
   - `make tag`: tag the release (`vX.Y.Z`, with `~` in the version swapped for `-` so it's a valid git ref)
   - `make upload`: upload the built `.rpm` to the `dnfLocal` Nexus repository
   - `make release`: run `rpmcl` → `commitcl` → `upload` → `tag`, in that order
   - `make clean`: remove the generated tarball and build tree
4. If you don't have a Nexus repository manager, or the `nxtools` uploader it relies on, building and distributing the package once made is up to you.

### ARCHBUILDER : Archlinux-based distros
1. In an Arch container or VM, cd to `__archlinux` and run `./install-build-deps.sh` (installs `base-devel` via pacman)
2. Run `make build` (calls `makepkg`, reading package metadata from `PKGBUILD`), `make release` (build, refresh the pacman `repo-add` database, then upload the package and database to the `archLocal` Nexus repository), `make clean`, or `make info` to print name/version/pkgrel/arch extracted from `PKGBUILD`

### macOS
There's no macOS package (no `.pkg`, no Homebrew formula) — `__macos/` just builds a plain `vmman4` binary. See [`__macos/README.md`](__macos/README.md) for the two-step process (`goget-macos.sh` to fetch the right Go toolchain, `build-macos.sh` to build). It's the same CGO build as everywhere else, so `libvirt` and `pkg-config` need to be available first (e.g. `brew install libvirt pkg-config`).

<a id="using-the-tool"></a>
# Using the tool

Once you have a `vmman4` binary on your `PATH` (or aliased as `vmman`), it talks to `qemu:///system` on the local machine by default — no setup needed if libvirtd is running locally and you're in the right group (e.g. `libvirt` / `kvm`) to talk to it:

```sh
vmman vm list
```

## Managing a remote hypervisor
To manage a remote hypervisor over SSH, create a named connection once:

```sh
vmman connection add
```

You'll be prompted for a connection name, a hostname, a username, and an optional comment. This is saved as `~/.config/JFG/vmman4/<name>.json`. From then on, point any command at it with `-c`/`--connectionfile`:

```sh
vmman -c myhypervisor vm list
```

This resolves to `qemu+ssh://<user>@<host>/system` under the hood (SSH key auth is assumed — there's no password field). To bypass named connections entirely, pass a raw libvirt URI with `-C`/`--connectionuri` instead.

Manage saved connections with `vmman connection ls`, `vmman connection info <name>`, and `vmman connection rm <name>`.

<a id="command-summary"></a>
# Command summary

<a id="global-flags"></a>
## Global flags
These apply to every subcommand:

| Flag | Short | Description |
| --- | --- | --- |
| `--quiet` | `-q` | Suppress non-essential output |
| `--debug` | `-D` | Enable debug mode |
| `--connectionfile` | `-c` | Named connection file to use (see [Using the tool](#using-the-tool)) |
| `--connectionuri` | `-C` | Raw libvirt connection URI (default `qemu:///system`) |

Other top-level commands: `vmman version` (prints the software and Go versions).

<a id="conn-ops"></a>
## Connection operations
`vmman connection <subcommand>` (alias: `vmman conn`)

| Subcommand | Aliases | Description |
| --- | --- | --- |
| `ls` | `list` | List saved connection files |
| `add` | `create` | Interactively create a connection file |
| `info NAME...` | `explain` | Print the details of one or more connection files |
| `rm NAME...` | `del`, `remove` | Delete one or more connection files |

<a id="vm-ops"></a>
## VM operations
`vmman vm <subcommand>`

| Subcommand | Aliases | Description |
| --- | --- | --- |
| `list` | `ls` | List all VMs known to the connection |
| `start VM...` | `up` | Start one or many VMs |
| `startall` | | Start all VMs at once |
| `stop VM...` | `down` | Stop one or many VMs |
| `stopall` | | Stop all VMs at once |
| `console VM` | | Open a console session on a VM (`-f`/`--force` to kick out a previous session) |
| `rename OLD NEW` | | Rename a VM. The VM's underlying disk keeps its old name; any existing snapshots must be removed first |

`list`, `start`, `startall`, `stop`, and `stopall` are also available directly off the root command (e.g. `vmman list` works the same as `vmman vm list`).

<a id="storage-ops"></a>
## Storage/pool operations

| Command | Aliases | Description |
| --- | --- | --- |
| `vmman pool list` | `vmman pool ls` | List storage pools with extended info |
| `vmman slist VM` | `vmman sls VM` | List the storage (disks/volumes) attached to a specific VM |

<a id="completion-ops"></a>
## Shell completion
`vmman completion bash` and `vmman completion zsh` generate completion scripts — see `vmman completion --help` for how to load them into your shell.

<a id="roadmap"></a>
# Roadmap

The following are planned but not yet wired up as CLI commands, even though groundwork for some of them already exists in the codebase:
- **VM creation** from a JSON spec file (disks, NICs, arch/machine) — the underlying logic exists, but there's no `vmman vm create` yet
- **Snapshot management** (list / create / remove / revert) — snapshot listing exists internally but isn't exposed as a command yet
- **Resource editing** (CPU, memory, storage resize) for existing VMs
