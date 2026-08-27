(() => {
    const nameEl = document.getElementById("selectedMasjidName");
    const locationEl = document.getElementById("selectedMasjidLocation");
    if (!nameEl || !locationEl) return;

    let catalogue = [];
    let lastMasjidID = null;

    async function getJSON(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`Request failed (${response.status})`);
        return response.json();
    }

    function render(masjidID) {
        lastMasjidID = masjidID || "";
        if (!masjidID) {
            nameEl.textContent = "No masjid selected";
            locationEl.textContent = "";
            return;
        }

        const selected = catalogue.find(item => item.id === masjidID);
        if (!selected) {
            nameEl.textContent = masjidID;
            locationEl.textContent = "";
            return;
        }

        nameEl.textContent = selected.name;
        locationEl.textContent = selected.location || "";
    }

    async function refresh() {
        try {
            const [status, streams] = await Promise.all([
                getJSON("/api/listen/status"),
                catalogue.length ? Promise.resolve(catalogue) : getJSON("/api/streams?kind=masjid")
            ]);
            catalogue = streams;
            render(status.masjid_id);
        } catch (_) {
            if (lastMasjidID === null) {
                nameEl.textContent = "Unable to load selected masjid";
                locationEl.textContent = "";
            }
        }
    }

    const streamInput = document.getElementById("stream");
    if (streamInput) {
        streamInput.addEventListener("change", () => {
            const selected = catalogue.find(item => item.id === streamInput.value);
            if (selected) {
                nameEl.textContent = selected.name;
                locationEl.textContent = selected.location || "";
            }
        });
    }

    refresh();
    setInterval(refresh, 2000);
})();
