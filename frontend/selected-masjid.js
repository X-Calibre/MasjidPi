(() => {
    const nameEl = document.getElementById("selectedMasjidName");
    const locationEl = document.getElementById("selectedMasjidLocation");
    if (!nameEl || !locationEl) return;

    let catalogue = [];
    let lastMasjidID = null;

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

    window.addEventListener("masjidpi:listen-status", event => render(event.detail.masjid_id));
    window.addEventListener("masjidpi:masjid-catalogue", event => {
        catalogue = event.detail;
        if (lastMasjidID !== null) render(lastMasjidID);
    });
})();
