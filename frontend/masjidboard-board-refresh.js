(() => {
    "use strict";

    const button = document.getElementById("refreshBoardsButton");
    if (!button) return;

    button.addEventListener("click", async () => {
        const original = button.textContent;
        button.disabled = true;
        button.textContent = "Refreshing…";

        try {
            const response = await fetch("/api/masjidboard/boards/refresh", {
                method: "POST",
                cache: "no-store",
            });
            if (!response.ok) {
                let message = `HTTP ${response.status}`;
                try {
                    const body = await response.json();
                    if (body && body.error) message = body.error;
                } catch (_) {}
                throw new Error(message);
            }

            // Reload the configuration view so the existing status renderer
            // immediately reflects the freshly fetched board data.
            window.location.reload();
        } catch (error) {
            const banner = document.getElementById("configBanner");
            if (banner) {
                banner.textContent = `Could not refresh selected boards: ${error.message}`;
                banner.className = "config-banner error";
            }
            button.disabled = false;
            button.textContent = original;
        }
    });
})();
