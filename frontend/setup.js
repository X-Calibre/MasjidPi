"use strict";

(() => {
    const boardURL = "/masjidboard.html?profile=appliance-720";
    const visualDemoStep = new URLSearchParams(window.location.search).get("visual-demo-step");
    document.body.classList.add("setup-720-layout");
    const networkStep = document.getElementById("networkStep");
    const passwordStep = document.getElementById("passwordStep");
    const successStep = document.getElementById("successStep");
    const locationStep = document.getElementById("locationStep");
    const masjidStep = document.getElementById("masjidStep");
    const networkList = document.getElementById("networkList");
    const networkStatus = document.getElementById("networkStatus");
    const selectedNetwork = document.getElementById("selectedNetwork");
    const ssid = document.getElementById("wifiSSID");
    const password = document.getElementById("wifiPassword");
    const hiddenNetworkFields = document.getElementById("hiddenNetworkFields");
    const passwordFields = document.getElementById("passwordFields");
    const connectStatus = document.getElementById("connectStatus");
    const keyboard = document.getElementById("keyboard");
    const keyboardRows = document.getElementById("keyboardRows");
    const connectButton = document.getElementById("connectButton");
    const countryButton = document.getElementById("countryButton");
    const regionButton = document.getElementById("regionButton");
    const cityButton = document.getElementById("cityButton");
    const findMasjidsButton = document.getElementById("findMasjidsButton");
    const masjidList = document.getElementById("masjidList");
    const finishSetupButton = document.getElementById("finishSetupButton");
    const pickerSheet = document.getElementById("pickerSheet");
    const pickerOptions = document.getElementById("pickerOptions");
    let currentNetwork = null;
    let hierarchy = null;
    let selectedMasjid = null;
    let selectedCountryName = "";
    let selectedRegionIndex = -1;
    let selectedCityName = "";
    let continueAction = () => window.location.replace(boardURL);
    let shifted = false;
    let symbols = false;
    let activeKeyboardInput = password;

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
            activeKeyboardInput.value = Array.from(activeKeyboardInput.value).slice(0, -1).join("");
            return;
        }
        if (key === "done") {
            setKeyboardOpen(false);
            return;
        }
        if (activeKeyboardInput.value.length >= activeKeyboardInput.maxLength) return;
        const value = key === "space" ? " " : (shifted ? key.toUpperCase() : key);
        activeKeyboardInput.value += value;
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
        currentNetwork = {...network, hidden: false};
        selectedNetwork.textContent = network.ssid;
        hiddenNetworkFields.hidden = true;
        passwordFields.hidden = network.security === "Open";
        document.getElementById("passwordHeading").textContent = network.security === "Open" ? "Connect to Wi-Fi" : "Enter Wi-Fi password";
        password.value = "";
        connectStatus.textContent = "";
        networkStep.hidden = true;
        passwordStep.hidden = false;
        activeKeyboardInput = password;
        setKeyboardOpen(network.security !== "Open");
    }

    function addHiddenNetwork() {
        currentNetwork = {ssid: "", security: "WPA2", hidden: true};
        ssid.value = "";
        password.value = "";
        selectedNetwork.textContent = "Hidden network";
        document.getElementById("passwordHeading").textContent = "Add hidden Wi-Fi network";
        hiddenNetworkFields.hidden = false;
        passwordFields.hidden = false;
        document.getElementById("hiddenSecurity").textContent = "Password protected";
        document.getElementById("hiddenSecurity").setAttribute("aria-pressed", "true");
        connectStatus.textContent = "";
        networkStep.hidden = true;
        passwordStep.hidden = false;
        activeKeyboardInput = ssid;
        setKeyboardOpen(true);
    }

    async function connect() {
        if (!currentNetwork) return;
        if (currentNetwork.hidden) currentNetwork.ssid = ssid.value.trim();
        if (!currentNetwork.ssid) {
            connectStatus.textContent = "Enter the hidden network name.";
            activeKeyboardInput = ssid;
            setKeyboardOpen(true);
            return;
        }
        setKeyboardOpen(false);
        connectButton.disabled = true;
        connectButton.textContent = "Connecting…";
        connectStatus.textContent = "This can take up to 20 seconds.";
        try {
            await jsonRequest("/api/setup/wifi/connect", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({ssid: currentNetwork.ssid, password: passwordFields.hidden ? "" : password.value, hidden: currentNetwork.hidden})
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
        await loadDeviceAccess();
        if (configured) {
            document.getElementById("continueButton").textContent = "Start MasjidFrame";
            continueAction = () => window.location.replace(boardURL);
        } else {
            document.getElementById("continueButton").textContent = "Choose your location";
            continueAction = showLocationStep;
        }
    }

    async function loadDeviceAccess() {
        const advancedSetup = document.getElementById("advancedSetup");
        const fqdnAccess = document.getElementById("fqdnAccess");
        const ipAccess = document.getElementById("ipAccess");
        fqdnAccess.hidden = true;
        ipAccess.hidden = true;
        try {
            const access = await jsonRequest("/api/setup/device-access");
            const port = window.location.port || "8080";
            if (access?.fqdn) {
                document.getElementById("fqdnURL").textContent = `http://${access.fqdn}:${port}`;
                fqdnAccess.hidden = false;
            }
            if (access?.ip_address) {
                document.getElementById("ipURL").textContent = `http://${access.ip_address}:${port}`;
                ipAccess.hidden = false;
            }
            advancedSetup.hidden = fqdnAccess.hidden && ipAccess.hidden;
        } catch (_) {
            advancedSetup.hidden = true;
        }
    }

    function countries() {
        return Array.isArray(hierarchy?.countries) ? hierarchy.countries : [];
    }

    function selectedCountry() {
        return countries().find((country) => country.name === selectedCountryName) || null;
    }

    function selectedRegion() {
        return selectedRegionIndex >= 0 ? (selectedCountry()?.regions || [])[selectedRegionIndex] || null : null;
    }

    function closePicker() {
        pickerSheet.hidden = true;
        pickerOptions.replaceChildren();
    }

    function openPicker(title, items, selectedValue, choose) {
        document.getElementById("pickerHeading").textContent = title;
        pickerOptions.replaceChildren();
        items.forEach((item) => {
            const button = document.createElement("button");
            button.type = "button";
            button.className = "picker-option";
            button.textContent = item.label;
            button.classList.toggle("selected", item.value === selectedValue);
            button.addEventListener("click", () => {
                choose(item);
                closePicker();
            });
            pickerOptions.append(button);
        });
        pickerSheet.hidden = false;
    }

    function populateCountries() {
        countryButton.disabled = countries().length === 0;
        const southAfrica = countries().find((country) => country.name === "South Africa");
        if (southAfrica) {
            selectedCountryName = southAfrica.name;
            countryButton.textContent = southAfrica.name;
            populateRegions();
        }
    }

    function populateRegions() {
        selectedRegionIndex = -1;
        selectedCityName = "";
        regionButton.textContent = "Select province or region…";
        cityButton.textContent = "Select town or city…";
        cityButton.disabled = true;
        findMasjidsButton.disabled = true;
        const country = selectedCountry();
        regionButton.disabled = !country;
    }

    function populateCities() {
        selectedCityName = "";
        cityButton.textContent = "Select town or city…";
        findMasjidsButton.disabled = true;
        const region = selectedRegion();
        cityButton.disabled = !region;
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
        document.body.classList.remove("masjid-step-open");
        successStep.hidden = true;
        masjidStep.hidden = true;
        locationStep.hidden = false;
        loadHierarchy().catch((error) => {
            document.getElementById("locationStatus").textContent = `Could not load locations: ${error.message}`;
        });
    }

    function locationValue() {
        return {country: selectedCountryName, region: selectedRegion()?.name || "", city: selectedCityName};
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
            document.body.classList.add("masjid-step-open");
            document.getElementById("masjidLocation").textContent = [location.city, location.region, location.country].filter(Boolean).join(", ");
            document.getElementById("masjidStatus").textContent = records.length
                ? `${records.length} MasjidBoard${records.length === 1 ? "" : "s"} found`
                : "No MasjidBoards were found for this location. Choose another location.";
            renderMasjids(records);
        } catch (error) {
            document.getElementById("locationStatus").textContent = `Could not find masjids: ${error.message}`;
        } finally {
            findMasjidsButton.disabled = !selectedCityName;
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
            document.body.classList.remove("masjid-step-open");
            masjidStep.hidden = true;
            successStep.hidden = false;
            document.getElementById("successHeading").textContent = "MasjidFrame is ready";
            document.getElementById("successNetwork").textContent = selectedMasjid.name;
            document.getElementById("continueButton").textContent = "Start MasjidFrame";
            continueAction = () => window.location.replace(boardURL);
            await loadDeviceAccess();
        } catch (error) {
            document.getElementById("masjidStatus").textContent = `Could not save this masjid: ${error.message}`;
            finishSetupButton.disabled = false;
            finishSetupButton.textContent = "Use selected masjid";
        }
    }

    function showVisualDemoStep(step) {
        document.body.classList.add("visual-demo");
        document.body.classList.remove("masjid-step-open");
        for (const section of [networkStep, passwordStep, successStep, locationStep, masjidStep]) section.hidden = true;
        pickerSheet.hidden = true;
        setKeyboardOpen(false);

        if (step === "password") {
            chooseNetwork({ssid:"MasjidFrame Demo Wi-Fi", security:"WPA2"});
            password.value = "visual-demo-password";
            return;
        }
        if (step === "hidden") {
            addHiddenNetwork();
            ssid.value = "Community Wi-Fi";
            password.value = "visual-demo-password";
            return;
        }
        if (step === "success") {
            successStep.hidden = false;
            document.getElementById("successHeading").textContent = "MasjidFrame is online";
            document.getElementById("successNetwork").textContent = "Connected to MasjidFrame Demo Wi-Fi";
            document.getElementById("fqdnURL").textContent = "http://masjidframe.example.test:8080";
            document.getElementById("ipURL").textContent = "http://192.0.2.25:8080";
            document.getElementById("fqdnAccess").hidden = false;
            document.getElementById("ipAccess").hidden = false;
            document.getElementById("advancedSetup").hidden = false;
            return;
        }
        if (step === "location" || step === "picker") {
            locationStep.hidden = false;
            document.getElementById("locationStatus").textContent = "Select the location nearest to your masjid.";
            countryButton.disabled = false;
            regionButton.disabled = false;
            cityButton.disabled = false;
            countryButton.textContent = "South Africa";
            regionButton.textContent = "Gauteng";
            cityButton.textContent = "Laudium";
            findMasjidsButton.disabled = false;
            if (step === "picker") {
                document.getElementById("pickerHeading").textContent = "Choose town or city";
                pickerOptions.replaceChildren();
                for (const name of ["Laudium", "Erasmia", "Centurion", "Pretoria"]) {
                    const button = document.createElement("button");
                    button.type = "button";
                    button.className = "picker-option";
                    button.textContent = name;
                    pickerOptions.append(button);
                }
                pickerSheet.hidden = false;
            }
            return;
        }
        if (step === "masjid") {
            masjidStep.hidden = false;
            document.body.classList.add("masjid-step-open");
            document.getElementById("masjidLocation").textContent = "Laudium, Gauteng, South Africa";
            document.getElementById("masjidStatus").textContent = "3 MasjidBoards found";
            const records = [
                {id:"demo-salaam", name:"Masjid us Salaam", city:"Laudium", region:"Gauteng", country:"South Africa"},
                {id:"demo-noor", name:"Masjid Al Noor", city:"Laudium", region:"Gauteng", country:"South Africa"},
                {id:"demo-hamza", name:"Masjid Hamza RA", city:"Erasmia", region:"Gauteng", country:"South Africa"}
            ];
            renderMasjids(records);
            const first = masjidList.querySelector(".masjid-option");
            first?.classList.add("selected");
            if (first) first.querySelector("input").checked = true;
            finishSetupButton.disabled = false;
            return;
        }

        networkStep.hidden = false;
        networkStatus.textContent = "3 networks found";
        networkList.replaceChildren();
        for (const network of [
            {ssid:"MasjidFrame Demo Wi-Fi", meta:"Secured · WPA2", signal:"Excellent"},
            {ssid:"Community Hall", meta:"Secured · WPA2", signal:"Good"},
            {ssid:"Guest Network", meta:"Open network", signal:"Fair"}
        ]) {
            const button = document.createElement("button");
            button.type = "button";
            button.className = "network-button";
            const details = document.createElement("span");
            const name = document.createElement("span");
            name.className = "network-name";
            name.textContent = network.ssid;
            const meta = document.createElement("span");
            meta.className = "network-meta";
            meta.textContent = network.meta;
            details.append(name, document.createElement("br"), meta);
            const signal = document.createElement("span");
            signal.className = "signal";
            signal.textContent = network.signal;
            button.append(details, signal);
            networkList.append(button);
        }
    }

    keyboardRows.addEventListener("click", (event) => {
        const key = event.target.closest("[data-key]");
        if (key) pressKey(key.dataset.key);
    });
    ssid.addEventListener("click", () => {
        activeKeyboardInput = ssid;
        setKeyboardOpen(true);
    });
    password.addEventListener("click", () => {
        activeKeyboardInput = password;
        setKeyboardOpen(true);
    });
    connectButton.addEventListener("click", connect);
    document.getElementById("addHiddenNetwork").addEventListener("click", addHiddenNetwork);
    document.getElementById("refreshNetworks").addEventListener("click", scanNetworks);
    document.getElementById("backToNetworks").addEventListener("click", () => {
        password.value = "";
        ssid.value = "";
        setKeyboardOpen(false);
        passwordStep.hidden = true;
        networkStep.hidden = false;
    });
    document.getElementById("hiddenSecurity").addEventListener("click", (event) => {
        const secured = event.currentTarget.getAttribute("aria-pressed") !== "true";
        event.currentTarget.setAttribute("aria-pressed", String(secured));
        event.currentTarget.textContent = secured ? "Password protected" : "Open network";
        passwordFields.hidden = !secured;
        currentNetwork.security = secured ? "WPA2" : "Open";
        activeKeyboardInput = ssid;
    });
    document.getElementById("togglePassword").addEventListener("click", (event) => {
        const showing = password.type === "text";
        password.type = showing ? "password" : "text";
        event.currentTarget.textContent = showing ? "Show" : "Hide";
        event.currentTarget.setAttribute("aria-pressed", String(!showing));
    });
    countryButton.addEventListener("click", () => openPicker(
        "Choose country",
        countries().map((country) => ({label: country.name, value: country.name})),
        selectedCountryName,
        (item) => {
            selectedCountryName = item.value;
            countryButton.textContent = item.label;
            populateRegions();
        }
    ));
    regionButton.addEventListener("click", () => openPicker(
        "Choose province or region",
        (selectedCountry()?.regions || []).map((region, index) => ({label: region.name || "Other areas", value: index})),
        selectedRegionIndex,
        (item) => {
            selectedRegionIndex = item.value;
            regionButton.textContent = item.label;
            populateCities();
        }
    ));
    cityButton.addEventListener("click", () => openPicker(
        "Choose town or city",
        (selectedRegion()?.cities || []).map((city) => ({label: city.name, value: city.name})),
        selectedCityName,
        (item) => {
            selectedCityName = item.value;
            cityButton.textContent = item.label;
            findMasjidsButton.disabled = false;
        }
    ));
    document.getElementById("closePicker").addEventListener("click", closePicker);
    findMasjidsButton.addEventListener("click", findMasjids);
    finishSetupButton.addEventListener("click", finishSetup);
    document.getElementById("backToLocation").addEventListener("click", showLocationStep);
    document.getElementById("continueButton").addEventListener("click", () => continueAction());

    const setupParams = new URLSearchParams(window.location.search);
    if (setupParams.get("return") === "board") {
        const returnToBoard = document.getElementById("returnToBoard");
        returnToBoard.href = boardURL;
        returnToBoard.hidden = false;
        document.getElementById("networkHeading").textContent = "Change Wi-Fi network";
    }
    const requestedStep = setupParams.get("step");
    if (visualDemoStep) {
        showVisualDemoStep(visualDemoStep);
        document.addEventListener("click", event => {
            event.preventDefault();
            event.stopImmediatePropagation();
        }, true);
    } else if (requestedStep === "location") showLocationStep();
    else scanNetworks();
})();
