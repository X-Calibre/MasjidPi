(() => {
    "use strict";

    const endpoint = "/api/masjidboard/layout";
    const select = document.getElementById("displayLayout");
    const themeInputs = Array.from(document.querySelectorAll('input[name="boardTheme"]'));
    const saveButton = document.getElementById("saveDisplayLayoutButton");
    const meta = document.getElementById("displayLayoutMeta");
    const banner = document.getElementById("configBanner");

    if (!select || !saveButton || !meta || themeInputs.length === 0) return;

    const supportedThemes = new Set(["emerald", "midnight", "slate", "ruby", "light", "black-white"]);
    const themeNames = {
        emerald: "Emerald", midnight: "Midnight", slate: "Slate",
        ruby: "Ruby", light: "Light", "black-white": "Black & White",
    };

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
            try { const body = await response.json(); if (body && body.error) message = body.error; } catch (_) {}
            throw new Error(message);
        }
        return response.json();
    }

    function selectedTheme() {
        return themeInputs.find((input) => input.checked)?.value || "emerald";
    }

    function setTheme(theme) {
        const value = supportedThemes.has(theme) ? theme : "emerald";
        for (const input of themeInputs) input.checked = input.value === value;
    }

    function describe(layout, theme) {
        const layoutName = layout === "detailed" ? "Detailed" : "Standard";
        return `${layoutName} layout with the ${themeNames[theme] || "Emerald"} theme will be used automatically on HDMI output.`;
    }

    async function load() {
        select.disabled = true; saveButton.disabled = true;
        for (const input of themeInputs) input.disabled = true;
        try {
            const state = await request();
            const layout = state && state.layout === "detailed" ? "detailed" : "standard";
            const theme = state && supportedThemes.has(state.theme) ? state.theme : "emerald";
            select.value = layout; setTheme(theme); meta.textContent = describe(layout, theme);
        } catch (error) {
            meta.textContent = `Could not load HDMI display settings: ${error.message}`;
        } finally {
            select.disabled = false; saveButton.disabled = false;
            for (const input of themeInputs) input.disabled = false;
        }
    }

    async function save() {
        saveButton.disabled = true; saveButton.textContent = "Saving…";
        try {
            const state = await request({
                method: "PUT",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({layout: select.value, theme: selectedTheme()}),
            });
            const layout = state && state.layout === "detailed" ? "detailed" : "standard";
            const theme = state && supportedThemes.has(state.theme) ? state.theme : "emerald";
            select.value = layout; setTheme(theme); meta.textContent = describe(layout, theme);
            showBanner(`HDMI display saved: ${layout === "detailed" ? "Detailed" : "Standard"}, ${themeNames[theme]}.`);
        } catch (error) {
            showBanner(`Could not save HDMI display settings: ${error.message}`, "error");
        } finally {
            saveButton.disabled = false; saveButton.textContent = "Save HDMI Display";
        }
    }

    select.addEventListener("change", () => { meta.textContent = describe(select.value, selectedTheme()); });
    for (const input of themeInputs) input.addEventListener("change", () => { meta.textContent = describe(select.value, selectedTheme()); });
    saveButton.addEventListener("click", save);
    load();
})();
