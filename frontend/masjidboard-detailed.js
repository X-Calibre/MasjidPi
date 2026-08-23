(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") === "portrait") return;
    const communityFixtureMode = params.get("notice-fixtures");
    const useCommunityFixtures = communityFixtureMode === "1" || communityFixtureMode === "new";
    const {collectCommunityItems, communityTypeLabel, orderedFields, plainText} = window.MasjidBoardCommunityUtils;

    document.body.classList.add("landscape-layout");
    const panel = document.getElementById("additionalTimes");
    const prayerGrid = document.getElementById("prayerGrid");
    const detailedGregorianDate = document.getElementById("detailedGregorianDate");
    const detailedIslamicDate = document.getElementById("detailedIslamicDate");
    const communityPanel = document.getElementById("detailedCommunityPanel");
    const communityCards = document.getElementById("detailedCommunityCards");
    const communityEmpty = document.getElementById("detailedCommunityEmpty");
    const refreshIntervalMs = 60_000;
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

    function fixtureCommunityItems() {
        const source = "Layout test fixture — not live";
        const items = [
            {
                type: "funeral",
                title: "Marhoom Abdullah Ismail",
                body: "",
                fields: {
                    relation: "Father of Yusuf Ismail",
                    salaah_time: "After Zuhr · 13:15",
                    salaah_venue: "Central Masjid",
                    pickup: "12:45 from the family residence",
                    cemetery: "Central Cemetery",
                    address: "10 Example Road",
                },
                source,
            },
            {
                type: "nikah",
                title: "Muhammad & Ayesha",
                body: "",
                fields: {
                    groom_relation: "Son of Ahmad & Fatima",
                    relation_two: "Daughter of Ismail & Maryam",
                    date: "Saturday, 29 August 2026",
                    time: "After Asr · 16:45",
                    venue: "Masjid Hall",
                },
                source,
            },
            {
                type: "eid",
                title: "Eid Salaah Notice",
                body: "",
                fields: {
                    date: "Monday, 25 May 2026",
                    venue: "Community Sports Ground",
                    address: "1 Example Field Road",
                    lecture: "Lecture · 07:00",
                    salaah: "Salaah · 07:30",
                },
                source,
            },
            {
                type: "announcement",
                title: "Important Access Notice",
                body: "Please use the northern entrance while maintenance work is under way. The main parking area will be closed after Maghrib. Elderly worshippers and families may use the reserved drop-off area near the hall entrance. Please follow the directions of the volunteers on duty.",
                fields: {},
                source,
            },
			{
				type: "dawah",
				title: "Dawah and Gasht",
				body: "",
				fields: {masjid_taleem: "Daily after Esha Salaah", gasht_out_day: "Thursday", gasht_out_time: "After Asr", gasht_in_day: "Monday", gasht_in_time: "After Maghrib"},
				source,
			},
			{
				type: "three_day_jamaat",
				title: "Three-Day Jamaat",
				body: "",
				fields: {first_location: "Hartbeespoort area", first_date: "4–6 September", second_location: "Pretoria West", second_date: "11–13 September"},
				source,
			},
			{
				type: "contribution",
				title: "Masjid Contributions — Lillah Only",
				body: "",
				fields: {bank: "Example Bank", account_name: "Example Masjid Trust", branch_code: "123456", account_number: "000 123 456", bsb: ""},
				source,
			},
            {
                type: "salaah_change",
                title: "Esha Time Change",
                body: "",
                fields: {prayer: "Esha", effective_date: "1 September", new_time: "19:45"},
                source,
            },
            {
                type: "programme",
                title: "Taleem Programme",
                body: "Wednesday 11:15–12:15\nResident's home",
                fields: {},
                source,
            },
            {
                type: "new_moon",
                title: "New Moon Information",
                body: "",
                fields: {birth_date: "23 August", birth: "05:27", visibility_date: "24 August", best_visibility: "18:37"},
                source,
            },
            {
                type: "announcement",
                title: "Masjid Announcement",
                body: "تذكير: سيكون البرنامج بعد صلاة العشاء بإذن الله. نرجو من الجميع الحضور في الوقت المحدد.",
                fields: {},
                source,
            },
            {
                type: "well_wishes",
                title: "Du'a Requested",
                body: "The community is requested to make du'a for those who are unwell and for their families.",
                fields: {},
                source,
            },
            {
                type: "announcement",
                title: "Weekly Programme",
                body: "The weekly community programme will take place after Esha on Thursday evening.",
                fields: {},
                source,
            },
        ];
        if (communityFixtureMode === "new") {
			const newTypes = new Set(["dawah", "three_day_jamaat", "contribution"]);
            return items.filter((item) => newTypes.has(item.type));
        }
        return items;
    }

    function isDetailedCommunityItem(item) {
        return plainText(item.body).length > 180 || orderedFields(item).length > 5;
    }

    function packCommunityPages(items) {
        const pages = [];
        let index = 0;
        while (index < items.length) {
            const remaining = items.length - index;
            if (remaining === 1) {
                pages.push({layout: "single", entries: [{item: items[index], span: 1}]});
                index += 1;
                continue;
            }
            if (remaining === 2) {
                pages.push({layout: "halves", entries: [
                    {item: items[index], span: 1}, {item: items[index + 1], span: 1},
                ]});
                index += 2;
                continue;
            }

            const first = items[index];
            const second = items[index + 1];
            const third = items[index + 2];
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
        return pages;
    }

    function renderCommunityCard(item, span) {
        const card = makeElement("article", `detailed-community-card detailed-community-${item.type} community-span-${span}`);
        card.append(makeElement("div", "detailed-community-type", communityTypeLabel(item.type)));
        card.append(makeElement("h2", "detailed-community-title", item.title));
        if (item.body) {
            const body = makeElement("p", "detailed-community-body", item.body);
            body.dir = "auto";
            card.append(body);
        }

        const fields = orderedFields(item);
        if (fields.length > 0) {
            const list = makeElement("div", "detailed-community-fields");
            for (const field of fields) {
                const row = makeElement("div", "detailed-community-field");
                row.append(makeElement("span", "detailed-community-field-label", field.label));
                const value = makeElement("span", "detailed-community-field-value", field.value);
                value.dir = "auto";
                row.append(value);
                list.append(row);
            }
            card.append(list);
        }
        card.append(makeElement("div", "detailed-community-source", `From ${item.source}`));
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
        const items = useCommunityFixtures ? fixtureCommunityItems() : collectCommunityItems(boards);
        const signature = JSON.stringify(items);
        if (signature !== communitySignature) {
            communitySignature = signature;
            communityItems = items;
            communityPages = packCommunityPages(items);
            communityPage = 0;
        }
        renderCommunityPage();
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
            renderCommunityContent(boards);
            addSharedPrayerLabels();
        } catch (error) {
            console.warn("Landscape MasjidBoard times refresh failed", error);
        }
    }

    refresh();
    window.setInterval(refresh, refreshIntervalMs);
    window.setInterval(() => {
        if (communityPages.length <= 1) return;
        communityPage = (communityPage + 1) % communityPages.length;
        renderCommunityPage();
    }, communityPageIntervalMs);
})();
