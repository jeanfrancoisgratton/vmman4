SPEC      := vman4.spec
NAME      := $(shell rpmspec -q --qf '%{name}\n'    $(SPEC) | head -1)
VERSION   := $(shell rpmspec -q --qf '%{version}\n' $(SPEC) | head -1)
RELEASE   := $(shell rpmspec -q --qf '%{release}\n' $(SPEC) | head -1)
SOURCES   := $(HOME)/rpmbuild/SOURCES
RPMS      := $(HOME)/rpmbuild/RPMS
DEST      := /data
REPOROOT  := $(shell git rev-parse --show-toplevel)
NEXUSREPO := dnfLocal

.PHONY: tarball rpm rpmcl clean upload changelog commitcl

# Intended workflow: make rpm → make upload → make commitcl

tarball:
        git -C $(REPOROOT) archive --prefix=$(NAME)-$(VERSION)/ HEAD \
          | gzip > $(SOURCES)/$(NAME)-$(VERSION).tar.gz

changelog: tarball
        $(REPOROOT)/__redhat/updateChangelog.sh

rpm: tarball
	sudo dnf builddep -y $(SPEC); \
	rpmbuild -bb $(SPEC); BUILD_STATUS=$$?; \
	sudo dnf history undo last -y; \
	find $(RPMS) -name '*.rpm' -exec cp {} $(DEST)/ \;; \
	exit $$BUILD_STATUS

rpmcl: changelog
	sudo dnf builddep -y $(SPEC); \
	rpmbuild -bb $(SPEC); BUILD_STATUS=$$?; \
	sudo dnf history undo last -y; \
	find $(RPMS) -name '*.rpm' -exec cp {} $(DEST)/ \;; \
	exit $$BUILD_STATUS

upload:
        nxtools upload -r $(NEXUSREPO) $(DEST)/$(NAME)-$(VERSION)-$(RELEASE).x86_64.rpm

commitcl:
        git add $(SPEC)
        git commit -m "chore: update changelog for $(VERSION)-$(RELEASE)"
		git push

clean:
        rm -f  $(SOURCES)/$(NAME)-$(VERSION).tar.gz
        rm -rf $(HOME)/rpmbuild/BUILD/$(NAME)-$(VERSION)
