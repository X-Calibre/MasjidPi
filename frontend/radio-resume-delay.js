(() => {
    const slider = document.getElementById("radioResumeDelaySlider");
    const value = document.getElementById("radioResumeDelayValue");
    const status = document.getElementById("radioResumeDelayStatus");

    if (!slider || !value || !status) return;

    let editing = false;

    function formatMinutes(minutes) {
        return `${minutes} ${minutes === 1 ? "minute" : "minutes"}`;
    }

    function renderValue(minutes) {
        value.textContent = formatMinutes(minutes);
    }

    function renderPending(data) {
        if (!data.radio_resume_pending || !data.radio_resume_at) {
            status.textContent = "";
            status.classList.add("hidden");
            return;
        }
        const remainingMs = new Date(data.radio_resume_at).getTime() - Date.now();
        const remainingSeconds = Math.max(0, Math.ceil(remainingMs / 1000));
        const minutes = Math.floor(remainingSeconds / 60);
        const seconds = remainingSeconds % 60;
        status.textContent = `Radio resumes in ${minutes}:${String(seconds).padStart(2, "0")}`;
        status.classList.remove("hidden");
    }

    async function refresh() {
        try {
            const response = await fetch("/api/listen/status");
            if (!response.ok) return;
            const data = await response.json();
            if (!editing && Number.isInteger(data.radio_resume_delay_minutes)) {
                slider.value = data.radio_resume_delay_minutes;
                renderValue(data.radio_resume_delay_minutes);
            }
            renderPending(data);
        } catch (_) {}
    }

    slider.addEventListener("input", () => {
        editing = true;
        renderValue(Number(slider.value));
    });

    slider.addEventListener("change", async () => {
        const minutes = Number(slider.value);
        try {
            const response = await fetch("/api/listen/radio-delay", {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ minutes })
            });
            if (!response.ok) {
                const body = await response.json().catch(() => ({}));
                throw new Error(body.error || `Request failed (${response.status})`);
            }
            window.MasjidPiUI?.notify?.(`Radio resume delay set to ${formatMinutes(minutes)}.`, "success");
        } catch (err) {
            window.MasjidPiUI?.notify?.(err.message, "error");
        } finally {
            editing = false;
            await refresh();
        }
    });

    renderValue(Number(slider.value));
    refresh();
    setInterval(refresh, 1000);
})();
