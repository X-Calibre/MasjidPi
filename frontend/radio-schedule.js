(() => {
    const enabled = document.getElementById("radioScheduleEnabled");
    const start = document.getElementById("radioScheduleStart");
    const stop = document.getElementById("radioScheduleStop");
    const times = document.getElementById("radioScheduleTimes");
    const saveButton = document.getElementById("radioScheduleSave");
    const status = document.getElementById("radioScheduleStatus");

    if (!enabled || !start || !stop || !times || !saveButton || !status) return;

    let saved = null;
    let saving = false;

    function values() {
        return { enabled: enabled.checked, start: start.value, stop: stop.value };
    }

    function isDirty() {
        const current = values();
        return saved && (current.enabled !== saved.enabled || current.start !== saved.start || current.stop !== saved.stop);
    }

    function updateControls() {
        start.disabled = !enabled.checked || enabled.disabled;
        stop.disabled = !enabled.checked || enabled.disabled;
        times.classList.toggle("radio-schedule-disabled", !enabled.checked);
        saveButton.disabled = saving || enabled.disabled || !isDirty();
        if (isDirty() && !saving) status.textContent = "Unsaved schedule changes.";
    }

    function renderStatus(data) {
        if (!data.radio_schedule_enabled) {
            status.textContent = "Daily Radio hours are off. In scheduled mode, Radio may play at any time.";
            return;
        }
        const hours = `${data.radio_schedule_start} to ${data.radio_schedule_stop}`;
        status.textContent = `Radio plays daily from ${hours} (appliance local time). It is currently ${data.radio_schedule_allows_now ? "within" : "outside"} these hours.`;
    }

    function refresh(data) {
        const dirty = isDirty();
        saved = {
            enabled: Boolean(data.radio_schedule_enabled),
            start: data.radio_schedule_start || "06:00",
            stop: data.radio_schedule_stop || "22:00"
        };
        if (!dirty || saving) {
            enabled.checked = saved.enabled;
            start.value = saved.start;
            stop.value = saved.stop;
            renderStatus(data);
        }
        updateControls();
    }

    async function save() {
        const next = values();
        if (next.enabled && (!next.start || !next.stop || next.start === next.stop)) {
            status.textContent = "Choose different start and stop times before saving.";
            return;
        }
        saving = true;
        updateControls();
        try {
            const response = await fetch("/api/listen/radio-schedule", {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(next)
            });
            const body = await response.json().catch(() => ({}));
            if (!response.ok) throw new Error(body.error || `Request failed (${response.status})`);
            saved = next;
            window.MasjidFrameUI?.notify?.("Radio schedule saved.", "success");
            renderStatus(body);
        } catch (err) {
            status.textContent = `Could not save schedule: ${err.message}`;
            window.MasjidFrameUI?.notify?.(err.message, "error");
        } finally {
            saving = false;
            updateControls();
            await window.MasjidFrameRefreshListenStatus?.();
        }
    }

    enabled.addEventListener("change", updateControls);
    for (const input of [start, stop]) input.addEventListener("change", updateControls);
    saveButton.addEventListener("click", save);
    updateControls();
    window.addEventListener("masjidframe:listen-status", event => refresh(event.detail));
})();
