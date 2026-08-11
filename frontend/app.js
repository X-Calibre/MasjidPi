let catalogue = [];
let filteredCatalogue = [];
let backendOnline = true;
let playerStatus = null;
let preferences = { last_stream_id: "", autoplay: false };
let favouriteIds = new Set();

const state = document.getElementById("state");
const statusDetail = document.getElementById("statusDetail");
const volume = document.getElementById("volume");
const volumeNote = document.getElementById("volumeNote");
const stream = document.getElementById("url");
const streamInput = document.getElementById("stream");
const streamSearch = document.getElementById("streamSearch");
const streamCount = document.getElementById("streamCount");
const favouritesSection = document.getElementById("favouritesSection");
const favourites = document.getElementById("favourites");
const favouriteButton = document.getElementById("favouriteButton");
const volumeSlider = document.getElementById("volumeSlider");
const volumeValue = document.getElementById("volumeValue");
const audioDevice = document.getElementById("audioDevice");
const playButton = document.getElementById("play");
const stopButton = document.getElementById("stop");
const updateCatalogueButton = document.getElementById("updateCatalogueButton");
const autoplay = document.getElementById("autoplay");

async function getStatus() {
    const response = await fetch("/api/player/status");
    if (!response.ok) throw new Error("Unable to get player status");
    return response.json();
}

async function getAudioDevices() {
    const response = await fetch("/api/player/volume?devices=1");
    if (!response.ok) throw new Error("Unable to get audio devices");
    return response.json();
}

async function setAudioDevice(name) {
    const response = await fetch("/api/player/volume", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ audio_device: name })
    });
    if (!response.ok) throw new Error("Unable to change audio device");
    return response.json();
}

async function getPreferences() {
    const response = await fetch("/api/preferences");
    if (!response.ok) throw new Error("Unable to load preferences");
    return response.json();
}

async function savePreferences(next) {
    const updated = { ...preferences, ...next };
    const response = await fetch("/api/preferences", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updated)
    });
    if (!response.ok) throw new Error("Unable to save preferences");
    preferences = await response.json();
}

async function loadFavourites() {
    const response = await fetch("/api/favourites");
    if (!response.ok) throw new Error("Unable to load favourites");
    const data = await response.json();
    favouriteIds = new Set(data.ids || []);
}

async function saveFavourites() {
    const response = await fetch("/api/favourites", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ids: [...favouriteIds] })
    });
    if (!response.ok) throw new Error("Unable to save favourites");
}

function streamMatchesQuery(item, query) {
    if (!query) return true;
    return [item.name, item.location, item.id]
        .filter(Boolean)
        .some(value => value.toLowerCase().includes(query));
}

function streamLabel(item) {
    return item.location ? `${item.name} — ${item.location}` : item.name;
}

function renderFavourites(query) {
    favourites.innerHTML = "";
    const favouriteStreams = catalogue.filter(item => favouriteIds.has(item.id) && streamMatchesQuery(item, query));
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

function renderStreams(preferredId = streamInput.value) {
    const query = streamSearch.value.trim().toLowerCase();
    filteredCatalogue = catalogue.filter(item => streamMatchesQuery(item, query));
    streamInput.innerHTML = "";

    for (const item of filteredCatalogue) {
        const option = document.createElement("option");
        option.value = item.id;
        option.textContent = streamLabel(item);
        streamInput.appendChild(option);
    }

    if (preferredId && filteredCatalogue.some(item => item.id === preferredId)) {
        streamInput.value = preferredId;
    } else if (filteredCatalogue.length > 0 && !streamInput.value) {
        streamInput.selectedIndex = 0;
    }

    streamCount.textContent = query ? `${filteredCatalogue.length} of ${catalogue.length} masjids` : `${catalogue.length} masjids`;
    streamCount.classList.toggle("hidden", catalogue.length === 0);
    renderFavourites(query);
    updateFavouriteButton();
}

async function loadStreams() {
    const currentSelection = streamInput.value;
    const response = await fetch("/api/streams");
    if (!response.ok) {
        console.error("Unable to load streams");
        return;
    }
    catalogue = await response.json();
    const preferred = currentSelection || preferences.last_stream_id;
    renderStreams(preferred);
}

async function loadAudioDevices() {
    try {
        const devices = await getAudioDevices();
        audioDevice.innerHTML = "";
        for (const device of devices) {
            const option = document.createElement("option");
            option.value = device.name;
            option.textContent = device.description ? device.description : device.name;
            audioDevice.appendChild(option);
        }
        audioDevice.disabled = !backendOnline || devices.length === 0;
    } catch (err) {
        console.error(err);
        audioDevice.innerHTML = "";
        const option = document.createElement("option");
        option.textContent = "Unable to detect audio devices";
        audioDevice.appendChild(option);
        audioDevice.disabled = true;
    }
}

async function playStream(id) {
    const response = await fetch("/api/player/play", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id })
    });
    if (!response.ok) throw new Error("Unable to play stream");
}

async function stopStream() {
    const response = await fetch("/api/player/stop", { method: "POST" });
    if (!response.ok) throw new Error("Unable to stop playback");
}

async function setVolume(level) {
    const response = await fetch("/api/player/volume", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ volume: level })
    });
    if (!response.ok) throw new Error("Unable to change volume");
}

async function updateCatalogue() {
    const response = await fetch("/api/catalogue/update", { method: "POST" });
    if (!response.ok) throw new Error("Unable to update catalogue");
    return response.json();
}

function showToast(message, type = "success") {
    const container = document.getElementById("toastContainer");
    const toast = document.createElement("div");
    toast.className = "toast toast-" + type;
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(() => {
        toast.style.opacity = "0";
        toast.style.transition = "opacity .3s";
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

function setBusy(button, busy, busyText, normalText) {
    button.disabled = busy;
    if (busy) {
        button.dataset.label = normalText;
        button.textContent = busyText;
    } else {
        button.textContent = button.dataset.label || normalText;
    }
}

function setOffline(offline) {
    const banner = document.getElementById("offlineBanner");
    banner.classList.toggle("hidden", !offline);
    playButton.disabled = offline;
    stopButton.disabled = offline;
    streamInput.disabled = offline;
    streamSearch.disabled = offline;
    favouriteButton.disabled = offline;
    volumeSlider.disabled = offline || !playerStatus?.volume_supported;
    audioDevice.disabled = offline || audioDevice.options.length === 0;
    updateCatalogueButton.disabled = offline;
}

function updateFavouriteButton() {
    const selectedId = streamInput.value;
    const selected = catalogue.find(item => item.id === selectedId);
    const isFavourite = selected && favouriteIds.has(selected.id);
    favouriteButton.disabled = !backendOnline || !selected;
    favouriteButton.textContent = isFavourite ? "★ Remove from Favourites" : "☆ Add to Favourites";
}

function updateControls(status) {
    if (!backendOnline) {
        setOffline(true);
        return;
    }
    const active = ["waiting", "connecting", "playing", "retrying"].includes(status.state);
    const selectedCurrentStream = active && status.stream_id && streamInput.value === status.stream_id;
    playButton.disabled = selectedCurrentStream;
    stopButton.disabled = !active;
    streamInput.disabled = false;
    streamSearch.disabled = false;
    volumeSlider.disabled = !status.volume_supported;
    volumeNote.textContent = status.volume_supported
        ? "Controls the selected audio device's hardware volume."
        : "The selected audio device does not provide a controllable hardware mixer.";
    audioDevice.disabled = audioDevice.options.length === 0;
    updateCatalogueButton.disabled = false;
    updateFavouriteButton();
}

function updateStatusDetail(status) {
    const detail = status.error || "";
    statusDetail.textContent = detail;
    statusDetail.className = detail ? "status-detail status-detail-" + status.state : "status-detail hidden";
}

function findStreamByURL(url) {
    return catalogue.find(stream => stream.url === url);
}

async function refreshStatus() {
    try {
        const status = await getStatus();
        playerStatus = status;
        if (!backendOnline) {
            backendOnline = true;
            setOffline(false);
            showToast("Connection to MasjidPi restored.", "success");
        }
        document.getElementById("version").textContent = "MasjidPi " + status.version;
        state.textContent = status.message || status.state;
        state.className = "status-badge status-" + status.state;
        updateStatusDetail(status);
        volume.textContent = status.volume + "%";
        volumeSlider.value = status.volume;
        volumeValue.textContent = status.volume + "%";
        if (status.audio_device && [...audioDevice.options].some(option => option.value === status.audio_device)) {
            audioDevice.value = status.audio_device;
        }
        updateControls(status);
        if (!status.url) {
            stream.textContent = "No stream playing";
        } else {
            const current = findStreamByURL(status.url);
            stream.textContent = current ? current.name + (current.location ? " — " + current.location : "") + "\n" + current.url : status.url;
        }
    } catch (err) {
        console.error(err);
        if (backendOnline) {
            backendOnline = false;
            setOffline(true);
            state.textContent = "Offline";
            state.className = "status-badge status-error";
            statusDetail.textContent = "Unable to reach MasjidPi.";
            statusDetail.className = "status-detail status-detail-error";
            showToast("Connection to MasjidPi lost.", "error");
        }
    }
}

streamSearch.addEventListener("input", () => {
    renderStreams(streamInput.value);
    if (playerStatus) updateControls(playerStatus);
});

streamInput.addEventListener("change", async () => {
    const id = streamInput.value;
    if (!id) return;
    try {
        await savePreferences({ last_stream_id: id });
    } catch (err) {
        showToast(err.message, "error");
    }
    updateFavouriteButton();
    if (playerStatus) updateControls(playerStatus);
});

favourites.addEventListener("click", async event => {
    const item = event.target.closest(".favourite-item");
    if (!item) return;
    const id = item.dataset.id;
    if (event.target.closest(".favourite-remove")) {
        favouriteIds.delete(id);
        try {
            await saveFavourites();
            showToast("Removed from favourites.", "success");
        } catch (err) {
            showToast(err.message, "error");
        }
        renderStreams(streamInput.value);
        if (playerStatus) updateControls(playerStatus);
        return;
    }
    if (catalogue.some(stream => stream.id === id)) {
        streamInput.value = id;
        try {
            await savePreferences({ last_stream_id: id });
        } catch (err) {
            showToast(err.message, "error");
        }
        updateFavouriteButton();
        if (playerStatus) updateControls(playerStatus);
    }
});

favouriteButton.addEventListener("click", async () => {
    const id = streamInput.value;
    if (!id) return;
    const wasFavourite = favouriteIds.has(id);
    if (wasFavourite) {
        favouriteIds.delete(id);
    } else {
        favouriteIds.add(id);
    }
    try {
        await saveFavourites();
        showToast(wasFavourite ? "Removed from favourites." : "Added to favourites.", "success");
    } catch (err) {
        if (wasFavourite) favouriteIds.add(id); else favouriteIds.delete(id);
        showToast(err.message, "error");
    }
    renderStreams(id);
    if (playerStatus) updateControls(playerStatus);
});

playButton.addEventListener("click", async () => {
    if (!streamInput.value) {
        showToast("Please select a masjid.", "warning");
        return;
    }
    setBusy(playButton, true, "Playing...", "▶ Play");
    try {
        await playStream(streamInput.value);
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        setBusy(playButton, false, "Playing...", "▶ Play");
        await refreshStatus();
    }
});

stopButton.addEventListener("click", async () => {
    setBusy(stopButton, true, "Stopping...", "■ Stop");
    try {
        await stopStream();
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        setBusy(stopButton, false, "Stopping...", "■ Stop");
        await refreshStatus();
    }
});

volumeSlider.addEventListener("input", async () => {
    const value = Number(volumeSlider.value);
    volumeValue.textContent = value + "%";
    try {
        await setVolume(value);
        await refreshStatus();
    } catch (err) {
        console.error(err);
        await refreshStatus();
    }
});

audioDevice.addEventListener("change", async () => {
    const selected = audioDevice.value;
    if (!selected) return;
    audioDevice.disabled = true;
    try {
        const status = await setAudioDevice(selected);
        playerStatus = status;
        updateControls(status);
        showToast("Audio output changed.", "success");
    } catch (err) {
        showToast(err.message, "error");
        await loadAudioDevices();
        await refreshStatus();
    }
});

autoplay.addEventListener("change", async () => {
    try {
        await savePreferences({ autoplay: autoplay.checked });
    } catch (err) {
        showToast(err.message, "error");
    }
});

updateCatalogueButton.addEventListener("click", async () => {
    try {
        setBusy(updateCatalogueButton, true, "Updating...", "🔄 Update Catalogue");
        await updateCatalogue();
        await loadStreams();
        await refreshStatus();
        showToast("Catalogue updated successfully.", "success");
    } catch (err) {
        showToast(err.message, "error");
    } finally {
        setBusy(updateCatalogueButton, false, "Updating...", "🔄 Update Catalogue");
    }
});

async function initialize() {
    try {
        preferences = await getPreferences();
        await loadFavourites();
    } catch (err) {
        console.error(err);
        showToast("Unable to load saved preferences.", "error");
    }

    autoplay.checked = preferences.autoplay === true;
    await loadStreams();
    await loadAudioDevices();
    await refreshStatus();
    setInterval(refreshStatus, 1000);
}

initialize();
