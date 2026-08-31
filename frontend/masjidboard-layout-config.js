(() => {
    "use strict";

    const endpoint = "/api/masjidboard/layout";
    const themeInputs = Array.from(document.querySelectorAll('input[name="boardTheme"]'));
    const saveStatus = document.getElementById("displaySaveStatus");
    const meta = document.getElementById("displayLayoutMeta");
    const slideDuration = document.getElementById("slideDuration");
    const slideDurationValue = document.getElementById("slideDurationValue");
    const showEconomicIndicators = document.getElementById("showEconomicIndicators");

    if (!saveStatus || !meta || !slideDuration || !slideDurationValue || !showEconomicIndicators || themeInputs.length === 0) return;

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

    function describe(theme) {
        return `The ${themeNames[theme] || "Emerald"} theme is shared by the automatically selected standard and appliance display profiles.`;
    }

    function updateDurationLabel() { slideDurationValue.textContent = `${slideDuration.value} seconds`; }

    function normaliseState(state) {
        return {
            theme: state && supportedThemes.has(state.theme) ? state.theme : "emerald",
            slide_duration_seconds: Number(state && state.slide_duration_seconds) || 15,
            show_economic_indicators: Boolean(state && state.show_economic_indicators),
        };
    }

    function currentState() {
        return normaliseState({
            theme: selectedTheme(),
            slide_duration_seconds: Number(slideDuration.value),
            show_economic_indicators: showEconomicIndicators.checked,
        });
    }

    function applyState(state) {
        setTheme(state.theme);
        slideDuration.value = String(state.slide_duration_seconds);
        showEconomicIndicators.checked = state.show_economic_indicators;
        updateDurationLabel();
        meta.textContent = describe(state.theme);
    }

    function setSaveStatus(message, className = "") {
        saveStatus.textContent = message;
        saveStatus.className = `config-save-status${className ? ` ${className}` : ""}`;
    }

    async function load() {
        slideDuration.disabled = true;
        showEconomicIndicators.disabled = true;
        for (const input of themeInputs) input.disabled = true;
        try {
            lastSavedState = normaliseState(await request());
            applyState(lastSavedState);
            setSaveStatus("Changes are saved automatically.");
        } catch (error) {
            meta.textContent = `Could not load display settings: ${error.message}`;
            setSaveStatus("Settings could not be loaded.", "error");
        } finally {
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
                showBanner(`Could not save display settings: ${error.message}`, "error");
            }
        }
        saving = false;
    }

    function saveAutomatically() {
        savePending = true;
        void drainSaves();
    }

    slideDuration.addEventListener("input", updateDurationLabel);
    slideDuration.addEventListener("change", saveAutomatically);
    showEconomicIndicators.addEventListener("change", saveAutomatically);
    for (const input of themeInputs) input.addEventListener("change", () => { meta.textContent = describe(selectedTheme()); saveAutomatically(); });
    load();
})();
