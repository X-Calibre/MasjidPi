(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") !== "detailed") return;

    document.body.classList.add("detailed-layout");
    const panel = document.getElementById("additionalTimes");
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

    function add(items, label, time) {
        if (validTime(time)) items.push({label, time});
    }

    function buildItems(board) {
        const astronomical = board && board.astronomical ? board.astronomical : {};
        const items = [];

        add(items, "Sehri Ends", astronomical.suhur);
        add(items, "Fajr Starts", astronomical.fajr_start);
        add(items, "Sunrise", astronomical.sunrise);
        add(items, "Ishraq", astronomical.ishraaq);

        // MasjidBoard Live can publish a caution/start/end sequence around
        // Istiwa and Zawaal. Show each supplied value in chronological order.
        add(items, "Istiwa Caution", astronomical.istiwa_caution);
        add(items, "Istiwa", astronomical.istiwa);
        add(items, "Zawaal Ends", astronomical.zawaal_end);

        // Prayer-start values are astronomical starts, not the selected
        // masjid's Adhan/Jamaah values. Fajr and Esha are explicit upstream;
        // Asr exposes both common calculation methods.
        add(items, "Asr Starts (Shafi'i)", astronomical.asr_shafii);
        add(items, "Asr Starts (Hanafi)", astronomical.asr_hanafi);
        add(items, "Maghrib Starts", astronomical.sunset);
        add(items, "Esha Starts", astronomical.esha_start);

        items.sort((left, right) => minutes(left.time) - minutes(right.time));
        return items;
    }

    function render(board) {
        panel.replaceChildren();

        const heading = document.createElement("div");
        heading.className = "additional-times-heading";
        heading.textContent = "Additional Times";
        panel.append(heading);

        const source = document.createElement("div");
        source.className = "additional-times-source";
        source.textContent = board ? `From ${board.name}` : "From first masjid";
        panel.append(source);

        const items = buildItems(board);
        if (items.length === 0) {
            const empty = document.createElement("div");
            empty.className = "additional-times-empty";
            empty.textContent = "No additional times available";
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
            value.textContent = format(item.time);
            row.append(label, value);
            list.append(row);
        }
        panel.append(list);
    }

    async function refresh() {
        try {
            const response = await fetch("/api/masjidboard/display", {cache: "no-store"});
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const view = await response.json();
            const boards = Array.isArray(view.boards) ? view.boards : [];
            render(boards[0] || null);
        } catch (error) {
            console.warn("Detailed MasjidBoard times refresh failed", error);
        }
    }

    refresh();
    window.setInterval(refresh, refreshIntervalMs);
})();
