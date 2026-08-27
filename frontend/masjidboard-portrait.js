(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") !== "portrait") return;
    const communityFixtureMode = params.get("notice-fixtures");
    const useCommunityFixtures = communityFixtureMode === "1" || communityFixtureMode === "new";

    document.body.classList.add("portrait-layout");
    const utils = window.MasjidBoardDisplayUtils;
    const dateUtils = window.MasjidBoardDate;
    const {collectCommunityItems, fixtureCommunityItems, formatNoticeDate, formatRand, formatUpdatedAt, orderedFields, plainText} = window.MasjidBoardCommunityUtils;
    const state = document.getElementById("portraitState");
    const slidesHost = document.getElementById("portraitSlides");
    const dotsHost = document.getElementById("portraitDots");
    const primaryName = document.getElementById("portraitPrimaryName");
    const clock = document.getElementById("portraitClock");
    const gregorianDate = document.getElementById("portraitGregorianDate");
    const islamicDate = document.getElementById("portraitIslamicDate");
    const nextName = document.getElementById("portraitNextName");
    const nextTime = document.getElementById("portraitNextTime");
    const countdown = document.getElementById("portraitCountdown");
    const prayerLabels = {fajr: "Fajr", dhuhr: "Dhuhr", asr: "Asr", maghrib: "Maghrib", esha: "Esha"};
    let latestView = null;
    let slides = [];
    let activeSlide = 0;
    let slideDurationSeconds = 15;
    let slideTimer = 0;
    let transitionTimer = 0;
    let gestureStart = null;

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

        const zawaalStart = validTime(astronomical.istiwa_caution)
            ? astronomical.istiwa_caution
            : astronomical.istiwa;
        const zawaalEnd = validTime(astronomical.zawaal_end)
            ? astronomical.zawaal_end
            : astronomical.istiwa;
        if (validTime(zawaalStart)) {
            const range = validTime(zawaalEnd) && minutes(zawaalEnd) !== minutes(zawaalStart)
                ? utils.formatClock(zawaalStart) + "–" + utils.formatClock(zawaalEnd)
                : utils.formatClock(zawaalStart);
            add("Zawaal / Istiwa", zawaalStart, range);
        }

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
        const slide = element("article", "portrait-slide");
        const heading = element("header", "portrait-slide-heading");
        heading.append(element("small", "", "SALAAH TIMES"), element("h2", "", board.name));
        slide.append(heading);
        const table = element("div", "portrait-times");
        const head = element("div", "portrait-time-row portrait-time-head");
        head.append(element("span", "", "Salaah"), element("span", "", "Adhan"), element("span", "", "Jamaah"));
        table.append(head);
        const now = utils.displayNow();
        const friday = utils.displayDate().getDay() === 5;
        const next = utils.nextEventForBoard(board, now, friday);
        for (const key of ["fajr", "dhuhr", "asr", "maghrib", "esha"]) {
            if (friday && key === "dhuhr") {
                const service = Array.isArray(board.jumuah) ? board.jumuah[0] : null;
                const items = service ? utils.jumuahItems(service) : [];
                const adhan = items.find((item) => item.kind === "adhan");
                // jumuahItems identifies an explicit Salaah/Jamaah as the
                // salaah item, or uses Khutbah when no such time is supplied.
                const jamaah = items.find((item) => item.kind === "salaah");
                const row = element("div", "portrait-time-row");
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
            const row = element("div", "portrait-time-row");
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
        const slide = element("article", "portrait-slide");
        const heading = element("header", "portrait-slide-heading");
        heading.append(element("small", "", "DAILY TIMES"), element("h2", "", board.name));
        slide.append(heading);
        const table = element("div", "portrait-times portrait-daily-times");
        for (const item of dailyItems(board)) {
            const row = element("div", "portrait-time-row portrait-daily-row");
            row.append(element("span", "", item.label), element("strong", "", item.valueText));
            table.append(row);
        }
        slide.append(table);
        return slide;
    }

    function communityCard(item, compact) {
        const card = element("article", "portrait-community-card portrait-community-" + item.type + (compact ? " compact" : ""));
        card.append(element("h2", "portrait-community-title", item.title));
        if (item.body) {
            const body = element("p", "portrait-community-body", item.body);
            body.dir = "auto";
            card.append(body);
        }

        if (item.type === "salaah_change") {
            const main = element("div", "portrait-salaah-change-main");
            main.append(
                element("div", "portrait-salaah-change-effective", "Effective from\n" + formatNoticeDate(item.fields.effective_date)),
                element("div", "portrait-salaah-change-time", plainText(item.fields.new_time))
            );
            card.append(main);
        }

        const fields = item.type === "salaah_change" ? [] : orderedFields(item);
        if (fields.length > 0) {
            const list = element("div", "portrait-community-fields");
            for (const field of fields) {
                const row = element("div", "portrait-community-field");
                row.append(element("span", "portrait-community-field-label", field.label));
                const value = element("span", "portrait-community-field-value", field.value);
                value.dir = "auto";
                row.append(value);
                list.append(row);
            }
            card.append(list);
        }
        card.append(element("div", "portrait-community-source", "From " + item.source));
        return card;
    }

    function communitySlide(entries) {
        const paired = entries.length === 2;
        const compactSingle = entries.length === 1 && isCompactCommunityItem(entries[0]);
        const layout = paired ? "paired" : compactSingle ? "single compact-single" : "single";
        const slide = element("article", "portrait-slide portrait-community-slide " + layout);
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
        const slide = element("article", "portrait-slide portrait-economic-slide");
        const heading = element("header", "portrait-slide-heading");
        heading.append(element("h2", "", "Islamic Economic Indicators"));
        slide.append(heading);
        const effectiveDate = new Date(`${indicators.effective_date}T12:00:00`);
        const dateText = Number.isNaN(effectiveDate.getTime()) ? indicators.effective_date : effectiveDate.toLocaleDateString("en-ZA", {day: "numeric", month: "long", year: "numeric"});
        slide.append(element("div", "portrait-economic-date", `Effective ${dateText}`));
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
        const grid = element("div", "portrait-economic-values");
        for (const [label, value] of values) {
            const row = element("div", "portrait-economic-value");
            const accountingValue = element("strong", "portrait-economic-accounting");
            accountingValue.append(
                element("span", "portrait-economic-currency", "R"),
                element("span", "portrait-economic-amount", value.slice(1))
            );
            row.append(element("span", "", label), accountingValue);
            grid.append(row);
        }
        const footer = element("footer", "portrait-economic-footer");
        const retrievedAt = formatUpdatedAt(indicators.fetched_at);
        if (retrievedAt) footer.append(element("div", "portrait-economic-retrieved", `Retrieved at ${retrievedAt}`));
        footer.append(element("div", "portrait-economic-source", `From ${indicators.source}`));
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
        if (slides.length < 2) return;
        slideTimer = window.setInterval(() => showSlide(activeSlide + 1, false), slideDurationSeconds * 1000);
    }

    function renderSlides(boards, indicators) {
        slidesHost.replaceChildren();
        dotsHost.replaceChildren();
        slides = boards.map(salaahSlide);
        if (boards[0] && dailyItems(boards[0]).length > 0) slides.push(dailySlide(boards[0]));
        const communityItems = useCommunityFixtures
            ? fixtureCommunityItems(communityFixtureMode)
            : collectCommunityItems(boards);
        slides.push(...communitySlides(communityItems));
        const indicatorsSlide = economicSlide(indicators);
        if (indicatorsSlide) slides.push(indicatorsSlide);
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
        const friday = utils.displayDate().getDay() === 5;
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
    }

    function render(view) {
        latestView = view;
        const boards = view && Array.isArray(view.boards) ? view.boards.slice(0, 3) : [];
        if (!view || !view.configured || boards.length === 0) return;
        primaryName.textContent = boards[0].name;
        renderSlides(boards, view.economic_indicators);
        updateHeader();
        state.classList.remove("hidden");
    }

    async function refresh() {
        try {
            const [displayResponse, layoutResponse] = await Promise.all([
                fetch("/api/masjidboard/display", {cache: "no-store"}),
                fetch("/api/masjidboard/layout", {cache: "no-store"}),
            ]);
            if (!displayResponse.ok) throw new Error(`HTTP ${displayResponse.status}`);
            render(await displayResponse.json());
            if (layoutResponse.ok) {
                const preference = await layoutResponse.json();
                const duration = Number(preference.slide_duration_seconds);
                slideDurationSeconds = duration >= 5 && duration <= 60 ? duration : 15;
                startTimer();
            }
        } catch (error) {
            console.warn("Portrait MasjidBoard refresh failed", error);
        }
    }

    state.addEventListener("pointerdown", (event) => { gestureStart = {x: event.clientX, y: event.clientY}; });
    state.addEventListener("pointerup", (event) => {
        if (!gestureStart) return;
        const dx = event.clientX - gestureStart.x;
        const dy = event.clientY - gestureStart.y;
        gestureStart = null;
        if (Math.abs(dx) > 55 && Math.abs(dx) > Math.abs(dy)) showSlide(activeSlide + (dx < 0 ? 1 : -1));
    });

    refresh();
    updateHeader();
    window.setInterval(updateHeader, 1000);
    window.setInterval(refresh, 60000);
})();
