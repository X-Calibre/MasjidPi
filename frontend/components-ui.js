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
