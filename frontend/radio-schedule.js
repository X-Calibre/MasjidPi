(() => {
    const enabled = document.getElementById("radioScheduleEnabled");
    const start = document.getElementById("radioScheduleStart");
    const stop = document.getElementById("radioScheduleStop");
    const times = document.getElementById("radioScheduleTimes");
    const status = document.getElementById("radioScheduleStatus");

    if (!enabled || !start || !stop || !times || !status) return;

    let editing = false;

    function renderEnabled() {
        start.disabled = !enabled.checked;
        stop.disabled = !enabled.checked;
        times.classList.toggle("radio-schedule-disabled", !enabled.checked);
    }

    function renderStatus(data) {
        if (!data.radio_schedule_enabled) {
            status.textContent = "Radio schedule disabled — radio may play at any time.";
            return;
        }
        status.textContent = data.radio_schedule_allows_now
            ? `Radio is currently allowed (${data.radio_schedule_start}–${data.radio_schedule_stop}).`
            : `Radio is currently silenced (${data.radio_schedule_start}–${data.radio_schedule_stop}).`;
    }

    async function save() {
        if (enabled.checked && start.value === stop.value) {
            window.MasjidPiUI?.notify?.("Radio start and stop times must differ.", "error");
            await window.MasjidPiRefreshListenStatus?.();
            return;
        }

        try {
            const response = await fetch("/api/listen/radio-schedule", {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    enabled: enabled.checked,
                    start: start.value,
                    stop: stop.value
                })
            });
            if (!response.ok) {
                const body = await response.json().catch(() => ({}));
                throw new Error(body.error || `Request failed (${response.status})`);
            }
            window.MasjidPiUI?.notify?.(
                enabled.checked ? `Radio schedule set to ${start.value}–${stop.value}.` : "Radio schedule disabled.",
                "success"
            );
        } catch (err) {
            window.MasjidPiUI?.notify?.(err.message, "error");
        } finally {
            editing = false;
            await window.MasjidPiRefreshListenStatus?.();
        }
    }

    function refresh(data) {
        if (!editing) {
            enabled.checked = Boolean(data.radio_schedule_enabled);
            if (data.radio_schedule_start) start.value = data.radio_schedule_start;
            if (data.radio_schedule_stop) stop.value = data.radio_schedule_stop;
            renderEnabled();
        }
        renderStatus(data);
    }

    enabled.addEventListener("change", async () => {
        editing = true;
        renderEnabled();
        await save();
    });

    for (const input of [start, stop]) {
        input.addEventListener("input", () => { editing = true; });
        input.addEventListener("change", save);
    }

    renderEnabled();
    window.addEventListener("masjidpi:listen-status", event => refresh(event.detail));
})();
