(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") !== "detailed") return;

    document.body.classList.add("detailed-layout");
    const panel = document.getElementById("additionalTimes");
    const prayerGrid = document.getElementById("prayerGrid");
    const detailedGregorianDate = document.getElementById("detailedGregorianDate");
    const detailedIslamicDate = document.getElementById("detailedIslamicDate");
    const refreshIntervalMs = 60_000;
    const islamicWeekdays = [
        "Al-Ahad",
        "Al-Ithnayn",
        "Ath-Thulatha",
        "Al-Arbi'a",
        "Al-Khamis",
        "Al-Jumu'ah",
        "As-Sabt",
    ];

    function validTime(time) {
        return time && Number.isInteger(time.hour) && Number.isInteger(time.minute);
    }

    function minutes(time) {
        return time.hour * 60 + time.minute;
    }

    function format(time) {
        return `${String(time.hour).padStart(2, "0")}:${String(time.minute).padStart(2, "0")}`;
    }

    function add(items, label, time, valueText = "") {
        if (validTime(time)) items.push({label, time, valueText: valueText || format(time)});
    }

    function findPrayer(board, key) {
        return board && Array.isArray(board.prayers)
            ? board.prayers.find((prayer) => prayer.key === key)
            : null;
    }

    function chooseAsrStart(board, astronomical) {
        const candidates = [astronomical.asr_shafii, astronomical.asr_hanafi].filter(validTime);
        if (candidates.length === 0) return null;

        const asr = findPrayer(board, "asr");
        const reference = asr && validTime(asr.adhan) ? minutes(asr.adhan) : null;
        if (reference !== null) {
            const beforeAdhan = candidates
                .filter((time) => minutes(time) <= reference)
                .sort((left, right) => minutes(right) - minutes(left));
            if (beforeAdhan.length > 0) return beforeAdhan[0];
        }

        return validTime(astronomical.asr_hanafi) ? astronomical.asr_hanafi : astronomical.asr_shafii;
    }

    function buildItems(board) {
        const astronomical = board && board.astronomical ? board.astronomical : {};
        const items = [];

        add(items, "Sehri Ends", astronomical.suhur);
        add(items, "Fajr Starts", astronomical.fajr_start);
        add(items, "Sunrise", astronomical.sunrise);
        add(items, "Ishraq", astronomical.ishraaq);

        const zawaalStart = validTime(astronomical.istiwa_caution)
            ? astronomical.istiwa_caution
            : astronomical.istiwa;
        const zawaalEnd = validTime(astronomical.zawaal_end)
            ? astronomical.zawaal_end
            : astronomical.istiwa;
        if (validTime(zawaalStart)) {
            const range = validTime(zawaalEnd) && minutes(zawaalEnd) !== minutes(zawaalStart)
                ? `${format(zawaalStart)}–${format(zawaalEnd)}`
                : format(zawaalStart);
            add(items, "Zawaal / Istiwa", zawaalStart, range);
        }

        add(items, "Asr Starts", chooseAsrStart(board, astronomical));
        add(items, "Sunset", astronomical.sunset);
        add(items, "Esha Starts", astronomical.esha_start);

        items.sort((left, right) => minutes(left.time) - minutes(right.time));
        return items;
    }

    function boardGregorianDate(board) {
        const value = board && board.date ? board.date.gregorian : "";
        if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return new Date();
        const [year, month, day] = value.split("-").map(Number);
        return new Date(year, month - 1, day);
    }

    function formatGregorianDate(board) {
        return boardGregorianDate(board).toLocaleDateString([], {
            weekday: "long",
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    function islamicWeekday(board) {
        const date = boardGregorianDate(board);
        const sunset = board && board.astronomical ? board.astronomical.sunset : null;
        if (validTime(sunset)) {
            const now = new Date();
            const rolloverSeconds = sunset.hour * 3600 + sunset.minute * 60 + 185;
            const nowSeconds = now.getHours() * 3600 + now.getMinutes() * 60 + now.getSeconds();
            if (nowSeconds >= rolloverSeconds) date.setDate(date.getDate() + 1);
        }
        return islamicWeekdays[date.getDay()];
    }

    function renderDetailedDates(board) {
        if (detailedGregorianDate) {
            detailedGregorianDate.textContent = formatGregorianDate(board);
        }
        if (!detailedIslamicDate) return;
        const islamicDate = board && board.date ? board.date.islamic : "";
        detailedIslamicDate.textContent = islamicDate ? `${islamicWeekday(board)}, ${islamicDate}` : "";
    }

    function render(board) {
        renderDetailedDates(board);
        panel.replaceChildren();

        const heading = document.createElement("div");
        heading.className = "additional-times-heading";
        heading.textContent = "Daily Times";
        panel.append(heading);

        const source = document.createElement("div");
        source.className = "additional-times-source";
        source.textContent = board ? `From ${board.name}` : "From first masjid";
        panel.append(source);

        const items = buildItems(board);
        if (items.length === 0) {
            const empty = document.createElement("div");
            empty.className = "additional-times-empty";
            empty.textContent = "No daily times available";
            panel.append(empty);
            return;
        }

        const list = document.createElement("div");
        list.className = "additional-times-list";
        for (const item of items) {
            const row = document.createElement("div");
            row.className = "additional-time-row";
            const label = document.createElement("span");
            label.className = "additional-time-label";
            label.textContent = item.label;
            const value = document.createElement("span");
            value.className = "additional-time-value";
            value.textContent = item.valueText;
            row.append(label, value);
            list.append(row);
        }
        panel.append(list);
    }

    function transformRow(row, preferredLabels) {
        if (row.querySelector(":scope > .shared-time-labels")) return;

        const cells = Array.from(row.children).filter((child) =>
            child.classList.contains("prayer-cell") || child.classList.contains("jumuah-cell")
        );
        if (cells.length === 0) return;

        const available = new Set();
        for (const cell of cells) {
            for (const label of cell.querySelectorAll(":scope > .time-line > .time-label")) {
                const text = label.textContent.trim();
                if (text) available.add(text);
            }
        }
        const labels = preferredLabels.filter((label) => available.has(label));
        if (labels.length === 0) return;

        row.style.setProperty("--time-slot-count", String(labels.length));
        const shared = document.createElement("div");
        shared.className = "shared-time-labels";

        for (const wanted of labels) {
            const label = document.createElement("div");
            label.className = `shared-time-label shared-${wanted.toLowerCase().replace(/[^a-z0-9]+/g, "-")}`;
            label.textContent = wanted;
            shared.append(label);
        }

        const prayerLabel = row.querySelector(":scope > .prayer-label-card");
        if (!prayerLabel) return;
        prayerLabel.insertAdjacentElement("afterend", shared);

        for (const cell of cells) {
            const existing = new Map(
                Array.from(cell.querySelectorAll(":scope > .time-line")).map((line) => [line.querySelector(".time-label")?.textContent.trim(), line])
            );
            const fragment = document.createDocumentFragment();
            for (const wanted of labels) {
                const line = existing.get(wanted);
                if (line) {
                    line.classList.add("shared-label-value");
                    fragment.append(line);
                } else {
                    const placeholder = document.createElement("div");
                    placeholder.className = "time-line shared-label-value shared-time-placeholder";
                    const dash = document.createElement("span");
                    dash.className = "shared-missing-value";
                    dash.textContent = "—";
                    placeholder.append(dash);
                    fragment.append(placeholder);
                }
            }
            cell.replaceChildren(fragment);
        }
    }

    function addSharedPrayerLabels() {
        if (!prayerGrid) return;

        for (const row of prayerGrid.querySelectorAll(".prayer-row:not(.jumuah-row)")) {
            transformRow(row, ["Adhan", "Jamaah"]);
        }
        for (const row of prayerGrid.querySelectorAll(".prayer-row.jumuah-row")) {
            transformRow(row, ["Adhan", "Sunan", "Khutbah"]);
        }
    }

    let transformPending = false;
    const observer = new MutationObserver(() => {
        if (transformPending) return;
        transformPending = true;
        window.requestAnimationFrame(() => {
            transformPending = false;
            addSharedPrayerLabels();
        });
    });
    if (prayerGrid) observer.observe(prayerGrid, {childList: true, subtree: true});

    async function refresh() {
        try {
            const response = await fetch("/api/masjidboard/display", {cache: "no-store"});
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const view = await response.json();
            const boards = Array.isArray(view.boards) ? view.boards.slice(0, 3) : [];
            document.body.dataset.boardCount = String(Math.max(1, boards.length));
            render(boards[0] || null);
            addSharedPrayerLabels();
        } catch (error) {
            console.warn("Detailed MasjidBoard times refresh failed", error);
        }
    }

    refresh();
    window.setInterval(refresh, refreshIntervalMs);
})();
