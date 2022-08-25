#!/usr/bin/env bash
set -Eeuo pipefail

DEFAULT_TUF_REPO_URL=${DEFAULT_TUF_REPO_URL:-"http://localhost:8080"}
REQUEST_ID=${REQUEST_ID:-$(uuidgen | tr "[:upper:]" "[:lower:]")}

usage() {
  cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-s server_uri] [-n namespace] [-r repo_uuid]

Returns list available TUF repositories.

Available options:

-h, --help      Print this help and exit
-v, --verbose   Print script debug info
-s, --server    TUF server bse URI (by default ${DEFAULT_TUF_REPO_URL})
-o, --offset    offset of request
-l, --limit     max number records in response
EOF
  exit
}

setup_colors() {
  if [[ -t 2 ]] && [[ -z "${NO_COLOR-}" ]] && [[ "${TERM-}" != "dumb" ]]; then
    NOFORMAT='\033[0m' RED='\033[0;31m' GREEN='\033[0;32m' ORANGE='\033[0;33m' BLUE='\033[0;34m' PURPLE='\033[0;35m' CYAN='\033[0;36m' YELLOW='\033[1;33m'
  else
    NOFORMAT='' RED='' GREEN='' ORANGE='' BLUE='' PURPLE='' CYAN='' YELLOW=''
  fi
}

msg() {
  echo >&2 -e "${1-}"
}

parse_params() {
  TUF_REPO_URL=${DEFAULT_TUF_REPO_URL}
  OFFSET=0
  LIMIT=5
  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) set -x ;;
    --no-color) NO_COLOR=1 ;;
    -s | --server) # example named parameter
      TUF_REPO_URL="${2-}"
      shift
      ;;
    -o | --offset) # example named parameter
      OFFSET="${2-}"
      shift
      ;;
    -l | --limit) # example named parameter
      LIMIT="${2-}"
      shift
      ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac
    shift
  done

  # check required params and arguments
  [[ -z "${TUF_REPO_URL-}" ]] && die "Missing required parameter: server"
  [[ -z "${LIMIT-}" ]] && die "Missing required parameter: limit"

  return 0
}

parse_response() {
  local response=${1}
  local http_code
  http_code=$(tail -n1 <<< "$response")  # get the last line
  local content
  content=$(sed '$d' <<< "$response")   # get all except the first and last lines
  local head=true
  local header=""
  local response_body=""
  while read -r line; do
    if $head; then
      if [[ $line = $'\r' ]]; then
          head=false
      else
          header="$header"$'\n\t'"$line"
      fi
    else
      response_body="$response_body"$'\n'"$line"
    fi
  done < <(echo "$content")

  if [[ "${http_code}" -ne 200 ]] ; then
    msg "${RED}HTTP response code: ${NOFORMAT}${http_code}"
  else
    msg "${BLUE}HTTP response code: ${NOFORMAT}${http_code}"
  fi
  msg "${RED}Headers:${NOFORMAT}$header"
  echo "${response_body}"
}

parse_params "$@"
setup_colors

URL="${TUF_REPO_URL}/api/v1/repos?offset=${OFFSET}&limit=${LIMIT}"

msg "${GREEN}RequestID:${NOFORMAT} ${REQUEST_ID}"
msg "${GREEN}URL      :${NOFORMAT} ${URL}"

response=$(curl -si -w "%{http_code}" \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: ${REQUEST_ID}" \
  "${URL}")

parse_response "${response}"
