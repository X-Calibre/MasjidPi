"use strict";

(() => {
    const networkStep = document.getElementById("networkStep");
    const passwordStep = document.getElementById("passwordStep");
    const successStep = document.getElementById("successStep");
    const networkList = document.getElementById("networkList");
    const networkStatus = document.getElementById("networkStatus");
    const selectedNetwork = document.getElementById("selectedNetwork");
    const password = document.getElementById("wifiPassword");
    const connectStatus = document.getElementById("connectStatus");
    const keyboard = document.getElementById("keyboard");
    const keyboardRows = document.getElementById("keyboardRows");
    const connectButton = document.getElementById("connectButton");
    let currentNetwork = null;
    let shifted = false;
    let symbols = false;

    const letterRows = [
        ["q", "w", "e", "r", "t", "y", "u", "i", "o", "p"],
        ["a", "s", "d", "f", "g", "h", "j", "k", "l"],
        ["shift", "z", "x", "c", "v", "b", "n", "m", "backspace"],
        ["symbols", "-", "_", "space", ".", "@", "done"]
    ];
    const symbolRows = [
        ["1", "2", "3", "4", "5", "6", "7", "8", "9", "0"],
        ["!", "@", "#", "$", "%", "^", "&", "*", "(", ")"],
        ["letters", "+", "=", "[", "]", "{", "}", "?", "backspace"],
        ["\\", "|", ";", "space", ":", "'", "\"", "/", "done"]
    ];

    function keyLabel(key) {
        const labels = {shift: "⇧", backspace: "⌫", symbols: "123", letters: "ABC", space: "Space", done: "Done"};
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

    async function scanNetworks() {
        networkStatus.textContent = "Looking for nearby networks…";
        networkList.replaceChildren();
        try {
            const response = await fetch("/api/setup/wifi/networks", {cache: "no-store"});
            const payload = await response.json();
            if (!response.ok) throw new Error(payload.error || "Could not scan for Wi-Fi networks");
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
            const response = await fetch("/api/setup/wifi/connect", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({ssid: currentNetwork.ssid, password: password.value})
            });
            const payload = await response.json();
            if (!response.ok) throw new Error(payload.error || "Could not connect to Wi-Fi");
            password.value = "";
            passwordStep.hidden = true;
            successStep.hidden = false;
            document.getElementById("successNetwork").textContent = `Connected to ${currentNetwork.ssid}`;
        } catch (error) {
            connectStatus.textContent = error.message;
            setKeyboardOpen(true);
        } finally {
            connectButton.disabled = false;
            connectButton.textContent = "Connect";
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
    document.getElementById("continueButton").addEventListener("click", () => {
        window.location.replace("/masjidboard-config.html");
    });

    scanNetworks();
})();
