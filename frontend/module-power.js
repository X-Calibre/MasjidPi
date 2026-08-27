(() => {
    const masjidSwitch = document.getElementById("masjidPowerSwitch");
    const radioSwitch = document.getElementById("radioPowerSwitch");
    const masjidStatus = document.getElementById("masjidPowerStatus");
    const radioStatus = document.getElementById("radioPowerStatus");

    if (!masjidSwitch || !radioSwitch || !masjidStatus || !radioStatus) return;

    const radioControls = [
        "radioModeSchedule",
        "radioModePlayNow",
        "radioModeStop",
        "radioVolumeSlider",
        "radioResumeDelaySlider",
        "radioScheduleEnabled",
        "radioScheduleStart",
        "radioScheduleStop",
        "radioStream"
    ].map(id => document.getElementById(id)).filter(Boolean);

    let updating = false;

    async function setPower(module, enabled) {
        const response = await fetch("/api/listen/power", {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ module, enabled })
        });
        if (!response.ok) {
            const body = await response.json().catch(() => ({}));
            throw new Error(body.error || `Request failed (${response.status})`);
        }
        return response.json();
    }

    function render(data) {
        updating = true;
        masjidSwitch.checked = Boolean(data.masjid_enabled);
        radioSwitch.checked = Boolean(data.radio_enabled);
        radioSwitch.disabled = !data.masjid_enabled;

        masjidStatus.textContent = data.masjid_enabled
            ? "Masjid module is powered on."
            : "Masjid module is powered off. Radio is also forced off.";

        if (!data.masjid_enabled) {
            radioStatus.textContent = "Radio cannot be powered on while Masjid is off.";
        } else {
            radioStatus.textContent = data.radio_enabled
                ? "Radio module is powered on."
                : "Radio module is powered off until switched back on.";
        }

        for (const control of radioControls) {
            if (control === radioSwitch) continue;
            control.disabled = !data.radio_enabled;
        }
        updating = false;
    }

    async function refresh() {
        try {
            const response = await fetch("/api/listen/status");
            if (!response.ok) return;
            render(await response.json());
        } catch (_) {}
    }

    masjidSwitch.addEventListener("change", async () => {
        if (updating) return;
        masjidSwitch.disabled = true;
        try {
            const data = await setPower("masjid", masjidSwitch.checked);
            render(data);
            window.MasjidPiUI?.notify?.(
                masjidSwitch.checked ? "Masjid powered on." : "Masjid and Radio powered off.",
                "success"
            );
        } catch (err) {
            window.MasjidPiUI?.notify?.(err.message, "error");
            await refresh();
        } finally {
            masjidSwitch.disabled = false;
        }
    });

    radioSwitch.addEventListener("change", async () => {
        if (updating) return;
        radioSwitch.disabled = true;
        try {
            const data = await setPower("radio", radioSwitch.checked);
            render(data);
            window.MasjidPiUI?.notify?.(
                radioSwitch.checked ? "Radio powered on." : "Radio powered off.",
                "success"
            );
        } catch (err) {
            window.MasjidPiUI?.notify?.(err.message, "error");
            await refresh();
        } finally {
            await refresh();
        }
    });

    refresh();
    setInterval(refresh, 2000);
})();
