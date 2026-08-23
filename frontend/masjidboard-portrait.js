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

    function plainText(value) {
        const source = String(value || "").replace(/<br\\s*\\/?\\s*>/gi, "\n");
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
            time: "Time", venue: "Venue", account_name: "Account Name", account_number: "Account Number",
            bank: "Bank", branch_code: "Branch Code", bsb: "BSB", masjid_taleem: "Masjid Taleem",
            gasht_out_day: "Gasht Out", gasht_out_time: "Out Time", gasht_in_day: "Gasht In",
            gasht_in_time: "In Time", first_location: "First Jamaat", first_date: "First Date",
            second_location: "Second Jamaat", second_date: "Second Date",
        };
        return labels[name] || name.replace(/_/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
    }

    function orderedFields(item) {
        const fields = item && item.fields && typeof item.fields === "object" ? item.fields : {};
        const preferred = {
            funeral: ["relation", "salaah_time", "salaah_venue", "pickup", "cemetery", "address"],
            nikah: ["groom_relation", "relation_one", "relation_two", "date", "time"],
            eid: ["date", "venue", "address", "lecture", "salaah"],
            salaah_change: ["effective_date", "new_time"],
            new_moon: ["birth_date", "birth", "best_visibility", "visibility_date", "first_moonset", "first_age"],
            dawah: ["masjid_taleem", "gasht_out_day", "gasht_out_time", "gasht_in_day", "gasht_in_time"],
            three_day_jamaat: ["first_location", "first_date", "second_location", "second_date"],
            contribution: ["bank", "account_name", "branch_code", "account_number", "bsb"],
        }[item.type] || [];
        const titleFields = new Set(item.type === "funeral" ? ["name"]
            : item.type === "nikah" ? ["name_one", "name_two", "bride"]
                : item.type === "salaah_change" ? ["prayer"] : []);
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
        return plainText(notice.title) || fieldLabel(notice.type || "general") + " Notice";
    }

    function communityTypeLabel(type) {
        const labels = {
            announcement: "Announcement", eid: "Eid Notice", funeral: "Funeral Notice",
            nikah: "Nikah Notice", well_wishes: "Well Wishes", salaah_change: "Salaah Time Change",
            programme: "Programme", new_moon: "New Moon", dawah: "Dawah / Gasht",
            three_day_jamaat: "Three-Day Jamaat", contribution: "Contributions",
        };
        return labels[type] || fieldLabel(type || "general") + " Notice";
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
            for (const programme of Array.isArray(board.programmes) ? board.programmes : []) {
                add({
                    type: "programme",
                    title: plainText(programme.title) || "Masjid Programme",
                    body: plainText(programme.content),
                    fields: {},
                    source: board.name,
                });
            }
            if (board.new_moon && board.new_moon.fields && typeof board.new_moon.fields === "object") {
                add({type: "new_moon", title: "New Moon Information", body: "", fields: board.new_moon.fields, source: board.name});
            }
            if (board.banking && board.banking.fields && typeof board.banking.fields === "object") {
                add({
                    type: "contribution",
                    title: plainText(board.banking.title) || "Masjid Contributions",
                    body: "",
                    fields: board.banking.fields,
                    source: board.name,
                });
            }
        }
        return items;
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
        card.append(element("div", "portrait-community-type", communityTypeLabel(item.type)));
        card.append(element("h2", "portrait-community-title", item.title));
        if (item.body) {
            const body = element("p", "portrait-community-body", item.body);
            body.dir = "auto";
            card.append(body);
        }

        const fields = orderedFields(item);
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
        slides.push(...communitySlides(collectCommunityItems(boards)));
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
