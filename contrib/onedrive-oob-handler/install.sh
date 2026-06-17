#!/usr/bin/env bash
#
# install.sh - build and register the rclone OneDrive OOB scheme handler.
#
# This builds the onedrive-oob-handler binary and registers it as the handler
# for the urn: URI scheme, so that when the OneDrive web_auth (oob) flow
# redirects the browser to urn:ietf:wg:oauth:2.0:oob?code=..., the code is
# extracted and copied to the clipboard.
#
# Usage:
#   ./install.sh            install the handler
#   ./install.sh --uninstall remove the handler
#
set -euo pipefail

APP_ID="rclone-onedrive-oob"
DESKTOP_FILE="${APP_ID}.desktop"
SCHEME="x-scheme-handler/urn"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="${XDG_BIN_HOME:-$HOME/.local/bin}"
APPS_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/applications"
BIN_PATH="${BIN_DIR}/onedrive-oob-handler"

uninstall() {
	echo "Removing ${APPS_DIR}/${DESKTOP_FILE}"
	rm -f "${APPS_DIR}/${DESKTOP_FILE}"
	echo "Removing ${BIN_PATH}"
	rm -f "${BIN_PATH}"
	if command -v update-desktop-database >/dev/null 2>&1; then
		update-desktop-database "${APPS_DIR}" >/dev/null 2>&1 || true
	fi
	echo "Uninstalled. You may want to reset the urn: handler with:"
	echo "  xdg-mime default '' ${SCHEME}"
}

if [[ "${1:-}" == "--uninstall" ]]; then
	uninstall
	exit 0
fi

echo "Building handler..."
mkdir -p "${BIN_DIR}" "${APPS_DIR}"
( cd "${SCRIPT_DIR}" && go build -o "${BIN_PATH}" . )
echo "Installed binary at ${BIN_PATH}"

echo "Installing desktop entry..."
sed "s|@HANDLER_PATH@|${BIN_PATH}|g" \
	"${SCRIPT_DIR}/${DESKTOP_FILE}.in" > "${APPS_DIR}/${DESKTOP_FILE}"

if command -v update-desktop-database >/dev/null 2>&1; then
	update-desktop-database "${APPS_DIR}" >/dev/null 2>&1 || true
fi

if command -v xdg-mime >/dev/null 2>&1; then
	xdg-mime default "${DESKTOP_FILE}" "${SCHEME}"
	echo "Registered ${DESKTOP_FILE} as the ${SCHEME} handler."
else
	echo "WARNING: xdg-mime not found; register manually:"
	echo "  xdg-mime default ${DESKTOP_FILE} ${SCHEME}"
fi

if ! command -v xclip >/dev/null 2>&1 \
	&& ! command -v xsel >/dev/null 2>&1 \
	&& ! command -v wl-copy >/dev/null 2>&1; then
	echo "WARNING: no clipboard tool found. Install one of: xclip, xsel, wl-clipboard"
fi

cat <<EOF

Done. Test it with:
  ${BIN_PATH} 'urn:ietf:wg:oauth:2.0:oob?code=test123'
(the string 'test123' should now be on your clipboard)

To use it, configure your OneDrive remote with:
  web_auth = true
EOF
