%define debug_package   %{nil}
%define _build_id_links none
%define _name vmman4
%define _prefix /opt
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
BuildRequires: gcc
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
libvirt client

%prep
%autosetup

%build
cd src
go mod download
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre
%install
install -Dpm 0755 %{_builddir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
# Bash completion — always install
vmman completion bash > %{_datadir}/bash-completion/completions/vmman

# Zsh completion — only if zsh is present
if command -v zsh > /dev/null 2>&1; then
    mkdir -p %{_datadir}/zsh/site-functions
    vmman completion zsh > %{_datadir}/zsh/site-functions/_vmman
fi

%preun

%postun
if [ $1 -eq 0 ]; then
    # $1 == 0 means this is a full uninstall, not an upgrade
    rm -f %{_datadir}/bash-completion/completions/vmman
    rm -f %{_datadir}/zsh/site-functions/_vmman
fi

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
