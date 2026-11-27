#!/usr/bin/env sh
# vmman4 — bump the package version, the release number and/or the Go
# version everywhere they need to be typed by hand.
#
#   -v <version>    package version      vmman4.json "versionnumber"
#                                        + every __*/ packaging file
#   -r <release>    release number       vmman4.json "defaultreleasenumber"
#                                        + every __*/ packaging file
#   -g <goversion>  Go version           vmman4.json "goversion" + go.version
#
# Any flag left out means "keep the current value". At least one flag is required.
#
# Version and release are different things (see docs/CHANGELOG.md history /
# the packaging Makefiles): release numbers are NOT kept in sync across
# distros day to day — a Debian-only packaging fix bumps only
# __debian/control, not the RPM spec. But a new *version* is a fresh package
# on every distro, so -v without -r resets every release counter to 1. If one
# distro needs to diverge afterward (e.g. a distro-specific repackaging issue),
# edit that one file's release field by hand — this script only sets the
# same value everywhere.
#
# The Go version is not embedded in any __*/ file: __debian/install_build_deps.sh
# and __redhat/rpmbuild-deps.sh read go.version when they run, so updating
# go.version (and the "goversion" field) is all that -g has to do.
#
# src/cmd/root.go needs no edit: it reads its version via -ldflags -X at
# build time from whichever file below is authoritative for that build.

set -eu

SCRIPT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
cd "$SCRIPT_DIR"

#PKG_TARGETS="__alpine/APKBUILD __archlinux/PKGBUILD __debian/control __redhat/vmman4.spec"
PKG_TARGETS="__alpine/APKBUILD __debian/control __redhat/vmman4.spec"
JSON_FILE="vmman4.json"
GO_FILE="go.version"

# --- output helpers ------------------------------------------------------
# Colours only when the stream is a terminal (and NO_COLOR is unset), so logs
# and pipes don't get escape sequences. Errors: red. Warnings/notices: yellow.

if [ -z "${NO_COLOR-}" ] && [ -t 2 ]; then
  ERR_RED=$(printf '\033[31m')
  ERR_YEL=$(printf '\033[33m')
  ERR_RST=$(printf '\033[0m')
else
  ERR_RED=''
  ERR_YEL=''
  ERR_RST=''
fi
if [ -z "${NO_COLOR-}" ] && [ -t 1 ]; then
  OUT_YEL=$(printf '\033[33m')
  OUT_RST=$(printf '\033[0m')
else
  OUT_YEL=''
  OUT_RST=''
fi

err() { printf '%sERROR: %s%s\n' "$ERR_RED" "$1" "$ERR_RST" >&2; }
warn() { printf '%sWARNING: %s%s\n' "$ERR_YEL" "$1" "$ERR_RST" >&2; }
notice() { printf '%sNOTICE: %s%s\n' "$OUT_YEL" "$1" "$OUT_RST"; }
die() {
  err "$1"
  exit "${2:-1}"
}

print_usage() {
  cat <<EOF
Usage: $0 [-v <version>] [-r <release>] [-g <goversion>]

At least one flag is required; a flag left out keeps its current value.

  -v <version>    New package version (e.g. 2.3.4). Updates "versionnumber" in
                  vmman4.json and every __*/ packaging file. Must
                  start with a digit and contain only [A-Za-z0-9._~-] (no "v"
                  prefix). Unless -r is also given, the release number is
                  reset to 1.
  -r <release>    New release number (digits only). Updates
                  "defaultreleasenumber" in vmman4.json and every
                  __*/ packaging file. With -v too, this value wins over the
                  automatic reset to 1.
  -g <goversion>  New Go version (e.g. 1.27.4). Updates "goversion" in
                  vmman4.json and the go.version file.
  -h, --help      Show this help.

Examples:
  $0 -v 2.3.4              new version, release reset to 1
  $0 -r 2                  release-only bump, version untouched
  $0 -g 1.27.5             Go version only
  $0 -v 2.3.4 -g 1.27.5    new version (release -> 1) and new Go version
EOF
}

# --- argument handling -------------------------------------------------

case "${1-}" in
  --help)
    print_usage
    exit 0
    ;;
esac

HAVE_V=''
HAVE_R=''
HAVE_G=''
NEWVERSION=''
NEWREL=''
NEWGO=''

while getopts ':v:r:g:h' opt; do
  case "$opt" in
    v)
      HAVE_V=1
      NEWVERSION="$OPTARG"
      ;;
    r)
      HAVE_R=1
      NEWREL="$OPTARG"
      ;;
    g)
      HAVE_G=1
      NEWGO="$OPTARG"
      ;;
    h)
      print_usage
      exit 0
      ;;
    :)
      die "option -$OPTARG requires an argument. See '$0 -h'." 2
      ;;
    *)
      die "unknown option -$OPTARG. See '$0 -h'." 2
      ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$HAVE_V$HAVE_R$HAVE_G" ]; then
  warn "no flag provided, nothing to do. Give at least one of -v, -r, -g (see '$0 -h')."
  if [ "$#" -gt 0 ]; then
    warn "a bare version argument is no longer accepted; use: $0 -v $1"
  fi
  exit 2
fi

[ "$#" -eq 0 ] || die "unexpected argument(s): $*. See '$0 -h'." 2

# --- validation ------------------------------------------------------------

if [ -n "$HAVE_V" ]; then
  case "$NEWVERSION" in
    -*)
      die "'$NEWVERSION' looks like an option, not a version (-v). See '$0 -h'." 2
      ;;
    [0-9]*) ;;
    *)
      die "version (-v) must start with a digit (e.g. 1.3.2), got '$NEWVERSION'"
      ;;
  esac
  case "$NEWVERSION" in
    *[!A-Za-z0-9._~-]*)
      die "'$NEWVERSION' contains characters not expected in a version string"
      ;;
  esac
fi

if [ -n "$HAVE_R" ]; then
  case "$NEWREL" in
    '' | *[!0-9]*)
      die "release number (-r) must be digits only (e.g. 2), got '$NEWREL'"
      ;;
  esac
fi

if [ -n "$HAVE_G" ]; then
  case "$NEWGO" in
    [0-9]*) ;;
    *)
      die "Go version (-g) must start with a digit (e.g. 1.27.4), got '$NEWGO'"
      ;;
  esac
  case "$NEWGO" in
    *[!A-Za-z0-9.]*)
      die "'$NEWGO' contains characters not expected in a Go version (e.g. 1.27.4, 1.28rc1)"
      ;;
  esac
fi

# A new version is a fresh package on every distro: release goes back to 1,
# unless the caller explicitly asked for another value with -r.
if [ -n "$HAVE_V" ] && [ -z "$HAVE_R" ]; then
  NEWREL=1
  notice "-v given without -r: release number reset to 1 (defaultreleasenumber and every distro's release counter)."
fi

# Do the packaging files need touching at all? (anything but a -g-only run)
BUMP_PKG=''
if [ -n "$HAVE_V$HAVE_R" ]; then
  BUMP_PKG=1
fi

# --- preconditions -------------------------------------------------------

TARGETS="$JSON_FILE"
[ -z "$BUMP_PKG" ] || TARGETS="$TARGETS $PKG_TARGETS"
[ -z "$HAVE_G" ] || TARGETS="$TARGETS $GO_FILE"

for f in $TARGETS; do
  [ -f "$f" ] || die "expected to find $f here — is this script still at the project root?"
done

# Refuse to run against a dirty tree for the files we're about to rewrite,
# so any mistake here (or in a future edit to this script) is always a
# trivial 'git checkout -- <file>' away from being undone.
# shellcheck disable=SC2086  # TARGETS is a deliberate space-separated list
dirty=$(git status --porcelain -- $TARGETS 2>/dev/null || true)
if [ -n "$dirty" ]; then
  die "the following files have uncommitted changes; commit or stash them first so this bump stays trivially revertible:
$dirty"
fi

# --- patch ----------------------------------------------------------------

# patch <file> <sed-script>  — writes via a temp file + mv so this works
# identically under GNU, BSD, and BusyBox sed (their -i flags don't agree),
# then restores the original file's permissions (mv onto a freshly created
# temp file would otherwise silently drop the executable bit).
patch() {
  file="$1"
  script="$2"
  mode=$(stat -c '%a' "$file" 2>/dev/null || stat -f '%Lp' "$file")
  sed "$script" "$file" > "$file.bumptmp"
  mv "$file.bumptmp" "$file"
  chmod "$mode" "$file"
}

# verify <file> <expected-literal-text>  — fails loudly instead of silently
# no-op'ing if a file's format ever drifts out from under this script's sed
# patterns.
verify() {
  grep -qF -- "$2" "$1" || die "expected to find '$2' in $1 after patching — its format may have changed; check 'git diff' before re-running."
}

# verify_line <file> <expected-full-line>  — same, but the whole line must match,
# so that e.g. 'pkgrel=1' is not satisfied by a leftover 'pkgrel=10'.
verify_line() {
  grep -qxF -- "$2" "$1" || die "expected to find the line '$2' in $1 after patching — its format may have changed; check 'git diff' before re-running."
}

echo "Bumping vmman4:"
[ -z "$HAVE_V" ] || echo "  package version -> $NEWVERSION"
[ -z "$BUMP_PKG" ] || echo "  release number  -> $NEWREL"
[ -z "$HAVE_G" ] || echo "  Go version      -> $NEWGO"

# --- vmman4.json ----------------------------------------------

if [ -n "$HAVE_V" ]; then
  patch "$JSON_FILE" \
    "s/\"versionnumber\": *\"[^\"]*\"/\"versionnumber\": \"$NEWVERSION\"/"
  verify "$JSON_FILE" "\"versionnumber\": \"$NEWVERSION\""
fi
if [ -n "$BUMP_PKG" ]; then
  patch "$JSON_FILE" \
    "s/\"defaultreleasenumber\": *\"[^\"]*\"/\"defaultreleasenumber\": \"$NEWREL\"/"
  verify "$JSON_FILE" "\"defaultreleasenumber\": \"$NEWREL\""
fi
if [ -n "$HAVE_G" ]; then
  patch "$JSON_FILE" \
    "s/\"goversion\": *\"[^\"]*\"/\"goversion\": \"$NEWGO\"/"
  verify "$JSON_FILE" "\"goversion\": \"$NEWGO\""
fi

# --- go.version -------------------------------------------------------------

if [ -n "$HAVE_G" ]; then
  patch "$GO_FILE" "1s/.*/$NEWGO/"
  verify_line "$GO_FILE" "$NEWGO"
fi

# --- __*/ packaging files ---------------------------------------------------

if [ -n "$BUMP_PKG" ]; then
  # Alpine
  if [ -n "$HAVE_V" ]; then
    patch __alpine/APKBUILD "s/^pkgver=.*/pkgver=$NEWVERSION/"
    verify_line __alpine/APKBUILD "pkgver=$NEWVERSION"
  fi
  patch __alpine/APKBUILD "s/^pkgrel=.*/pkgrel=$NEWREL/"
  verify_line __alpine/APKBUILD "pkgrel=$NEWREL"

  # ArchLinux
  #if [ -n "$HAVE_V" ]; then
  #  patch __archlinux/PKGBUILD "s/^pkgver=.*/pkgver=$NEWVERSION/"
  #  verify_line __archlinux/PKGBUILD "pkgver=$NEWVERSION"
  #fi
  #patch __archlinux/PKGBUILD "s/^pkgrel=.*/pkgrel=$NEWREL/"
  #verify_line __archlinux/PKGBUILD "pkgrel=$NEWREL"

  # Red Hat
  if [ -n "$HAVE_V" ]; then
    patch __redhat/vmman4.spec "s/^%define _version .*/%define _version $NEWVERSION/"
    verify_line __redhat/vmman4.spec "%define _version $NEWVERSION"
  fi
  patch __redhat/vmman4.spec "s/^%define _rel .*/%define _rel $NEWREL/"
  verify_line __redhat/vmman4.spec "%define _rel $NEWREL"

  # Debian: control's "Version:" is <upstream>-<debian revision>, so both flags
  # land on the same line. For a release-only bump, keep whatever upstream
  # version the file already carries.
  if [ -n "$HAVE_V" ]; then
    debver="$NEWVERSION"
  else
    debver=$(sed -n 's/^Version: \(.*\)-[^-]*$/\1/p' __debian/control)
    [ -n "$debver" ] || die "could not find a 'Version: <upstream>-<revision>' line in __debian/control — its format may have changed."
  fi
  patch __debian/control "s/^Version: .*/Version: $debver-$NEWREL/"
  verify_line __debian/control "Version: $debver-$NEWREL"
fi

echo
echo "Changed:"
grep -H '"versionnumber"\|"defaultreleasenumber"\|"goversion"' "$JSON_FILE"
if [ -n "$HAVE_G" ]; then
  grep -H '' "$GO_FILE"
fi
if [ -n "$BUMP_PKG" ]; then
  grep -H '^pkgver=\|^pkgrel=' __alpine/APKBUILD
  #grep -H '^pkgver=\|^pkgrel=' __archlinux/PKGBUILD
  grep -H '^Version:' __debian/control
  grep -H '^%define _version\|^%define _rel' __redhat/vmman4.spec
fi
echo
if [ -n "$BUMP_PKG" ]; then
  echo "Nothing else to edit by hand — src/cmd/root.go picks this up via -ldflags at build time."
fi
if [ -n "$HAVE_G" ]; then
  echo "Go version: the builder scripts read go.version when they run; nothing else to edit."
fi
