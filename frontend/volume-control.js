(() => {
    // app.js historically handled every slider input as a persistent volume
    // change. Capture the event before that handler so slider movement changes
    // live hardware volume without rewriting persistent state.
    const slider = document.getElementById("volumeSlider");
    const valueLabel = document.getElementById("volumeValue");
    if (!slider) return;

    let timer = null;
    let requestSerial = 0;

    async function sendVolume(value, persist) {
        const response = await fetch("/api/player/volume", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ volume: value, persist })
        });
        if (!response.ok) throw new Error("Unable to change volume");
        return response.json();
    }

    function scheduleLiveVolume(value) {
        requestSerial += 1;
        const serial = requestSerial;
        clearTimeout(timer);
        timer = setTimeout(async () => {
            try {
                const status = await sendVolume(value, false);
                if (serial === requestSerial && window.playerStatus !== undefined) {
                    window.playerStatus = status;
                }
            } catch (err) {
                console.error(err);
            }
        }, 60);
    }

    slider.addEventListener("input", event => {
        event.stopImmediatePropagation();
        const value = Number(slider.value);
        valueLabel.textContent = value + "%";
        scheduleLiveVolume(value);
    }, true);

    slider.addEventListener("change", async event => {
        event.stopImmediatePropagation();
        const value = Number(slider.value);
        clearTimeout(timer);
        requestSerial += 1;
        try {
            const status = await sendVolume(value, true);
            if (window.playerStatus !== undefined) {
                window.playerStatus = status;
            }
        } catch (err) {
            console.error(err);
        }
    }, true);
})();
