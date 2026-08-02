async function refreshStatus() {
    const response = await fetch("/api/player/status");

    if (!response.ok) {
        console.error("Unable to get player status");
        return;
    }

    const status = await response.json();

    // ----- State badge -----

    const state = document.getElementById("state");

    state.textContent = status.state.toUpperCase();

    state.className =
        status.state === "playing"
            ? "status-playing"
            : "status-stopped";

    // ----- Volume -----

    document.getElementById("volume").textContent =
        status.volume + "%";

    document.getElementById("volumeSlider").value =
        status.volume;

    document.getElementById("volumeValue").textContent =
        status.volume + "%";

    // ----- Current Stream -----

    const stream = document.getElementById("url");

    if (!status.url) {
        stream.textContent = "No stream playing";
    } else if (status.url.includes("activetakbeer")) {
        stream.textContent = "🕌 Active Takbeer";
    } else {
        stream.textContent = status.url;
    }
}

refreshStatus();

setInterval(refreshStatus, 1000);

document.getElementById("play").addEventListener("click", async () => {
    const url = document.getElementById("stream").value;

    if (!url) {
        alert("Please enter a stream URL");
        return;
    }

    const response = await fetch("/api/player/play", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            url: url
        })
    });

    if (!response.ok) {
        alert("Unable to play stream");
        return;
    }

document.getElementById("stream").value = "";

await refreshStatus();
});

document.getElementById("stop").addEventListener("click", async () => {
    const response = await fetch("/api/player/stop", {
        method: "POST"
    });

    if (!response.ok) {
        alert("Unable to stop playback");
        return;
    }

    await refreshStatus();
});

document.getElementById("volumeSlider").addEventListener("input", async (event) => {
    const volume = Number(event.target.value);

    document.getElementById("volumeValue").textContent = volume + "%";

    const response = await fetch("/api/player/volume", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            volume: volume
        })
    });

    if (!response.ok) {
        console.error("Unable to set volume");
        return;
    }

    await refreshStatus();
});