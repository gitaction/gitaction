#!/bin/bash
set -e
set -u

GA_ENV=${GA_ENV:-local}
PARAMETER=

BASE_DIR=$(dirname $0)
SCRIPT_PATH="$( cd "${BASE_DIR}" && pwd -P )"

exit_err() {
  echo "ERROR: ${1}" >&2
  exit 1
}

# Usage: -h, --help
# Description: Show help text
option_help() {
  printf "Usage: %s [options...] COMMAND <parameter> \n\n" "${0}"
  printf "Default command: --help\n\n"

  echo "Options:"
  grep -e '^[[:space:]]*# Usage:' -e '^[[:space:]]*# Description:' -e '^option_.*()[[:space:]]*{' "${0}" | while read -r usage; read -r description; read -r option; do
    if [[ ! "${usage}" =~ Usage ]] || [[ ! "${description}" =~ Description ]] || [[ ! "${option}" =~ ^option_ ]]; then
      exit_err "Error generating help text."
    fi
    printf " %-32s %s\n" "${usage##"# Usage: "}" "${description##"# Description: "}"
  done

  printf "\n"
  echo "Commands:"
  grep -e '^[[:space:]]*# Command Usage:' -e '^[[:space:]]*# Command Description:' -e '^command_.*()[[:space:]]*{' "${0}" | while read -r usage; read -r description; read -r command; do
    if [[ ! "${usage}" =~ Usage ]] || [[ ! "${description}" =~ Description ]] || [[ ! "${command}" =~ ^command_ ]]; then
      exit_err "Error generating help text."
    fi
    printf " %-32s %s\n" "${usage##"# Command Usage: "}" "${description##"# Command Description: "}"
  done
}

# Usage: -p, --prod
# Description: Set the GA env to production (default local)
option_prod() {
  GA_ENV=prod
  load_env
}

# Command Usage: run
# Command Description: Gradle project bootRun
command_run() {
  echo "run..."
}

# Command Usage: test <unit|function|integration>
# Command Description: Run unit|function|integration test
command_test() {
  echo "test..."
}

# Command Usage: up
# Command Description: Docker compose start up brand new database container
command_up() {
  docker-compose -f ./docker-compose.yml up -d
  check_msg "Docker container up"
  echo "Database provision..."
  ./mssql-provision/up.sh
}

# Command Usage: down
# Command Description: Docker compose remove database instance totally
command_down() {
  docker-compose -f ./docker-compose.yml down
}

# Command Usage: push
# Command Description: Docker push image to AWS ECR
command_push() {
  echo "push..."
}

# Command Usage: deploy
# Command Description: Deploy application
function command_deploy() {
  echo ">>> deploy $APP_NAME <<<"
}

command_build() {
  echo "build..."
}

check_msg() {
  printf "\xE2\x9C\x94 ${1}\n"
}

main() {
  [[ -z "${@}" ]] && eval set -- "--help"

  local theCommand=

  set_command() {
    [[ -z "${theCommand}" ]] || exit_err "Only one command at a time!"
    theCommand="${1}"
  }

  while [[ "$#" -gt 0 ]]; do
    case "$1" in

      --help|-h)
        option_help
        exit 0
        ;;

      --prod|-p)
        option_prod
        ;;

      run|test|clean|up|build|deploy|down|console|push)
        set_command "${1}"
        ;;

      *)
        PARAMETER="${1}"
        ;;
    esac

    shift 1
  done

  [[ ! -z "${theCommand}" ]] || exit_err "Command not found!"

  case "${theCommand}" in
    run) command_run;;
    test) command_test;;
    clean) command_clean;;
    up) command_up;;
    build) command_build;;
    deploy) command_deploy;;
    down) command_down;;
    console) command_console;;
    push) command_push;;

    *) option_help; exit 1;;
  esac
}

main "${@-}"
