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
uuid=${1:-$(uuidgen | tr "[:upper:]" "[:lower:]")}
URL="${TUF_REPO_URL}/api/v1/root/${uuid}"
body="{\"keyType\":\"ed25519\"}"
curl -H "Content-Type: application/json" -X "POST" --data "${body}" "${URL}"
