(() => {
    "use strict";

    const endpoint = "/api/masjidboard/layout";
    const refreshIntervalMs = 2_000;
    const supportedThemes = new Set(["emerald", "midnight", "slate", "ruby", "light", "black-white"]);
    const params = new URLSearchParams(window.location.search);
    const themeOverride = params.get("theme");
    const hasThemeOverride = supportedThemes.has(themeOverride);
    let refreshPending = false;

    function applyTheme(theme) {
        const value = supportedThemes.has(theme) ? theme : "emerald";
        if (document.body.dataset.boardTheme !== value) {
            document.body.dataset.boardTheme = value;
        }
    }

    function currentLayout() {
        const layout = new URLSearchParams(window.location.search).get("layout");
        return ["detailed", "portrait"].includes(layout) ? layout : "standard";
    }

    function applyLayout(layout) {
        const wanted = ["detailed", "portrait"].includes(layout) ? layout : "standard";
        if (wanted === currentLayout()) return false;

        const url = new URL(window.location.href);
        if (wanted !== "standard") url.searchParams.set("layout", wanted);
        else url.searchParams.delete("layout");

        // Preserve development/test parameters such as date, time and theme.
        window.location.replace(url.toString());
        return true;
    }

    async function refresh() {
        if (refreshPending) return;
        refreshPending = true;
        try {
            const response = await fetch(endpoint, {cache: "no-store"});
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const state = await response.json();

            // A layout change needs the matching Standard/Detailed scripts to be
            // loaded, so navigate the existing board page rather than restarting
            // the Cog display process.
            if (applyLayout(state && state.layout)) return;

            // Explicit URL themes remain useful as development/test overrides.
            // Normal appliance display follows the persisted WebUI preference.
            if (hasThemeOverride) applyTheme(themeOverride);
            else applyTheme(state && state.theme);
        } catch (error) {
            if (hasThemeOverride) applyTheme(themeOverride);
            console.warn("Could not refresh MasjidBoard HDMI display settings", error);
        } finally {
            refreshPending = false;
        }
    }

    applyTheme(hasThemeOverride ? themeOverride : "emerald");
    refresh();
    window.setInterval(refresh, refreshIntervalMs);
})();
