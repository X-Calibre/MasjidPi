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

    syncAll();
    window.addEventListener("masjidpi:listen-status", syncAll);
})();
