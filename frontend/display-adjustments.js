(() => {
    "use strict";

    const filters = {
        mild: "0.96 0 0 0 0  0 0.985 0 0 0  0 0 1.02 0 0  0 0 0 1 0",
        medium: "0.92 0 0 0 0  0 0.97 0 0 0  0 0 1.04 0 0  0 0 0 1 0",
        strong: "0.88 0 0 0 0  0 0.95 0 0 0  0 0 1.07 0 0  0 0 0 1 0"
    };
    const enabled = new URLSearchParams(window.location.search).get("profile") === "appliance-720";

    function ensureFilters() {
        if (document.getElementById("masjidpiDisplayFilters")) return;
        const host = document.createElementNS("http://www.w3.org/2000/svg", "svg");
        host.id = "masjidpiDisplayFilters";
        host.setAttribute("width", "0");
        host.setAttribute("height", "0");
        host.setAttribute("aria-hidden", "true");
        host.style.position = "absolute";
        for (const [name, values] of Object.entries(filters)) {
            const filter = document.createElementNS(host.namespaceURI, "filter");
            filter.id = `masjidpi-cool-${name}`;
            filter.setAttribute("color-interpolation-filters", "sRGB");
            const matrix = document.createElementNS(host.namespaceURI, "feColorMatrix");
            matrix.setAttribute("type", "matrix");
            matrix.setAttribute("values", values);
            filter.append(matrix);
            host.append(filter);
        }
        document.body.append(host);
    }

    function apply(colorTemperature = "off") {
        if (!enabled) colorTemperature = "off";
        ensureFilters();
        const value = filters[colorTemperature] ? `url(#masjidpi-cool-${colorTemperature})` : "";
        document.querySelectorAll(".board-shell,.setup-shell,.keyboard,.picker-sheet")
            .forEach(element => { element.style.filter = value; });
        document.documentElement.dataset.displayTemperature = colorTemperature;
    }

    async function request(settings) {
        const options = settings ? {
            method: "PUT",
            headers: {"Content-Type":"application/json"},
            body: JSON.stringify(settings)
        } : {};
        const response = await fetch("/api/display/settings", options);
        if (!response.ok) throw new Error(`Display settings request failed (${response.status})`);
        const result = await response.json();
        apply(result.color_temperature);
        window.dispatchEvent(new CustomEvent("masjidpi:display-settings", {detail:result}));
        return result;
    }

    window.MasjidPiDisplaySettings = {apply, load:() => request(), save:request};
    if (enabled) request().catch(() => apply("off"));
})();
