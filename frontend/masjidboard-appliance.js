(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("profile") !== "appliance") return;
    const communityFixtureMode = params.get("notice-fixtures");
    const useCommunityFixtures = communityFixtureMode === "1" || communityFixtureMode === "new";
    const useJumuahKhateebFixture = params.get("jumuah-fixture") === "khateeb";
    const useDuaAfterAdhanFixture = params.get("dua-fixture") === "1";

    document.body.classList.add("appliance-layout");
    const utils = window.MasjidBoardDisplayUtils;
    const dateUtils = window.MasjidBoardDate;
    const warningUtils = window.MasjidBoardWarningUtils;
    const {communityPriorityGroups, dailyIslamicItems, duaAfterAdhanItem, duaAfterAdhanWindowMinutes, fixtureCommunityItems, formatNoticeDate, formatRand, formatUpdatedAt, orderedCommunityItemGroups, orderedFields, plainText, specialDhuhrItem} = window.MasjidBoardCommunityUtils;
    const state = document.getElementById("applianceState");
    const slidesHost = document.getElementById("applianceSlides");
    const dotsHost = document.getElementById("applianceDots");
    const primaryName = document.getElementById("appliancePrimaryName");
    const clock = document.getElementById("applianceClock");
    const gregorianDate = document.getElementById("applianceGregorianDate");
    const islamicDate = document.getElementById("applianceIslamicDate");
    const nextName = document.getElementById("applianceNextName");
    const nextTime = document.getElementById("applianceNextTime");
    const countdown = document.getElementById("applianceCountdown");
    const prayerLabels = {fajr: "Fajr", dhuhr: "Dhuhr", asr: "Asr", maghrib: "Maghrib", esha: "Esha"};
    let latestView = null;
    let slides = [];
    let activeSlide = 0;
    let slideDurationSeconds = 15;
    let slideTimer = 0;
    let transitionTimer = 0;
    let gestureStart = null;
    let duaAfterAdhanVisible = false;

    function element(tag, className, text) {
        const node = document.createElement(tag);
        if (className) node.className = className;
        if (text !== undefined) node.textContent = text;
        return node;
    }

    function validTime(time) {
        return time && Number.isInteger(time.hour) && Number.isInteger(time.minute);
    }

    function minutes(time) { return time.hour * 60 + time.minute; }

    function formatEvent(event) {
        if (!event) return {name: "No upcoming event", time: null};
        if (event.kind === "jumuah") return {name: `Jumu'ah\n${event.label || "Salaah"}`, time: event.time};
        const prayer = prayerLabels[event.key] || event.key;
        return {name: `${prayer}\n${event.event === "jamaah" ? "Jamaah" : "Adhan"}`, time: event.time};
    }

    function dailyItems(board) {
        const astronomical = board && board.astronomical ? board.astronomical : {};
        const result = [];
        const add = (label, time, valueText = "") => {
            if (validTime(time)) result.push({label, time, valueText: valueText || utils.formatClock(time)});
        };

        add("Sehri Ends", astronomical.suhur);
        add("Fajr Starts", astronomical.fajr_start);
        add("Sunrise", astronomical.sunrise);
        add("Ishraq", astronomical.ishraaq);
        add("Duha / Chaasht", astronomical.duha);

        const {start: zawaalStart, end: zawaalEnd} = warningUtils.zawaalWindow(astronomical);
        if (validTime(zawaalStart)) {
            const range = validTime(zawaalEnd) && minutes(zawaalEnd) !== minutes(zawaalStart)
                ? utils.formatClock(zawaalStart) + "–" + utils.formatClock(zawaalEnd)
                : utils.formatClock(zawaalStart);
            add("Zawaal / Istiwa", zawaalStart, range);
        }

        const specialDhuhr = specialDhuhrItem(board);
        if (specialDhuhr) add(specialDhuhr.label, specialDhuhr.time);

        add("Asr Shafi‘i", astronomical.asr_shafii);
        add("Asr Hanafi", astronomical.asr_hanafi);
        add("Sunset", astronomical.sunset);
        add("Esha Starts", astronomical.esha_start);
        return result.sort((left, right) => minutes(left.time) - minutes(right.time));
    }

    function isCompactCommunityItem(item) {
        return plainText(item.body).length <= 80 && orderedFields(item).length <= 2;
    }

    function salaahSlide(board) {
        const slide = element("article", "appliance-slide");
        const heading = element("header", "appliance-slide-heading");
        heading.append(element("small", "", "SALAAH TIMES"), element("h2", "", board.name));
        slide.append(heading);
        const table = element("div", "appliance-times");
        const head = element("div", "appliance-time-row appliance-time-head");
        head.append(element("span", "", "Salaah"), element("span", "", "Adhan"), element("span", "", "Jamaah"));
        table.append(head);
        const now = utils.displayNow();
        const friday = dateUtils.isIslamicFriday(board, now);
        const next = utils.nextEventForBoard(board, now, friday);
        for (const key of ["fajr", "dhuhr", "asr", "maghrib", "esha"]) {
            if (friday && key === "dhuhr") {
                const service = Array.isArray(board.jumuah) ? board.jumuah[0] : null;
                const items = service ? utils.jumuahItems(service) : [];
                const adhan = items.find((item) => item.kind === "adhan");
                // jumuahItems identifies an explicit Salaah/Jamaah as the
                // salaah item, or uses Khutbah when no such time is supplied.
                const jamaah = items.find((item) => item.kind === "salaah");
                const row = element("div", "appliance-time-row");
                if (next && next.kind === "jumuah") row.classList.add("upcoming");
                row.append(
                    element("span", "", "Jumu'ah"),
                    element("strong", "", adhan ? utils.formatClock(adhan.time) : "-"),
                    element("strong", "", jamaah ? utils.formatClock(jamaah.time) : "-"),
                );
                table.append(row);
                continue;
            }
            const prayer = utils.findPrayer(board, key);
            const row = element("div", "appliance-time-row");
            if (next && next.kind === "prayer" && next.key === key) row.classList.add("upcoming");
            row.append(
                element("span", "", prayerLabels[key]),
                element("strong", "", prayer && prayer.adhan ? utils.formatClock(prayer.adhan) : "-"),
                element("strong", "", prayer && prayer.jamaah ? utils.formatClock(prayer.jamaah) : "-"),
            );
            table.append(row);
        }
        slide.append(table);
        return slide;
    }

    function dailySlide(board) {
        const slide = element("article", "appliance-slide");
        const heading = element("header", "appliance-slide-heading");
        heading.append(element("small", "", "DAILY TIMES"), element("h2", "", board.name));
        slide.append(heading);
        const table = element("div", "appliance-times appliance-daily-times");
        for (const item of dailyItems(board)) {
            const row = element("div", "appliance-time-row appliance-daily-row");
            row.append(element("span", "", item.label), element("strong", "", item.valueText));
            table.append(row);
        }
        slide.append(table);
        return slide;
    }

    function communityCard(item, compact) {
        const card = element("article", "appliance-community-card appliance-community-" + item.type + (compact ? " compact" : ""));
		const contentLength = plainText(item.title).length + plainText(item.body).length + orderedFields(item).reduce((total, field) => total + field.value.length, 0);
		if (contentLength > 300) card.classList.add("content-long");
		if (contentLength > 600) card.classList.add("content-very-long");
		if (item.typeLabel) {
			const typeLabel = element("div", "appliance-community-type", item.typeLabel);
			typeLabel.dir = "auto";
			card.append(typeLabel);
		}
		const title = element("h2", "appliance-community-title", item.title);
		title.dir = "auto";
		card.append(title);
		if (item.type === "jumuah_schedule") {
			const schedule = element("div", "appliance-jumuah-schedule");
			for (const event of item.schedule || []) {
				const column = element("div", "appliance-jumuah-event");
				column.append(element("span", "", event.heading), element("strong", "", event.time));
				schedule.append(column);
			}
			card.append(schedule);
		}
		if (item.type === "daily_ayah" && plainText(item.fields.ayah_number)) {
			card.append(element("div", "appliance-daily-ayah-number", plainText(item.fields.ayah_number)));
		}
        if (item.body) {
            const body = element("p", "appliance-community-body", item.body);
            body.dir = "auto";
            card.append(body);
        }

        if (item.type === "salaah_change") {
            const main = element("div", "appliance-salaah-change-main");
            main.append(
                element("div", "appliance-salaah-change-effective", "Effective from\n" + formatNoticeDate(item.fields.effective_date)),
                element("div", "appliance-salaah-change-time", plainText(item.fields.new_time))
            );
            card.append(main);
        }

        const fields = item.type === "salaah_change" ? [] : orderedFields(item).filter((field) =>
            item.type !== "daily_ayah" || field.label !== "Ayah"
        );
        if (fields.length > 0) {
            const list = element("div", "appliance-community-fields");
            for (const field of fields) {
                const row = element("div", "appliance-community-field");
                row.append(element("span", "appliance-community-field-label", field.label));
                const value = element("span", "appliance-community-field-value", field.value);
                value.dir = "auto";
                row.append(value);
                list.append(row);
            }
            card.append(list);
        }
        if (item.type !== "dua_after_adhan") {
            const source = element("div", "appliance-community-source", "Source: " + item.source);
            source.dir = "auto";
            card.append(source);
        }
        return card;
    }

    function communitySlide(entries) {
        const paired = entries.length === 2;
        const compactSingle = entries.length === 1 && isCompactCommunityItem(entries[0]);
        const layout = paired ? "paired" : compactSingle ? "single compact-single" : "single";
        const slide = element("article", "appliance-slide appliance-community-slide " + layout);
        for (const item of entries) slide.append(communityCard(item, isCompactCommunityItem(item)));
        return slide;
    }

    function communitySlides(items) {
        const result = [];
        let pendingCompact = null;
        function flushCompact() {
            if (!pendingCompact) return;
            result.push(communitySlide([pendingCompact]));
            pendingCompact = null;
        }

        for (const item of items) {
            if (!isCompactCommunityItem(item)) {
                flushCompact();
                result.push(communitySlide([item]));
                continue;
            }
            if (pendingCompact) {
                result.push(communitySlide([pendingCompact, item]));
                pendingCompact = null;
            } else {
                pendingCompact = item;
            }
        }
        flushCompact();
        return result;
    }

    function economicSlide(indicators) {
        if (!indicators) return null;
        const slide = element("article", "appliance-slide appliance-economic-slide");
        const heading = element("header", "appliance-slide-heading");
        heading.append(element("h2", "", "Islamic Economic Indicators"));
        slide.append(heading);
        const effectiveDate = new Date(`${indicators.effective_date}T12:00:00`);
        const dateText = Number.isNaN(effectiveDate.getTime()) ? indicators.effective_date : effectiveDate.toLocaleDateString("en-ZA", {day: "numeric", month: "long", year: "numeric"});
        slide.append(element("div", "appliance-economic-date", `Effective ${dateText}`));
        const values = [
            ["Rand/Dollar", formatRand(indicators.rand_dollar)],
            ["Nisaab", formatRand(indicators.nisaab)],
            ["Krugerrand", formatRand(indicators.krugerrand)],
            ["Gold 24 ct / g", formatRand(indicators.gold_24_carat_per_gram)],
            ["Gold 22 ct / g", formatRand(indicators.gold_22_carat_per_gram)],
            ["Gold 18 ct / g", formatRand(indicators.gold_18_carat_per_gram)],
            ["Gold 14 ct / g", formatRand(indicators.gold_14_carat_per_gram)],
            ["Gold 9 ct / g", formatRand(indicators.gold_9_carat_per_gram)],
            ["Silver / g", formatRand(indicators.silver_per_gram)],
            ["Minimum Mahr", formatRand(indicators.minimum_mahr)],
            ["Mahr Faatimi", formatRand(indicators.mahr_faatimi)],
        ];
        const grid = element("div", "appliance-economic-values");
        for (const [label, value] of values) {
            const row = element("div", "appliance-economic-value");
            const accountingValue = element("strong", "appliance-economic-accounting");
            accountingValue.append(
                element("span", "appliance-economic-currency", "R"),
                element("span", "appliance-economic-amount", value.slice(1))
            );
            row.append(element("span", "", label), accountingValue);
            grid.append(row);
        }
        const footer = element("footer", "appliance-economic-footer");
        const retrievedAt = formatUpdatedAt(indicators.fetched_at);
        if (retrievedAt) footer.append(element("div", "appliance-economic-retrieved", `Retrieved at ${retrievedAt}`));
        footer.append(element("div", "appliance-economic-source", `From ${indicators.source}`));
        slide.append(grid, footer);
        return slide;
    }

    function showSlide(index, restart = true) {
        if (slides.length === 0) return;

        const nextIndex = (index + slides.length) % slides.length;
        const previousIndex = activeSlide;
        const direction = index < previousIndex ? -1 : 1;
        const previousSlide = slides[previousIndex];
        const nextSlide = slides[nextIndex];

        window.clearTimeout(transitionTimer);
        for (const slide of slides) {
            slide.classList.remove("enter-from-left", "enter-from-right", "exit-to-left", "exit-to-right");
        }

        if (previousSlide && nextSlide && nextIndex !== previousIndex) {
            const enterClass = direction > 0 ? "enter-from-right" : "enter-from-left";
            const exitClass = direction > 0 ? "exit-to-left" : "exit-to-right";
            previousSlide.classList.remove("active");
            previousSlide.classList.add(exitClass);
            nextSlide.classList.add("active", enterClass);
            void nextSlide.offsetWidth;
            window.requestAnimationFrame(() => nextSlide.classList.remove(enterClass));
            transitionTimer = window.setTimeout(() => previousSlide.classList.remove(exitClass), 320);
        } else if (nextSlide) {
            nextSlide.classList.add("active");
        }

        activeSlide = nextIndex;
        Array.from(dotsHost.children).forEach((dot, itemIndex) => dot.classList.toggle("active", itemIndex === activeSlide));
        if (restart) startTimer();
    }

    function startTimer() {
        window.clearInterval(slideTimer);
        if (state.classList.contains("listen-panel-open")) return;
        if (slides.length < 2) return;
        slideTimer = window.setInterval(() => showSlide(activeSlide + 1, false), slideDurationSeconds * 1000);
    }

    function renderSlides(boards, indicators, dailyContent, showDuaAfterAdhan) {
        slidesHost.replaceChildren();
        dotsHost.replaceChildren();
        slides = [];
        const primaryBoard = boards[0];
        const duaItem = duaAfterAdhanItem(primaryBoard ? [primaryBoard] : [], utils.displayNow(), showDuaAfterAdhan || useDuaAfterAdhanFixture, duaAfterAdhanWindowMinutes, useDuaAfterAdhanFixture);
        duaAfterAdhanVisible = Boolean(duaItem);
        if (duaItem) {
            slides = [communitySlide([duaItem])];
        } else {
            boards.forEach((board, boardIndex) => {
                slides.push(salaahSlide(board));
                if (boardIndex === 0 && dailyItems(board).length > 0) slides.push(dailySlide(board));
                const communityGroups = useCommunityFixtures && boardIndex === 0
                    ? communityPriorityGroups(fixtureCommunityItems(communityFixtureMode))
                    : useCommunityFixtures ? [] : orderedCommunityItemGroups(board, utils.displayNow(), dateUtils.isIslamicFriday);
                if (useJumuahKhateebFixture) {
                    for (const item of communityGroups.flat().filter((entry) => entry.type === "jumuah_schedule")) item.body = "Khateeb: To be announced";
                }
                for (const group of communityGroups) slides.push(...communitySlides(group));
            });
            for (const item of dailyIslamicItems(dailyContent)) slides.push(communitySlide([item]));
            const indicatorsSlide = economicSlide(indicators);
            if (indicatorsSlide) slides.push(indicatorsSlide);
        }
        slides.forEach((slide, index) => {
            slidesHost.append(slide);
            const dot = element("button", "", "");
            dot.type = "button";
            dot.setAttribute("aria-label", `Show slide ${index + 1}`);
            dot.addEventListener("click", () => showSlide(index));
            dotsHost.append(dot);
        });
        showSlide(Math.min(activeSlide, slides.length - 1));
    }

    function updateHeader() {
        if (!latestView || !latestView.boards || !latestView.boards[0]) return;
        const board = latestView.boards[0];
        const now = utils.displayNow();
        const friday = dateUtils.isIslamicFriday(board, now);
        const event = utils.nextEventForBoard(board, now, friday);
        const formatted = formatEvent(event);
        clock.textContent = now.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit", hour12: false});
        gregorianDate.textContent = dateUtils.formatGregorianDate(board);
        islamicDate.textContent = board.date && board.date.islamic
            ? `${dateUtils.islamicWeekday(board, now)}, ${board.date.islamic}`
            : "";
        nextName.textContent = formatted.name;
        nextTime.textContent = formatted.time ? utils.formatClock(formatted.time) : "--:--";
        countdown.textContent = formatted.time ? utils.countdownText(formatted.time, now, event.dayOffset || 0) : "";

        const boards = latestView.boards.slice(0, 3);
        const duaVisibleNow = Boolean(duaAfterAdhanItem(
            boards[0] ? [boards[0]] : [],
            now,
            latestView.show_dua_after_adhan || useDuaAfterAdhanFixture,
            duaAfterAdhanWindowMinutes,
            useDuaAfterAdhanFixture,
        ));
        if (duaVisibleNow !== duaAfterAdhanVisible) {
            renderSlides(boards, latestView.economic_indicators, latestView.daily_islamic_content, latestView.show_dua_after_adhan);
        }
    }

    function render(view) {
        latestView = view;
        const boards = view && Array.isArray(view.boards) ? view.boards.slice(0, 3) : [];
        if (!view || !view.configured || boards.length === 0) return;
        primaryName.textContent = boards[0].name;
        renderSlides(boards, view.economic_indicators, view.daily_islamic_content, view.show_dua_after_adhan);
        updateHeader();
        state.classList.remove("hidden");
    }

    function refresh(view) {
        render(view);
        const duration = Number(view.slide_duration_seconds);
        slideDurationSeconds = duration >= 5 && duration <= 60 ? duration : 15;
        startTimer();
    }

    state.addEventListener("pointerdown", (event) => { gestureStart = {x: event.clientX, y: event.clientY}; });
    state.addEventListener("pointerup", (event) => {
        if (!gestureStart) return;
        const dx = event.clientX - gestureStart.x;
        const dy = event.clientY - gestureStart.y;
        gestureStart = null;
        if (Math.abs(dx) > 55 && Math.abs(dx) > Math.abs(dy)) showSlide(activeSlide + (dx < 0 ? 1 : -1));
    });

    window.addEventListener("masjidpi:board-view", event => refresh(event.detail));
    if (window.MasjidBoardCurrentView) refresh(window.MasjidBoardCurrentView);
    window.addEventListener("masjidpi:appliance-listen-panel", event => {
        state.classList.toggle("listen-panel-open", Boolean(event.detail && event.detail.open));
        if (state.classList.contains("listen-panel-open")) window.clearInterval(slideTimer);
        else startTimer();
    });
    updateHeader();
    window.setInterval(updateHeader, 1000);
})();
