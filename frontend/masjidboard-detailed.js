(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") !== "detailed") return;
    const useCommunityFixtures = params.get("notice-fixtures") === "1";

    document.body.classList.add("detailed-layout");
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

    function renderDetailedDates(board) {
        if (detailedGregorianDate) {
            detailedGregorianDate.textContent = window.MasjidBoardDate.formatGregorianDate(board);
        }
        if (!detailedIslamicDate) return;
        const islamicDate = board && board.date ? board.date.islamic : "";
        detailedIslamicDate.textContent = islamicDate ? `${window.MasjidBoardDate.islamicWeekday(board)}, ${islamicDate}` : "";
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

    function makeElement(tag, className, text) {
        const element = document.createElement(tag);
        if (className) element.className = className;
        if (text !== undefined) element.textContent = text;
        return element;
    }

    function plainText(value) {
        const source = String(value || "").replace(/<br\s*\/?\s*>/gi, "\n");
        if (!source) return "";
        const parsed = new DOMParser().parseFromString(source, "text/html");
        return (parsed.body.textContent || "").replace(/\r/g, "").replace(/\n{3,}/g, "\n\n").trim();
    }

    function fieldLabel(name) {
        const labels = {
            address: "Address", bride: "Bride", cemetery: "Cemetery", date: "Date",
            groom_relation: "Groom", lecture: "Lecture", name: "Name", name_one: "Name",
            name_two: "Name", pickup: "Pickup", relation: "Relation", relation_one: "Family",
            relation_two: "Family", salaah: "Salaah", salaah_time: "Janazah", salaah_venue: "Venue",
            time: "Time", venue: "Venue",
        };
        return labels[name] || name.replace(/_/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
    }

    function orderedFields(item) {
        const fields = item && item.fields && typeof item.fields === "object" ? item.fields : {};
        const preferred = {
            funeral: ["relation", "salaah_time", "salaah_venue", "pickup", "cemetery", "address"],
            nikah: ["groom_relation", "relation_one", "relation_two", "date", "time"],
            eid: ["date", "venue", "address", "lecture", "salaah"],
        }[item.type] || [];
        const titleFields = new Set(item.type === "funeral" ? ["name"] : item.type === "nikah" ? ["name_one", "name_two", "bride"] : []);
        const names = [...preferred, ...Object.keys(fields).sort()].filter((name, index, all) =>
            !titleFields.has(name) && all.indexOf(name) === index && plainText(fields[name])
        );
        return names.slice(0, 6).map((name) => ({label: fieldLabel(name), value: plainText(fields[name])}));
    }

    function noticeTitle(notice) {
        const fields = notice && notice.fields && typeof notice.fields === "object" ? notice.fields : {};
        if (notice.type === "funeral") return plainText(fields.name) || plainText(notice.title) || "Funeral Notice";
        if (notice.type === "nikah") {
            const names = [fields.name_one, fields.bride || fields.name_two].map(plainText).filter(Boolean);
            if (names.length > 0) return names.join(" & ");
        }
        return plainText(notice.title) || `${fieldLabel(notice.type || "general")} Notice`;
    }

    function collectCommunityItems(boards) {
        const items = [];
        const seen = new Set();
        function add(item) {
            const key = JSON.stringify([item.type, item.title, item.body, item.fields]);
            if (!item.title && !item.body && Object.keys(item.fields || {}).length === 0) return;
            if (seen.has(key)) return;
            seen.add(key);
            items.push(item);
        }

        for (const board of boards) {
            for (const notice of Array.isArray(board.notices) ? board.notices : []) {
                const fields = notice && notice.fields && typeof notice.fields === "object" ? notice.fields : {};
                add({
                    type: plainText(notice.type).toLowerCase() || "general",
                    title: noticeTitle(notice),
                    body: Object.keys(fields).length === 0 ? plainText(notice.content) : "",
                    fields,
                    source: board.name,
                });
            }
            for (const announcement of Array.isArray(board.announcements) ? board.announcements : []) {
                add({
                    type: "announcement",
                    title: plainText(announcement.title) || "Masjid Announcement",
                    body: plainText(announcement.content),
                    fields: {},
                    source: board.name,
                });
            }
        }
        return items;
    }

    function fixtureCommunityItems() {
        const source = "Layout test fixture — not live";
        return [
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
                type: "announcement",
                title: "Masjid Announcement",
                body: "تذكير: سيكون البرنامج بعد صلاة العشاء بإذن الله. نرجو من الجميع الحضور في الوقت المحدد.",
                fields: {},
                source,
            },
        ];
    }

    function communityTypeLabel(type) {
        return {announcement: "Announcement", eid: "Eid Notice", funeral: "Funeral Notice", nikah: "Nikah Notice", well_wishes: "Well Wishes"}[type]
            || `${fieldLabel(type || "general")} Notice`;
    }

    function renderCommunityCard(item) {
        const card = makeElement("article", `detailed-community-card detailed-community-${item.type}`);
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

        const start = (communityPage * 2) % communityItems.length;
        const count = Math.min(2, communityItems.length - start);
        for (let offset = 0; offset < count; offset += 1) {
            communityCards.append(renderCommunityCard(communityItems[start + offset]));
        }
        communityCards.classList.toggle("single", count === 1);
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
            console.warn("Detailed MasjidBoard times refresh failed", error);
        }
    }

    refresh();
    window.setInterval(refresh, refreshIntervalMs);
    window.setInterval(() => {
        if (communityItems.length <= 2) return;
        communityPage = (communityPage + 1) % Math.ceil(communityItems.length / 2);
        renderCommunityPage();
    }, communityPageIntervalMs);
})();
