let masjidCatalogue = [];
let radioCatalogue = [];
let backendOnline = true;
let listenStatus = null;
let favouriteIds = new Set();
let renderedMasjidID = null;
let renderedRadioID = null;

function publishListenStatus(status) {
    window.dispatchEvent(new CustomEvent("masjidpi:listen-status", {detail: status}));
}

function publishMasjidCatalogue(items) {
    window.dispatchEvent(new CustomEvent("masjidpi:masjid-catalogue", {detail: items}));
}

const state = document.getElementById("state");
const statusDetail = document.getElementById("statusDetail");
const currentStream = document.getElementById("url");
const sourceExplanation = document.getElementById("sourceExplanation");
const streamInput = document.getElementById("stream");
const radioInput = document.getElementById("radioStream");
const streamSearch = document.getElementById("streamSearch");
const streamCount = document.getElementById("streamCount");
const favouritesSection = document.getElementById("favouritesSection");
const favourites = document.getElementById("favourites");
const favouriteButton = document.getElementById("favouriteButton");
const masjidVolumeSlider = document.getElementById("masjidVolumeSlider");
const masjidVolumeValue = document.getElementById("masjidVolumeValue");
const radioVolumeSlider = document.getElementById("radioVolumeSlider");
const radioVolumeValue = document.getElementById("radioVolumeValue");
const volumeSlider = document.getElementById("volumeSlider");
const volumeValue = document.getElementById("volumeValue");
const volumeNote = document.getElementById("volumeNote");
const audioDevice = document.getElementById("audioDevice");
const refreshAudioDevicesButton = document.getElementById("refreshAudioDevices");
const playButton = document.getElementById("play");
const stopButton = document.getElementById("stop");
const updateCatalogueButton = document.getElementById("updateCatalogueButton");

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

function jsonOptions(method, body) {
    return {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    };
}

async function getListenStatus() {
    return requestJSON("/api/listen/status");
}

async function getStreams(kind) {
    return requestJSON(`/api/streams?kind=${encodeURIComponent(kind)}`);
}

async function setSelection(body) {
    return requestJSON("/api/listen/selection", jsonOptions("PUT", body));
}

async function setSourceVolume(source, volume) {
    return requestJSON("/api/listen/volume", jsonOptions("PUT", { source, volume }));
}

async function startListening() {
    return requestJSON("/api/listen/start", { method: "POST" });
}

async function stopListening() {
    return requestJSON("/api/listen/stop", { method: "POST" });
}

async function getAudioDevices() {
    return requestJSON("/api/player/volume?devices=1");
}

async function setAudioDevice(name) {
    return requestJSON("/api/player/volume", jsonOptions("POST", { audio_device: name }));
}

async function loadFavourites() {
    const data = await requestJSON("/api/favourites");
    favouriteIds = new Set(data.ids || []);
}

async function saveFavourites() {
    await requestJSON("/api/favourites", jsonOptions("PUT", { ids: [...favouriteIds] }));
}

async function updateCatalogue() {
    return requestJSON("/api/catalogue/update", { method: "POST" });
}

function showToast(message, type = "success") {
    window.MasjidPiUI.notify(message, type);
}

function setBusy(button, busy, busyText, normalText) {
    button.disabled = busy;
    button.textContent = busy ? busyText : normalText;
}

function streamLabel(item) {
    return item.location ? `${item.name} — ${item.location}` : item.name;
}

function streamMatchesQuery(item, query) {
    if (!query) return true;
    return [item.name, item.location, item.id]
        .filter(Boolean)
        .some(value => value.toLowerCase().includes(query));
}

function renderMasjids(preferredId = listenStatus?.masjid_id || streamInput.value) {
    const query = streamSearch.value.trim().toLowerCase();
    const filtered = masjidCatalogue.filter(item => streamMatchesQuery(item, query));
    streamInput.innerHTML = "";

    for (const item of filtered) {
        const option = document.createElement("option");
        option.value = item.id;
        option.textContent = streamLabel(item);
        streamInput.appendChild(option);
    }

    if (preferredId && filtered.some(item => item.id === preferredId)) {
        streamInput.value = preferredId;
    }

    streamCount.textContent = query
        ? `${filtered.length} of ${masjidCatalogue.length} masjids`
        : `${masjidCatalogue.length} masjids`;
    streamCount.classList.toggle("hidden", masjidCatalogue.length === 0);
    renderFavourites();
    updateFavouriteButton();
}

function renderRadios(preferredId = listenStatus?.radio_id || radioInput.value) {
    radioInput.innerHTML = "";
    for (const item of radioCatalogue) {
        const option = document.createElement("option");
        option.value = item.id;
        option.textContent = streamLabel(item);
        radioInput.appendChild(option);
    }
    if (preferredId && radioCatalogue.some(item => item.id === preferredId)) {
        radioInput.value = preferredId;
    }
}

function renderFavourites() {
    favourites.innerHTML = "";
    const favouriteStreams = masjidCatalogue.filter(item => favouriteIds.has(item.id));
    favouritesSection.classList.toggle("hidden", favouriteStreams.length === 0);

    for (const item of favouriteStreams) {
        const button = document.createElement("button");
        button.type = "button";
        button.className = "favourite-item";
        button.dataset.id = item.id;

        const label = document.createElement("span");
        label.textContent = `★ ${streamLabel(item)}`;

        const remove = document.createElement("span");
        remove.className = "favourite-remove";
        remove.textContent = "×";
        remove.title = "Remove from favourites";
        remove.setAttribute("aria-label", `Remove ${item.name} from favourites`);

        button.append(label, remove);
        favourites.appendChild(button);
    }
}

function updateFavouriteButton() {
    const selected = masjidCatalogue.find(item => item.id === streamInput.value);
    const isFavourite = selected && favouriteIds.has(selected.id);
    favouriteButton.disabled = !backendOnline || !selected;
    favouriteButton.textContent = isFavourite ? "★ Remove from Favourites" : "☆ Add to Favourites";
}

async function loadStreams() {
    const [masjids, radios] = await Promise.all([
        getStreams("masjid"),
        getStreams("radio")
    ]);
    masjidCatalogue = masjids;
    radioCatalogue = radios;
    publishMasjidCatalogue(masjids);
    renderMasjids();
    renderRadios();
}

async function loadAudioDevices(options = {}) {
    try {
        const previous = audioDevice.value;
        const devices = await getAudioDevices();
        audioDevice.innerHTML = "";
        for (const device of devices) {
            const option = document.createElement("option");
            option.value = device.name;
            option.textContent = (device.description || device.name) + (device.unavailable ? " — unavailable" : "");
            option.disabled = Boolean(device.unavailable);
            audioDevice.appendChild(option);
        }
        if (previous && [...audioDevice.options].some(option => option.value === previous)) {
            audioDevice.value = previous;
        }
        audioDevice.disabled = !backendOnline || devices.length === 0;
    } catch (err) {
        console.error(err);
        audioDevice.innerHTML = "";
        const option = document.createElement("option");
        option.textContent = "Unable to detect audio devices";
        audioDevice.appendChild(option);
        audioDevice.disabled = true;
        if (options.throwOnError) throw err;
    }
}

function findStream(id) {
    return masjidCatalogue.find(item => item.id === id) || radioCatalogue.find(item => item.id === id);
}

function setOffline(offline) {
    const radioPowered = Boolean(listenStatus?.radio_enabled);
    document.getElementById("offlineBanner").classList.toggle("hidden", !offline);
    streamInput.disabled = offline;
    radioInput.disabled = offline || !radioPowered;
    streamSearch.disabled = offline;
    favouriteButton.disabled = offline;
    masjidVolumeSlider.disabled = offline;
    radioVolumeSlider.disabled = offline || !radioPowered;
    volumeSlider.disabled = offline || !listenStatus?.master_volume_supported;
    audioDevice.disabled = offline || audioDevice.options.length === 0;
    refreshAudioDevicesButton.disabled = offline;
    playButton.disabled = offline || Boolean(listenStatus?.listening);
    stopButton.disabled = offline || !listenStatus?.listening;
    updateCatalogueButton.disabled = offline;
    document.getElementById("masjidPowerSwitch").disabled = offline;
    document.getElementById("radioPowerSwitch").disabled = offline;
}

function updateControlAvailability(status) {
    const radioEnabled = backendOnline && Boolean(status.radio_enabled);
    for (const id of [
        "radioModeSchedule", "radioModePlayNow", "radioModeStop", "radioVolumeSlider",
        "radioResumeDelaySlider", "radioScheduleEnabled", "radioStream"
    ]) {
        const control = document.getElementById(id);
        if (control) control.disabled = !radioEnabled;
    }
    for (const id of ["radioScheduleStart", "radioScheduleStop"]) {
        const control = document.getElementById(id);
        if (control) control.disabled = !radioEnabled || !document.getElementById("radioScheduleEnabled")?.checked;
    }
}

function renderStatus(status) {
    listenStatus = status;

    masjidVolumeSlider.value = status.masjid_volume;
    masjidVolumeValue.textContent = status.masjid_volume + "%";
    radioVolumeSlider.value = status.radio_volume;
    radioVolumeValue.textContent = status.radio_volume + "%";
    volumeSlider.value = status.master_volume;
    volumeValue.textContent = status.master_volume + "%";

    if (status.audio_device && [...audioDevice.options].some(option => option.value === status.audio_device)) {
        audioDevice.value = status.audio_device;
    }

    volumeSlider.disabled = !status.master_volume_supported;
    radioInput.disabled = !status.radio_enabled;
    radioVolumeSlider.disabled = !status.radio_enabled;
    volumeNote.textContent = status.master_volume_supported
        ? "Controls the selected audio device's hardware volume."
        : "The selected audio device does not provide a controllable hardware mixer.";

    if (status.masjid_id !== renderedMasjidID) {
        renderedMasjidID = status.masjid_id || "";
        if (status.masjid_id && [...streamInput.options].some(option => option.value === status.masjid_id)) {
            streamInput.value = status.masjid_id;
        }
    }
    if (status.radio_id !== renderedRadioID) {
        renderedRadioID = status.radio_id || "";
        if (status.radio_id && [...radioInput.options].some(option => option.value === status.radio_id)) {
            radioInput.value = status.radio_id;
        }
    }

    playButton.disabled = status.listening;
    stopButton.disabled = !status.listening;
    updateFavouriteButton();
    publishListenStatus(status);
    updateControlAvailability(status);

    const detail = status.error || "";
    statusDetail.textContent = detail;
    statusDetail.className = detail ? "status-detail status-detail-error" : "status-detail hidden";

    if (!status.listening) {
        state.textContent = "Stopped";
        state.className = "status-badge status-idle";
        currentStream.textContent = "Listening is stopped";
        sourceExplanation.textContent = "Start Listening to use Masjid priority with Radio as the secondary source.";
        return;
    }

    if (status.active_source === "masjid") {
        const current = findStream(status.active_stream_id);
        state.textContent = status.playback_state === "playing" ? "Masjid" : status.playback_state;
        state.className = "status-badge status-playing";
        currentStream.textContent = current ? streamLabel(current) : "Selected masjid";
        const radio = findStream(status.radio_id);
        sourceExplanation.textContent = radio
            ? `${radio.name} is standing by and will resume when the masjid goes offline.`
            : "The selected masjid has priority.";
        return;
    }

    if (status.active_source === "radio") {
        const current = findStream(status.active_stream_id);
        state.textContent = status.playback_state === "playing" ? "Radio" : status.playback_state;
        state.className = "status-badge status-playing";
        currentStream.textContent = current ? streamLabel(current) : "Selected radio station";
        sourceExplanation.textContent = status.masjid_id
            ? "Radio is playing while the selected masjid is offline. It will stop automatically when the masjid comes online."
            : "Radio is playing. Select a primary masjid to enable automatic priority switching.";
        return;
    }

    state.textContent = "Waiting";
    state.className = "status-badge status-waiting";
    currentStream.textContent = "No source currently playing";
    if (status.radio_resume_pending) {
        sourceExplanation.textContent = "The masjid broadcast has ended. Radio will resume after the configured delay.";
        return;
    }
    sourceExplanation.textContent = status.masjid_id
        ? "Waiting for the selected masjid. Select a radio station for continuous secondary audio."
        : "Select a primary masjid and/or secondary radio station.";
}

async function refreshStatus() {
    try {
        const status = await getListenStatus();
        if (!backendOnline) {
            backendOnline = true;
            showToast("Connection to MasjidPi restored.", "success");
        }
        renderStatus(status);
        setOffline(false);
    } catch (err) {
        console.error(err);
        if (backendOnline) {
            backendOnline = false;
            showToast("Connection to MasjidPi lost.", "error");
        }
        state.textContent = "Offline";
        state.className = "status-badge status-error";
        statusDetail.textContent = "Unable to reach MasjidPi.";
        statusDetail.className = "status-detail status-detail-error";
        setOffline(true);
    }
}

window.MasjidPiRefreshListenStatus = refreshStatus;

function activateTab(name) {
    document.querySelectorAll("[data-listen-tab]").forEach(button => {
        const active = button.dataset.listenTab === name;
        button.classList.toggle("active", active);
        button.setAttribute("aria-selected", active ? "true" : "false");
        button.tabIndex = active ? 0 : -1;
    });
    document.querySelectorAll("[data-listen-panel]").forEach(panel => {
        panel.classList.toggle("hidden", panel.dataset.listenPanel !== name);
    });
    sessionStorage.setItem("masjidpi-listen-tab", name);
}

document.querySelectorAll("[data-listen-tab]").forEach(button => {
    button.addEventListener("click", () => activateTab(button.dataset.listenTab));
    button.addEventListener("keydown", event => {
        if (event.key !== "ArrowLeft" && event.key !== "ArrowRight") return;
        const tabs = Array.from(document.querySelectorAll("[data-listen-tab]"));
        const direction = event.key === "ArrowRight" ? 1 : -1;
        const next = tabs[(tabs.indexOf(button) + direction + tabs.length) % tabs.length];
        activateTab(next.dataset.listenTab);
        next.focus();
    });
});

const savedListenTab = sessionStorage.getItem("masjidpi-listen-tab");
if (["masjid", "radio", "audio"].includes(savedListenTab)) activateTab(savedListenTab);

streamSearch.addEventListener("input", () => renderMasjids(streamInput.value));

streamInput.addEventListener("change", async () => {
    const id = streamInput.value;
    if (!id) return;
    try {
        listenStatus = await setSelection({ masjid_id: id });
        renderedMasjidID = id;
        updateFavouriteButton();
        await refreshStatus();
        showToast("Primary masjid updated.", "success");
    } catch (err) {
        showToast(err.message, "error");
        await refreshStatus();
    }
});

radioInput.addEventListener("change", async () => {
    const id = radioInput.value;
    if (!id) return;
    try {
        listenStatus = await setSelection({ radio_id: id });
        renderedRadioID = id;
        await refreshStatus();
        showToast("Secondary radio updated.", "success");
    } catch (err) {
        showToast(err.message, "error");
        await refreshStatus();
    }
});

favourites.addEventListener("click", async event => {
    const item = event.target.closest(".favourite-item");
    if (!item) return;
    const id = item.dataset.id;

    if (event.target.closest(".favourite-remove")) {
        favouriteIds.delete(id);
        try {
            await saveFavourites();
            renderFavourites();
            updateFavouriteButton();
            showToast("Removed from favourites.", "success");
        } catch (err) {
            favouriteIds.add(id);
            showToast(err.message, "error");
        }
        return;
    }

    streamSearch.value = "";
    renderMasjids(id);
    streamInput.value = id;
    try {
        await setSelection({ masjid_id: id });
        renderedMasjidID = id;
        await refreshStatus();
    } catch (err) {
        showToast(err.message, "error");
    }
});

favouriteButton.addEventListener("click", async () => {
    const id = streamInput.value;
    if (!id) return;
    const wasFavourite = favouriteIds.has(id);
    if (wasFavourite) favouriteIds.delete(id); else favouriteIds.add(id);

    try {
        await saveFavourites();
        renderFavourites();
        updateFavouriteButton();
        showToast(wasFavourite ? "Removed from favourites." : "Added to favourites.", "success");
    } catch (err) {
        if (wasFavourite) favouriteIds.add(id); else favouriteIds.delete(id);
        renderFavourites();
        updateFavouriteButton();
        showToast(err.message, "error");
    }
});

masjidVolumeSlider.addEventListener("input", () => {
    masjidVolumeValue.textContent = masjidVolumeSlider.value + "%";
});

masjidVolumeSlider.addEventListener("change", async () => {
    try {
        await setSourceVolume("masjid", Number(masjidVolumeSlider.value));
        await refreshStatus();
    } catch (err) {
        showToast(err.message, "error");
        await refreshStatus();
    }
});

radioVolumeSlider.addEventListener("input", () => {
    radioVolumeValue.textContent = radioVolumeSlider.value + "%";
});

radioVolumeSlider.addEventListener("change", async () => {
    try {
        await setSourceVolume("radio", Number(radioVolumeSlider.value));
        await refreshStatus();
    } catch (err) {
        showToast(err.message, "error");
        await refreshStatus();
    }
});

playButton.addEventListener("click", async () => {
    setBusy(playButton, true, "Starting...", "▶ Start Listening");
    try {
        await startListening();
        await refreshStatus();
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        playButton.textContent = "▶ Start Listening";
    }
});

stopButton.addEventListener("click", async () => {
    setBusy(stopButton, true, "Stopping...", "■ Stop");
    try {
        await stopListening();
        await refreshStatus();
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        stopButton.textContent = "■ Stop";
    }
});

audioDevice.addEventListener("change", async () => {
    const selected = audioDevice.value;
    if (!selected) return;
    audioDevice.disabled = true;
    try {
        await setAudioDevice(selected);
        showToast("Audio output changed.", "success");
        await refreshStatus();
    } catch (err) {
        showToast(err.message, "error");
        await loadAudioDevices();
        await refreshStatus();
    } finally {
        audioDevice.disabled = false;
    }
});

refreshAudioDevicesButton.addEventListener("click", async () => {
    setBusy(refreshAudioDevicesButton, true, "Refreshing...", "Refresh Devices");
    try {
        await loadAudioDevices({throwOnError: true});
        await refreshStatus();
        showToast("Audio devices refreshed.", "success");
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        setBusy(refreshAudioDevicesButton, false, "Refreshing...", "Refresh Devices");
    }
});

updateCatalogueButton.addEventListener("click", async () => {
    setBusy(updateCatalogueButton, true, "Updating...", "🔄 Update Masjid Catalogue");
    try {
        await updateCatalogue();
        await loadStreams();
        await refreshStatus();
        showToast("Masjid catalogue updated successfully.", "success");
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        updateCatalogueButton.textContent = "🔄 Update Masjid Catalogue";
    }
});

async function initialize() {
    try {
        await Promise.all([loadFavourites(), loadAudioDevices()]);
        listenStatus = await getListenStatus();
        await loadStreams();
        renderStatus(listenStatus);
        setOffline(false);
    } catch (err) {
        console.error(err);
        backendOnline = false;
        setOffline(true);
    }

    const scheduleRefresh = () => {
        window.setTimeout(async () => {
            await refreshStatus();
            scheduleRefresh();
        }, document.hidden ? 5000 : 1000);
    };
    scheduleRefresh();
}

initialize();
