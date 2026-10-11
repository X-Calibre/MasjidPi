(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("profile") === "appliance-720") return;
    const communityFixtureMode = params.get("notice-fixtures");
    const useCommunityFixtures = communityFixtureMode === "1" || communityFixtureMode === "new";
    const useJumuahKhateebFixture = params.get("jumuah-fixture") === "khateeb";
    const useDuaAfterAdhanFixture = params.get("dua-fixture") === "1";
    const displayUtils = window.MasjidBoardDisplayUtils;
    const dateUtils = window.MasjidBoardDate;
    const warningUtils = window.MasjidBoardWarningUtils;
    const {communityPriorityGroups, dailyIslamicItems, duaAfterAdhanItem, duaAfterAdhanWindowMinutes, fixtureCommunityItems, formatNoticeDate, formatRand, formatUpdatedAt, orderedCommunityItemGroups, orderedFields, plainText, specialDhuhrItem} = window.MasjidBoardCommunityUtils;

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
    let communityPageNodes = [];
    let communitySignature = "";
    let communityPage = 0;
    let latestCommunityView = null;
    let duaAfterAdhanVisible = false;

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

        const {start: zawaalStart, end: zawaalEnd} = warningUtils.zawaalWindow(astronomical);
        if (validTime(zawaalStart)) {
            const range = validTime(zawaalEnd) && minutes(zawaalEnd) !== minutes(zawaalStart)
                ? `${format(zawaalStart)}–${format(zawaalEnd)}`
                : format(zawaalStart);
            add(items, "Zawaal / Istiwa", zawaalStart, range);
        }

        const specialDhuhr = specialDhuhrItem(board);
        if (specialDhuhr) add(items, specialDhuhr.label, specialDhuhr.time);

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
        list.style.setProperty("--daily-time-columns", String(items.length));
        if (items.length > 10) list.classList.add("is-crowded");
        if (items.length > 11) list.classList.add("is-very-crowded");
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
        return item.type === "jumuah_schedule" || item.type === "dua_after_adhan" || plainText(item.body).length > 180 || orderedFields(item).length > 5;
    }

    function isRefinedNotice(item) {
        return ["funeral", "nikah", "eid", "announcement", "urgent_announcement", "programme", "class_change", "class_time_change", "well_wishes", "dawah", "three_day_jamaat", "new_moon"].includes(item.type);
    }

    function landscapeNoticeFields(item) {
        const fields = orderedFields(item);
        const groups = {
            dawah: [["Masjid Taleem"], ["Gasht Out", "Out Time"], ["Gasht In", "In Time"]],
            three_day_jamaat: [["First Jamaat", "First Date"], ["Second Jamaat", "Second Date"]],
            new_moon: [["Birth Date", "Birth"], ["Best Visibility", "Visibility Date"]],
        }[item.type];
        if (!groups) return fields;
        const used = new Set();
        const grouped = groups.map((labels) => {
            const selected = labels.flatMap((label) => fields.filter((field) => field.label === label));
            selected.forEach((field) => used.add(field));
            return {label: labels[0], value: selected.map((field) => field.value).join("\n")};
        }).filter((field) => field.value);
        return [...grouped, ...fields.filter((field) => !used.has(field))];
    }

    function isFullHeightNotice(item) {
        return ["contribution", "funeral", "nikah", "eid", "jumuah_schedule", "dua_after_adhan"].includes(item.type) || isDetailedCommunityItem(item);
    }

    function packCommunityPages(items) {
        const pages = [];
        const dailyItems = items.filter(isDailyItem);
        const regularItems = items.filter((item) => item.type !== "economic" && !isDailyItem(item));
        let index = 0;
        while (index < regularItems.length) {
            const first = regularItems[index];
            const second = regularItems[index + 1];
            if (isFullHeightNotice(first) || !second || isFullHeightNotice(second)) {
                pages.push({layout: "single", entries: [{item: first, span: 1}]});
                index += 1;
            } else {
                pages.push({layout: "halves", entries: [{item: first, span: 1}, {item: second, span: 1}]});
                index += 2;
            }
        }
        if (dailyItems.length) {
            pages.push({layout: "daily", entries: dailyItems.map((item) => ({item, span: 1}))});
        }
        for (const item of items.filter((entry) => entry.type === "economic")) {
            pages.push({layout: "single", entries: [{item, span: 1}]});
        }
        return pages;
    }

    function isDailyItem(item) {
        return ["daily_ayah", "daily_hadith", "daily_sunnah"].includes(item.type);
    }

    function renderCommunityCard(item, span) {
        const card = makeElement("article", `detailed-community-card detailed-community-${item.type} community-span-${span}`);
        if (isRefinedNotice(item)) card.classList.add("detailed-refined-notice");
		const contentLength = plainText(item.title).length + plainText(item.body).length + orderedFields(item).reduce((total, field) => total + field.value.length, 0);
		if (contentLength > 300) card.classList.add("content-long");
		if (contentLength > 600) card.classList.add("content-very-long");
		if (item.typeLabel) {
			const typeLabel = makeElement("div", "detailed-community-type", item.typeLabel);
			typeLabel.dir = "auto";
			card.append(typeLabel);
		}
		const title = makeElement("h2", "detailed-community-title", item.title);
		title.dir = "auto";
		card.append(title);
		if (item.type === "jumuah_schedule") {
			const schedule = makeElement("div", "detailed-jumuah-schedule");
			for (const event of item.schedule || []) {
				const column = makeElement("div", "detailed-jumuah-event");
				column.append(makeElement("span", "", event.heading), makeElement("strong", "", event.time));
				schedule.append(column);
			}
			card.append(schedule);
		}
		if (item.type === "daily_ayah" && plainText(item.fields.ayah_number)) {
			card.append(makeElement("div", "detailed-daily-ayah-number", plainText(item.fields.ayah_number)));
		}
        if (item.body) {
            const body = makeElement("p", "detailed-community-body", item.body);
            if (isDailyItem(item)) {
                body.replaceChildren(makeElement("span", "detailed-daily-text", item.body));
                card.classList.add("detailed-daily-card");
            }
            if (isRefinedNotice(item)) {
                body.replaceChildren(makeElement("span", "detailed-notice-text", item.body));
                if (/[\u0600-\u06ff]/.test(item.body)) body.classList.add("detailed-notice-arabic");
            }
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

        const fields = item.type === "salaah_change" ? [] : (isRefinedNotice(item) ? landscapeNoticeFields(item) : orderedFields(item)).filter((field) =>
            (item.type !== "economic" || field.label !== "Updated at") &&
            (item.type !== "daily_ayah" || field.label !== "Ayah")
        );
        if (fields.length > 0) {
            const list = makeElement("div", "detailed-community-fields");
            if (item.type === "contribution" && fields.length > 4) {
                list.classList.add("detailed-contribution-fields-many");
            }
            for (const field of fields) {
                const row = makeElement("div", "detailed-community-field");
                if ((item.type === "funeral" && field.label === "Relation") || (item.type === "nikah" && ["Groom", "Family"].includes(field.label))) {
                    row.classList.add("detailed-notice-family");
                }
                if ((item.type === "funeral" && ["Janazah", "Venue"].includes(field.label)) || (item.type === "eid" && field.label === "Salaah") || (item.type === "nikah" && field.label === "Time")) {
                    row.classList.add("detailed-notice-primary");
                }
                if (item.type === "contribution" && field.label === "Account Number") {
                    row.classList.add("detailed-contribution-account-number");
                }
                if (item.type === "dua_after_adhan") {
                    row.classList.add("detailed-dua-" + field.label.toLowerCase());
                }
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
                    if (item.type === "dua_after_adhan" && field.label === "Arabic") {
                        value.dir = "rtl";
                        value.lang = "ar";
                    }
                    if (item.type === "dua_after_adhan" && ["Arabic", "Translation"].includes(field.label)) {
                        value.replaceChildren(makeElement("span", "detailed-dua-text", field.value));
                    }
                }
                row.append(value);
                list.append(row);
            }
            card.append(list);
        }
        if (item.type === "economic") {
            const footer = makeElement("footer", "detailed-community-footer");
            if (item.notes) {
                const notes = makeElement("p", "detailed-economic-note", item.notes);
                notes.dir = "auto";
                footer.append(notes);
            }
            if (item.fields.updated_at) {
                footer.append(makeElement("div", "detailed-community-retrieved", `Updated at ${item.fields.updated_at}`));
            }
            footer.append(makeElement("div", "detailed-community-source", `From ${item.source}`));
            card.append(footer);
        } else if (item.type !== "dua_after_adhan") {
            const source = makeElement("div", "detailed-community-source", `Source: ${item.source}`);
            source.dir = "auto";
            card.append(source);
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

        const pageIndex = communityPage % communityPages.length;
        const page = communityPages[pageIndex];
        communityCards.className = `detailed-community-cards community-layout-${page.layout}`;
        communityCards.replaceChildren(...communityPageNodes[pageIndex]);
        communityCards.style.gridTemplateRows = page.layout === "daily"
            ? page.entries.map(({item}) => `${Math.max(1, Math.min(2.5, Math.sqrt(plainText(item.body).length / 100)))}fr`).join(" ")
            : "";
        scheduleDailyFit();
    }

    let dailyFitPending = false;
    let dailyFitSignature = "";
    function scheduleDailyFit() {
        if (dailyFitPending) return;
        dailyFitPending = true;
        window.requestAnimationFrame(() => {
            dailyFitPending = false;
            if (!communityCards || !communityCards.clientHeight || !communityPages.length) return;
            const pageIndex = communityPage % communityPages.length;
            const signature = `${communitySignature}:${pageIndex}:${communityCards.clientWidth}:${communityCards.clientHeight}`;
            if (signature === dailyFitSignature) return;
            dailyFitSignature = signature;
            // The root font follows the 1920x1080 reference design (16px)
            // and doubles at 4K. Keep measured font limits on the same scale.
            const rootFontSize = Number.parseFloat(window.getComputedStyle(document.documentElement).fontSize) || 16;
            const fontScale = rootFontSize / 16;
            const scaledFont = (size) => Math.round(size * fontScale);
            for (const row of communityCards.querySelectorAll(".detailed-dua-arabic, .detailed-dua-translation")) {
                const value = row.querySelector(".detailed-community-field-value");
                const text = row.querySelector(".detailed-dua-text");
                if (!text) continue;
                const arabic = row.classList.contains("detailed-dua-arabic");
                let low = scaledFont(arabic ? 28 : 18);
                let high = scaledFont(arabic ? 48 : 30);
                const fitsAt = (size) => {
                    value.style.fontSize = `${size}px`;
                    return text.scrollHeight <= value.clientHeight && text.scrollWidth <= value.clientWidth;
                };
                while (high - low > 1) {
                    const middle = Math.floor((low + high) / 2);
                    if (fitsAt(middle)) low = middle;
                    else high = middle;
                }
                fitsAt(fitsAt(high) ? high : low);
            }
            let noticesFit = true;
            for (const card of communityCards.querySelectorAll(".detailed-refined-notice")) {
                const body = card.querySelector(".detailed-community-body");
                const text = card.querySelector(".detailed-notice-text");
                const fitsAt = (size) => {
                    card.style.setProperty("--notice-value-size", `${size}px`);
                    if (body) body.style.fontSize = `${size}px`;
                    return card.scrollHeight <= card.clientHeight && (!text || text.scrollHeight <= body.clientHeight && text.scrollWidth <= body.clientWidth);
                };
                let low = scaledFont(18);
                let high = scaledFont(body && !card.querySelector(".detailed-community-fields") ? 32 : 28);
                if (!fitsAt(low)) noticesFit = false;
                while (high - low > 1) {
                    const middle = Math.floor((low + high) / 2);
                    if (fitsAt(middle)) low = middle;
                    else high = middle;
                }
                fitsAt(fitsAt(high) ? high : low);
            }
            let fits = true;
            for (const card of communityCards.querySelectorAll(".detailed-daily-card")) {
                const body = card.querySelector(".detailed-community-body");
                const text = card.querySelector(".detailed-daily-text");
                const fitsAt = (size) => {
                    body.style.fontSize = `${size}px`;
                    return text.scrollHeight <= body.clientHeight && text.scrollWidth <= body.clientWidth && card.scrollHeight <= card.clientHeight;
                };
                let low = scaledFont(16);
                let high = scaledFont(32);
                if (!fitsAt(low)) fits = false;
                while (high - low > 1) {
                    const middle = Math.floor((low + high) / 2);
                    if (fitsAt(middle)) low = middle;
                    else high = middle;
                }
                fitsAt(fitsAt(high) ? high : low);
            }
            const page = communityPages[pageIndex];
            if (page.entries.length > 1 && ((!fits && page.layout === "daily") || !noticesFit)) {
                communityPages.splice(pageIndex, 1, ...page.entries.map((entry) => ({layout: page.layout === "daily" ? "daily" : "single", entries: [entry]})));
                dailyFitSignature = "";
                buildCommunityPageNodes();
                renderCommunityPage();
            }
        });
    }

    window.addEventListener("resize", () => { dailyFitSignature = ""; scheduleDailyFit(); });
    if (typeof ResizeObserver !== "undefined" && communityCards) {
        new ResizeObserver(() => scheduleDailyFit()).observe(communityCards);
    }
    if (document.fonts && document.fonts.ready) {
        document.fonts.ready.then(() => { dailyFitSignature = ""; scheduleDailyFit(); });
    }

    function buildCommunityPageNodes() {
        communityPageNodes = communityPages.map((page) =>
            page.entries.map((entry) => renderCommunityCard(entry.item, entry.span))
        );
    }

    function renderCommunityContent(boards, dailyContent, indicators, showDuaAfterAdhan) {
        if (!communityPanel) return;
        if (boards.length === 0) {
            communityPanel.classList.add("hidden");
            return;
        }
		let itemGroups = boards.flatMap((board, boardIndex) => useCommunityFixtures
			? (boardIndex === 0 ? communityPriorityGroups(fixtureCommunityItems(communityFixtureMode)) : [])
			: orderedCommunityItemGroups(board, displayUtils.displayNow(), dateUtils.isIslamicFriday));
		if (useJumuahKhateebFixture) {
			for (const item of itemGroups.flat().filter((entry) => entry.type === "jumuah_schedule")) item.body = "Khateeb: To be announced";
		}
		const primaryBoard = boards[0];
		const duaItem = duaAfterAdhanItem(primaryBoard ? [primaryBoard] : [], displayUtils.displayNow(), showDuaAfterAdhan || useDuaAfterAdhanFixture, duaAfterAdhanWindowMinutes, useDuaAfterAdhanFixture);
		duaAfterAdhanVisible = Boolean(duaItem);
		if (duaItem) {
			itemGroups = [[duaItem]];
		} else {
			const sharedItems = dailyIslamicItems(dailyContent);
			const economicItem = economicCommunityItem(indicators);
			if (economicItem) sharedItems.push(economicItem);
			itemGroups.push(sharedItems);
		}
		const items = itemGroups.flat();
        const signature = JSON.stringify(items);
        if (signature !== communitySignature) {
            communitySignature = signature;
            communityItems = items;
            communityPages = itemGroups.flatMap((group) => packCommunityPages(group));
            buildCommunityPageNodes();
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
                gold_21: formatRand(indicators.gold_21_carat_per_gram),
                gold_18: formatRand(indicators.gold_18_carat_per_gram),
                gold_14: formatRand(indicators.gold_14_carat_per_gram),
                gold_9: formatRand(indicators.gold_9_carat_per_gram),
                silver: formatRand(indicators.silver_per_gram),
                minimum_mahr: formatRand(indicators.minimum_mahr),
                mahr_faatimi: formatRand(indicators.mahr_faatimi),
                updated_at: formatUpdatedAt(indicators.updated_at),
            },
            notes: plainText(indicators.notes),
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
            latestCommunityView = view;
            const boards = Array.isArray(view.boards) ? view.boards.slice(0, 3) : [];
            document.body.dataset.boardCount = String(Math.max(1, boards.length));
            render(boards[0] || null);
			renderCommunityContent(boards, view.daily_islamic_content, view.economic_indicators, view.show_dua_after_adhan);
            addSharedPrayerLabels();
    }

    window.addEventListener("masjidframe:board-view", event => refresh(event.detail));
    if (window.MasjidBoardCurrentView) refresh(window.MasjidBoardCurrentView);
    window.setInterval(() => {
        if (!latestCommunityView) return;
        const boards = Array.isArray(latestCommunityView.boards) ? latestCommunityView.boards.slice(0, 3) : [];
        const duaVisibleNow = Boolean(duaAfterAdhanItem(
            boards[0] ? [boards[0]] : [],
            displayUtils.displayNow(),
            latestCommunityView.show_dua_after_adhan || useDuaAfterAdhanFixture,
            duaAfterAdhanWindowMinutes,
            useDuaAfterAdhanFixture,
        ));
        if (duaVisibleNow !== duaAfterAdhanVisible) {
            renderCommunityContent(boards, latestCommunityView.daily_islamic_content, latestCommunityView.economic_indicators, latestCommunityView.show_dua_after_adhan);
        }
    }, 1000);
    window.setInterval(() => {
        if (communityPages.length <= 1) return;
        communityPage = (communityPage + 1) % communityPages.length;
        renderCommunityPage();
    }, communityPageIntervalMs);
})();
