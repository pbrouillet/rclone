// Command onedrive-oob-handler receives an OAuth2 out-of-band (OOB) redirect
// URI and copies the authorization code it contains to the clipboard.
//
// When rclone's OneDrive "web_auth" flow signs in, Microsoft redirects the
// browser to
//
//	urn:ietf:wg:oauth:2.0:oob?code=<authorization code>
//
// Some Linux desktops cannot render the urn: scheme and instead hand the URI
// to an external handler (via xdg-open). This program is that handler: it
// extracts the code, copies it to the clipboard and shows a desktop
// notification so the code can be pasted straight into rclone's prompt.
//
// It is intentionally independent of rclone itself. See README.md for how to
// register it as the urn: scheme handler.
package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
)

// extractCode parses an OAuth redirect URI and returns the authorization code.
//
// If the URI carries an OAuth error (error / error_description) it is returned
// as a Go error instead.
func extractCode(rawURI string) (string, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return "", fmt.Errorf("could not parse redirect URI %q: %w", rawURI, err)
	}
	q := u.Query()
	if e := q.Get("error"); e != "" {
		if desc := q.Get("error_description"); desc != "" {
			return "", fmt.Errorf("authorization failed: %s: %s", e, desc)
		}
		return "", fmt.Errorf("authorization failed: %s", e)
	}
	code := q.Get("code")
	if code == "" {
		return "", fmt.Errorf("no authorization code found in redirect URI %q", rawURI)
	}
	return code, nil
}

// notify shows a best-effort desktop notification. Any failure (e.g.
// notify-send not installed) is ignored.
func notify(title, body string) {
	path, err := exec.LookPath("notify-send")
	if err != nil {
		return
	}
	_ = exec.Command(path, title, body).Run()
}

func run(args []string) error {
	if len(args) != 1 || args[0] == "" {
		return fmt.Errorf("usage: onedrive-oob-handler <redirect-uri>")
	}

	code, err := extractCode(args[0])
	if err != nil {
		notify("rclone OneDrive sign-in failed", err.Error())
		return err
	}

	if clipboard.Unsupported {
		notify("rclone OneDrive", "Clipboard unavailable - copy the code manually")
		// Still print the code so it can be copied from a terminal.
		fmt.Println(code)
		return fmt.Errorf("clipboard not supported on this system (install xclip, xsel or wl-clipboard); code: %s", code)
	}

	if err := clipboard.WriteAll(code); err != nil {
		notify("rclone OneDrive", "Failed to copy code to clipboard")
		fmt.Println(code)
		return fmt.Errorf("failed to copy code to clipboard: %w (code: %s)", err, code)
	}

	notify("rclone OneDrive", "Authorization code copied to clipboard - paste it into rclone")
	fmt.Fprintln(os.Stderr, "Authorization code copied to clipboard. Paste it into rclone's prompt.")
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "onedrive-oob-handler:", err)
		os.Exit(1)
	}
}
