(() => {
    "use strict";

    const endpoint = "/api/masjidboard/layout";
    const select = document.getElementById("displayLayout");
    const saveButton = document.getElementById("saveDisplayLayoutButton");
    const meta = document.getElementById("displayLayoutMeta");
    const banner = document.getElementById("configBanner");

    if (!select || !saveButton || !meta) return;

    function showBanner(message, kind = "success") {
        if (!banner) return;
        banner.textContent = message;
        banner.className = `config-banner ${kind}`;
        window.clearTimeout(showBanner.timer);
        showBanner.timer = window.setTimeout(() => banner.classList.add("hidden"), 6000);
    }

    async function request(options = {}) {
        const response = await fetch(endpoint, {cache: "no-store", ...options});
        if (!response.ok) {
            let message = `HTTP ${response.status}`;
            try {
                const body = await response.json();
                if (body && body.error) message = body.error;
            } catch (_) {}
            throw new Error(message);
        }
        return response.json();
    }

    function describe(layout) {
        return layout === "detailed"
            ? "Detailed will be used automatically on HDMI output."
            : "Standard will be used automatically on HDMI output.";
    }

    async function load() {
        select.disabled = true;
        saveButton.disabled = true;
        try {
            const state = await request();
            const layout = state && state.layout === "detailed" ? "detailed" : "standard";
            select.value = layout;
            meta.textContent = describe(layout);
        } catch (error) {
            meta.textContent = `Could not load HDMI layout: ${error.message}`;
        } finally {
            select.disabled = false;
            saveButton.disabled = false;
        }
    }

    async function save() {
        saveButton.disabled = true;
        saveButton.textContent = "Saving…";
        try {
            const state = await request({
                method: "PUT",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({layout: select.value}),
            });
            const layout = state && state.layout === "detailed" ? "detailed" : "standard";
            select.value = layout;
            meta.textContent = describe(layout);
            showBanner(`HDMI display layout saved as ${layout === "detailed" ? "Detailed" : "Standard"}.`);
        } catch (error) {
            showBanner(`Could not save HDMI display layout: ${error.message}`, "error");
        } finally {
            saveButton.disabled = false;
            saveButton.textContent = "Save HDMI Layout";
        }
    }

    saveButton.addEventListener("click", save);
    load();
})();
