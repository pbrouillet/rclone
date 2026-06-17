# rclone OneDrive OOB handler

A small, independent helper for rclone's OneDrive **`web_auth`** flow.

## Why

When you authenticate OneDrive for Business / SharePoint with `web_auth`,
rclone uses a Microsoft first-party app and a browser sign-in. After you
sign in, Microsoft redirects the browser to:

```text
urn:ietf:wg:oauth:2.0:oob?code=<authorization code>
```

Some Linux desktops can't render the `urn:` scheme and hand the URI to
`xdg-open`, where nothing is listening — so the code is hard to retrieve. This
handler registers itself for the `urn:` scheme, extracts the `code`, and
copies it to your clipboard with a desktop notification.

## Requirements

- Go (to build the handler).
- A clipboard tool: `xclip`, `xsel`, or `wl-clipboard`.
- Optional: `notify-send` (from `libnotify-bin`) for desktop notifications.

## Install

```sh
cd contrib/onedrive-oob-handler
./install.sh
```

This builds `onedrive-oob-handler` into `~/.local/bin`, installs a `.desktop`
entry, and registers it as the `x-scheme-handler/urn` handler.

Test it:

```sh
~/.local/bin/onedrive-oob-handler 'urn:ietf:wg:oauth:2.0:oob?code=test123'
# "test123" is now on your clipboard
```

## Use

Configure your OneDrive remote with `web_auth = true` (this uses the OOB
redirect). Then run `rclone config` (or `rclone config reconnect <remote>:`),
open the sign-in URL, authenticate, and the code is copied to your clipboard
automatically. Paste it into rclone's "Verification code" prompt.

The first time the browser triggers the handler it may ask you to confirm
opening the external application — allow it.

## Uninstall

```sh
./install.sh --uninstall
```

## Notes / caveats

- This claims the **entire** `urn:` URI scheme for your user. In practice few
  applications use `urn:` links, but be aware of it.
- The handler is intentionally standalone and does not depend on a running
  rclone process; it only puts the code on the clipboard.
