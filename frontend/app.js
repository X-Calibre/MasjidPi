// ---------- DOM ----------

const state = document.getElementById("state");
const volume = document.getElementById("volume");
const stream = document.getElementById("url");

const streamInput = document.getElementById("stream");

const volumeSlider = document.getElementById("volumeSlider");
const volumeValue = document.getElementById("volumeValue");

const playButton = document.getElementById("play");
const stopButton = document.getElementById("stop");
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
    const response = await fetch("/api/streams");

    if (!response.ok) {
        console.error("Unable to load streams");
        return;
    }

    const streams = await response.json();

    const select = document.getElementById("stream");

    select.innerHTML = "";

    for (const stream of streams) {
        const option = document.createElement("option");

        option.value = stream.id;
        option.textContent = stream.name;

        select.appendChild(option);
    }

    const last = localStorage.getItem("lastStream");

        if (last) {
        select.value = last;
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
        } else if (status.url.includes("activetakbeer")) {
            stream.textContent = "🕌 Active Takbeer";
        } else {
            stream.textContent = status.url;
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

        streamInput.value = "";

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

// ---------- Startup ----------

async function initialize() {

    await loadStreams();

    await refreshStatus();

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