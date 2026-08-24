(() => {
    "use strict";

    async function loadComponents() {
        const response = await fetch("/api/components", {cache: "no-store"});
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    }

    function applyVisibility(components) {
        document.querySelectorAll("[data-component]").forEach((element) => {
            const component = element.dataset.component;
            element.classList.toggle("hidden", !Boolean(components[component]));
        });
    }

    function enforcePage(components) {
        const page = document.body.dataset.componentPage;
        if (!page || components[page]) return;

        if (components.board) {
            window.location.replace("/masjidboard-config.html");
            return;
        }
        if (components.listen) {
            window.location.replace("/index.html");
        }
    }

    async function updateVersionFooter() {
        const footer = document.getElementById("version");
        if (!footer) return;
        try {
            const response = await fetch("/api/player/status", {cache: "no-store"});
            if (!response.ok) return;
            const status = await response.json();
            if (status.version) footer.textContent = "MasjidPi " + status.version;
        } catch (_) {
            // Keep the product name when version status is unavailable.
        }
    }

    updateVersionFooter();

    loadComponents()
        .then((components) => {
            applyVisibility(components);
            enforcePage(components);
            document.documentElement.dataset.componentsReady = "true";
        })
        .catch((error) => {
            console.warn("Unable to load installed MasjidPi components", error);
        });
})();
