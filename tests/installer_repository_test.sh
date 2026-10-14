#!/usr/bin/env bash
set -euo pipefail
root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

# Check both download entry points without running a production installation.
for script in install-latest.sh github.sh; do
    (
        # Read only repository constants from each entry point.
        # shellcheck disable=SC1090
        source <(awk '/^info\(\)|^update_repository\(\)/ {exit} {print}' "$root_dir/scripts/$script")
        test "$REPO_OWNER/$REPO_NAME" = "X-Calibre/MasjidPi"
        test "$RELEASE_API" = "https://api.github.com/repos/X-Calibre/MasjidPi/releases/latest"
        if [[ "$script" = github.sh ]]; then
            test "$REPO_URL" = "https://github.com/X-Calibre/MasjidPi.git"
        fi
    )
done
printf '%s\n' 'Installer repository tests passed'
