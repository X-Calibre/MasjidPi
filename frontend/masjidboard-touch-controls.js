(() => {
    "use strict";

    const profile = new URLSearchParams(window.location.search).get("profile");
    if (profile !== "appliance-720") return;

    const state = document.getElementById("applianceState");
    const bottomPanel = document.getElementById("applianceListenPanel");
    const quickPanel = document.getElementById("applianceQuickPanel");
    if (!state || !bottomPanel || !quickPanel) return;

    const byID = id => document.getElementById(id);
    const changeWiFi = byID("applianceChangeWiFi");
    if (changeWiFi) changeWiFi.href = `/setup.html?return=board&profile=${profile}`;

    const connections = [byID("applianceListenConnection"), byID("applianceQuickConnection")].filter(Boolean);
    const statusBadge = byID("applianceListenState");
    const nowPlaying = byID("applianceListenNowPlaying");
    const detail = byID("applianceListenDetail");
    const favouriteHost = byID("applianceFavouriteMasjids");
    const radioHost = byID("applianceRadioStations");
    const masjidSelection = byID("applianceMasjidSelection");
    const radioSelection = byID("applianceRadioSelection");
    const volumeControls = {
        master: byID("applianceMasterVolume"),
        masjid: byID("applianceMasjidVolume"),
        radio: byID("applianceRadioVolume")
    };
    const volumeOutputs = {
        master: byID("applianceMasterVolumeValue"),
        masjid: byID("applianceMasjidVolumeValue"),
        radio: byID("applianceRadioVolumeValue")
    };
    const playMasjid = byID("appliancePlayMasjid");
    const stopListening = byID("applianceStopListening");
    const radioModeButtons = [...state.querySelectorAll("[data-radio-mode]")];
    const radioModeDetail = byID("applianceRadioModeDetail");
    const themeHost = byID("applianceThemeChoices");
    const networkFQDNRow = byID("applianceNetworkFQDNRow");
    const networkFQDN = byID("applianceNetworkFQDN");
    const networkIPRow = byID("applianceNetworkIPRow");
    const networkIP = byID("applianceNetworkIP");
    const networkUnavailable = byID("applianceNetworkUnavailable");
    const brightness = byID("applianceBrightness");
    const brightnessValue = byID("applianceBrightnessValue");
    const brightnessUnavailable = byID("applianceBrightnessUnavailable");
    const themes = [
        ["emerald","Emerald","MasjidPi green"],["midnight","Midnight","Deep blue"],
        ["slate","Slate","Neutral gold"],["ruby","Ruby","Warm red"],
        ["light","Light Gold","Warm gold"],["ivory","Ivory","Ivory & emerald"],
        ["sage","Sage","Soft forest green"],["sky","Sky","Cool blue"],
        ["rose","Rose","Blush & burgundy"],["black-white","Black & White","Maximum contrast"]
    ];

    let openPanel = "";
    let status = null;
    let favouriteMasjids = [];
    let radios = [];
    let selectedMasjidID = "";
    let selectedRadioID = "";
    let currentTheme = document.body.dataset.boardTheme || "emerald";
    let displaySettings = null;
    let refreshTimer = 0;
    let inactivityTimer = 0;
    let gestureStart = null;
    let closeGestureStart = null;
    let busy = false;
    let brightnessSaveTimer = 0;
    const volumeSaveTimers = {master:0,masjid:0,radio:0};
    const volumeSaveSerials = {master:0,masjid:0,radio:0};
    const pendingVolumes = {master:null,masjid:null,radio:null};
    const inactivityTimeout = 60000;

    const jsonOptions = (method, body) => ({method,headers:{"Content-Type":"application/json"},body:JSON.stringify(body)});

    async function requestJSON(url, options = {}) {
        const response = await fetch(url, options);
        if (!response.ok) {
            let message = `Request failed (${response.status})`;
            try { message = (await response.json()).error || message; } catch (_) {}
            throw new Error(message);
        }
        return response.json();
    }

    const label = stream => stream?.location ? `${stream.name} — ${stream.location}` : stream?.name || "Unknown source";
    function formatResumeCountdown(resumeAt) {
        if (!resumeAt) return "";
        const seconds = Math.max(0, Math.ceil((new Date(resumeAt).getTime() - Date.now()) / 1000));
        return `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2,"0")}`;
    }

    function setConnectionError(message = "") {
        for (const connection of connections) {
            connection.textContent = message;
            connection.classList.toggle("hidden", !message);
        }
    }

    function renderNetworkAccess(access = null) {
        const port = window.location.port || "8080";
        const fqdn = access?.fqdn || "";
        const ipAddress = access?.ip_address || "";
        networkFQDN.textContent = fqdn ? `http://${fqdn}:${port}` : "";
        networkIP.textContent = ipAddress ? `http://${ipAddress}:${port}` : "";
        networkFQDNRow.classList.toggle("hidden", !fqdn);
        networkIPRow.classList.toggle("hidden", !ipAddress);
        networkUnavailable.classList.toggle("hidden", Boolean(fqdn || ipAddress));
    }

    function renderDisplaySettings() {
        const available = Boolean(displaySettings?.brightness_available);
        brightness.disabled = busy || !available;
        brightness.value = displaySettings?.brightness_percent ?? 100;
        brightnessValue.textContent = available ? `${brightness.value}%` : "Unavailable";
        brightnessUnavailable.classList.toggle("hidden", available);
    }

    function resetInactivityTimer() {
        window.clearTimeout(inactivityTimer);
        if (openPanel) inactivityTimer = window.setTimeout(() => setOpenPanel(""), inactivityTimeout);
    }

    function setOpenPanel(name) {
        openPanel = name;
        bottomPanel.classList.toggle("hidden", name !== "bottom");
        quickPanel.classList.toggle("hidden", name !== "quick");
        bottomPanel.setAttribute("aria-hidden", name === "bottom" ? "false" : "true");
        quickPanel.setAttribute("aria-hidden", name === "quick" ? "false" : "true");
        window.dispatchEvent(new CustomEvent("masjidpi:appliance-listen-panel", {detail:{open:Boolean(name)}}));
        window.clearTimeout(refreshTimer);
        window.clearTimeout(inactivityTimer);
        if (name) {
            resetInactivityTimer();
            loadPanel();
        }
    }

    function activateTab(name) {
        bottomPanel.querySelectorAll("[data-touch-tab]").forEach(button => {
            const active = button.dataset.touchTab === name;
            button.classList.toggle("active", active);
            button.setAttribute("aria-selected", active ? "true" : "false");
        });
        bottomPanel.querySelectorAll("[data-touch-panel]").forEach(section =>
            section.classList.toggle("hidden", section.dataset.touchPanel !== name));
    }

    function renderSources() {
        favouriteHost.replaceChildren();
        if (!favouriteMasjids.length) {
            const empty = document.createElement("div");
            empty.className = "appliance-source-empty";
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
            button.setAttribute("role","option");
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
            button.setAttribute("role","option");
            button.setAttribute("aria-selected", item.id === selectedRadioID ? "true" : "false");
            radioHost.append(button);
        }
        if (!radios.length) {
            const empty = document.createElement("div");
            empty.className = "appliance-source-empty";
            empty.textContent = "No Radio stations are available.";
            radioHost.append(empty);
        }
        masjidSelection.textContent = favouriteMasjids.find(item => item.id === selectedMasjidID)?.name || "";
        radioSelection.textContent = radios.find(item => item.id === selectedRadioID)?.name || "";
    }

    function renderThemes() {
        themeHost.replaceChildren();
        for (const [value,name,description] of themes) {
            const button = document.createElement("button");
            button.type = "button";
            button.className = "appliance-theme-option";
            button.dataset.theme = value;
            button.classList.toggle("active", value === currentTheme);
            button.setAttribute("role","radio");
            button.setAttribute("aria-checked", value === currentTheme ? "true" : "false");
            button.disabled = busy;
            const swatch = document.createElement("span");
            swatch.className = "appliance-theme-swatch";
            swatch.setAttribute("aria-hidden","true");
            const text = document.createElement("span");
            const strong = document.createElement("strong");
            strong.textContent = name;
            const small = document.createElement("span");
            small.className = "appliance-theme-description";
            small.textContent = description;
            text.append(strong,small);
            button.append(swatch,text);
            themeHost.append(button);
        }
    }

    function renderStatus() {
        if (!status) {
            statusBadge.textContent = "Offline";
            statusBadge.className = "appliance-listen-badge error";
            nowPlaying.textContent = "Listen is unavailable";
            detail.textContent = "Check that the Listen component is installed and running.";
        } else if (!status.listening) {
            statusBadge.textContent = "Stopped";
            statusBadge.className = "appliance-listen-badge stopped";
            nowPlaying.textContent = "Listening is stopped";
            detail.textContent = "Choose a source and start playback.";
        } else if (status.active_source === "masjid") {
            statusBadge.textContent = "Masjid";
            statusBadge.className = "appliance-listen-badge";
            nowPlaying.textContent = status.active_stream_name || status.masjid_name || "Selected Masjid";
            detail.textContent = status.radio_name ? `${status.radio_name} is standing by.` : "Radio will resume when the Masjid broadcast ends.";
        } else if (status.active_source === "radio") {
            statusBadge.textContent = "Radio";
            statusBadge.className = "appliance-listen-badge";
            nowPlaying.textContent = status.active_stream_name || status.radio_name || "Selected Radio station";
            detail.textContent = "Radio will yield automatically when the Masjid comes online.";
        } else {
            statusBadge.textContent = "Waiting";
            statusBadge.className = "appliance-listen-badge waiting";
            nowPlaying.textContent = status.radio_resume_pending
                ? `${status.radio_name || "Radio"} resumes in ${formatResumeCountdown(status.radio_resume_at)}`
                : "No source currently playing";
            detail.textContent = status.radio_resume_pending ? "Post-Masjid Radio delay is active." : "Waiting for an available source.";
        }

        if (status) {
            const values = {master:status.master_volume,masjid:status.masjid_volume,radio:status.radio_volume};
            for (const [name,control] of Object.entries(volumeControls)) {
                const displayed = pendingVolumes[name] ?? values[name];
                if (document.activeElement !== control) control.value = displayed;
                volumeOutputs[name].textContent = `${control.value}%`;
            }
            volumeControls.master.disabled = busy || !status.master_volume_supported;
            volumeOutputs.master.textContent = status.master_volume_supported ? `${volumeControls.master.value}%` : "Unavailable";
            volumeControls.masjid.disabled = busy;
            volumeControls.radio.disabled = busy || !status.radio_enabled;
            selectedMasjidID ||= status.masjid_id || "";
            selectedRadioID ||= status.radio_id || "";
        } else {
            Object.values(volumeControls).forEach(control => { control.disabled = true; });
        }

        const selectedMasjidPlaying = status?.active_source === "masjid" && status.active_stream_id === selectedMasjidID;
        playMasjid.disabled = busy || !selectedMasjidID || selectedMasjidPlaying;
        playMasjid.textContent = selectedMasjidPlaying ? "Masjid Playing" : "▶ Play Masjid";
        stopListening.disabled = busy || !status?.listening;
        for (const button of radioModeButtons) {
            const mode = button.dataset.radioMode;
            button.disabled = busy || !selectedRadioID || !status || (mode === "stopped" && !status.radio_enabled);
            button.classList.toggle("active", status?.radio_mode === mode);
        }
        for (const stepButton of state.querySelectorAll("[data-volume-step]")) {
            stepButton.disabled = volumeControls[stepButton.dataset.volumeStep.split(":")[0]].disabled;
        }
        radioModeDetail.textContent = status?.radio_mode === "play_now" ? "Play Now override active"
            : status?.radio_mode === "stopped" ? "Radio remains stopped until another mode is selected"
            : status?.radio_schedule_enabled ? `Scheduled ${status.radio_schedule_start}–${status.radio_schedule_stop}`
            : "Scheduled mode · Radio may play whenever the Masjid is offline";
        renderSources();
        renderThemes();
        renderDisplaySettings();
    }

    function setBusy(value) { busy = value; renderStatus(); }

    async function refreshStatus() {
        if (!openPanel) return;
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
            if (openPanel) refreshTimer = window.setTimeout(refreshStatus,1000);
        }
    }

    async function loadPanel() {
        setConnectionError();
        const results = await Promise.allSettled([
            requestJSON("/api/listen/status"),requestJSON("/api/streams?kind=masjid"),
            requestJSON("/api/streams?kind=radio"),requestJSON("/api/favourites"),
            requestJSON("/api/masjidboard/layout"),requestJSON("/api/setup/device-access"),
            requestJSON("/api/display/settings")
        ]);
        const boardLayout = results[4].status === "fulfilled" ? results[4].value : null;
        if (boardLayout) currentTheme = boardLayout.theme || "emerald";
        renderNetworkAccess(results[5].status === "fulfilled" ? results[5].value : null);
        if (results[6].status === "fulfilled") displaySettings = results[6].value;
        if (results.slice(0,4).every(result => result.status === "fulfilled")) {
            const [newStatus,masjids,radioItems,favourites] = results.map(result => result.value);
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
            const failure = results.slice(0,4).find(result => result.status === "rejected");
            setConnectionError(failure?.reason?.message || "Listen controls are unavailable.");
        }
        renderStatus();
        window.clearTimeout(refreshTimer);
        if (openPanel) refreshTimer = window.setTimeout(refreshStatus,1000);
    }

    async function runAction(action, refreshListen = true) {
        if (busy) return;
        setBusy(true);
        setConnectionError();
        try {
            await action();
            if (refreshListen) status = await requestJSON("/api/listen/status");
        } catch (error) { setConnectionError(error.message); }
        finally { setBusy(false); resetInactivityTimer(); }
    }

    async function ensureSelection(source,id) {
        if (!id) throw new Error(`Select a ${source === "masjid" ? "favourite Masjid" : "Radio station"} first.`);
        await requestJSON("/api/listen/selection",jsonOptions("PUT",{[`${source}_id`]:id}));
    }

    async function selectSource(source,id) {
        if (busy || !id) return;
        const previous = source === "masjid" ? selectedMasjidID : selectedRadioID;
        if (source === "masjid") selectedMasjidID = id; else selectedRadioID = id;
        renderSources();
        setBusy(true);
        setConnectionError();
        try {
            await ensureSelection(source,id);
            status = await requestJSON("/api/listen/status");
        } catch (error) {
            if (source === "masjid") selectedMasjidID = previous; else selectedRadioID = previous;
            setConnectionError(error.message);
        } finally { setBusy(false); }
    }

    favouriteHost.addEventListener("click", event => {
        const button = event.target.closest("button[data-source-id]");
        if (button) selectSource("masjid",button.dataset.sourceId);
    });
    radioHost.addEventListener("click", event => {
        const button = event.target.closest("button[data-source-id]");
        if (button) selectSource("radio",button.dataset.sourceId);
    });
    themeHost.addEventListener("click", event => {
        const button = event.target.closest("button[data-theme]");
        if (!button) return;
        runAction(async () => {
            const saved = await requestJSON("/api/masjidboard/layout",jsonOptions("PUT",{theme:button.dataset.theme}));
            currentTheme = saved.theme || button.dataset.theme;
            document.body.dataset.boardTheme = currentTheme;
            renderThemes();
        },false);
    });

    brightness.addEventListener("input", () => {
        brightnessValue.textContent = `${brightness.value}%`;
        window.clearTimeout(brightnessSaveTimer);
        brightnessSaveTimer = window.setTimeout(async () => {
            try {
                displaySettings = await requestJSON("/api/display/settings",jsonOptions("PUT",{brightness_percent:Number(brightness.value)}));
                setConnectionError();
            } catch (error) { setConnectionError(error.message); }
            renderDisplaySettings();
        },150);
    });

    bottomPanel.querySelectorAll("[data-touch-tab]").forEach(button =>
        button.addEventListener("click",() => activateTab(button.dataset.touchTab)));
    bottomPanel.querySelectorAll("[data-listen-close]").forEach(button =>
        button.addEventListener("click",() => setOpenPanel("")));
    quickPanel.querySelectorAll("[data-quick-close]").forEach(button =>
        button.addEventListener("click",() => setOpenPanel("")));

    async function saveVolume(name,value) {
        window.clearTimeout(volumeSaveTimers[name]);
        const serial = ++volumeSaveSerials[name];
        pendingVolumes[name] = value;
        volumeOutputs[name].textContent = `${value}%`;
        try {
            if (name === "master") await requestJSON("/api/player/volume",jsonOptions("POST",{volume:value,persist:true}));
            else await requestJSON("/api/listen/volume",jsonOptions("PUT",{source:name,volume:value}));
            const refreshed = await requestJSON("/api/listen/status");
            if (serial === volumeSaveSerials[name]) { status = refreshed; setConnectionError(); }
        } catch (error) {
            if (serial === volumeSaveSerials[name]) setConnectionError(error.message);
        } finally {
            if (serial === volumeSaveSerials[name]) { pendingVolumes[name] = null; renderStatus(); }
        }
    }

    function scheduleVolumeSave(name) {
        const value = Number(volumeControls[name].value);
        pendingVolumes[name] = value;
        volumeOutputs[name].textContent = `${value}%`;
        window.clearTimeout(volumeSaveTimers[name]);
        volumeSaveTimers[name] = window.setTimeout(() => saveVolume(name,value),120);
    }

    for (const [name,control] of Object.entries(volumeControls)) {
        control.addEventListener("input",() => scheduleVolumeSave(name));
        control.addEventListener("change",() => saveVolume(name,Number(control.value)));
    }
    state.addEventListener("click", event => {
        const stepButton = event.target.closest("[data-volume-step]");
        if (!stepButton) return;
        const [name,amount] = stepButton.dataset.volumeStep.split(":");
        const control = volumeControls[name];
        control.value = Math.max(Number(control.min),Math.min(Number(control.max),Number(control.value) + Number(amount)));
        saveVolume(name,Number(control.value));
    });

    for (const sheet of state.querySelectorAll(".appliance-listen-sheet,.appliance-quick-sheet")) {
        for (const eventName of ["pointerdown","keydown","input","change"]) sheet.addEventListener(eventName,resetInactivityTimer);
        sheet.addEventListener("scroll",resetInactivityTimer,true);
    }

    playMasjid.addEventListener("click",() => runAction(async () => {
        await ensureSelection("masjid",selectedMasjidID);
        if (!status?.masjid_enabled) await requestJSON("/api/listen/power",jsonOptions("PUT",{module:"masjid",enabled:true}));
        await requestJSON("/api/listen/start",{method:"POST"});
    }));
    stopListening.addEventListener("click",() => runAction(() => requestJSON("/api/listen/stop",{method:"POST"})));

    for (const button of radioModeButtons) {
        button.addEventListener("click",() => runAction(async () => {
            const mode = button.dataset.radioMode;
            if (mode !== "stopped") await ensureSelection("radio",selectedRadioID);
            if (!status?.radio_enabled && mode !== "stopped")
                await requestJSON("/api/listen/power",jsonOptions("PUT",{module:"radio",enabled:true}));
            if (!status?.listening && mode !== "stopped") await requestJSON("/api/listen/start",{method:"POST"});
            await requestJSON("/api/listen/radio-mode",jsonOptions("PUT",{mode}));
        }));
    }

    state.addEventListener("pointerdown", event => {
        if (openPanel) return;
        gestureStart = {x:event.clientX,y:event.clientY,height:state.clientHeight};
    });
    state.addEventListener("pointerup", event => {
        if (!gestureStart || openPanel) return;
        const start = gestureStart;
        gestureStart = null;
        const dx = event.clientX - start.x;
        const dy = event.clientY - start.y;
        if (Math.abs(dy) <= Math.abs(dx)) return;
        if (start.y <= 120 && dy > 70) setOpenPanel("quick");
        else if (start.y >= start.height - 120 && dy < -70) setOpenPanel("bottom");
    });

    function wireCloseGesture(panel,selector,direction) {
        panel.addEventListener("pointerdown",event => event.stopPropagation());
        panel.addEventListener("pointermove",event => {
            if (!closeGestureStart || closeGestureStart.panel !== panel || closeGestureStart.pointerId !== event.pointerId) return;
            const dx = event.clientX - closeGestureStart.x;
            const dy = event.clientY - closeGestureStart.y;
            if (direction * dy > 45 && Math.abs(dy) > Math.abs(dx)) {
                const target = closeGestureStart.target;
                closeGestureStart = null;
                if (target.hasPointerCapture?.(event.pointerId)) target.releasePointerCapture(event.pointerId);
                setOpenPanel("");
            }
            event.stopPropagation();
        });
        for (const eventName of ["pointerup","pointercancel"]) panel.addEventListener(eventName,event => {
            closeGestureStart = null;
            event.stopPropagation();
        });
        panel.querySelectorAll(selector).forEach(target => target.addEventListener("pointerdown",event => {
            if (event.target.closest("button,a,input")) return;
            closeGestureStart = {x:event.clientX,y:event.clientY,pointerId:event.pointerId,target,panel};
            target.setPointerCapture?.(event.pointerId);
            event.preventDefault();
            event.stopPropagation();
        }));
    }
    wireCloseGesture(bottomPanel,".appliance-listen-handle,.appliance-listen-heading",1);
    wireCloseGesture(quickPanel,".appliance-quick-handle,.appliance-quick-heading",-1);
    document.addEventListener("keydown",event => { if (openPanel && event.key === "Escape") setOpenPanel(""); });
})();