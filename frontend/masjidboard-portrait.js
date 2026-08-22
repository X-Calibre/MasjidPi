(() => {
    "use strict";

    if (new URLSearchParams(window.location.search).get("layout") !== "portrait") return;

    document.body.classList.add("portrait-layout");
    const utils = window.MasjidBoardDisplayUtils;
    const dateUtils = window.MasjidBoardDate;
    const state = document.getElementById("portraitState");
    const slidesHost = document.getElementById("portraitSlides");
    const dotsHost = document.getElementById("portraitDots");
    const primaryName = document.getElementById("portraitPrimaryName");
    const boardStatus = document.getElementById("portraitBoardStatus");
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
        if (event.kind === "jumuah") return {name: `Jumuâah Â· ${event.label || "Salaah"}`, time: event.time};
        const prayer = prayerLabels[event.key] || event.key;
        return {name: `${prayer} Â· ${event.event === "jamaah" ? "Jamaah" : "Adhan"}`, time: event.time};
    }

    function dailyItems(board) {
        const astronomical = board && board.astronomical ? board.astronomical : {};
        const result = [];
        const add = (label, time) => { if (validTime(time)) result.push({label, time}); };
        add("Sehri Ends", astronomical.suhur);
        add("Fajr Starts", astronomical.fajr_start);
        add("Sunrise", astronomical.sunrise);
        add("Ishraq", astronomical.ishraaq);
        add("Zawaal / Istiwa", astronomical.istiwa_caution || astronomical.istiwa);
        add("Sunset", astronomical.sunset);
        add("Esha Starts", astronomical.esha_start);
        return result.sort((left, right) => minutes(left.time) - minutes(right.time));
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
            const prayer = utils.findPrayer(board, key);
            const row = element("div", "portrait-time-row");
            if (next && next.kind === "prayer" && next.key === key) row.classList.add("upcoming");
            row.append(
                element("span", "", prayerLabels[key]),
                element("strong", "", prayer && prayer.adhan ? utils.formatClock(prayer.adhan) : "â"),
                element("strong", "", prayer && prayer.jamaah ? utils.formatClock(prayer.jamaah) : "â"),
            );
            table.append(row);
        }
        slide.append(table);
        return slide;
    }

    function dailySlide(board) {
        const slide = element("article", "portrait-slide");
        const heading = element("header", "portrait-slide-heading");
        heading.append(element("small", "", board.name), element("h2", "", "Other daily times"));
        slide.append(heading);
        const grid = element("div", "portrait-daily-grid");
        for (const item of dailyItems(board)) {
            const card = element("div", "portrait-daily-card");
            card.append(element("span", "", item.label), element("strong", "", utils.formatClock(item.time)));
            grid.append(card);
        }
        slide.append(grid);
        return slide;
    }

    function showSlide(index, restart = true) {
        if (slides.length === 0) return;
        activeSlide = (index + slides.length) % slides.length;
        slides.forEach((slide, itemIndex) => slide.classList.toggle("active", itemIndex === activeSlide));
        Array.from(dotsHost.children).forEach((dot, itemIndex) => dot.classList.toggle("active", itemIndex === activeSlide));
        if (restart) startTimer();
    }

    function startTimer() {
        window.clearInterval(slideTimer);
        if (slides.length < 2) return;
        slideTimer = window.setInterval(() => showSlide(activeSlide + 1, false), slideDurationSeconds * 1000);
    }

    function renderSlides(boards) {
        slidesHost.replaceChildren();
        dotsHost.replaceChildren();
        slides = boards.map(salaahSlide);
        if (boards[0] && dailyItems(boards[0]).length > 0) slides.push(dailySlide(boards[0]));
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
        boardStatus.textContent = boards[0].status === "stale" ? "STALE" : boards[0].status === "unavailable" ? "OFFLINE" : "ONLINE";
        boardStatus.dataset.status = boards[0].status || "online";
        renderSlides(boards);
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
