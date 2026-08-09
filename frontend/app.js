let catalogue = [];
let backendOnline = true;

const state = document.getElementById("state");
const statusDetail = document.getElementById("statusDetail");
const volume = document.getElementById("volume");
const stream = document.getElementById("url");
const streamInput = document.getElementById("stream");
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

async function loadStreams() {
    const currentSelection = streamInput.value;
    const response = await fetch("/api/streams");
    if (!response.ok) {
        console.error("Unable to load streams");
        return;
    }

    catalogue = await response.json();
    streamInput.innerHTML = "";

    for (const stream of catalogue) {
        const option = document.createElement("option");
        option.value = stream.id;
        option.textContent = stream.name;
        streamInput.appendChild(option);
    }

    const preferred = currentSelection || localStorage.getItem("lastStream");
    if (preferred && [...streamInput.options].some(option => option.value === preferred)) {
        streamInput.value = preferred;
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
    volumeSlider.disabled = offline;
    updateCatalogueButton.disabled = offline;
}

function updateControls(status) {
    if (!backendOnline) {
        setOffline(true);
        return;
    }

    const active = ["waiting", "connecting", "playing", "retrying"].includes(status.state);
    playButton.disabled = active;
    stopButton.disabled = !active;
    streamInput.disabled = false;
    volumeSlider.disabled = false;
    updateCatalogueButton.disabled = false;
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
            stream.textContent = current ? current.name + "\n" + current.url : status.url;
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

streamInput.addEventListener("change", () => {
    localStorage.setItem("lastStream", streamInput.value);
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