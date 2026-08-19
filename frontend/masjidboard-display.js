(() => {
    "use strict";

    const refreshIntervalMs = 60_000;

    const displayState = document.getElementById("displayState");
    const unconfiguredState = document.getElementById("unconfiguredState");
    const loadErrorState = document.getElementById("loadErrorState");
    const boardHeaders = document.getElementById("boardHeaders");
    const prayerGrid = document.getElementById("prayerGrid");
    const jumuahSection = document.getElementById("jumuahSection");
    const jumuahGrid = document.getElementById("jumuahGrid");
    const currentTime = document.getElementById("currentTime");
    const currentDate = document.getElementById("currentDate");
    const connectionState = document.getElementById("connectionState");

    function formatClock(time) {
        if (!time || !Number.isInteger(time.hour) || !Number.isInteger(time.minute)) {
            return "";
        }
        return `${String(time.hour).padStart(2, "0")}:${String(time.minute).padStart(2, "0")}`;
    }

    function updateClock() {
        const now = new Date();
        currentTime.textContent = now.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        });
        currentDate.textContent = now.toLocaleDateString([], {
            weekday: "long",
            day: "numeric",
            month: "long",
        });
    }

    function showOnly(element) {
        for (const candidate of [displayState, unconfiguredState, loadErrorState]) {
            candidate.classList.toggle("hidden", candidate !== element);
        }
    }

    function setGridCount(count) {
        const value = String(Math.max(1, count));
        boardHeaders.style.setProperty("--board-count", value);
        jumuahGrid.style.setProperty("--board-count", value);
    }

    function makeElement(tag, className, text) {
        const element = document.createElement(tag);
        if (className) {
            element.className = className;
        }
        if (text !== undefined) {
            element.textContent = text;
        }
        return element;
    }

    function boardStateClass(board) {
        if (board.status === "stale") {
            return "stale";
        }
        if (board.status === "unavailable") {
            return "unavailable";
        }
        return "";
    }

    function renderHeaders(boards) {
        boardHeaders.replaceChildren();
        boardHeaders.append(makeElement("div", "board-header-spacer"));

        for (const board of boards) {
            const header = makeElement("article", `board-header ${boardStateClass(board)}`.trim());
            header.append(makeElement("h2", "board-name", board.name));

            const status = makeElement("div", `board-status ${boardStateClass(board)}`.trim());
            if (board.status === "stale") {
                status.textContent = "Using last updated timetable";
            } else if (board.status === "unavailable") {
                status.textContent = "Timetable unavailable";
            }
            header.append(status);
            boardHeaders.append(header);
        }
    }

    function findPrayer(board, key) {
        return Array.isArray(board.prayers)
            ? board.prayers.find((prayer) => prayer.key === key)
            : undefined;
    }

    function appendTimeLine(cell, label, time, dominant, single) {
        if (!time) {
            return;
        }
        const line = makeElement(
            "div",
            `time-line ${dominant ? "dominant" : "secondary"}${single ? " single-time" : ""}`,
        );
        line.append(makeElement("span", "time-label", label));
        line.append(makeElement("span", "time-value", formatClock(time)));
        cell.append(line);
    }

    function renderPrayerCell(board, prayer) {
        const cell = makeElement("div", `prayer-cell ${boardStateClass(board)}`.trim());

        if (board.status === "unavailable" && !prayer) {
            cell.append(makeElement("div", "unavailable-copy", "No timetable data"));
            return cell;
        }

        if (!prayer) {
            cell.append(makeElement("div", "unavailable-copy", "Not provided"));
            return cell;
        }

        const hasAdhan = Boolean(prayer.adhan);
        const hasJamaah = Boolean(prayer.jamaah);
        const onlyOne = Number(hasAdhan) + Number(hasJamaah) === 1;

        // Reading order is always Adhan then Jamaah. Jamaah is dominant when
        // both exist; a sole available value inherits the dominant style.
        appendTimeLine(cell, "Adhan", prayer.adhan, onlyOne, onlyOne);
        appendTimeLine(cell, "Jamaah", prayer.jamaah, hasJamaah, onlyOne);

        if (!hasAdhan && !hasJamaah) {
            cell.append(makeElement("div", "unavailable-copy", "Not provided"));
        }
        return cell;
    }

    function renderPrayers(boards) {
        prayerGrid.replaceChildren();
        const prayers = [
            ["fajr", "Fajr"],
            ["dhuhr", "Dhuhr"],
            ["asr", "Asr"],
            ["maghrib", "Maghrib"],
            ["esha", "Esha"],
        ];

        for (const [key, label] of prayers) {
            const row = makeElement("div", "prayer-row");
            row.style.setProperty("--board-count", String(Math.max(1, boards.length)));
            row.append(makeElement("div", "prayer-label", label));
            for (const board of boards) {
                row.append(renderPrayerCell(board, findPrayer(board, key)));
            }
            prayerGrid.append(row);
        }
    }

    function eventTime(service, heading) {
        if (!service || !Array.isArray(service.events)) {
            return null;
        }
        const event = service.events.find((item) => item.heading === heading && item.time);
        return event ? event.time : null;
    }

    function renderJumuahCell(board) {
        const cell = makeElement("div", `jumuah-cell ${boardStateClass(board)}`.trim());
        const service = Array.isArray(board.jumuah) ? board.jumuah[0] : null;

        if (!service) {
            cell.append(makeElement("div", "unavailable-copy", "No Jumu’ah time provided"));
            return cell;
        }

        const adhan = service.adhan || eventTime(service, "Adhan");
        const salaah = service.effective_salaah || service.jamaah || eventTime(service, "Khutbah");
        const hasAdhan = Boolean(adhan);
        const hasSalaah = Boolean(salaah);
        const onlyOne = Number(hasAdhan) + Number(hasSalaah) === 1;

        appendTimeLine(cell, "Adhan", adhan, onlyOne, onlyOne);
        appendTimeLine(cell, "Salaah", salaah, hasSalaah, onlyOne);

        if (!hasAdhan && !hasSalaah) {
            const timedEvents = Array.isArray(service.events)
                ? service.events.filter((event) => event.time)
                : [];
            if (timedEvents.length === 0) {
                cell.append(makeElement("div", "unavailable-copy", "Times not provided"));
            }
        }

        if (Array.isArray(service.events)) {
            for (const event of service.events) {
                if (!event.time || event.heading === "Adhan" || event.heading === "Khutbah") {
                    continue;
                }
                const row = makeElement("div", "jumuah-event");
                row.append(makeElement("span", "", event.heading));
                row.append(makeElement("strong", "", formatClock(event.time)));
                cell.append(row);
            }
        }
        return cell;
    }

    function renderJumuah(boards) {
        const hasAnyJumuah = boards.some((board) => Array.isArray(board.jumuah) && board.jumuah.length > 0);
        jumuahSection.classList.toggle("hidden", !hasAnyJumuah);
        if (!hasAnyJumuah) {
            jumuahGrid.replaceChildren();
            return;
        }

        jumuahGrid.replaceChildren();
        jumuahGrid.append(makeElement("div", "jumuah-label-spacer"));
        for (const board of boards) {
            jumuahGrid.append(renderJumuahCell(board));
        }
    }

    function render(view) {
        if (!view || !view.configured) {
            showOnly(unconfiguredState);
            return;
        }

        const boards = Array.isArray(view.boards) ? view.boards.slice(0, 3) : [];
        if (boards.length === 0) {
            showOnly(unconfiguredState);
            return;
        }

        setGridCount(boards.length);
        renderHeaders(boards);
        renderPrayers(boards);
        renderJumuah(boards);
        showOnly(displayState);
    }

    async function refresh() {
        try {
            const response = await fetch("/api/masjidboard/display", {cache: "no-store"});
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            const view = await response.json();
            render(view);
            connectionState.textContent = "";
            connectionState.classList.remove("warning");
        } catch (error) {
            connectionState.textContent = "Connection interrupted";
            connectionState.classList.add("warning");
            if (displayState.classList.contains("hidden")) {
                showOnly(loadErrorState);
            }
            console.warn("MasjidBoard display refresh failed", error);
        }
    }

    updateClock();
    window.setInterval(updateClock, 1_000);
    refresh();
    window.setInterval(refresh, refreshIntervalMs);
})();
