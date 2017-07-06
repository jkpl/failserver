#!/bin/sh

log_json() {
    if hash jq 2>/dev/null; then
        jq < "$1"
    else
        cat "$1"
    fi
}

print_help() {
    cat <<EOF >&2
./run.sh [command]

Commands:
load  = run the load test and exit
help  = print this help
build = (re)build local images

No command:
Run all the services
EOF
}

run_load_test() {
    docker-compose -f docker-compose.load.yml up --abort-on-container-exit
    log_json ./results/results.json
}

run_service() {
    docker-compose -f docker-compose.yml up
}

build_images() {
    docker-compose -f docker-compose.yml \
                   -f docker-compose.load.yml \
                   build
}

main() {
    local command="$1"
    shift
    case "$command" in
        help) print_help "$@" ;;
        load) run_load_test "$@" ;;
        build) build_images "$@" ;;
        *) run_service "$@" ;;
    esac
}

main "$@"
