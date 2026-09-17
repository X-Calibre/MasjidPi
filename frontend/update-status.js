(() => {
    "use strict";

    const widgets = [...document.querySelectorAll("[data-update-widget]")];
    if (widgets.length === 0) return;

    const formatDate = value => {
        if (!value) return "Not yet";
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return value;
        return new Intl.DateTimeFormat("en-ZA", {
            day: "numeric",
            month: "long",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit"
        }).format(date);
    };

    const formatDateOnly = value => {
        if (!value) return "Not scheduled";
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return value;
        return new Intl.DateTimeFormat("en-ZA", {
            day: "numeric",
            month: "long",
            year: "numeric"
        }).format(date);
    };

    function setText(widget, selector, value) {
        const element = widget.querySelector(selector);
        if (element) element.textContent = value;
    }

    function touchSummary(state, release, installation) {
        const runningTarget = installation?.version &&
            ["probation", "installed"].includes(installation.status);
        const releaseInstalled = installation?.status === "installed" &&
            installation.version === release?.version;
        const availableVersion = release && !releaseInstalled
            ? release.version
            : "None";
        const installationDate = release && !releaseInstalled
            ? installation?.staged_at || installation?.started_at ||
                state?.approved_at || state?.approval_deadline
            : null;

        return {
            installedVersion: runningTarget
                ? installation.version
                : state?.current_version || "Unknown",
            availableVersion,
            installationDate: formatDateOnly(installationDate)
        };
    }

    function render(widget, state) {
        widget.updateState = state;
        const release = state?.available_release || null;
        const download = state?.download || null;
        const installation = state?.installation || null;
        const status = widget.querySelector("[data-update-state]");
        const link = widget.querySelector("[data-update-link]");
        const actions = widget.querySelector("[data-update-actions]");
        const summary = touchSummary(state, release, installation);

        setText(widget, "[data-update-current]", state?.current_version || "Unknown");
        setText(widget, "[data-update-installed-version]", summary.installedVersion);
        setText(widget, "[data-update-available]", summary.availableVersion);
        setText(widget, "[data-update-install-date]", summary.installationDate);
        setText(widget, "[data-update-last-checked]", formatDate(state?.last_checked_at));
        setText(widget, "[data-update-deadline]", release ? formatDate(state?.approval_deadline) : "Not applicable");
        let downloadText = "Not started";
        if (download?.status === "downloading") downloadText = "Downloading…";
        if (download?.status === "failed") downloadText = "Download failed; retry pending";
        if (download?.status === "verified") downloadText = "Downloaded and verified";
        setText(widget, "[data-update-download]", release ? downloadText : "Not applicable");
        let installationText = "Not started";
        if (installation?.status === "installing") installationText = "Installing…";
        if (installation?.status === "reboot_pending") installationText = "Restart pending";
        if (installation?.status === "probation") installationText = "Health probation in progress";
        if (installation?.status === "installed") installationText = "Installed and confirmed";
        if (installation?.status === "rolled_back") installationText = "Rolled back safely";
        if (installation?.status === "failed") installationText = "Installation failed";
        setText(widget, "[data-update-installation]", release ? installationText : "Not applicable");

        if (release) {
            if (status) {
                if (state.approved_at) {
                    status.textContent = `${release.version} is approved`;
                } else if (state.postponed_until) {
                    status.textContent = `${release.version} postponed until ${formatDate(state.postponed_until)}`;
                } else {
                    status.textContent = `${release.version} is available`;
                }
                status.className = "update-state update-state-available";
            }
            setText(widget, "[data-update-published]", formatDate(release.published_at));
            if (link) {
                link.href = release.page_url;
                link.classList.remove("hidden");
            }
            actions?.classList.remove("hidden");
            const approve = widget.querySelector("[data-update-approve]");
            if (approve) approve.disabled = Boolean(state.approved_at);
            const install = widget.querySelector("[data-update-install]");
            if (install) {
                const ready = download?.status === "verified" &&
                    !["installing", "reboot_pending", "probation", "installed"].includes(installation?.status);
                install.classList.toggle("hidden", !ready);
            }
        } else {
            if (status) {
                status.textContent = "No stable update available";
                status.className = "update-state update-state-current";
            }
            setText(widget, "[data-update-published]", "Not applicable");
            if (link) {
                link.removeAttribute("href");
                link.classList.add("hidden");
            }
            actions?.classList.add("hidden");
            widget.querySelector("[data-update-install]")?.classList.add("hidden");
        }

        const error = widget.querySelector("[data-update-error]");
        if (error) {
            const message = state?.last_check_error || download?.last_error || installation?.last_error || "";
            error.textContent = message;
            error.classList.toggle("hidden", !message);
        }
    }

    function renderFailure(widget, message) {
        const status = widget.querySelector("[data-update-state]");
        if (status) {
            status.textContent = "Update status unavailable";
            status.className = "update-state update-state-error";
        }
        setText(widget, "[data-update-available]", "Unavailable");
        setText(widget, "[data-update-install-date]", "Unavailable");
        const error = widget.querySelector("[data-update-error]");
        if (error) {
            error.textContent = message;
            error.classList.remove("hidden");
        }
    }

    async function requestState(method) {
        const response = await fetch(
            method === "POST" ? "/api/update/check" : "/api/update/status",
            {method}
        );
        let state = null;
        try { state = await response.json(); } catch (_) {}
        if (!response.ok && !state?.schema_version) {
            throw new Error(state?.error || `Request failed (${response.status})`);
        }
        return {state, ok: response.ok};
    }

    async function requestDecision(path, body) {
        const response = await fetch(path, {
            method: "POST",
            headers: body ? {"Content-Type": "application/json"} : undefined,
            body: body ? JSON.stringify(body) : undefined
        });
        const payload = await response.json();
        if (!response.ok) throw new Error(payload?.error || `Request failed (${response.status})`);
        return payload;
    }

    async function decide(widget, action) {
        const buttons = [...widget.querySelectorAll("[data-update-actions] button")];
        buttons.forEach(button => { button.disabled = true; });
        try {
            const sevenDays = Date.now() + 7 * 24 * 60 * 60 * 1000;
            const deadline = new Date(widget.updateState?.approval_deadline || sevenDays).getTime();
            const state = action === "approve"
                ? await requestDecision("/api/update/approve")
                : await requestDecision("/api/update/postpone", {
                    until: new Date(Math.min(sevenDays, deadline)).toISOString()
                });
            render(widget, state);
        } catch (error) {
            renderFailure(widget, error.message);
        } finally {
            buttons.forEach(button => { button.disabled = false; });
			const approve = widget.querySelector("[data-update-approve]");
			if (approve) approve.disabled = Boolean(widget.updateState?.approved_at);
        }
    }

    async function installNow(widget) {
        const button = widget.querySelector("[data-update-install]");
        const confirmed = window.confirm(
            "Install this verified update now? MasjidPi may stop active audio and will restart automatically."
        );
        if (!confirmed) return;
        button.disabled = true;
        button.textContent = "Installing…";
        try {
            const state = await requestDecision("/api/update/install", {
                interrupt_playback: true
            });
            render(widget, state);
        } catch (error) {
            renderFailure(widget, error.message);
        } finally {
            button.disabled = false;
            button.textContent = "Install now";
        }
    }

    async function refresh(widget, checkNow = false) {
        const button = widget.querySelector("[data-update-check]");
        if (button) {
            button.disabled = true;
            button.textContent = checkNow ? "Checking…" : "Loading…";
        }
        try {
            const result = await requestState(checkNow ? "POST" : "GET");
            render(widget, result.state);
        } catch (error) {
            renderFailure(widget, error.message);
        } finally {
            if (button) {
                button.disabled = false;
                button.textContent = "Check now";
            }
        }
    }

    for (const widget of widgets) {
        widget.querySelector("[data-update-check]")?.addEventListener("click", () => refresh(widget, true));
        widget.querySelector("[data-update-approve]")?.addEventListener("click", () => decide(widget, "approve"));
        widget.querySelector("[data-update-postpone]")?.addEventListener("click", () => decide(widget, "postpone"));
        widget.querySelector("[data-update-install]")?.addEventListener("click", () => installNow(widget));
        refresh(widget);
    }
})();
