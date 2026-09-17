%define debug_package   %{nil}
%define _build_id_links none
%define _name vmman4
%define _prefix /opt
%define _bindir %{_prefix}/bin
%define _version 1.2.0
%define _rel 1
%define _arch x86_64
%define _binaryname vmman

Name:       vmman4
Version:    %{_version}
Release:    %{_rel}
Summary:    Libvirt client

Group:      Virtualization
License:    GPL-3.0-or-later
URL:        https://git.famillegratton.net:3000/devops/vmman4.git

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc, pkg-config, libvirt-devel
Requires: libvirt-libs
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
Virtual Machine Manager

%prep
%autosetup

%build
cd src
go mod download
# rpmbuild runs %build under set -e, so both of these fast-fail the package
# build; go test exits 0 for packages with no test files and only fails on
# an actual test failure.
PATH=$PATH:/opt/go/bin CGO_ENABLED=1 go vet -tags libvirt_dlopen ./...
PATH=$PATH:/opt/go/bin CGO_ENABLED=1 go test -tags libvirt_dlopen ./...
PATH=$PATH:/opt/go/bin CGO_ENABLED=1 go build -tags libvirt_dlopen -trimpath -ldflags="-s -w -buildid= -X vmman4/cmd.buildVersion=%{_version} -X vmman4/cmd.buildDate=%(date +%%Y.%%m.%%d)" -o %{_builddir}/%{name}-%{version}/%{_binaryname} .

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
* Thu Sep 17 2026 Binary package builder <builder@famillegratton.net> 1.2.0-1
- Version bump
- bug: fixed field in manifest
- feature: dynamic version numbering, sleep between vm stop/start
- chore: update changelog for 1.1.0-1

* Mon Sep 14 2026 Binary package builder <builder@famillegratton.net> 1.1.0-1
- chore: doc update, version bump
- feat: new volume attach command
- enhancement: more help in all commands -h
- feat: added the cluster subcommand
- Merge branch 'main' into develop
- chore: update changelog for 1.0.0-1

* Mon Sep 14 2026 Binary package builder <builder@famillegratton.net> 1.0.0-1
- Merge branch 'develop'
- feat: completed pool and vol subcommands
- chore: update changelog for 0.8.1-1

* Mon Sep 14 2026 Binary package builder <builder@famillegratton.net> 0.8.1-1
- Merge branch 'develop'
- Version bump
- enhancement: sample files are now located in the config dir
- chore: update changelog for 0.8.0-1

* Mon Sep 14 2026 Binary package builder <builder@famillegratton.net> 0.8.0-1
- Merge branch 'develop'
- version bump before testing
- completed vm create and vm provision
- chore: update changelog for 0.7.0-1

* Sun Sep 13 2026 Binary package builder <builder@famillegratton.net> 0.7.0-1
- Merge branch 'develop'
- chore: version bump
- task: completed vm management
- task: completed snapshot management
- chore (vm_mgt) : renamed functions
- chore: renamed connection_mgt functions
- chore: update changelog for 0.6.0-2
- Merge branch 'develop'
- bug(BUILDERS): static builds workaround for missing lib functions
- chore: update changelog for 0.6.0-1

* Sat Sep 12 2026 Binary package builder <builder@famillegratton.net> 0.6.0-2
- Merge branch 'develop'
- bug(BUILDERS): static builds workaround for missing lib functions
- chore: update changelog for 0.6.0-1

* Fri Sep 11 2026 Binary package builder <builder@famillegratton.net> 0.6.0-1
- Merge branch 'develop'
- Merge branch 'vmmanagement' into develop
- chore: version bump
- feat: Completed vmInfo(), enhancements to vmInventory()
- completed the DumpXML command
- Added SetVMem() and SetVCPUs()
- added Reset() / ResetAll()
- chore: yet another package refactoring
- chore: package refactoring
- bug(vm remove): reordered tasks in vm rm command
- chore: update changelog for 0.5.0-1

* Thu Sep 10 2026 Binary package builder <builder@famillegratton.net> 0.5.0-1
- Merge branch 'develop'
- Merge branch 'vmmanagement' into develop
- feat: completed vm remove
- Merge branch 'vmmanagement' into develop
- Completed all vm subcommands except Console
- chore: update changelog for 0.4.0-4
- fixed deps issue
- Merge branch 'vmmanagement' into develop
- chore: block windows build
- chore: update changelog for 0.4.0-3
- chore: doc update
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-3
- bug(RPMBUILDER): fixed missing builddep
- chore: CL update
- chore: update changelog for 0.4.0-2
- bug(RPMBUILDER): wrong package name
- bug(DEBBUILDER): fixed endline issue in control file
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-1

* Wed Sep 09 2026 Binary package builder <builder@famillegratton.net> 0.4.0-4
- fixed deps issue
- Merge branch 'vmmanagement' into develop
- chore: block windows build
- chore: update changelog for 0.4.0-3
- chore: doc update
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-3
- bug(RPMBUILDER): fixed missing builddep
- chore: CL update
- chore: update changelog for 0.4.0-2
- bug(RPMBUILDER): wrong package name
- bug(DEBBUILDER): fixed endline issue in control file
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-1

* Wed Sep 09 2026 Binary package builder <builder@famillegratton.net> 0.4.0-3
- chore: doc update
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-3
- bug(RPMBUILDER): fixed missing builddep
- chore: CL update
- chore: update changelog for 0.4.0-2
- bug(RPMBUILDER): wrong package name
- bug(DEBBUILDER): fixed endline issue in control file
- Merge branch 'main' into develop
- chore: update changelog for 0.4.0-1

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

