#!/usr/bin/env bash
#
HELPER=helper_go
echo "Launch helper${1:-} binary..."
# need this because binary will detect exit in parent shell and die otherwise
screen -S "session${1:-}" -d -m "$HELPER_DIR/$HELPER" "${HELPER}${1:-}"
