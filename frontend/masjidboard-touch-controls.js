(() => {
    "use strict";

    const params = new URLSearchParams(window.location.search);
    if (params.get("layout") !== "portrait") return;

    const state = document.getElementById("portraitState");
    const panel = document.getElementById("portraitListenPanel");
    if (!state || !panel) return;

    const connection = document.getElementById("portraitListenConnection");
    const statusBadge = document.getElementById("portraitListenState");
    const nowPlaying = document.getElementById("portraitListenNowPlaying");
    const detail = document.getElementById("portraitListenDetail");
    const favouriteHost = document.getElementById("portraitFavouriteMasjids");
    const radioHost = document.getElementById("portraitRadioStations");
    const masjidSelection = document.getElementById("portraitMasjidSelection");
    const radioSelection = document.getElementById("portraitRadioSelection");
    const masterVolume = document.getElementById("portraitMasterVolume");
    const masjidVolume = document.getElementById("portraitMasjidVolume");
    const radioVolume = document.getElementById("portraitRadioVolume");
    const volumeControls = {master: masterVolume, masjid: masjidVolume, radio: radioVolume};
    const volumeOutputs = {
        master: document.getElementById("portraitMasterVolumeValue"),
        masjid: document.getElementById("portraitMasjidVolumeValue"),
        radio: document.getElementById("portraitRadioVolumeValue")
    };
    const playMasjid = document.getElementById("portraitPlayMasjid");
    const stopListening = document.getElementById("portraitStopListening");
    const radioModeButtons = {
        schedule: document.getElementById("portraitRadioSchedule"),
        play_now: document.getElementById("portraitRadioPlayNow"),
        stopped: document.getElementById("portraitRadioStop")
    };
    const radioModeDetail = document.getElementById("portraitRadioModeDetail");
    const themeHost = document.getElementById("portraitThemeChoices");
    const themes = [
        ["emerald", "Emerald", "MasjidPi green"],
        ["midnight", "Midnight", "Deep blue"],
        ["slate", "Slate", "Neutral gold"],
        ["ruby", "Ruby", "Warm red"],
        ["light", "Light", "Bright display"],
        ["black-white", "Black & White", "Maximum contrast"]
    ];

    let open = false;
    let status = null;
    let favouriteMasjids = [];
    let radios = [];
    let selectedMasjidID = "";
    let selectedRadioID = "";
    let refreshTimer = 0;
    let gestureStart = null;
    let busy = false;
    let currentTheme = document.body.dataset.boardTheme || "emerald";
    const volumeSaveTimers = {master:0, masjid:0, radio:0};
    const volumeSaveSerials = {master:0, masjid:0, radio:0};
    const pendingVolumes = {master:null, masjid:null, radio:null};

    function jsonOptions(method, body) {
        return {method, headers:{"Content-Type":"application/json"}, body:JSON.stringify(body)};
    }

    async function requestJSON(url, options = {}) {
        const response = await fetch(url, options);
        if (!response.ok) {
            let message = `Request failed (${response.status})`;
            try {
                const body = await response.json();
                if (body.error) message = body.error;
            } catch (_) {}
            throw new Error(message);
        }
        return response.json();
    }

    function label(stream) {
        return stream && stream.location ? `${stream.name} — ${stream.location}` : (stream && stream.name) || "Unknown source";
    }

    function formatResumeCountdown(resumeAt) {
        if (!resumeAt) return "";
        const seconds = Math.max(0, Math.ceil((new Date(resumeAt).getTime() - Date.now()) / 1000));
        const minutes = Math.floor(seconds / 60);
        return `${minutes}:${String(seconds % 60).padStart(2, "0")}`;
    }

    function setConnectionError(message = "") {
        connection.textContent = message;
        connection.classList.toggle("hidden", !message);
    }

    function setBusy(value) {
        busy = value;
        renderStatus();
    }

    function setOpen(value) {
        open = value;
        panel.classList.toggle("hidden", !open);
        panel.setAttribute("aria-hidden", open ? "false" : "true");
        window.dispatchEvent(new CustomEvent("masjidpi:portrait-listen-panel", {detail:{open}}));
        window.clearTimeout(refreshTimer);
        if (open) {
            loadPanel();
            panel.querySelector(".portrait-listen-close")?.focus();
        }
    }

    function activateTab(name) {
        panel.querySelectorAll("[data-touch-tab]").forEach(button => {
            const active = button.dataset.touchTab === name;
            button.classList.toggle("active", active);
            button.setAttribute("aria-selected", active ? "true" : "false");
        });
        panel.querySelectorAll("[data-touch-panel]").forEach(section => section.classList.toggle("hidden", section.dataset.touchPanel !== name));
    }

    function renderSources() {
        favouriteHost.replaceChildren();
        if (favouriteMasjids.length === 0) {
            const empty = document.createElement("div");
            empty.className = "portrait-source-empty";
            empty.textContent = "No favourite masjids. Add favourites through the full MasjidPi Web UI.";
            favouriteHost.append(empty);
        }
        for (const item of favouriteMasjids) {
            const button = document.createElement("button");
            button.type = "button";
            button.dataset.sourceId = item.id;
            button.textContent = item.name || "Unknown Masjid";
            button.classList.toggle("selected", item.id === selectedMasjidID);
            button.classList.toggle("playing", status?.active_source === "masjid" && item.id === status.active_stream_id);
            button.setAttribute("role", "option");
            button.setAttribute("aria-selected", item.id === selectedMasjidID ? "true" : "false");
            favouriteHost.append(button);
        }

        radioHost.replaceChildren();
        for (const item of radios) {
            const button = document.createElement("button");
            button.type = "button";
            button.dataset.sourceId = item.id;
            button.textContent = label(item);
            button.classList.toggle("selected", item.id === selectedRadioID);
            button.classList.toggle("playing", status?.active_source === "radio" && item.id === status.active_stream_id);
            button.setAttribute("role", "option");
            button.setAttribute("aria-selected", item.id === selectedRadioID ? "true" : "false");
            radioHost.append(button);
        }
        if (radios.length === 0) {
            const empty = document.createElement("div");
            empty.className = "portrait-source-empty";
            empty.textContent = "No Radio stations are available.";
            radioHost.append(empty);
        }
        masjidSelection.textContent = favouriteMasjids.find(item => item.id === selectedMasjidID)?.name || "";
        radioSelection.textContent = radios.find(item => item.id === selectedRadioID)?.name || "";
    }

    function renderThemes() {
        themeHost.replaceChildren();
        for (const [value, name, description] of themes) {
            const button = document.createElement("button");
            button.type = "button";
            button.className = "portrait-theme-option";
            button.dataset.theme = value;
            button.classList.toggle("active", value === currentTheme);
            button.setAttribute("role", "radio");
            button.setAttribute("aria-checked", value === currentTheme ? "true" : "false");
            button.disabled = busy;
            const swatch = document.createElement("span");
            swatch.className = "portrait-theme-swatch";
            swatch.setAttribute("aria-hidden", "true");
            const text = document.createElement("span");
            const strong = document.createElement("strong");
            strong.textContent = name;
            const small = document.createElement("span");
            small.className = "portrait-theme-description";
            small.textContent = description;
            text.append(strong, small);
            button.append(swatch, text);
            themeHost.append(button);
        }
    }

    function renderStatus() {
        if (!status) {
            statusBadge.textContent = "Offline";
            statusBadge.className = "portrait-listen-badge error";
            nowPlaying.textContent = "Listen is unavailable";
            detail.textContent = "Check that the Listen component is installed and running.";
        } else if (!status.listening) {
            statusBadge.textContent = "Stopped";
            statusBadge.className = "portrait-listen-badge stopped";
            nowPlaying.textContent = "Listening is stopped";
            detail.textContent = "Choose a source and start playback.";
        } else if (status.active_source === "masjid") {
            statusBadge.textContent = "Masjid";
            statusBadge.className = "portrait-listen-badge";
            nowPlaying.textContent = status.active_stream_name || status.masjid_name || "Selected Masjid";
            detail.textContent = status.radio_name
                ? `${status.radio_name} is standing by.`
                : "Radio will resume when the Masjid broadcast ends.";
        } else if (status.active_source === "radio") {
            statusBadge.textContent = "Radio";
            statusBadge.className = "portrait-listen-badge";
            nowPlaying.textContent = status.active_stream_name || status.radio_name || "Selected Radio station";
            detail.textContent = "Radio will yield automatically when the Masjid comes online.";
        } else {
            statusBadge.textContent = "Waiting";
            statusBadge.className = "portrait-listen-badge waiting";
            nowPlaying.textContent = status.radio_resume_pending
                ? `${status.radio_name || "Radio"} resumes in ${formatResumeCountdown(status.radio_resume_at)}`
                : "No source currently playing";
            detail.textContent = status.radio_resume_pending ? "Post-Masjid Radio delay is active." : "Waiting for an available source.";
        }

        if (status) {
            const values = {master:status.master_volume, masjid:status.masjid_volume, radio:status.radio_volume};
            for (const [name, control] of Object.entries(volumeControls)) {
                const displayed = pendingVolumes[name] ?? values[name];
                if (document.activeElement !== control) control.value = displayed;
                volumeOutputs[name].textContent = `${control.value}%`;
            }
            masterVolume.disabled = busy || !status.master_volume_supported;
            volumeOutputs.master.textContent = status.master_volume_supported ? `${masterVolume.value}%` : "Unavailable";
            masjidVolume.disabled = busy;
            radioVolume.disabled = busy || !status.radio_enabled;
            selectedMasjidID ||= status.masjid_id || "";
            selectedRadioID ||= status.radio_id || "";
        } else {
            for (const control of Object.values(volumeControls)) control.disabled = true;
        }

        const selectedMasjidPlaying = status?.active_source === "masjid" && status.active_stream_id === selectedMasjidID;
        playMasjid.disabled = busy || !selectedMasjidID || selectedMasjidPlaying;
        playMasjid.textContent = selectedMasjidPlaying ? "Masjid Playing" : "▶ Play Masjid";
        stopListening.disabled = busy || !status?.listening;
        for (const [mode, button] of Object.entries(radioModeButtons)) {
            button.disabled = busy || !selectedRadioID || !status || (mode === "stopped" && !status.radio_enabled);
            button.classList.toggle("active", status?.radio_mode === mode);
        }
        for (const stepButton of panel.querySelectorAll("[data-volume-step]")) {
            const name = stepButton.dataset.volumeStep.split(":")[0];
            stepButton.disabled = volumeControls[name].disabled;
        }
        const modeText = status?.radio_mode === "play_now" ? "Play Now override active"
            : status?.radio_mode === "stopped" ? "Radio remains stopped until another mode is selected"
            : status?.radio_schedule_enabled ? `Scheduled ${status.radio_schedule_start}–${status.radio_schedule_stop}`
            : "Scheduled mode · Radio may play whenever the Masjid is offline";
        radioModeDetail.textContent = modeText;
        renderSources();
        renderThemes();
    }

    async function refreshStatus() {
        if (!open) return;
        try {
            status = await requestJSON("/api/listen/status");
            setConnectionError();
            renderStatus();
        } catch (error) {
            status = null;
            setConnectionError(error.message);
            renderStatus();
        } finally {
            window.clearTimeout(refreshTimer);
            if (open) refreshTimer = window.setTimeout(refreshStatus, 1000);
        }
    }

    async function loadPanel() {
        setConnectionError();
        try {
            const results = await Promise.allSettled([
                requestJSON("/api/listen/status"),
                requestJSON("/api/streams?kind=masjid"),
                requestJSON("/api/streams?kind=radio"),
                requestJSON("/api/favourites"),
                requestJSON("/api/masjidboard/layout")
            ]);
            const boardLayout = results[4].status === "fulfilled" ? results[4].value : null;
            if (boardLayout) currentTheme = boardLayout.theme || "emerald";
            if (results.slice(0, 4).every(result => result.status === "fulfilled")) {
                const [newStatus, masjids, radioItems, favourites] = results.map(result => result.value);
                const favouriteIDs = new Set(favourites.ids || []);
                status = newStatus;
                favouriteMasjids = masjids.filter(item => favouriteIDs.has(item.id));
                radios = radioItems;
                selectedMasjidID = status.masjid_id || favouriteMasjids[0]?.id || "";
                selectedRadioID = status.radio_id || radios[0]?.id || "";
                setConnectionError();
            } else {
                status = null;
                favouriteMasjids = [];
                radios = [];
                const failure = results.slice(0, 4).find(result => result.status === "rejected");
                setConnectionError(failure?.reason?.message || "Listen controls are unavailable.");
            }
            renderStatus();
        } catch (error) {
            setConnectionError(error.message);
        }
        window.clearTimeout(refreshTimer);
        if (open) refreshTimer = window.setTimeout(refreshStatus, 1000);
    }

    async function runAction(action, refreshListen = true) {
        if (busy) return;
        setBusy(true);
        setConnectionError();
        try {
            await action();
            if (refreshListen) status = await requestJSON("/api/listen/status");
        } catch (error) {
            setConnectionError(error.message);
        } finally {
            setBusy(false);
        }
    }

    async function ensureSelection(source, id) {
        if (!id) throw new Error(`Select a ${source === "masjid" ? "favourite Masjid" : "Radio station"} first.`);
        await requestJSON("/api/listen/selection", jsonOptions("PUT", {[`${source}_id`]:id}));
    }

    async function selectSource(source, id) {
        if (busy || !id) return;
        const previous = source === "masjid" ? selectedMasjidID : selectedRadioID;
        if (source === "masjid") selectedMasjidID = id;
        else selectedRadioID = id;
        renderSources();
        setBusy(true);
        setConnectionError();
        try {
            await ensureSelection(source, id);
            status = await requestJSON("/api/listen/status");
        } catch (error) {
            if (source === "masjid") selectedMasjidID = previous;
            else selectedRadioID = previous;
            setConnectionError(error.message);
        } finally {
            setBusy(false);
        }
    }

    favouriteHost.addEventListener("click", event => {
        const button = event.target.closest("button[data-source-id]");
        if (!button) return;
        selectSource("masjid", button.dataset.sourceId);
    });
    radioHost.addEventListener("click", event => {
        const button = event.target.closest("button[data-source-id]");
        if (!button) return;
        selectSource("radio", button.dataset.sourceId);
    });
    themeHost.addEventListener("click", event => {
        const button = event.target.closest("button[data-theme]");
        if (!button) return;
        const theme = button.dataset.theme;
        runAction(async () => {
            const saved = await requestJSON("/api/masjidboard/layout", jsonOptions("PUT", {theme}));
            currentTheme = saved.theme || theme;
            document.body.dataset.boardTheme = currentTheme;
            renderThemes();
        }, false);
    });

    for (const button of panel.querySelectorAll("[data-touch-tab]")) button.addEventListener("click", () => activateTab(button.dataset.touchTab));
    for (const button of panel.querySelectorAll("[data-listen-close]")) button.addEventListener("click", () => setOpen(false));

    async function saveVolume(name, value) {
        window.clearTimeout(volumeSaveTimers[name]);
        const serial = ++volumeSaveSerials[name];
        pendingVolumes[name] = value;
        volumeOutputs[name].textContent = `${value}%`;
        try {
            if (name === "master") {
                await requestJSON("/api/player/volume", jsonOptions("POST", {volume:value, persist:true}));
            } else {
                await requestJSON("/api/listen/volume", jsonOptions("PUT", {source:name, volume:value}));
            }
            const refreshed = await requestJSON("/api/listen/status");
            if (serial === volumeSaveSerials[name]) {
                status = refreshed;
                setConnectionError();
            }
        } catch (error) {
            if (serial === volumeSaveSerials[name]) setConnectionError(error.message);
        } finally {
            if (serial === volumeSaveSerials[name]) {
                pendingVolumes[name] = null;
                renderStatus();
            }
        }
    }

    function scheduleVolumeSave(name) {
        const value = Number(volumeControls[name].value);
        pendingVolumes[name] = value;
        volumeOutputs[name].textContent = `${value}%`;
        window.clearTimeout(volumeSaveTimers[name]);
        volumeSaveTimers[name] = window.setTimeout(() => saveVolume(name, value), 120);
    }

    for (const [name, control] of Object.entries(volumeControls)) {
        control.addEventListener("input", () => scheduleVolumeSave(name));
        control.addEventListener("change", () => saveVolume(name, Number(control.value)));
    }
    panel.addEventListener("click", event => {
        const stepButton = event.target.closest("[data-volume-step]");
        if (!stepButton) return;
        const [name, amount] = stepButton.dataset.volumeStep.split(":");
        const control = volumeControls[name];
        control.value = Math.max(Number(control.min), Math.min(Number(control.max), Number(control.value) + Number(amount)));
        saveVolume(name, Number(control.value));
    });

    playMasjid.addEventListener("click", () => runAction(async () => {
        await ensureSelection("masjid", selectedMasjidID);
        if (!status?.masjid_enabled) await requestJSON("/api/listen/power", jsonOptions("PUT", {module:"masjid", enabled:true}));
        await requestJSON("/api/listen/start", {method:"POST"});
    }));
    stopListening.addEventListener("click", () => runAction(() => requestJSON("/api/listen/stop", {method:"POST"})));

    for (const [mode, button] of Object.entries(radioModeButtons)) {
        button.addEventListener("click", () => runAction(async () => {
            if (mode !== "stopped") await ensureSelection("radio", selectedRadioID);
            if (!status?.radio_enabled && mode !== "stopped") await requestJSON("/api/listen/power", jsonOptions("PUT", {module:"radio", enabled:true}));
            if (!status?.listening && mode !== "stopped") await requestJSON("/api/listen/start", {method:"POST"});
            await requestJSON("/api/listen/radio-mode", jsonOptions("PUT", {mode}));
        }));
    }

    state.addEventListener("pointerdown", event => {
        if (open) return;
        gestureStart = {x:event.clientX, y:event.clientY};
    });
    state.addEventListener("pointerup", event => {
        if (!gestureStart || open) return;
        const dx = event.clientX - gestureStart.x;
        const dy = event.clientY - gestureStart.y;
        gestureStart = null;
        if (dy < -70 && Math.abs(dy) > Math.abs(dx)) setOpen(true);
    });
    panel.addEventListener("pointerdown", event => event.stopPropagation());
    panel.addEventListener("pointerup", event => event.stopPropagation());
    for (const target of panel.querySelectorAll(".portrait-listen-handle,.portrait-listen-heading")) {
        target.addEventListener("pointerdown", event => {
            gestureStart = {x:event.clientX, y:event.clientY};
            event.stopPropagation();
        });
        target.addEventListener("pointerup", event => {
            if (!gestureStart) return;
            const dx = event.clientX - gestureStart.x;
            const dy = event.clientY - gestureStart.y;
            gestureStart = null;
            if (dy > 70 && Math.abs(dy) > Math.abs(dx)) setOpen(false);
            event.stopPropagation();
        });
    }
    document.addEventListener("keydown", event => { if (open && event.key === "Escape") setOpen(false); });
})();
