(() => {
    "use strict";

    const supportedThemes = new Set(["emerald", "midnight", "slate", "ruby", "light", "black-white"]);
    const supportedFontPreviews = new Set(["ibm-plex", "source-sans"]);
    const params = new URLSearchParams(window.location.search);
    const themeOverride = params.get("theme");
    const hasThemeOverride = supportedThemes.has(themeOverride);
    const fontPreview = params.get("font");

    function applyTheme(theme) {
        const value = supportedThemes.has(theme) ? theme : "emerald";
        if (document.body.dataset.boardTheme !== value) {
            document.body.dataset.boardTheme = value;
        }
    }

    function applyFontPreview(font) {
        if (supportedFontPreviews.has(font)) document.documentElement.dataset.boardFont = font;
        else delete document.documentElement.dataset.boardFont;
    }

    function refresh(state) {
        if (hasThemeOverride) applyTheme(themeOverride);
        else applyTheme(state && state.theme);
    }

    applyTheme(hasThemeOverride ? themeOverride : "emerald");
    applyFontPreview(fontPreview);
    window.addEventListener("masjidpi:board-view", event => refresh(event.detail));
    if (window.MasjidBoardCurrentView) refresh(window.MasjidBoardCurrentView);
})();
