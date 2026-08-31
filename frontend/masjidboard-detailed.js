(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("profile") === "appliance") return;
    const communityFixtureMode = params.get("notice-fixtures");
    const useCommunityFixtures = communityFixtureMode === "1" || communityFixtureMode === "new";
    const {collectCommunityItems, fixtureCommunityItems, formatNoticeDate, formatRand, formatUpdatedAt, orderedFields, plainText} = window.MasjidBoardCommunityUtils;

    document.documentElement.classList.add("landscape-layout");
    document.body.classList.add("landscape-layout");
    const panel = document.getElementById("additionalTimes");
    const prayerGrid = document.getElementById("prayerGrid");
    const detailedGregorianDate = document.getElementById("detailedGregorianDate");
    const detailedIslamicDate = document.getElementById("detailedIslamicDate");
    const communityPanel = document.getElementById("detailedCommunityPanel");
    const communityCards = document.getElementById("detailedCommunityCards");
    const communityEmpty = document.getElementById("detailedCommunityEmpty");
    const communityPageIntervalMs = 15_000;
    let communityItems = [];
    let communityPages = [];
    let communitySignature = "";
    let communityPage = 0;

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

    function buildItems(board) {
        const astronomical = board && board.astronomical ? board.astronomical : {};
        const items = [];

        add(items, "Sehri Ends", astronomical.suhur);
        add(items, "Fajr Starts", astronomical.fajr_start);
        add(items, "Sunrise", astronomical.sunrise);
        add(items, "Ishraq", astronomical.ishraaq);
        add(items, "Duha / Chaasht", astronomical.duha);

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

        add(items, "Asr Shafi‘i", astronomical.asr_shafii);
        add(items, "Asr Hanafi", astronomical.asr_hanafi);
        add(items, "Sunset", astronomical.sunset);
        add(items, "Esha Starts", astronomical.esha_start);

        items.sort((left, right) => minutes(left.time) - minutes(right.time));
        return items;
    }

    function renderLandscapeDates(board) {
        if (detailedGregorianDate) {
            detailedGregorianDate.textContent = window.MasjidBoardDate.formatGregorianDate(board);
        }
        if (!detailedIslamicDate) return;
        const islamicDate = board && board.date ? board.date.islamic : "";
        detailedIslamicDate.textContent = islamicDate ? `${window.MasjidBoardDate.islamicWeekday(board)}, ${islamicDate}` : "";
    }

    function render(board) {
        renderLandscapeDates(board);
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

    function makeElement(tag, className, text) {
        const element = document.createElement(tag);
        if (className) element.className = className;
        if (text !== undefined) element.textContent = text;
        return element;
    }

    function isDetailedCommunityItem(item) {
        return plainText(item.body).length > 180 || orderedFields(item).length > 5;
    }

    function packCommunityPages(items) {
        const pages = [];
        const regularItems = items.filter((item) => item.type !== "economic");
        let index = 0;
        while (index < regularItems.length) {
            const remaining = regularItems.length - index;
            if (remaining === 1) {
                pages.push({layout: "single", entries: [{item: regularItems[index], span: 1}]});
                index += 1;
                continue;
            }
            if (remaining === 2) {
                pages.push({layout: "halves", entries: [
                    {item: regularItems[index], span: 1}, {item: regularItems[index + 1], span: 1},
                ]});
                index += 2;
                continue;
            }

            const first = regularItems[index];
            const second = regularItems[index + 1];
            const third = regularItems[index + 2];
            const firstDetailed = isDetailedCommunityItem(first);
            const secondDetailed = isDetailedCommunityItem(second);
            const thirdDetailed = isDetailedCommunityItem(third);

            if (firstDetailed && !secondDetailed) {
                pages.push({layout: "thirds", entries: [{item: first, span: 2}, {item: second, span: 1}]});
                index += 2;
            } else if (!firstDetailed && secondDetailed) {
                pages.push({layout: "thirds", entries: [{item: first, span: 1}, {item: second, span: 2}]});
                index += 2;
            } else if (!firstDetailed && !secondDetailed && !thirdDetailed) {
                pages.push({layout: "thirds", entries: [
                    {item: first, span: 1}, {item: second, span: 1}, {item: third, span: 1},
                ]});
                index += 3;
            } else {
                pages.push({layout: "halves", entries: [{item: first, span: 1}, {item: second, span: 1}]});
                index += 2;
            }
        }
        for (const item of items.filter((entry) => entry.type === "economic")) {
            pages.push({layout: "single", entries: [{item, span: 1}]});
        }
        return pages;
    }

    function renderCommunityCard(item, span) {
        const card = makeElement("article", `detailed-community-card detailed-community-${item.type} community-span-${span}`);
        card.append(makeElement("h2", "detailed-community-title", item.title));
        if (item.body) {
            const body = makeElement("p", "detailed-community-body", item.body);
            body.dir = "auto";
            card.append(body);
        }

        if (item.type === "salaah_change") {
            const main = makeElement("div", "salaah-change-main");
            main.append(
                makeElement("div", "salaah-change-effective", `Effective from\n${formatNoticeDate(item.fields.effective_date)}`),
                makeElement("div", "salaah-change-time", plainText(item.fields.new_time))
            );
            card.append(main);
        }

        const fields = item.type === "salaah_change" ? [] : orderedFields(item).filter((field) =>
            item.type !== "economic" || field.label !== "Retrieved at"
        );
        if (fields.length > 0) {
            const list = makeElement("div", "detailed-community-fields");
            for (const field of fields) {
                const row = makeElement("div", "detailed-community-field");
                row.append(makeElement("span", "detailed-community-field-label", field.label));
                const value = makeElement("span", "detailed-community-field-value");
                if (item.type === "economic" && field.value.startsWith("R")) {
                    value.classList.add("economic-accounting-value");
                    value.append(
                        makeElement("span", "economic-currency-symbol", "R"),
                        makeElement("span", "economic-amount", field.value.slice(1))
                    );
                } else {
                    value.textContent = field.value;
                    value.dir = "auto";
                }
                row.append(value);
                list.append(row);
            }
            card.append(list);
        }
        if (item.type === "economic") {
            const footer = makeElement("footer", "detailed-community-footer");
            if (item.fields.retrieved_at) {
                footer.append(makeElement("div", "detailed-community-retrieved", `Retrieved at ${item.fields.retrieved_at}`));
            }
            footer.append(makeElement("div", "detailed-community-source", `From ${item.source}`));
            card.append(footer);
        } else {
            card.append(makeElement("div", "detailed-community-source", `Source: ${item.source}`));
        }
        return card;
    }

    function renderCommunityPage() {
        if (!communityPanel || !communityCards || !communityEmpty) return;
        communityPanel.classList.remove("hidden");
        communityCards.replaceChildren();
        const hasItems = communityItems.length > 0;
        communityCards.classList.toggle("hidden", !hasItems);
        communityEmpty.classList.toggle("hidden", hasItems);
        if (!hasItems) return;

        const page = communityPages[communityPage % communityPages.length];
        communityCards.className = `detailed-community-cards community-layout-${page.layout}`;
        for (const entry of page.entries) {
            communityCards.append(renderCommunityCard(entry.item, entry.span));
        }
    }

    function renderCommunityContent(boards) {
        if (!communityPanel) return;
        if (boards.length === 0) {
            communityPanel.classList.add("hidden");
            return;
        }
        const items = useCommunityFixtures ? fixtureCommunityItems(communityFixtureMode) : collectCommunityItems(boards);
        const signature = JSON.stringify(items);
        if (signature !== communitySignature) {
            communitySignature = signature;
            communityItems = items;
            communityPages = packCommunityPages(items);
            communityPage = 0;
        }
        renderCommunityPage();
    }

    function economicCommunityItem(indicators) {
        if (!indicators) return null;
        const effectiveDate = new Date(`${indicators.effective_date}T12:00:00`);
        const dateText = Number.isNaN(effectiveDate.getTime()) ? indicators.effective_date : effectiveDate.toLocaleDateString("en-ZA", {day: "numeric", month: "short", year: "numeric"});
        return {
            type: "economic",
            title: "Islamic Economic Indicators",
            body: `Effective ${dateText}`,
            source: indicators.source,
            fields: {
                rand_dollar: formatRand(indicators.rand_dollar),
                nisaab: formatRand(indicators.nisaab),
                krugerrand: formatRand(indicators.krugerrand),
                gold_24: formatRand(indicators.gold_24_carat_per_gram),
                gold_22: formatRand(indicators.gold_22_carat_per_gram),
                gold_18: formatRand(indicators.gold_18_carat_per_gram),
                gold_14: formatRand(indicators.gold_14_carat_per_gram),
                gold_9: formatRand(indicators.gold_9_carat_per_gram),
                silver: formatRand(indicators.silver_per_gram),
                minimum_mahr: formatRand(indicators.minimum_mahr),
                mahr_faatimi: formatRand(indicators.mahr_faatimi),
                retrieved_at: formatUpdatedAt(indicators.fetched_at),
            },
        };
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
        const singleMaghribAdhan = prayerLabel.textContent.trim() === "Maghrib" &&
            labels.length === 1 && labels[0] === "Adhan";
        row.classList.toggle("single-adhan-as-jamaah", singleMaghribAdhan);
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

    function refresh(view) {
            const boards = Array.isArray(view.boards) ? view.boards.slice(0, 3) : [];
            document.body.dataset.boardCount = String(Math.max(1, boards.length));
            render(boards[0] || null);
            const economicItem = economicCommunityItem(view.economic_indicators);
            if (economicItem) {
                const originalFixtureMode = useCommunityFixtures;
                const baseItems = originalFixtureMode ? fixtureCommunityItems(communityFixtureMode) : collectCommunityItems(boards);
                const items = [...baseItems, economicItem];
                const signature = JSON.stringify(items);
                if (signature !== communitySignature) {
                    communitySignature = signature;
                    communityItems = items;
                    communityPages = packCommunityPages(items);
                    communityPage = 0;
                }
                renderCommunityPage();
            } else {
                renderCommunityContent(boards);
            }
            addSharedPrayerLabels();
    }

    window.addEventListener("masjidpi:board-view", event => refresh(event.detail));
    window.setInterval(() => {
        if (communityPages.length <= 1) return;
        communityPage = (communityPage + 1) % communityPages.length;
        renderCommunityPage();
    }, communityPageIntervalMs);
})();