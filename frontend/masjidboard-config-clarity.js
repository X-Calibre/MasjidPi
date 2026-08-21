(() => {
    "use strict";

    const statusEndpoint = "/api/masjidboard/status";
    const catalogueEndpoint = "/api/masjidboard/catalogue";
    const selectionEndpoint = "/api/masjidboard/selection";
    const statusRoot = document.getElementById("masjidBoardStatus");
    const banner = document.getElementById("configBanner");

    if (!statusRoot) return;

    let bannerBusy = false;

    function text(tag, className, value) {
        const element = document.createElement(tag);
        if (className) element.className = className;
        element.textContent = value;
        return element;
    }

    function stateLabel(value) {
        switch (value) {
        case "current": return "Current";
        case "stale": return "Stale";
        case "unavailable": return "Unavailable";
        default: return "Unknown";
        }
    }

    function formatTimestamp(value) {
        if (!value) return "Not yet";
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return "Unknown";

        const now = new Date();
        const sameDay = date.getFullYear() === now.getFullYear()
            && date.getMonth() === now.getMonth()
            && date.getDate() === now.getDate();

        const time = date.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"});
        if (sameDay) return `Today, ${time}`;

        return date.toLocaleString([], {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function detailRow(label, value, className = "") {
        const row = document.createElement("div");
        row.className = `status-board-detail${className ? ` ${className}` : ""}`;
        row.append(text("span", "status-board-detail-label", label));
        row.append(text("span", "status-board-detail-value", value));
        return row;
    }

    function renderStatus(status) {
        const boards = Array.isArray(status.boards) ? status.boards : [];
        statusRoot.replaceChildren();

        const summary = document.createElement("div");
        summary.className = "status-summary";
        summary.append(detailRow("Configured", status.configured ? "Yes" : "No"));
        summary.append(detailRow("Selected Masjids", String(boards.length)));
        statusRoot.append(summary);

        if (!boards.length) return;

        const list = document.createElement("div");
        list.className = "status-board-list status-board-list-detailed";

        for (const board of boards) {
            const value = board.status || "unknown";
            const item = document.createElement("article");
            item.className = `status-board-card status-${value}`;

            const heading = document.createElement("div");
            heading.className = "status-board-heading";
            heading.append(text("strong", "status-board-name", board.name || board.catalogue_id || "MasjidBoard"));
            heading.append(text("span", `status-board-state ${value}`, stateLabel(value)));
            item.append(heading);

            const details = document.createElement("div");
            details.className = "status-board-details";

            if (board.using_cached_data || value === "stale") {
                details.append(text("div", "status-board-note warning", "Using cached timetable data"));
            }
            if (value === "unavailable" && !board.board) {
                details.append(text("div", "status-board-note error", "No timetable data is currently available"));
            }

            if (board.last_successful_update) {
                details.append(detailRow("Last successful update", formatTimestamp(board.last_successful_update)));
            }
            if (board.last_attempt) {
                details.append(detailRow("Last attempt", formatTimestamp(board.last_attempt)));
            }
            if (board.update_error) {
                details.append(detailRow("Update error", board.update_error, "error"));
            }
            if (board.persistence_error) {
                details.append(detailRow("Cache error", board.persistence_error, "error"));
            }

            item.append(details);
            list.append(item);
        }

        statusRoot.append(list);
    }

    async function refreshDetailedStatus() {
        try {
            const response = await fetch(statusEndpoint, {cache: "no-store"});
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            renderStatus(await response.json());
        } catch (_) {
            // Leave the primary configuration script's error/status output intact.
        }
    }

    function normalizeRefreshButtons() {
        const hierarchy = document.getElementById("refreshHierarchyButton");
        const catalogue = document.getElementById("refreshCatalogueButton");
        const boards = document.getElementById("refreshBoardsButton");
        if (hierarchy && hierarchy.textContent === "Refresh Locations") hierarchy.textContent = "Refresh Location List";
        if (catalogue && catalogue.textContent === "Refresh Masjids") catalogue.textContent = "Refresh Masjid List";
        if (boards && boards.textContent === "Refresh Boards") boards.textContent = "Refresh Timetables";
    }

    async function updateSuccessBanner() {
        if (!banner || bannerBusy) return;
        const message = banner.textContent.trim();
        if (!message) return;

        bannerBusy = true;
        try {
            if (message === "Locations saved and MasjidBoard catalogue updated.") {
                const response = await fetch(catalogueEndpoint, {cache: "no-store"});
                if (response.ok) {
                    const data = await response.json();
                    const records = Array.isArray(data.records) ? data.records.filter(item => item.status === "active") : [];
                    banner.textContent = `Locations saved · ${records.length} Masjid${records.length === 1 ? "" : "s"} available`;
                }
            } else if (message === "Selected Masjids saved. The display has been updated.") {
                const response = await fetch(selectionEndpoint, {cache: "no-store"});
                if (response.ok) {
                    const data = await response.json();
                    const boards = Array.isArray(data.boards) ? data.boards : [];
                    banner.textContent = `${boards.length} Masjid${boards.length === 1 ? "" : "s"} saved · display updated`;
                }
            } else if (message === "MasjidBoard catalogue refreshed.") {
                const response = await fetch(catalogueEndpoint, {cache: "no-store"});
                if (response.ok) {
                    const data = await response.json();
                    const records = Array.isArray(data.records) ? data.records.filter(item => item.status === "active") : [];
                    banner.textContent = `Masjid list refreshed · ${records.length} Masjid${records.length === 1 ? "" : "s"} available`;
                }
            } else if (message === "Location hierarchy refreshed.") {
                banner.textContent = "Location list refreshed.";
            }
        } catch (_) {
            // The original success message is still valid if enrichment fails.
        } finally {
            bannerBusy = false;
        }
    }

    const statusObserver = new MutationObserver(() => {
        if (!statusRoot.querySelector(".status-board-card") && !statusRoot.querySelector(".status-summary")) {
            window.setTimeout(refreshDetailedStatus, 0);
        }
    });
    statusObserver.observe(statusRoot, {childList: true, subtree: true});

    if (banner) {
        const bannerObserver = new MutationObserver(() => {
            normalizeRefreshButtons();
            updateSuccessBanner();
        });
        bannerObserver.observe(banner, {childList: true, characterData: true, subtree: true, attributes: true});
    }

    for (const id of ["refreshHierarchyButton", "refreshCatalogueButton", "refreshBoardsButton"]) {
        const button = document.getElementById(id);
        if (!button) continue;
        new MutationObserver(normalizeRefreshButtons).observe(button, {childList: true, characterData: true, subtree: true});
    }

    normalizeRefreshButtons();
    window.setTimeout(refreshDetailedStatus, 50);
})();
