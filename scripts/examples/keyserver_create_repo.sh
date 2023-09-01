#!/usr/bin/env bash

set -eo pipefail
DEBUG=${DEBUG:-false}
[[ ${DEBUG} = true ]] && set -x

DEFAULT_TUF_REPO_URL=${DEFAULT_TUF_REPO_URL:-"http://localhost:8080"}
DEFAULT_NAMESPACE=${DEFAULT_NAMESPACE:-"default"}
REQUEST_ID=${REQUEST_ID:-$(uuidgen | tr "[:upper:]" "[:lower:]")}
ALLOWED_KEY_TYPES=("RSA" "ED25519")
DEFAULT_KEY_TYPE="RSA"


usage() {
  cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-s server_uri] [-n namespace] [-r repo_uuid]

Creates TUF repository with provided parameters.

Available options:

-h, --help      Print this help and exit
-v, --verbose   Print script debug info
-s, --server    TUF server bse URI (by default ${DEFAULT_TUF_REPO_URL})
-n, --namespace TUF repo namespace (by default ${DEFAULT_NAMESPACE})
-r, --repo      TUF repo UUID (generated in runtime by default)
-k, --key       TUF repo key type(supported types: ${ALLOWED_KEY_TYPES[*]})
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

die() {
  local msg=$1
  local code=${2-1} # default exit status 1
  msg "$msg"
  exit "$code"
}

parse_params() {
  TUF_REPO_URL=${DEFAULT_TUF_REPO_URL}
  NAMESPACE=${DEFAULT_NAMESPACE}
  KEY_TYPE=${DEFAULT_KEY_TYPE}
  REPO_ID="$(uuidgen | tr "[:upper:]" "[:lower:]")"
  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) set -x ;;
    --no-color) NO_COLOR=1 ;;
    -s | --server)
      TUF_REPO_URL="${2-}"
      shift
      ;;
    -n | --namespace)
      NAMESPACE="${2-}"
      shift
      ;;
    -k | --key)
      KEY_TYPE="${2-}"
      shift
      ;;
    -r | --repo)
      REPO_ID="${2-}"
      shift
      ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac
    shift
  done

  # check required params and arguments
  [[ -z "${TUF_REPO_URL-}" ]] && die "Missing required parameter: server"
  [[ -z "${NAMESPACE-}" ]] && die "Missing required parameter: namespace"
  [[ -z "${KEY_TYPE-}" ]] && die "Missing required parameter: key type"
  [[ "${ALLOWED_KEY_TYPES[*]}" =~ (^|[[:space:]])"$KEY_TYPE"($|[[:space:]]) ]] || die "parameter ${YELLOW}key${NOFORMAT} has incorrect value"
  [[ -z "${REPO_ID-}" ]] && die "Missing required parameter: repo ID"

  return 0
}

parse_response() {
  local response=${1}
  local http_code
  local content
  local head=true
  local header=""
  local response_body=""
  http_code=$(tail -c4 <<< "$response")  # get the last line
  content=$(sed '1d' <<< "$response")   # get all except the first and last lines

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
  done < <(echo "${content:0:${#content}-3}")

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

URL="${TUF_REPO_URL}/api/v1/root/${REPO_ID}"
BODY="{\"keyType\":\"${KEY_TYPE}\", \"threshold\": 2}"
echo "${BODY}"

msg "${GREEN}RequestID:${NOFORMAT} ${REQUEST_ID}"
msg "${GREEN}URL      :${NOFORMAT} ${KEY_TYPE}"
msg "${GREEN}KEY      :${NOFORMAT} ${URL}"

response=$(curl -si -w "%{http_code}" \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: ${REQUEST_ID}" \
  -H "x-ats-tuf-force-sync: keys" \
  --data "${BODY}" \
  -X "POST" \
  "${URL}")

parse_response "${response}"
