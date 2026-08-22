(() => {
    "use strict";

    const supported = new Set(["emerald", "midnight", "slate", "ruby", "light", "black-white"]);
    const override = new URLSearchParams(window.location.search).get("theme");

    function apply(theme) {
        const value = supported.has(theme) ? theme : "emerald";
        document.body.dataset.boardTheme = value;
    }

    async function refresh() {
        if (override && supported.has(override)) {
            apply(override);
            return;
        }
        try {
            const response = await fetch("/api/masjidboard/layout", {cache: "no-store"});
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const state = await response.json();
            apply(state && state.theme);
        } catch (error) {
            apply("emerald");
            console.warn("Could not load MasjidBoard colour theme", error);
        }
    }

    apply(override || "emerald");
    refresh();
    window.setInterval(refresh, 60_000);
})();
