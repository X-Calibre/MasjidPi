(() => {
    const scheduleButton = document.getElementById("radioModeSchedule");
    const playNowButton = document.getElementById("radioModePlayNow");
    const stopButton = document.getElementById("radioModeStop");
    const status = document.getElementById("radioModeStatus");

    if (!scheduleButton || !playNowButton || !stopButton || !status) return;

    const buttons = {
        schedule: scheduleButton,
        play_now: playNowButton,
        stopped: stopButton
    };

    function render(data) {
        const mode = data.radio_mode || "schedule";
        for (const [name, button] of Object.entries(buttons)) {
            const active = Boolean(data.radio_enabled) && name === mode;
            button.classList.toggle("active", active);
            button.setAttribute("aria-pressed", active ? "true" : "false");
        }
        if (!data.radio_enabled) {
            status.textContent = "Radio is off. Choose Play on Schedule or Play Now to enable it.";
            return;
        }

        if (mode === "stopped") {
            status.textContent = "Radio is stopped until you choose Play on Schedule or Play Now.";
            return;
        }
        if (mode === "play_now") {
            status.textContent = data.masjid_online
                ? "Play Now is waiting because the masjid has priority."
                : "Play Now override is active until the next masjid broadcast or schedule boundary.";
            return;
        }
        if (data.radio_schedule_enabled && !data.radio_schedule_allows_now) {
            status.textContent = "Scheduled mode is active; Radio is currently in quiet time.";
        } else if (data.radio_resume_pending) {
            status.textContent = "Scheduled mode is active; Radio is waiting for the configured post-masjid delay.";
        } else {
            status.textContent = "Scheduled mode is active.";
        }
    }

    async function setMode(mode) {
        for (const button of Object.values(buttons)) button.disabled = true;
        try {
            if (mode !== "stopped") {
                const current = await fetch("/api/listen/status");
                if (!current.ok) throw new Error(`Could not read Listen status (${current.status})`);
                if (!(await current.json()).radio_enabled) {
                    const power = await fetch("/api/listen/power", {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ module: "radio", enabled: true })
                    });
                    if (!power.ok) throw new Error((await power.json().catch(() => ({}))).error || `Could not enable Radio (${power.status})`);
                }
            }
            const response = await fetch("/api/listen/radio-mode", {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ mode })
            });
            const body = await response.json().catch(() => ({}));
            if (!response.ok) throw new Error(body.error || `Request failed (${response.status})`);
            render(body);
            const labels = { schedule: "Radio set to Play on Schedule.", play_now: "Radio Play Now enabled.", stopped: "Radio stopped." };
            window.MasjidFrameUI?.notify?.(labels[mode], "success");
        } catch (err) {
            window.MasjidFrameUI?.notify?.(err.message, "error");
        } finally {
            for (const button of Object.values(buttons)) button.disabled = false;
            await refresh();
        }
    }

    const refresh = () => window.MasjidFrameRefreshListenStatus?.();

    scheduleButton.addEventListener("click", () => setMode("schedule"));
    playNowButton.addEventListener("click", () => setMode("play_now"));
    stopButton.addEventListener("click", () => setMode("stopped"));

    window.addEventListener("masjidframe:listen-status", event => render(event.detail));
})();
