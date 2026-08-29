(() => {
    "use strict";

    const endpoint = "/api/masjidboard/layout";
    const select = document.getElementById("displayLayout");
    const themeInputs = Array.from(document.querySelectorAll('input[name="boardTheme"]'));
    const saveStatus = document.getElementById("displaySaveStatus");
    const meta = document.getElementById("displayLayoutMeta");
    const slideDuration = document.getElementById("slideDuration");
    const slideDurationValue = document.getElementById("slideDurationValue");
    const showEconomicIndicators = document.getElementById("showEconomicIndicators");
    const previewLink = document.getElementById("displayPreviewLink");

    if (!select || !saveStatus || !meta || !slideDuration || !slideDurationValue || !showEconomicIndicators || themeInputs.length === 0) return;

    let lastSavedState = null;
    let savePending = false;
    let saving = false;

    const supportedThemes = new Set(["emerald", "midnight", "slate", "ruby", "light", "black-white"]);
    const themeNames = {
        emerald: "Emerald", midnight: "Midnight", slate: "Slate",
        ruby: "Ruby", light: "Light", "black-white": "Black & White",
    };

    function showBanner(message, kind = "success") {
        window.MasjidPiUI.notify(message, kind);
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

    function normaliseLayout(layout) {
        return layout === "portrait" ? "portrait" : "landscape";
    }

    function describe(layout, theme) {
        const layoutName = {landscape: "Landscape (1920 × 1080)", portrait: "Portrait (600 × 1024)"}[normaliseLayout(layout)];
        return `${layoutName} layout with the ${themeNames[theme] || "Emerald"} theme will be used automatically on HDMI output.`;
    }

    function updateDurationLabel() { slideDurationValue.textContent = `${slideDuration.value} seconds`; }

    function normaliseState(state) {
        return {
            layout: normaliseLayout(state && state.layout),
            theme: state && supportedThemes.has(state.theme) ? state.theme : "emerald",
            slide_duration_seconds: Number(state && state.slide_duration_seconds) || 15,
            show_economic_indicators: Boolean(state && state.show_economic_indicators),
        };
    }

    function currentState() {
        return normaliseState({
            layout: select.value,
            theme: selectedTheme(),
            slide_duration_seconds: Number(slideDuration.value),
            show_economic_indicators: showEconomicIndicators.checked,
        });
    }

    function applyState(state) {
        select.value = state.layout;
        setTheme(state.theme);
        slideDuration.value = String(state.slide_duration_seconds);
        showEconomicIndicators.checked = state.show_economic_indicators;
        updateDurationLabel();
        updatePreviewLink(state.layout);
        meta.textContent = describe(state.layout, state.theme);
    }

    function setSaveStatus(message, className = "") {
        saveStatus.textContent = message;
        saveStatus.className = `config-save-status${className ? ` ${className}` : ""}`;
    }

    function updatePreviewLink(layout) {
        if (!previewLink) return;
        previewLink.href = normaliseLayout(layout) === "portrait" ? "masjidboard.html?layout=portrait" : "masjidboard.html";
        previewLink.setAttribute("aria-label", `Open ${normaliseLayout(layout)} display preview`);
    }

    async function load() {
        select.disabled = true;
        slideDuration.disabled = true;
        showEconomicIndicators.disabled = true;
        for (const input of themeInputs) input.disabled = true;
        try {
            lastSavedState = normaliseState(await request());
            applyState(lastSavedState);
            setSaveStatus("Changes are saved automatically.");
        } catch (error) {
            meta.textContent = `Could not load HDMI display settings: ${error.message}`;
            setSaveStatus("Settings could not be loaded.", "error");
        } finally {
            select.disabled = false;
            slideDuration.disabled = false;
            showEconomicIndicators.disabled = false;
            for (const input of themeInputs) input.disabled = false;
        }
    }

    async function drainSaves() {
        if (saving) return;
        saving = true;
        while (savePending) {
            savePending = false;
            const candidate = currentState();
            setSaveStatus("Saving…", "saving");
            try {
                lastSavedState = normaliseState(await request({
                    method: "PUT",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify(candidate),
                }));
                if (!savePending) applyState(lastSavedState);
                setSaveStatus(savePending ? "Saving…" : "Saved", savePending ? "saving" : "saved");
            } catch (error) {
                savePending = false;
                if (lastSavedState) applyState(lastSavedState);
                setSaveStatus("Could not save changes.", "error");
                showBanner(`Could not save HDMI display settings: ${error.message}`, "error");
            }
        }
        saving = false;
    }

    function saveAutomatically() {
        savePending = true;
        void drainSaves();
    }

    select.addEventListener("change", () => { meta.textContent = describe(select.value, selectedTheme()); updatePreviewLink(select.value); saveAutomatically(); });
    slideDuration.addEventListener("input", updateDurationLabel);
    slideDuration.addEventListener("change", saveAutomatically);
    showEconomicIndicators.addEventListener("change", saveAutomatically);
    for (const input of themeInputs) input.addEventListener("change", () => { meta.textContent = describe(select.value, selectedTheme()); saveAutomatically(); });
    load();
})();
