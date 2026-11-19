#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ACTION="${1:-build}"

case "${ACTION}" in
  lint)
    python3 -m py_compile "${SCRIPT_DIR}/plugins/modules/vmman.py"
    python3 -m py_compile "${SCRIPT_DIR}/plugins/modules/vmman_info.py"
    echo "Python syntax OK"
    ;;
  test)
    "${SCRIPT_DIR}/build.sh" lint
    ANSIBLE_LIBRARY="${SCRIPT_DIR}/plugins/modules" ansible-playbook --syntax-check "${SCRIPT_DIR}/tests/test_local.yaml"
    ;;
  build)
    "${SCRIPT_DIR}/build.sh" clean
    ansible-galaxy collection build --force "${SCRIPT_DIR}" --output-path "${SCRIPT_DIR}"
    echo "Collection built successfully."
    ;;
  install)
    "${SCRIPT_DIR}/build.sh" build
    ansible-galaxy collection install "${SCRIPT_DIR}"/devops-vmman4-*.tar.gz --force
    ;;
  clean)
    rm -f "${SCRIPT_DIR}"/*.tar.gz
    ;;
  *)
    echo "Usage: $0 [build|install|lint|test|clean]"
    exit 1
    ;;
esac
