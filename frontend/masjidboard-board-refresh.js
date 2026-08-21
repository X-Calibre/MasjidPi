(() => {
    "use strict";

    const button = document.getElementById("refreshBoardsButton");
    if (!button) return;

    const successKey = "masjidboard-refresh-success";
    const banner = document.getElementById("configBanner");

    if (sessionStorage.getItem(successKey) === "1") {
        sessionStorage.removeItem(successKey);
        if (banner) {
            banner.textContent = "Timetables refreshed.";
            banner.className = "config-banner success";
            window.setTimeout(() => banner.classList.add("hidden"), 6000);
        }
    }

    button.addEventListener("click", async () => {
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

            sessionStorage.setItem(successKey, "1");
            window.location.reload();
        } catch (error) {
            if (banner) {
                banner.textContent = `Could not refresh timetables: ${error.message}`;
                banner.className = "config-banner error";
            }
            button.disabled = false;
            button.textContent = "Refresh Timetables";
        }
    });
})();
