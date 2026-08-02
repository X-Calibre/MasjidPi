async function refreshStatus() {
    const response = await fetch("/api/player/status");

    if (!response.ok) {
        console.error("Unable to get player status");
        return;
    }

    const status = await response.json();

    document.getElementById("state").textContent = status.state;
    document.getElementById("volume").textContent = status.volume;
    document.getElementById("url").textContent = status.url;
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

    refreshStatus();
});

document.getElementById("stop").addEventListener("click", async () => {

    const response = await fetch("/api/player/stop", {
        method: "POST"
    });

    if (!response.ok) {
        alert("Unable to stop playback");
        return;
    }

    refreshStatus();
});