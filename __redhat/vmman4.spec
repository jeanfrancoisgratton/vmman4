%define debug_package   %{nil}
%define _build_id_links none
%define _name vmman4
%define _prefix /opt
%define _version 0.4.0~DEBUG
%define _rel 3
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

%preun

%postun

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Fri Aug 14 2026 Binary package builder <builder@famillegratton.net> 0.4.0~DEBUG-1
- repackaged to follow SemVer
- RPMBUILDER: fixed tag using reserved characted issue
- chore: update changelog for 0.30.00~DEBUG-1
- BUILDERS: added missing build dependencies everywhere, doc update
- removed shell completion scripts, builddeps updates, go update -> v1.26.6
- fixed dependencies in arch pkg building
- fixed fmt.Sprintf() issue that broke builds on archlinux
- updated build-deps, second attempt
- updated build-deps
- chore: update changelog for 0.20.00~DEBUG-0
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

* Fri Aug 14 2026 Binary package builder <builder@famillegratton.net> 0.30.00~DEBUG-1
- BUILDERS: added missing build dependencies everywhere, doc update
- removed shell completion scripts, builddeps updates, go update -> v1.26.6
- fixed dependencies in arch pkg building
- fixed fmt.Sprintf() issue that broke builds on archlinux
- updated build-deps, second attempt
- updated build-deps
- chore: update changelog for 0.20.00~DEBUG-0
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

