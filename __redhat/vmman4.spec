%define debug_package   %{nil}
%define _build_id_links none
%define _name vmman4
%define _prefix /opt
%define _bash_completionsdir /usr/share/bash-completion/completions
%define _zsh_completionsdir  /usr/share/zsh/site-functions
%define _version 0.10.00~DEBUG
%define _rel 0
%define _arch x86_64
%define _binaryname vmman

Name:       vmman4
Version:    %{_version}
Release:    %{_rel}
Summary:    libvirt client

Group:      Virtualization
License:    GPL2.0
URL:        https://git.famillegratton.net:9722/devops/vmman4.git

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc, pkg-config, libvirt-devel
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
libvirt client

%prep
%autosetup

%build
cd src
go mod download
PATH=$PATH:/opt/go/bin CGO_ENABLED=1 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre
%install
install -Dpm 0755 %{_builddir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
# Bash completion — always install
/opt/bin/vmman completion bash > %{_bash_completionsdir}/vmman

# Zsh completion — only if zsh is present
if command -v zsh > /dev/null 2>&1; then
    mkdir -p %{_zsh_completionsdir}/zsh/site-functions
    /opt/bin/vmman completion zsh > %{_zsh_completionsdir}/_vmman
    mkdir -p %{_zsh_completionsdir}/zsh/site-functions
    /opt/bin/vmman completion zsh > %{_zsh_completionsdir}/_vmman
    zsh -c 'autoload -Uz compinit && compinit' 2>/dev/null || true
fi

%preun

%postun
if [ $1 -eq 0 ]; then
    # $1 == 0 means this is a full uninstall, not an upgrade
    rm -f %{_bash_completionsdir}/vmman
    rm -f %{_zsh_completionsdir}/_vmman
fi

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Thu May 28 2026 Binary package builder <builder@famillegratton.net> 0.10.00~DEBUG-0
- fixed paths in script
- archlinux build fixes
- disabled some file so compile will not fail
- fixed error in PKGBUILD
- Fixed output for Stop/Start, added the Rename subcommand
- completed the console
- fixed vm stop[all]
- Fixed issue where vm list failed if a VM is not fully booted
- stubbed most vm subcommands
- completed start/startall commands
- fixed various cosmetic issues
- Completed vm list and the conn subcommand
- first attempt at wiring everything
- completed package building wiring
- refreshed the whole structure
- initial stub
- interim submit

