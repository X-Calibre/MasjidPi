# MasjidPi Release Candidate Soak Monitoring

This document describes a reusable lightweight monitoring service for release-candidate and post-release validation on Raspberry Pi hardware.

The monitor is intended for temporary 24–48 hour soak testing. It records service health, process resource usage, system load, memory, disk usage, Raspberry Pi temperature/throttling state, Listen status, MasjidBoard status summaries, and recent MasjidPi warnings/errors.

The monitor runs under `systemd`, so it continues after an SSH or terminal session disconnects and restarts automatically after a reboot.

## What it records

Every five minutes the monitor writes a sample containing:

- MasjidPi version
- uptime and load average
- memory usage
- root filesystem usage
- Raspberry Pi temperature and throttling state
- `masjidpi.service` state and restart count
- `masjidpi-display.service` state and restart count
- CPU, memory, RSS and elapsed time for MasjidPi, mpv, Cog and WPE processes
- `/api/listen/status`
- a compact `/api/masjidboard/status` summary
- recent MasjidPi warnings and errors

The log is written to:

```text
/var/log/masjidpi-rc-monitor/soak.log
```

## 1. Install the monitoring script

Create `/usr/local/bin/masjidpi-rc-monitor.sh`:

```bash
sudo tee /usr/local/bin/masjidpi-rc-monitor.sh >/dev/null <<'EOF'
#!/usr/bin/env bash

set -u

LOG_DIR="/var/log/masjidpi-rc-monitor"
LOG_FILE="${LOG_DIR}/soak.log"

mkdir -p "$LOG_DIR"

while true; do
    {
        echo "============================================================"
        date --iso-8601=seconds

        echo
        echo "=== VERSION ==="
        curl -fsS http://localhost:8080/api/version 2>/dev/null || echo "version endpoint unavailable"

        echo
        echo "=== UPTIME / LOAD ==="
        uptime

        echo
        echo "=== MEMORY ==="
        free -m

        echo
        echo "=== DISK ==="
        df -h /

        echo
        echo "=== TEMPERATURE / THROTTLING ==="
        if command -v vcgencmd >/dev/null 2>&1; then
            vcgencmd measure_temp
            vcgencmd get_throttled
        else
            echo -n "Temperature: "
            awk '{printf "%.1f°C\n", $1/1000}' \
                /sys/class/thermal/thermal_zone0/temp 2>/dev/null \
                || echo "unavailable"
        fi

        echo
        echo "=== SERVICE STATE ==="
        systemctl show masjidpi \
            -p ActiveState \
            -p SubState \
            -p NRestarts \
            -p ExecMainStartTimestamp

        systemctl show masjidpi-display \
            -p ActiveState \
            -p SubState \
            -p NRestarts \
            -p ExecMainStartTimestamp

        echo
        echo "=== PROCESSES ==="
        ps -eo pid,comm,%cpu,%mem,rss,etime \
            | grep -E 'masjidpi|mpv|cog|WPEWebProcess' \
            | grep -v grep \
            || true

        echo
        echo "=== PROCESS MEMORY DETAILS ==="
        for process_name in mpv WPEWebProcess; do
            process_pid="$(pgrep -xo "$process_name" 2>/dev/null || true)"
            if [ -z "$process_pid" ]; then
                echo "$process_name: not running"
                continue
            fi

            echo "$process_name PID=$process_pid"
            if [ -r "/proc/$process_pid/smaps_rollup" ]; then
                grep -E \
                    '^(Rss|Pss|Pss_Anon|Pss_File|Private_Clean|Private_Dirty|Swap):' \
                    "/proc/$process_pid/smaps_rollup" \
                    || true
            else
                echo "smaps_rollup unavailable"
            fi
        done

        echo
        echo "=== LISTEN STATUS ==="
        curl -fsS http://localhost:8080/api/listen/status 2>/dev/null \
            || echo "listen endpoint unavailable"

        echo
        echo "=== BOARD STATUS SUMMARY ==="
        curl -fsS http://localhost:8080/api/masjidboard/status 2>/dev/null \
            | jq -c '
                if .boards then
                    .boards[] |
                    {
                        name,
                        status,
                        using_cached_data,
                        update_failed,
                        last_successful_update,
                        update_error
                    }
                else
                    .
                end
            ' 2>/dev/null \
            || echo "board endpoint unavailable"

        echo
        echo "=== RECENT WARNINGS / ERRORS ==="
        journalctl -u masjidpi \
            --since "6 minutes ago" \
            --no-pager -o short-iso \
            | grep -E 'WARN|ERROR|failed|Failed|panic|fatal' \
            || echo "none"

        echo
    } >> "$LOG_FILE" 2>&1

    sleep 300
done
EOF

sudo chmod +x /usr/local/bin/masjidpi-rc-monitor.sh
```

## 2. Install the systemd service

Create `/etc/systemd/system/masjidpi-rc-monitor.service`:

```bash
sudo tee /etc/systemd/system/masjidpi-rc-monitor.service >/dev/null <<'EOF'
[Unit]
Description=MasjidPi RC Soak Monitor
After=network-online.target masjidpi.service
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/masjidpi-rc-monitor.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
```

Reload systemd and enable the monitor:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now masjidpi-rc-monitor.service
```

Because the service is enabled, it will start again automatically if the Raspberry Pi is rebooted during the soak test.

## 3. Verify that monitoring is running

Check the service:

```bash
systemctl status masjidpi-rc-monitor --no-pager
```

Confirm that the log exists and contains samples:

```bash
sudo tail -n 100 /var/log/masjidpi-rc-monitor/soak.log
```

## Reading the logs

Follow the monitor log live:

```bash
sudo tail -f /var/log/masjidpi-rc-monitor/soak.log
```

Read the most recent samples:

```bash
sudo tail -n 300 /var/log/masjidpi-rc-monitor/soak.log
```

Search for warnings, failures, throttling or unavailable endpoints:

```bash
sudo grep -Ei 'WARN|ERROR|failed|panic|fatal|unavailable|throttled' \
    /var/log/masjidpi-rc-monitor/soak.log
```

The normal Raspberry Pi throttling value is:

```text
throttled=0x0
```

A non-zero value should be investigated because it can indicate current or historical undervoltage, frequency capping or thermal throttling.

## Starting, stopping and restarting the monitor

Start the monitor:

```bash
sudo systemctl start masjidpi-rc-monitor.service
```

Stop the monitor without disabling it at boot:

```bash
sudo systemctl stop masjidpi-rc-monitor.service
```

Restart the monitor:

```bash
sudo systemctl restart masjidpi-rc-monitor.service
```

Show its current state:

```bash
systemctl status masjidpi-rc-monitor --no-pager
```

## Stop monitoring after RC acceptance

When the soak test is complete, stop the service and prevent it from starting on future boots:

```bash
sudo systemctl disable --now masjidpi-rc-monitor.service
```

The existing log is intentionally left in place at:

```text
/var/log/masjidpi-rc-monitor/soak.log
```

This allows the results to be reviewed after the monitor has been stopped.

## Remove the monitor completely

After the release candidate has been accepted and the log has been archived or reviewed, the temporary monitoring service can be removed:

```bash
sudo systemctl disable --now masjidpi-rc-monitor.service 2>/dev/null || true
sudo rm -f /etc/systemd/system/masjidpi-rc-monitor.service
sudo rm -f /usr/local/bin/masjidpi-rc-monitor.sh
sudo systemctl daemon-reload
```

To also delete the recorded soak data:

```bash
sudo rm -rf /var/log/masjidpi-rc-monitor
```

Do not remove the log directory until any required RC diagnostic data has been collected.

## What to look for during a soak test

A healthy 24–48 hour run should generally show:

- `masjidpi.service` and `masjidpi-display.service` staying active
- no unexpected increases in `NRestarts`
- no sustained growth in MasjidPi, mpv, Cog or WPE RSS usage
- adequate free memory and no unexpected swap pressure
- `throttled=0x0`
- no repeated mpv reconnect loops
- no unexpected Masjid/Radio source oscillation
- no repeated Board update failures affecting multiple boards
- Listen and Board HTTP endpoints remaining available

A single MasjidBoard provider warning can be valid when upstream data is malformed or unavailable, provided MasjidPi correctly retains last-known-good cached data. Repeated failures across multiple boards should be investigated.

## SD-card write consideration

This monitor intentionally writes a sample every five minutes and is intended for short RC validation periods, not permanent production operation.

For a normal 24–48 hour release-candidate soak this is acceptable. If the monitor is kept running for substantially longer, configure log rotation or stop the service to avoid unnecessary SD-card writes and unbounded log growth.
