# Makefile for __ansible (devops.vmman4 collection)

SHELL := /bin/bash
COLLECTION_NAME := devops-vmman4

.PHONY: help build clean install test lint

help:
	@echo "Available targets:"
	@echo "  build    - Build Ansible collection tarball (.tar.gz)"
	@echo "  clean    - Remove built collection tarballs"
	@echo "  install  - Install built collection into local environment (~/.ansible/collections)"
	@echo "  lint     - Run python syntax compilation check on modules"
	@echo "  test     - Run syntax-check on test playbook using local module library"

build: clean
	ansible-galaxy collection build --force

clean:
	rm -f *.tar.gz

install: build
	ansible-galaxy collection install $(COLLECTION_NAME)-*.tar.gz --force

lint:
	python3 -m py_compile plugins/modules/vmman.py
	python3 -m py_compile plugins/modules/vmman_info.py
	@echo "Syntax OK"

test: lint
	ANSIBLE_LIBRARY=plugins/modules ansible-playbook --syntax-check tests/test_local.yaml
