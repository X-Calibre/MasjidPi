(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") !== "detailed") return;

    document.body.classList.add("detailed-layout");
    const panel = document.getElementById("additionalTimes");
    const prayerGrid = document.getElementById("prayerGrid");
    const refreshIntervalMs = 60_000;

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

    function render(board) {
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

    function addSharedPrayerLabels() {
        if (!prayerGrid) return;

        for (const row of prayerGrid.querySelectorAll(".prayer-row:not(.jumuah-row)")) {
            if (row.querySelector(":scope > .shared-time-labels")) continue;

            const cells = Array.from(row.querySelectorAll(":scope > .prayer-cell"));
            const labels = ["Adhan", "Jamaah"].filter((wanted) =>
                cells.some((cell) => Array.from(cell.querySelectorAll(".time-label")).some((label) => label.textContent === wanted))
            );
            if (labels.length === 0) continue;

            row.style.setProperty("--time-slot-count", String(labels.length));
            const shared = document.createElement("div");
            shared.className = "shared-time-labels";

            for (const wanted of labels) {
                const label = document.createElement("div");
                label.className = `shared-time-label shared-${wanted.toLowerCase()}`;
                label.textContent = wanted;
                shared.append(label);
            }

            const prayerLabel = row.querySelector(":scope > .prayer-label-card");
            prayerLabel.insertAdjacentElement("afterend", shared);

            for (const cell of cells) {
                const existing = new Map(
                    Array.from(cell.querySelectorAll(":scope > .time-line")).map((line) => [line.querySelector(".time-label")?.textContent, line])
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
                        fragment.append(placeholder);
                    }
                }
                cell.replaceChildren(fragment);
            }
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
            const boards = Array.isArray(view.boards) ? view.boards : [];
            render(boards[0] || null);
            addSharedPrayerLabels();
        } catch (error) {
            console.warn("Detailed MasjidBoard times refresh failed", error);
        }
    }

    refresh();
    window.setInterval(refresh, refreshIntervalMs);
})();
