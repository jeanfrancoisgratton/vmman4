%define debug_package   %{nil}
%define _build_id_links none
%define _name vmman4
%define _prefix /opt
%define _bindir %{_prefix}/bin
%define _version 0.4.0
%define _rel 2
%define _arch x86_64
%define _binaryname vmman4

Name:       vmman4
Version:    %{_version}
Release:    %{_rel}
Summary:    Libvirt client

Group:      Virtualization
License:    GPL-3.0-or-later
URL:        https://git.famillegratton.net:3000/devops/vmman4.git

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
Virtual Machine Manager

%prep
%autosetup

%build
cd src
go mod download
PATH=$PATH:/opt/go/bin CGO_ENABLED=1 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{name}-%{version}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre

%install
rm -rf %{buildroot}
install -Dpm 0755 %{_builddir}/%{name}-%{version}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post

%preun

%postun

%files
%defattr(0755,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Wed Sep 09 2026 Binary package builder <builder@famillegratton.net> 0.4.0-2
- bug(RPMBUILDER): wrong package name
- bug(DEBBUILDER): fixed endline issue in control file
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-1

* Wed Sep 09 2026 Binary package builder <builder@famillegratton.net> 0.4.0-1
- chore: perm fix
- removed windows support, CGO build is now consistent
- removed dontexec
- completed repo resync
- Merge remote-tracking branch 'refs/remotes/origin/develop' into develop
- Added windows and macos support, go version bump, software version numbering now SemVer-aligned
- chore: update changelog for 0.4.0~DEBUG-4
- updated build deps
- RPMBUILDER: fixed build order
- RPMBUILDER: dummy build to test the new build order
- added a new script (to be integrated in stubber) to manage jenkins builds
- Merge remote-tracking branch 'refs/remotes/origin/main'
- Enable builds

* Tue Jun 09 2026 Binary package builder <builder@famillegratton.net> 0.20.00~DEBUG-0
- version bump
- completed pool list
- Merge remote-tracking branch 'refs/remotes/origin/develop' into develop
- full resynch
- minor doc update
- post-install fixes
- stubbed poolmanagement and snapshotmanagement
- stubbed vmCreate
- fixed path for autocompletion
- moved vm list to vmManagement
- removed the installation of un-needed deps
- updated build deps scripts
- fixed whitespace issue in Makefile
- fixed whitespace issue in Makefile
- fixed typo
- Merge remote-tracking branch 'refs/remotes/origin/develop' into develop
- Makefile fixes
- chore: update changelog for 0.10.00~DEBUG-0
- Merge remote-tracking branch 'refs/remotes/origin/develop' into develop
- Makefile fixes so I can cleanup the builddeps post-build
- chore: update changelog for 0.10.00~DEBUG-0
- added build dep
- chore: update changelog for 0.10.00~DEBUG-0
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

