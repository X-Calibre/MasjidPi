(() => {
    const masjidSwitch = document.getElementById("masjidPowerSwitch");
    const radioSwitch = document.getElementById("radioPowerSwitch");
    const masjidStatus = document.getElementById("masjidPowerStatus");
    const radioStatus = document.getElementById("radioPowerStatus");

    if (!masjidSwitch || !radioSwitch || !masjidStatus || !radioStatus) return;

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
        radioSwitch.disabled = false;

        masjidStatus.textContent = data.masjid_enabled
            ? "Masjid module is powered on."
            : "Masjid module is powered off. Radio is also off.";

        if (!data.masjid_enabled) {
            radioStatus.textContent = "Turning Radio on will also power Masjid on automatically.";
        } else {
            radioStatus.textContent = data.radio_enabled
                ? "Radio module is powered on."
                : "Radio module is powered off until switched back on.";
        }

        updating = false;
    }

    const refresh = () => window.MasjidPiRefreshListenStatus?.();

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
        } finally {
            await refresh();
        }
    });

    radioSwitch.addEventListener("change", async () => {
        if (updating) return;
        const enabling = radioSwitch.checked;
        const masjidWasOff = !masjidSwitch.checked;
        radioSwitch.disabled = true;
        try {
            const data = await setPower("radio", enabling);
            render(data);
            if (enabling && masjidWasOff && data.masjid_enabled) {
                window.MasjidPiUI?.notify?.("Radio powered on. Masjid was also powered on automatically.", "success");
            } else {
                window.MasjidPiUI?.notify?.(
                    enabling ? "Radio powered on." : "Radio powered off.",
                    "success"
                );
            }
        } catch (err) {
            window.MasjidPiUI?.notify?.(err.message, "error");
        } finally {
            await refresh();
        }
    });

    window.addEventListener("masjidpi:listen-status", event => render(event.detail));
})();
