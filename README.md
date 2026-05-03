![Dynamic XML Badge](https://img.shields.io/badge/dynamic/xml?url=https%3A%2F%2Fspetix.github.io%2Fotp_blocket%2Fcoverage.xml&query=round(100*%2F%2Fcoverage%2F%40line-rate)&suffix=%25&style=plastic&label=Code%20Coverage)

It provides a blocklet that shows your favourite's site/application otp.

## Installation

* Copy blocklet to blocklets directory (e.g. `$HOME/.config/i3blocks/blocklets`)
* Add blocklet to i3blocks config (e.g. `$HOME/.config/i3blocks/config`)

Here's a sample configuration, `${BLOCKLETS_DIR}` has to be set in i3 configuration or env:
```ini
[otp_blocklet]
command = ${BLOCKLETS_DIR}/otp_blocklet
interval = 1
fg_color = #ff0000
bg_color =
label = 🎄
```

## Downloads

* [linux amd64 version](otp_bloclet_linux_amd64)
* [linux arm64 version](otp_blocklet_linux_arm64)


# Build details

* [Coverage](code-coverage-results.md)

# Folders description

`internal/` contains the blocklet's core OTP logic and runtime helper packages.

- `internal/onetimepass`
  - Defines `PassCode`, the OTP result object with code text, expiration, remaining lifetime, and error state.
  - Provides validation for generated one-time passwords.
- `internal/provider`
  - Implements `OtpProvider`, which generates and refreshes TOTP codes from a configured key.
  - Uses a background refresher and atomic code storage so the current OTP is always available.
- `internal/render`
  - Contains the render adapter that exposes OTP values to the bar/blocklet framework.
  - Formats short/long text, label, and colors for display.
- `internal/secretmanager`
  - Manages loading and saving OTP key URLs from disk.
  - Handles path expansion (`~`) and JSON persistence for secret storage.

## License

MIT

## Author

[spetix](https://github.com/spetix)

