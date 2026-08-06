let catalogue = [];

// ---------- DOM ----------

const state = document.getElementById("state");
const volume = document.getElementById("volume");
const stream = document.getElementById("url");

const streamInput = document.getElementById("stream");

const volumeSlider = document.getElementById("volumeSlider");
const volumeValue = document.getElementById("volumeValue");

const playButton = document.getElementById("play");
const stopButton = document.getElementById("stop");
const updateCatalogueButton = document.getElementById("updateCatalogueButton");
const autoplay = document.getElementById("autoplay");

// ---------- API ----------

async function getStatus() {
    const response = await fetch("/api/player/status");

    if (!response.ok) {
        throw new Error("Unable to get player status");
    }

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

    const preferred =
        currentSelection ||
        localStorage.getItem("lastStream");

    if (preferred) {

        const exists = [...streamInput.options].some(
            option => option.value === preferred
        );

        if (exists) {
            streamInput.value = preferred;
        }
    }
}

async function playStream(id) {
    const response = await fetch("/api/player/play", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            id: id
        })
    });

    if (!response.ok) {
        throw new Error("Unable to play stream");
    }
}

async function stopStream() {
    const response = await fetch("/api/player/stop", {
        method: "POST"
    });

    if (!response.ok) {
        throw new Error("Unable to stop playback");
    }
}

async function setVolume(level) {
    const response = await fetch("/api/player/volume", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            volume: level
        })
    });

    if (!response.ok) {
        throw new Error("Unable to change volume");
    }
}

async function updateCatalogue() {

    const response = await fetch(
        "/api/catalogue/update",
        {
            method: "POST"
        }
    );

    if (!response.ok) {
        throw new Error("Unable to update catalogue");
    }

    return response.json();
}

async function loadVersion() {

    const response = await fetch("/api/version");

    if (!response.ok) {
        return;
    }

    const version = await response.json();

    document.getElementById("version").textContent =
        version.version;
}

// ---------- Helpers ----------

function findStreamByURL(url) {
    return catalogue.find(stream => stream.url === url);
}

// ---------- UI ----------

async function refreshStatus() {

    try {

        const status = await getStatus();

        state.textContent = status.state.toUpperCase();

        state.className =
            status.state === "playing"
                ? "status-playing"
                : "status-stopped";

        volume.textContent = status.volume + "%";

        volumeSlider.value = status.volume;
        volumeValue.textContent = status.volume + "%";

        if (!status.url) {

            stream.textContent = "No stream playing";

        } else {

            const current = findStreamByURL(status.url);

            if (current) {

                stream.textContent =
                    current.name + "\n" +
                    current.url;

            } else {

                stream.textContent = status.url;

            }
        }

    } catch (err) {

        console.error(err);

    }
}

// ---------- Events ----------

streamInput.addEventListener("change", () => {

    localStorage.setItem(
        "lastStream",
        streamInput.value
    );

});

playButton.addEventListener("click", async () => {

    if (!streamInput.value) {
        alert("Please select a masjid");
        return;
    }

    try {

        await playStream(streamInput.value);

        await refreshStatus();

    } catch (err) {

        alert(err.message);

    }

});

stopButton.addEventListener("click", async () => {

    try {

        await stopStream();

        await refreshStatus();

    } catch (err) {

        alert(err.message);

    }

});

volumeSlider.addEventListener("input", async () => {

    const value = Number(volumeSlider.value);

    volumeValue.textContent = value + "%";

    if (value > 100) {
        volumeValue.style.color = "#ffb347";
    } else {
        volumeValue.style.color = "#32c36c";
    }

    try {

        await setVolume(value);

        await refreshStatus();

    } catch (err) {

        console.error(err);

    }

});

autoplay.addEventListener("change", () => {

    localStorage.setItem(
        "autoplay",
        autoplay.checked
    );

});

updateCatalogueButton.addEventListener("click", async () => {

    try {

        updateCatalogueButton.disabled = true;
        updateCatalogueButton.textContent = "Updating...";

        await updateCatalogue();
        await loadStreams();
        await refreshStatus();

        alert("Catalogue updated successfully.");

    } catch (err) {

        alert(err.message);

    } finally {

        updateCatalogueButton.disabled = false;
        updateCatalogueButton.textContent = "🔄 Update Catalogue";

    }

});

// ---------- Startup ----------

async function initialize() {

    await loadStreams();

    await refreshStatus();

    await loadVersion();

    autoplay.checked =
        localStorage.getItem("autoplay") === "true";

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