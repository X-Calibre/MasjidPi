(() => {
    const controls = [
        {
            slider: document.getElementById("masjidVolumeSlider"),
            value: document.getElementById("masjidVolumeValue")
        },
        {
            slider: document.getElementById("radioVolumeSlider"),
            value: document.getElementById("radioVolumeValue")
        }
    ].filter(control => control.slider && control.value);

    function syncControl(control) {
        const boosted = Number(control.slider.value) > 100;
        control.slider.classList.toggle("volume-boost", boosted);
        control.value.classList.toggle("volume-boost", boosted);
    }

    function syncAll() {
        controls.forEach(syncControl);
    }

    controls.forEach(control => {
        control.slider.addEventListener("input", () => syncControl(control));
        control.slider.addEventListener("change", () => syncControl(control));
    });

    // app.js refreshes slider values from backend status once per second.
    // Keep the visual boost state aligned with those programmatic updates too.
    syncAll();
    setInterval(syncAll, 500);
})();
