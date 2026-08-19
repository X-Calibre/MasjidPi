(() => {
    "use strict";

    const refreshIntervalMs = 60_000;

    const displayState = document.getElementById("displayState");
    const unconfiguredState = document.getElementById("unconfiguredState");
    const loadErrorState = document.getElementById("loadErrorState");
    const boardHeaders = document.getElementById("boardHeaders");
    const prayerGrid = document.getElementById("prayerGrid");
    const currentTime = document.getElementById("currentTime");
    const currentDate = document.getElementById("currentDate");
    const connectionState = document.getElementById("connectionState");

    function displayDate() {
        const value = new URLSearchParams(window.location.search).get("date");
        if (!value || !/^\d{4}-\d{2}-\d{2}$/.test(value)) {
            return new Date();
        }

        const [year, month, day] = value.split("-").map(Number);
        const date = new Date(year, month - 1, day);
        if (
            date.getFullYear() !== year ||
            date.getMonth() !== month - 1 ||
            date.getDate() !== day
        ) {
            console.warn(`Ignoring invalid MasjidBoard display date override: ${value}`);
            return new Date();
        }
        return date;
    }

    function formatClock(time) {
        if (!time || !Number.isInteger(time.hour) || !Number.isInteger(time.minute)) {
            return "";
        }
        return `${String(time.hour).padStart(2, "0")}:${String(time.minute).padStart(2, "0")}`;
    }

    function minutesSinceMidnight(time) {
        return time.hour * 60 + time.minute;
    }

    function updateClock() {
        const now = new Date();
        const date = displayDate();
        currentTime.textContent = now.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        });
        currentDate.textContent = date.toLocaleDateString([], {
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
        boardHeaders.style.setProperty("--board-count", String(Math.max(1, count)));
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

    function appendTimeLine(cell, label, time, dominant, single, extraClass = "") {
        if (!time) {
            return;
        }
        const line = makeElement(
            "div",
            `time-line ${dominant ? "dominant" : "secondary"}${single ? " single-time" : ""}${extraClass ? ` ${extraClass}` : ""}`,
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
            return cell;
        }

        const hasAdhan = Boolean(prayer.adhan);
        const hasJamaah = Boolean(prayer.jamaah);
        const onlyOne = Number(hasAdhan) + Number(hasJamaah) === 1;

        appendTimeLine(cell, "Adhan", prayer.adhan, onlyOne, onlyOne);
        appendTimeLine(cell, "Jamaah", prayer.jamaah, hasJamaah, onlyOne);
        return cell;
    }

    function eventTime(service, heading) {
        if (!service || !Array.isArray(service.events)) {
            return null;
        }
        const event = service.events.find((item) => item.heading === heading && item.time);
        return event ? event.time : null;
    }

    function sameTime(left, right) {
        return left && right && left.hour === right.hour && left.minute === right.minute;
    }

    function jumuahItems(service) {
        const adhan = service.adhan || eventTime(service, "Adhan");
        const salaah = service.effective_salaah || service.jamaah || eventTime(service, "Khutbah");
        const items = [];

        if (adhan) {
            items.push({label: "Adhan", time: adhan, kind: "adhan"});
        }

        if (Array.isArray(service.events)) {
            for (const event of service.events) {
                if (!event.time || event.heading === "Adhan" || event.heading === "Khutbah") {
                    continue;
                }
                if (sameTime(event.time, adhan) || sameTime(event.time, salaah)) {
                    continue;
                }
                items.push({label: event.heading, time: event.time, kind: "event"});
            }
        }

        if (salaah) {
            items.push({label: "Salaah", time: salaah, kind: "salaah"});
        }

        items.sort((left, right) => minutesSinceMidnight(left.time) - minutesSinceMidnight(right.time));
        return items;
    }

    function renderJumuahCell(board) {
        const cell = makeElement("div", `jumuah-cell ${boardStateClass(board)}`.trim());
        const service = Array.isArray(board.jumuah) ? board.jumuah[0] : null;

        if (!service) {
            return cell;
        }

        for (const item of jumuahItems(service)) {
            appendTimeLine(
                cell,
                item.label,
                item.time,
                item.kind === "salaah",
                false,
                `jumuah-${item.kind}`,
            );
        }
        return cell;
    }

    function appendPrayerRow(boards, key, label) {
        const row = makeElement("div", "prayer-row");
        row.style.setProperty("--board-count", String(Math.max(1, boards.length)));
        row.append(makeElement("div", "prayer-label", label));
        for (const board of boards) {
            row.append(renderPrayerCell(board, findPrayer(board, key)));
        }
        prayerGrid.append(row);
    }

    function appendJumuahRow(boards) {
        const row = makeElement("div", "prayer-row jumuah-row");
        row.style.setProperty("--board-count", String(Math.max(1, boards.length)));
        row.append(makeElement("div", "prayer-label", "Jumu’ah"));
        for (const board of boards) {
            row.append(renderJumuahCell(board));
        }
        prayerGrid.append(row);
    }

    function renderPrayers(boards) {
        prayerGrid.replaceChildren();

        const friday = displayDate().getDay() === 5;
        prayerGrid.classList.toggle("friday", friday);

        appendPrayerRow(boards, "fajr", "Fajr");
        if (friday) {
            appendJumuahRow(boards);
        } else {
            appendPrayerRow(boards, "dhuhr", "Dhuhr");
        }
        appendPrayerRow(boards, "asr", "Asr");
        appendPrayerRow(boards, "maghrib", "Maghrib");
        appendPrayerRow(boards, "esha", "Esha");
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
