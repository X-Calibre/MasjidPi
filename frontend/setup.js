"use strict";

(() => {
    const networkStep = document.getElementById("networkStep");
    const passwordStep = document.getElementById("passwordStep");
    const successStep = document.getElementById("successStep");
    const locationStep = document.getElementById("locationStep");
    const masjidStep = document.getElementById("masjidStep");
    const networkList = document.getElementById("networkList");
    const networkStatus = document.getElementById("networkStatus");
    const selectedNetwork = document.getElementById("selectedNetwork");
    const password = document.getElementById("wifiPassword");
    const connectStatus = document.getElementById("connectStatus");
    const keyboard = document.getElementById("keyboard");
    const keyboardRows = document.getElementById("keyboardRows");
    const connectButton = document.getElementById("connectButton");
    const countrySelect = document.getElementById("countrySelect");
    const regionSelect = document.getElementById("regionSelect");
    const citySelect = document.getElementById("citySelect");
    const findMasjidsButton = document.getElementById("findMasjidsButton");
    const masjidList = document.getElementById("masjidList");
    const finishSetupButton = document.getElementById("finishSetupButton");
    let currentNetwork = null;
    let hierarchy = null;
    let selectedMasjid = null;
    let continueAction = () => window.location.replace("/masjidboard.html?profile=appliance");
    let shifted = false;
    let symbols = false;

    const letterRows = [
        ["1", "2", "3", "4", "5", "6", "7", "8", "9", "0"],
        ["q", "w", "e", "r", "t", "y", "u", "i", "o", "p"],
        ["a", "s", "d", "f", "g", "h", "j", "k", "l"],
        ["shift", "z", "x", "c", "v", "b", "n", "m", "backspace"],
        ["symbols", "-", "_", "space", ".", "@", "done"]
    ];
    const symbolRows = [
        ["~", "`", "!", "@", "#", "$", "%", "^", "&", "*"],
        ["(", ")", "-", "+", "=", "[", "]", "{", "}"],
        ["<", ">", "_", ":", ";", "'", "\"", "?", "backspace"],
        ["letters", "\\", "|", ",", "space", ".", "/", "done"]
    ];

    function keyLabel(key) {
        const labels = {shift: "⇧", backspace: "⌫", symbols: "#+=", letters: "ABC", space: "Space", done: "Done"};
        if (labels[key]) return labels[key];
        return shifted && /^[a-z]$/.test(key) ? key.toUpperCase() : key;
    }

    function renderKeyboard() {
        keyboardRows.replaceChildren();
        (symbols ? symbolRows : letterRows).forEach((row) => {
            const rowElement = document.createElement("div");
            rowElement.className = "keyboard-row";
            row.forEach((key) => {
                const button = document.createElement("button");
                button.type = "button";
                button.className = "key";
                if (["shift", "backspace", "symbols", "letters", "done"].includes(key)) button.classList.add("key-wide");
                if (key === "space") button.classList.add("key-space");
                if ((key === "shift" && shifted) || key === "done") button.classList.add("key-accent");
                button.dataset.key = key;
                button.textContent = keyLabel(key);
                button.setAttribute("aria-label", keyLabel(key));
                rowElement.appendChild(button);
            });
            keyboardRows.appendChild(rowElement);
        });
    }

    function setKeyboardOpen(open) {
        keyboard.hidden = !open;
        document.body.classList.toggle("keyboard-open", open);
        if (open) renderKeyboard();
    }

    function pressKey(key) {
        if (key === "shift") {
            shifted = !shifted;
            renderKeyboard();
            return;
        }
        if (key === "symbols" || key === "letters") {
            symbols = key === "symbols";
            shifted = false;
            renderKeyboard();
            return;
        }
        if (key === "backspace") {
            password.value = Array.from(password.value).slice(0, -1).join("");
            return;
        }
        if (key === "done") {
            setKeyboardOpen(false);
            return;
        }
        if (password.value.length >= password.maxLength) return;
        const value = key === "space" ? " " : (shifted ? key.toUpperCase() : key);
        password.value += value;
        if (shifted) {
            shifted = false;
            renderKeyboard();
        }
    }

    function signalLabel(value) {
        if (value >= 75) return "Excellent";
        if (value >= 55) return "Good";
        if (value >= 35) return "Fair";
        return "Weak";
    }

    async function jsonRequest(url, options = {}) {
        const response = await fetch(url, {cache: "no-store", ...options});
        let payload = null;
        try {
            payload = await response.json();
        } catch (_) {}
        if (!response.ok) throw new Error(payload?.error || `Request failed (${response.status})`);
        return payload;
    }

    async function scanNetworks() {
        networkStatus.textContent = "Looking for nearby networks…";
        networkList.replaceChildren();
        try {
            const payload = await jsonRequest("/api/setup/wifi/networks");
            if (!payload.networks.length) {
                networkStatus.textContent = "No 2.4 GHz Wi-Fi networks were found. Move MasjidFrame closer to your router and refresh.";
                return;
            }
            networkStatus.textContent = `${payload.networks.length} network${payload.networks.length === 1 ? "" : "s"} found`;
            payload.networks.forEach((network) => {
                const button = document.createElement("button");
                button.type = "button";
                button.className = "network-button";
                button.innerHTML = `<span><span class="network-name"></span><br><span class="network-meta"></span></span><span class="signal"></span>`;
                button.querySelector(".network-name").textContent = network.ssid;
                button.querySelector(".network-meta").textContent = network.security === "Open" ? "Open network" : `Secured · ${network.security}`;
                button.querySelector(".signal").textContent = signalLabel(network.signal);
                button.addEventListener("click", () => chooseNetwork(network));
                networkList.appendChild(button);
            });
        } catch (error) {
            networkStatus.textContent = error.message;
        }
    }

    function chooseNetwork(network) {
        currentNetwork = network;
        selectedNetwork.textContent = network.ssid;
        password.value = "";
        connectStatus.textContent = "";
        networkStep.hidden = true;
        passwordStep.hidden = false;
        setKeyboardOpen(network.security !== "Open");
    }

    async function connect() {
        if (!currentNetwork) return;
        setKeyboardOpen(false);
        connectButton.disabled = true;
        connectButton.textContent = "Connecting…";
        connectStatus.textContent = "This can take up to 20 seconds.";
        try {
            await jsonRequest("/api/setup/wifi/connect", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({ssid: currentNetwork.ssid, password: password.value})
            });
            password.value = "";
            passwordStep.hidden = true;
            await showConnectedStep();
        } catch (error) {
            connectStatus.textContent = error.message;
            setKeyboardOpen(true);
        } finally {
            connectButton.disabled = false;
            connectButton.textContent = "Connect";
        }
    }

    async function showConnectedStep() {
        let configured = false;
        try {
            const selection = await jsonRequest("/api/masjidboard/selection");
            configured = selection?.configured === true;
        } catch (_) {}

        successStep.hidden = false;
        document.getElementById("successHeading").textContent = "MasjidFrame is online";
        document.getElementById("successNetwork").textContent = `Connected to ${currentNetwork.ssid}`;
        if (configured) {
            document.getElementById("continueButton").textContent = "Start MasjidFrame";
            continueAction = () => window.location.replace("/masjidboard.html?profile=appliance");
        } else {
            document.getElementById("continueButton").textContent = "Choose your location";
            continueAction = showLocationStep;
        }
    }

    function option(value, label) {
        const item = document.createElement("option");
        item.value = value;
        item.textContent = label;
        return item;
    }

    function countries() {
        return Array.isArray(hierarchy?.countries) ? hierarchy.countries : [];
    }

    function selectedCountry() {
        return countries().find((country) => country.name === countrySelect.value) || null;
    }

    function selectedRegion() {
        const index = Number.parseInt(regionSelect.value, 10);
        return Number.isInteger(index) ? (selectedCountry()?.regions || [])[index] || null : null;
    }

    function populateCountries() {
        countrySelect.replaceChildren(option("", "Select country…"));
        countries().forEach((country) => countrySelect.append(option(country.name, country.name)));
        countrySelect.disabled = countries().length === 0;
        const southAfrica = countries().find((country) => country.name === "South Africa");
        if (southAfrica) {
            countrySelect.value = southAfrica.name;
            populateRegions();
        }
    }

    function populateRegions() {
        regionSelect.replaceChildren(option("", "Select province or region…"));
        citySelect.replaceChildren(option("", "Select town or city…"));
        citySelect.disabled = true;
        findMasjidsButton.disabled = true;
        const country = selectedCountry();
        regionSelect.disabled = !country;
        (country?.regions || []).forEach((region, index) => {
            const label = region.name || "Other areas";
            regionSelect.append(option(String(index), label));
        });
    }

    function populateCities() {
        citySelect.replaceChildren(option("", "Select town or city…"));
        findMasjidsButton.disabled = true;
        const region = selectedRegion();
        citySelect.disabled = !region;
        (region?.cities || []).forEach((city) => citySelect.append(option(city.name, city.name)));
    }

    async function loadHierarchy() {
        const status = document.getElementById("locationStatus");
        status.textContent = "Loading locations…";
        hierarchy = await jsonRequest("/api/masjidboard/hierarchy");
        if (!countries().length) {
            status.textContent = "Downloading the location list…";
            await jsonRequest("/api/masjidboard/hierarchy/refresh", {method: "POST"});
            hierarchy = await jsonRequest("/api/masjidboard/hierarchy");
        }
        populateCountries();
        status.textContent = countries().length ? "Select the location nearest to your masjid." : "No locations are currently available.";
    }

    function showLocationStep() {
        successStep.hidden = true;
        masjidStep.hidden = true;
        locationStep.hidden = false;
        loadHierarchy().catch((error) => {
            document.getElementById("locationStatus").textContent = `Could not load locations: ${error.message}`;
        });
    }

    function locationValue() {
        return {country: countrySelect.value, region: selectedRegion()?.name || "", city: citySelect.value};
    }

    function renderMasjids(records) {
        masjidList.replaceChildren();
        selectedMasjid = null;
        finishSetupButton.disabled = true;
        records.forEach((record) => {
            const label = document.createElement("label");
            label.className = "masjid-option";
            const radio = document.createElement("input");
            radio.type = "radio";
            radio.name = "masjid";
            radio.value = record.id;
            const details = document.createElement("span");
            const name = document.createElement("strong");
            name.textContent = record.name;
            const location = document.createElement("small");
            location.textContent = [record.city, record.region, record.country].filter(Boolean).join(" · ");
            details.append(name, location);
            label.append(radio, details);
            label.addEventListener("click", () => {
                selectedMasjid = record;
                masjidList.querySelectorAll(".masjid-option").forEach((item) => item.classList.toggle("selected", item === label));
                finishSetupButton.disabled = false;
            });
            masjidList.append(label);
        });
    }

    async function findMasjids() {
        const location = locationValue();
        if (!location.country || !location.city) return;
        findMasjidsButton.disabled = true;
        findMasjidsButton.textContent = "Finding masjids…";
        document.getElementById("locationStatus").textContent = "Saving your location and retrieving MasjidBoards…";
        try {
            await jsonRequest("/api/masjidboard/scope", {
                method: "PUT",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({locations: [location]})
            });
            await jsonRequest("/api/masjidboard/catalogue/refresh", {method: "POST"});
            const catalogue = await jsonRequest("/api/masjidboard/catalogue?status=active");
            const records = Array.isArray(catalogue?.records) ? catalogue.records : [];
            locationStep.hidden = true;
            masjidStep.hidden = false;
            document.getElementById("masjidLocation").textContent = [location.city, location.region, location.country].filter(Boolean).join(", ");
            document.getElementById("masjidStatus").textContent = records.length
                ? `${records.length} MasjidBoard${records.length === 1 ? "" : "s"} found`
                : "No MasjidBoards were found for this location. Choose another location.";
            renderMasjids(records);
        } catch (error) {
            document.getElementById("locationStatus").textContent = `Could not find masjids: ${error.message}`;
        } finally {
            findMasjidsButton.disabled = !citySelect.value;
            findMasjidsButton.textContent = "Find masjids";
        }
    }

    async function finishSetup() {
        if (!selectedMasjid) return;
        finishSetupButton.disabled = true;
        finishSetupButton.textContent = "Setting up MasjidFrame…";
        document.getElementById("masjidStatus").textContent = "Downloading the first timetable…";
        try {
            await jsonRequest("/api/masjidboard/selection", {
                method: "PUT",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({catalogue_ids: [selectedMasjid.id]})
            });
            masjidStep.hidden = true;
            successStep.hidden = false;
            document.getElementById("successHeading").textContent = "MasjidFrame is ready";
            document.getElementById("successNetwork").textContent = selectedMasjid.name;
            document.getElementById("continueButton").textContent = "Start MasjidFrame";
            continueAction = () => window.location.replace("/masjidboard.html?profile=appliance");
        } catch (error) {
            document.getElementById("masjidStatus").textContent = `Could not save this masjid: ${error.message}`;
            finishSetupButton.disabled = false;
            finishSetupButton.textContent = "Use selected masjid";
        }
    }

    keyboardRows.addEventListener("click", (event) => {
        const key = event.target.closest("[data-key]");
        if (key) pressKey(key.dataset.key);
    });
    password.addEventListener("click", () => setKeyboardOpen(true));
    connectButton.addEventListener("click", connect);
    document.getElementById("refreshNetworks").addEventListener("click", scanNetworks);
    document.getElementById("backToNetworks").addEventListener("click", () => {
        password.value = "";
        setKeyboardOpen(false);
        passwordStep.hidden = true;
        networkStep.hidden = false;
    });
    document.getElementById("togglePassword").addEventListener("click", (event) => {
        const showing = password.type === "text";
        password.type = showing ? "password" : "text";
        event.currentTarget.textContent = showing ? "Show" : "Hide";
        event.currentTarget.setAttribute("aria-pressed", String(!showing));
    });
    countrySelect.addEventListener("change", populateRegions);
    regionSelect.addEventListener("change", populateCities);
    citySelect.addEventListener("change", () => { findMasjidsButton.disabled = !citySelect.value; });
    findMasjidsButton.addEventListener("click", findMasjids);
    finishSetupButton.addEventListener("click", finishSetup);
    document.getElementById("backToLocation").addEventListener("click", showLocationStep);
    document.getElementById("continueButton").addEventListener("click", () => continueAction());

    const requestedStep = new URLSearchParams(window.location.search).get("step");
    if (requestedStep === "location") showLocationStep();
    else scanNetworks();
})();
