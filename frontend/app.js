let catalogue = [];
let filteredCatalogue = [];
let backendOnline = true;
let playerStatus = null;
let favouriteIds = new Set(JSON.parse(localStorage.getItem("favourites") || "[]"));

const state = document.getElementById("state");
const statusDetail = document.getElementById("statusDetail");
const volume = document.getElementById("volume");
const stream = document.getElementById("url");
const streamInput = document.getElementById("stream");
const streamSearch = document.getElementById("streamSearch");
const streamCount = document.getElementById("streamCount");
const favouritesSection = document.getElementById("favouritesSection");
const favourites = document.getElementById("favourites");
const favouriteButton = document.getElementById("favouriteButton");
const volumeSlider = document.getElementById("volumeSlider");
const volumeValue = document.getElementById("volumeValue");
const playButton = document.getElementById("play");
const stopButton = document.getElementById("stop");
const updateCatalogueButton = document.getElementById("updateCatalogueButton");
const autoplay = document.getElementById("autoplay");

async function getStatus() {
    const response = await fetch("/api/player/status");
    if (!response.ok) throw new Error("Unable to get player status");
    return response.json();
}

function saveFavourites() {
    localStorage.setItem("favourites", JSON.stringify([...favouriteIds]));
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

    const favouriteStreams = catalogue.filter(item =>
        favouriteIds.has(item.id) && streamMatchesQuery(item, query)
    );

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

    streamCount.textContent = query
        ? `${filteredCatalogue.length} of ${catalogue.length} masjids`
        : `${catalogue.length} masjids`;

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

    const preferred = currentSelection || localStorage.getItem("lastStream");
    renderStreams(preferred);
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
    volumeSlider.disabled = offline;
    updateCatalogueButton.disabled = offline;
}

function updateFavouriteButton() {
    const selectedId = streamInput.value;
    const selected = catalogue.find(item => item.id === selectedId);
    const isFavourite = selected && favouriteIds.has(selected.id);

    favouriteButton.disabled = !backendOnline || !selected;
    favouriteButton.textContent = isFavourite
        ? "★ Remove from Favourites"
        : "☆ Add to Favourites";
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
    volumeSlider.disabled = false;
    updateCatalogueButton.disabled = false;
    updateFavouriteButton();
}

function updateStatusDetail(status) {
    const detail = status.error || "";
    statusDetail.textContent = detail;
    statusDetail.className = detail
        ? "status-detail status-detail-" + status.state
        : "status-detail hidden";
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
        updateControls(status);

        if (!status.url) {
            stream.textContent = "No stream playing";
        } else {
            const current = findStreamByURL(status.url);
            stream.textContent = current
                ? current.name + (current.location ? " — " + current.location : "") + "\n" + current.url
                : status.url;
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

streamInput.addEventListener("change", () => {
    localStorage.setItem("lastStream", streamInput.value);
    updateFavouriteButton();
    if (playerStatus) updateControls(playerStatus);
});

favourites.addEventListener("click", event => {
    const item = event.target.closest(".favourite-item");
    if (!item) return;

    const id = item.dataset.id;
    if (event.target.closest(".favourite-remove")) {
        favouriteIds.delete(id);
        saveFavourites();
        renderStreams(streamInput.value);
        if (playerStatus) updateControls(playerStatus);
        return;
    }

    if (catalogue.some(stream => stream.id === id)) {
        streamInput.value = id;
        localStorage.setItem("lastStream", id);
        updateFavouriteButton();
        if (playerStatus) updateControls(playerStatus);
    }
});

favouriteButton.addEventListener("click", () => {
    const id = streamInput.value;
    if (!id) return;

    if (favouriteIds.has(id)) {
        favouriteIds.delete(id);
        showToast("Removed from favourites.", "success");
    } else {
        favouriteIds.add(id);
        showToast("Added to favourites.", "success");
    }

    saveFavourites();
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
    volumeValue.style.color = value > 100 ? "#ffb347" : "#32c36c";

    try {
        await setVolume(value);
        await refreshStatus();
    } catch (err) {
        console.error(err);
    }
});

autoplay.addEventListener("change", () => {
    localStorage.setItem("autoplay", autoplay.checked);
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
    await loadStreams();
    await refreshStatus();

    autoplay.checked = localStorage.getItem("autoplay") === "true";
    if (autoplay.checked && streamInput.value) {
        try {
            await playStream(streamInput.value);
            await refreshStatus();
        } catch (err) {
            console.error(err);
        }
    }

    setInterval(refreshStatus, 1000);
}

initialize();
