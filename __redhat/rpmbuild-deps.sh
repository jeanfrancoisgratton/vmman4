#!/usr/bin/env bash
set -euo pipefail

# One-time provisioning of the RPMBUILDER container. This installs the specfile's
# BuildRequires into the running container; they persist there. rpmbuild only
# *checks* BuildRequires, it never installs (nor removes) anything itself.

echo "Installing BuildRequires dependencies"; echo

sudo dnf install -y dnf-plugins-core rpm-build rpmdevtools

# libvirt-devel lives in CodeReady Builder on RHEL and its clones, and that repo
# is disabled by default. It is called 'crb' on 9+, 'powertools' on 8, and does
# not exist at all on Fedora -- hence the fallback chain.
sudo dnf config-manager --set-enabled crb 2>/dev/null \
  || sudo dnf config-manager --set-enabled powertools 2>/dev/null \
  || echo "NOTE: no crb/powertools repo to enable (Fedora?), continuing."

# Reads BuildRequires straight from the specfile, so there is a single source of
# truth and versioned constraints are honoured.
sudo dnf builddep -y vmman4.spec

echo; echo; echo "Done. Now installing the Go binaries"

VER=$(cat ../go.version)
ARCH=${1:-amd64}

echo "Fetching archive..."
sudo wget -q "https://go.dev/dl/go${VER}.linux-${ARCH}.tar.gz" -O /opt/go.tar.gz

echo "Unarchiving..."
cd /opt && sudo rm -rf go && sudo tar zxf go.tar.gz && sudo rm -f go.tar.gz

echo "Completed."
