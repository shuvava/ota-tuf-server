#!/usr/bin/env bash
set -Eeuo pipefail

BRed='\033[1;31m'
BGreen='\033[1;32m'
BWhite='\033[1;37m'
Color_Off='\033[0m'
print() {
  color=${2:-$BWhite}
  echo -e "${color}$1${Color_Off}"
}

TUF_REPO_URL=${TUF_REPO_URL:-"http://localhost:8080"}
URL="${TUF_REPO_URL}/healthz"
print "  - Checking TUF repo is alive at ${URL}..."
curl -i "${URL}"
URL="${TUF_REPO_URL}/readyz"
print "  - Checking TUF repo is readiness at ${URL}..."
curl -i "${URL}"
