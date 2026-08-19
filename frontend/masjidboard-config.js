(() => {
    "use strict";

    const maxLocations = 3;
    const hierarchyEndpoint = "/api/masjidboard/hierarchy";
    const hierarchyRefreshEndpoint = "/api/masjidboard/hierarchy/refresh";
    const scopeEndpoint = "/api/masjidboard/scope";
    const catalogueRefreshEndpoint = "/api/masjidboard/catalogue/refresh";
    const statusEndpoint = "/api/masjidboard/status";

    const locationConfiguration = document.getElementById("locationConfiguration");
    const addLocationButton = document.getElementById("addLocationButton");
    const saveLocationsButton = document.getElementById("saveLocationsButton");
    const refreshHierarchyButton = document.getElementById("refreshHierarchyButton");
    const hierarchyMeta = document.getElementById("hierarchyMeta");
    const configBanner = document.getElementById("configBanner");
    const masjidBoardStatus = document.getElementById("masjidBoardStatus");
    const themeToggle = document.getElementById("themeToggle");

    let hierarchy = null;
    let locationRows = [];

    function setupTheme() {
        const key = "masjidpi-theme";
        const validThemes = ["system", "light", "dark"];
        let theme = validThemes.includes(localStorage.getItem(key)) ? localStorage.getItem(key) : "system";

        function applyTheme() {
            if (theme === "system") document.documentElement.removeAttribute("data-theme");
            else document.documentElement.dataset.theme = theme;
            themeToggle.textContent = "Theme: " + theme.charAt(0).toUpperCase() + theme.slice(1);
        }

        applyTheme();
        themeToggle.addEventListener("click", () => {
            theme = validThemes[(validThemes.indexOf(theme) + 1) % validThemes.length];
            localStorage.setItem(key, theme);
            applyTheme();
        });
    }

    async function jsonRequest(url, options = {}) {
        const response = await fetch(url, {cache: "no-store", ...options});
        if (!response.ok) {
            let message = `HTTP ${response.status}`;
            try {
                const body = await response.json();
                if (body && body.error) message = body.error;
            } catch (_) {}
            throw new Error(message);
        }
        return response.json();
    }

    function showBanner(message, kind = "success") {
        configBanner.textContent = message;
        configBanner.className = `config-banner ${kind}`;
        window.clearTimeout(showBanner.timer);
        showBanner.timer = window.setTimeout(() => configBanner.classList.add("hidden"), 6000);
    }

    function option(value, text, selected = false) {
        const element = document.createElement("option");
        element.value = value;
        element.textContent = text;
        element.selected = selected;
        return element;
    }

    function countries() {
        return hierarchy && Array.isArray(hierarchy.countries) ? hierarchy.countries : [];
    }

    function countryByName(name) {
        return countries().find(item => item.name === name) || null;
    }

    function regionByName(country, name) {
        if (!country || !Array.isArray(country.regions)) return null;
        return country.regions.find(item => (item.name || "") === name) || null;
    }

    function populateCountry(select, selectedValue = "") {
        select.replaceChildren(option("", "Select country…"));
        for (const country of countries()) {
            const suffix = Number.isInteger(country.count) ? ` (${country.count})` : "";
            select.append(option(country.name, country.name + suffix, country.name === selectedValue));
        }
    }

    function populateRegion(row, selectedValue = "") {
        const country = countryByName(row.country.value);
        row.region.replaceChildren(option("", country ? "Select province / region…" : "Select country first…"));
        row.region.disabled = !country;
        row.city.replaceChildren(option("", "Select province / region first…"));
        row.city.disabled = true;
        if (!country) return;

        for (const region of country.regions || []) {
            const value = region.name || "";
            const label = value || "No province / region";
            const suffix = Number.isInteger(region.count) ? ` (${region.count})` : "";
            row.region.append(option(value, label + suffix, value === selectedValue));
        }

        if (selectedValue !== "" || (country.regions || []).some(region => (region.name || "") === "" && selectedValue === "")) {
            row.region.value = selectedValue;
            populateCity(row, row.pendingCity || "");
            row.pendingCity = "";
        }
    }

    function populateCity(row, selectedValue = "") {
        const country = countryByName(row.country.value);
        const region = regionByName(country, row.region.value);
        row.city.replaceChildren(option("", region ? "Select town / city…" : "Select province / region first…"));
        row.city.disabled = !region;
        if (!region) return;

        for (const city of region.cities || []) {
            const suffix = Number.isInteger(city.count) ? ` (${city.count})` : "";
            row.city.append(option(city.name, city.name + suffix, city.name === selectedValue));
        }
        if (selectedValue) row.city.value = selectedValue;
    }

    function createSelectField(labelText, className) {
        const field = document.createElement("div");
        field.className = "location-field";
        const label = document.createElement("label");
        label.textContent = labelText;
        const select = document.createElement("select");
        select.className = className;
        field.append(label, select);
        return {field, select};
    }

    function addLocation(location = {}) {
        if (locationRows.length >= maxLocations) return;

        const wrapper = document.createElement("div");
        wrapper.className = "location-row";
        const header = document.createElement("div");
        header.className = "location-row-header";
        const title = document.createElement("span");
        title.className = "location-row-title";
        const remove = document.createElement("button");
        remove.type = "button";
        remove.className = "location-remove";
        remove.textContent = "Remove";
        header.append(title, remove);

        const fields = document.createElement("div");
        fields.className = "location-fields";
        const countryField = createSelectField("Country", "location-country");
        const regionField = createSelectField("Province / Region", "location-region");
        const cityField = createSelectField("Town / City", "location-city");
        fields.append(countryField.field, regionField.field, cityField.field);
        wrapper.append(header, fields);
        locationConfiguration.append(wrapper);

        const row = {
            wrapper,
            title,
            country: countryField.select,
            region: regionField.select,
            city: cityField.select,
            pendingCity: location.city || "",
        };
        locationRows.push(row);

        populateCountry(row.country, location.country || "");
        populateRegion(row, location.region || "");

        row.country.addEventListener("change", () => {
            row.pendingCity = "";
            populateRegion(row);
        });
        row.region.addEventListener("change", () => populateCity(row));
        remove.addEventListener("click", () => {
            row.wrapper.remove();
            locationRows = locationRows.filter(item => item !== row);
            ensureOneLocation();
            updateLocationUI();
        });

        updateLocationUI();
    }

    function ensureOneLocation() {
        if (locationRows.length === 0) addLocation();
    }

    function updateLocationUI() {
        locationRows.forEach((row, index) => {
            row.title.textContent = `Location ${index + 1}`;
            const remove = row.wrapper.querySelector(".location-remove");
            remove.disabled = locationRows.length === 1;
        });
        addLocationButton.disabled = locationRows.length >= maxLocations;
        addLocationButton.textContent = locationRows.length >= maxLocations ? "Maximum 3 Locations" : "+ Add Location";
    }

    function readLocations() {
        return locationRows.map(row => ({
            country: row.country.value.trim(),
            region: row.region.value.trim(),
            city: row.city.value.trim(),
        }));
    }

    function validateLocations(locations) {
        for (let index = 0; index < locations.length; index += 1) {
            const location = locations[index];
            if (!location.country || !location.city) return `Complete Location ${index + 1} before saving.`;
        }
        const keys = new Set();
        for (const location of locations) {
            const key = [location.country, location.region, location.city].join("\u0000").toLowerCase();
            if (keys.has(key)) return "The same location cannot be selected more than once.";
            keys.add(key);
        }
        return "";
    }

    async function loadConfiguration() {
        const [hierarchyState, scopeState] = await Promise.all([
            jsonRequest(hierarchyEndpoint),
            jsonRequest(scopeEndpoint),
        ]);
        hierarchy = hierarchyState;
        locationConfiguration.replaceChildren();
        locationRows = [];
        const locations = Array.isArray(scopeState.locations) && scopeState.locations.length ? scopeState.locations : [{}];
        for (const location of locations.slice(0, maxLocations)) addLocation(location);
        ensureOneLocation();
        updateHierarchyMeta();
    }

    function updateHierarchyMeta() {
        if (!hierarchy) {
            hierarchyMeta.textContent = "";
            return;
        }
        const count = countries().length;
        const stamp = hierarchy.validated_at || hierarchy.retrieved_at;
        const dateText = stamp ? new Date(stamp).toLocaleString() : "unknown time";
        hierarchyMeta.textContent = `${count} countries available · hierarchy last validated ${dateText}`;
    }

    async function saveLocations() {
        const locations = readLocations();
        const validationError = validateLocations(locations);
        if (validationError) {
            showBanner(validationError, "error");
            return;
        }

        saveLocationsButton.disabled = true;
        saveLocationsButton.textContent = "Saving…";
        try {
            await jsonRequest(scopeEndpoint, {
                method: "PUT",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({locations}),
            });
            showBanner("Locations saved. Updating the scoped MasjidBoard catalogue…", "success");
            try {
                await jsonRequest(catalogueRefreshEndpoint, {method: "POST"});
                showBanner("Locations saved and MasjidBoard catalogue updated.", "success");
            } catch (error) {
                showBanner(`Locations saved, but catalogue refresh failed: ${error.message}`, "warning");
            }
        } catch (error) {
            showBanner(`Could not save locations: ${error.message}`, "error");
        } finally {
            saveLocationsButton.disabled = false;
            saveLocationsButton.textContent = "Save Locations";
        }
    }

    async function refreshHierarchy() {
        refreshHierarchyButton.disabled = true;
        refreshHierarchyButton.textContent = "Refreshing…";
        try {
            await jsonRequest(hierarchyRefreshEndpoint, {method: "POST"});
            const selected = readLocations();
            hierarchy = await jsonRequest(hierarchyEndpoint);
            locationConfiguration.replaceChildren();
            locationRows = [];
            for (const location of selected) addLocation(location);
            ensureOneLocation();
            updateHierarchyMeta();
            showBanner("Location hierarchy refreshed.", "success");
        } catch (error) {
            showBanner(`Could not refresh locations: ${error.message}`, "error");
        } finally {
            refreshHierarchyButton.disabled = false;
            refreshHierarchyButton.textContent = "Refresh Locations";
        }
    }

    async function loadStatus() {
        try {
            const status = await jsonRequest(statusEndpoint);
            const boards = Array.isArray(status.boards) ? status.boards : [];
            masjidBoardStatus.replaceChildren();

            const configured = document.createElement("div");
            configured.className = "status-row";
            configured.innerHTML = `<span>Configured</span><strong>${status.configured ? "Yes" : "No"}</strong>`;
            const count = document.createElement("div");
            count.className = "status-row";
            count.innerHTML = `<span>Selected boards</span><strong>${boards.length}</strong>`;
            masjidBoardStatus.append(configured, count);

            if (boards.length) {
                const list = document.createElement("div");
                list.className = "status-board-list";
                for (const board of boards) {
                    const item = document.createElement("div");
                    item.className = "status-board-item";
                    const name = document.createElement("span");
                    name.textContent = board.name || board.catalogue_id || "MasjidBoard";
                    const state = document.createElement("span");
                    const value = board.status || "unknown";
                    state.className = `status-board-state ${value}`;
                    state.textContent = value;
                    item.append(name, state);
                    list.append(item);
                }
                masjidBoardStatus.append(list);
            }
        } catch (error) {
            masjidBoardStatus.innerHTML = `<p class="status-detail status-detail-error">Unable to load MasjidBoard status: ${error.message}</p>`;
        }
    }

    setupTheme();
    addLocationButton.addEventListener("click", () => addLocation());
    saveLocationsButton.addEventListener("click", saveLocations);
    refreshHierarchyButton.addEventListener("click", refreshHierarchy);

    loadConfiguration().catch(error => {
        locationConfiguration.innerHTML = `<p class="status-detail status-detail-error">Unable to load locations: ${error.message}</p>`;
        showBanner("MasjidBoard location configuration could not be loaded.", "error");
    });
    loadStatus();
})();
