# Wi-Fi recovery and validation

This note records the supported recovery procedure for appliance Wi-Fi profiles and the evidence behind the v1.6.1 credential-validation hardening.

## Normal connection path

The appliance setup UI sends the selected SSID and password to the loopback-only setup API. The backend passes a non-empty password to `nmcli --ask` on standard input so credentials are not exposed in process arguments, API responses, or application logs.

For WPA/WPA2 Personal, a non-empty credential must be either:

- an 8–63 character passphrase; or
- exactly 64 hexadecimal characters representing a raw PSK.

Open networks continue to use an empty password.

## Stale or malformed NetworkManager profiles

Do not automatically delete an existing NetworkManager profile merely because an authentication attempt fails. A failed attempt can be caused by a mistyped password, and preserving known-good profiles allows NetworkManager to fall back to another configured network.

If a profile is suspected to be malformed:

1. Preserve the profile file before changing it.
2. Inspect only non-secret metadata and the PSK length; do not print or log the PSK.
3. Delete the affected connection through `nmcli connection delete "<SSID>"`.
4. Recreate it through the MasjidFrame setup UI.
5. Verify the resulting connection and reboot once to confirm autoconnect.

Example length-only inspection:

```sh
sudo python3 - <<'PY'
from pathlib import Path

path = Path("/etc/NetworkManager/system-connections/<SSID>.nmconnection")
psk = next((line[4:] for line in path.read_text().splitlines() if line.startswith("psk=")), None)
print("PSK stored:", psk is not None)
if psk is not None:
    print("PSK length:", len(psk))
PY
```

Never include the PSK itself in bug reports, logs, screenshots, shell history, or diagnostic archives.

## v1.6.0 field investigation

During Raspberry Pi 3 appliance testing on 26 September 2026, one saved profile contained an invalid 329-character PSK and a stale temporary profile contained a 66-character value. The source of those malformed values was not reproduced.

The current connection path was tested independently at each boundary:

- direct `nmcli --ask` stored an exact 17-character controlled test password;
- the loopback setup API stored the same controlled password unchanged;
- the touchscreen setup UI stored the same controlled password unchanged;
- the real 11-character WPA2-Personal credential connected successfully; and
- after reboot, NetworkManager automatically reconnected to the saved SSID, obtained DHCP configuration, and reached `CONNECTED_GLOBAL`.

The evidence therefore did not justify automatic profile deletion or a change to the `nmcli --ask` transport. v1.6.1 instead adds credential validation at both the API and NetworkManager boundaries so malformed non-empty PSKs cannot be submitted by those paths.

## Release validation

For v1.6.1, verify:

- invalid 1–7 character passwords are rejected without invoking NetworkManager;
- non-hexadecimal 64-character passwords are rejected;
- 8-character and 63-character passphrases are accepted;
- 64-character hexadecimal raw PSKs are accepted;
- open networks still permit an empty password;
- visible and hidden WPA2-Personal networks still connect from the touchscreen;
- a failed replacement attempt does not remove an unrelated known-good profile; and
- a successful replacement reconnects automatically after reboot.
