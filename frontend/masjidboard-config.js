(() => {
    "use strict";

    const maxLocations = 3;
    const maxBoards = 3;
    const hierarchyEndpoint = "/api/masjidboard/hierarchy";
    const hierarchyRefreshEndpoint = "/api/masjidboard/hierarchy/refresh";
    const scopeEndpoint = "/api/masjidboard/scope";
    const catalogueEndpoint = "/api/masjidboard/catalogue";
    const catalogueRefreshEndpoint = "/api/masjidboard/catalogue/refresh";
    const selectionEndpoint = "/api/masjidboard/selection";
    const statusEndpoint = "/api/masjidboard/status";

    const locationConfiguration = document.getElementById("locationConfiguration");
    const addLocationButton = document.getElementById("addLocationButton");
    const saveLocationsButton = document.getElementById("saveLocationsButton");
    const locationSaveStatus = document.getElementById("locationSaveStatus");
    const refreshHierarchyButton = document.getElementById("refreshHierarchyButton");
    const hierarchyMeta = document.getElementById("hierarchyMeta");
    const masjidBoardStatus = document.getElementById("masjidBoardStatus");
    const themeToggle = document.getElementById("themeToggle");

    const boardSelection = document.getElementById("boardSelection");
    const boardSearch = document.getElementById("boardSearch");
    const availableBoards = document.getElementById("availableBoards");
    const availableBoardMeta = document.getElementById("availableBoardMeta");
    const addBoardButton = document.getElementById("addBoardButton");
    const boardSaveStatus = document.getElementById("boardSaveStatus");
    const refreshCatalogueButton = document.getElementById("refreshCatalogueButton");
    const refreshBoardsButton = document.getElementById("refreshBoardsButton");

    let hierarchy = null;
    let locationRows = [];
    let catalogueRecords = [];
    let selectedBoards = [];
    let lastSavedBoards = [];
    let boardSaveTimer = 0;
    let boardSavePending = false;
    let boardSaving = false;

    function setSaveStatus(element, message, className = "") {
        if (!element) return;
        element.textContent = message;
        element.className = `config-save-status${className ? ` ${className}` : ""}`;
    }

    function markLocationsUnsaved() {
        setSaveStatus(locationSaveStatus, "Unsaved location changes", "unsaved");
    }

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
        window.MasjidPiUI.notify(message, kind);
    }

    function option(value, text, selected = false) {
        const element = document.createElement("option");
        element.value = value;
        element.textContent = text;
        element.selected = selected;
        return element;
    }

    function textElement(tag, className, value) {
        const element = document.createElement(tag);
        if (className) element.className = className;
        element.textContent = value;
        return element;
    }

    function statusDetail(label, value, className = "") {
        const row = document.createElement("div");
        row.className = `status-board-detail${className ? ` ${className}` : ""}`;
        row.append(textElement("span", "status-board-detail-label", label));
        row.append(textElement("span", "status-board-detail-value", value));
        return row;
    }

    function statusLabel(value) {
        return {current: "Current", stale: "Stale", unavailable: "Unavailable"}[value] || "Unknown";
    }

    function formatTimestamp(value) {
        if (!value) return "Not yet";
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return "Unknown";
        const now = new Date();
        const sameDay = date.toDateString() === now.toDateString();
        const time = date.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"});
        return sameDay ? `Today, ${time}` : date.toLocaleString([], {
            year: "numeric", month: "short", day: "numeric", hour: "2-digit", minute: "2-digit",
        });
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
            markLocationsUnsaved();
        });
        row.region.addEventListener("change", () => { populateCity(row); markLocationsUnsaved(); });
        row.city.addEventListener("change", markLocationsUnsaved);
        remove.addEventListener("click", () => {
            row.wrapper.remove();
            locationRows = locationRows.filter(item => item !== row);
            ensureOneLocation();
            updateLocationUI();
            markLocationsUnsaved();
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
            setSaveStatus(locationSaveStatus, "Locations saved", "saved");
            showBanner("Locations saved. Updating the scoped MasjidBoard catalogue…", "success");
            try {
                await jsonRequest(catalogueRefreshEndpoint, {method: "POST"});
                await loadBoardConfiguration();
                const count = activeCatalogueRecords().length;
                showBanner(`Locations saved · ${count} Masjid${count === 1 ? "" : "s"} available`, "success");
            } catch (error) {
                showBanner(`Locations saved, but catalogue refresh failed: ${error.message}`, "warning");
            }
        } catch (error) {
            setSaveStatus(locationSaveStatus, "Locations could not be saved", "error");
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
            showBanner("Location list refreshed.", "success");
        } catch (error) {
            showBanner(`Could not refresh locations: ${error.message}`, "error");
        } finally {
            refreshHierarchyButton.disabled = false;
            refreshHierarchyButton.textContent = "Refresh Location List";
        }
    }

    function activeCatalogueRecords() {
        return catalogueRecords.filter(record => record.status === "active");
    }

    function recordByID(id) {
        return catalogueRecords.find(record => record.id === id) || null;
    }

    function boardLabel(record) {
        const location = [record.city, record.region, record.country].filter(Boolean).join(", ");
        return location ? `${record.name} — ${location}` : record.name;
    }

    function selectedIDSet() {
        return new Set(selectedBoards.map(board => board.catalogue_id));
    }

    function renderAvailableBoards() {
        const query = boardSearch.value.trim().toLowerCase();
        const selected = selectedIDSet();
        const records = activeCatalogueRecords()
            .filter(record => !selected.has(record.id))
            .filter(record => !query || [record.name, record.city, record.region, record.country]
                .some(value => String(value || "").toLowerCase().includes(query)))
            .sort((left, right) => boardLabel(left).localeCompare(boardLabel(right)));

        availableBoards.replaceChildren();
        for (const record of records) {
            availableBoards.append(option(record.id, boardLabel(record)));
        }
        availableBoardMeta.textContent = `${records.length} available MasjidBoard${records.length === 1 ? "" : "s"}`;
        addBoardButton.disabled = selectedBoards.length >= maxBoards || records.length === 0;
        addBoardButton.textContent = selectedBoards.length >= maxBoards ? "Maximum 3 Masjids" : "+ Add Selected Masjid";
    }

    function renderSelectedBoards() {
        boardSelection.replaceChildren();

        if (selectedBoards.length === 0) {
            const empty = document.createElement("p");
            empty.className = "status-detail";
            empty.textContent = "No MasjidBoards selected yet. Add at least one from the list below.";
            boardSelection.append(empty);
        }

        selectedBoards.forEach((board, index) => {
            const row = document.createElement("div");
            row.className = "selected-board-row";

            const order = document.createElement("div");
            order.className = "selected-board-order";
            order.textContent = String(index + 1);

            const details = document.createElement("div");
            details.className = "selected-board-details";
            const name = document.createElement("strong");
            name.textContent = board.name || board.catalogue_id;
            const record = recordByID(board.catalogue_id);
            const meta = document.createElement("small");
            if (record && record.status === "active") {
                meta.textContent = [record.city, record.region, record.country].filter(Boolean).join(" · ");
            } else {
                meta.className = "board-outside-scope";
                meta.textContent = "Outside the current location catalogue — remove or replace before saving";
            }
            const role = document.createElement("small");
            role.className = index === 0 ? "board-role primary" : "board-role";
            role.textContent = index === 0
                ? "Primary · supplies additional daily times"
                : "Secondary";
            const jumuahOption = document.createElement("label");
            jumuahOption.className = "selected-board-option";
            const jumuahToggle = document.createElement("input");
            jumuahToggle.type = "checkbox";
            jumuahToggle.checked = board.show_detailed_jumuah !== false;
            jumuahToggle.setAttribute("aria-label", `Show detailed Jumu'ah schedule for ${board.name || "masjid"}`);
            jumuahToggle.addEventListener("change", () => {
                board.show_detailed_jumuah = jumuahToggle.checked;
                scheduleBoardSave();
            });
            jumuahOption.append(jumuahToggle, document.createTextNode(" Detailed Jumu’ah schedule"));
            details.append(name, role, meta, jumuahOption);

            const actions = document.createElement("div");
            actions.className = "selected-board-actions";
            const up = document.createElement("button");
            up.type = "button";
            up.className = "secondary board-order-button";
            up.textContent = "↑";
            up.title = "Move up";
            up.setAttribute("aria-label", `Move ${board.name || "masjid"} up`);
            up.disabled = index === 0;
            up.addEventListener("click", () => moveSelectedBoard(index, -1));

            const down = document.createElement("button");
            down.type = "button";
            down.className = "secondary board-order-button";
            down.textContent = "↓";
            down.title = "Move down";
            down.setAttribute("aria-label", `Move ${board.name || "masjid"} down`);
            down.disabled = index === selectedBoards.length - 1;
            down.addEventListener("click", () => moveSelectedBoard(index, 1));

            const remove = document.createElement("button");
            remove.type = "button";
            remove.className = "secondary board-remove-button";
            remove.textContent = "Remove";
            remove.disabled = selectedBoards.length === 1;
            remove.setAttribute("aria-label", `Remove ${board.name || "masjid"}`);
            remove.addEventListener("click", () => {
                selectedBoards.splice(index, 1);
                renderBoardConfiguration();
                scheduleBoardSave();
            });

            actions.append(up, down, remove);
            row.append(order, details, actions);
            boardSelection.append(row);
        });

    }

    function renderBoardConfiguration() {
        renderSelectedBoards();
        renderAvailableBoards();
    }

    function moveSelectedBoard(index, direction) {
        const target = index + direction;
        if (target < 0 || target >= selectedBoards.length) return;
        const [board] = selectedBoards.splice(index, 1);
        selectedBoards.splice(target, 0, board);
        renderBoardConfiguration();
        scheduleBoardSave();
    }

    function addSelectedBoard() {
        const id = availableBoards.value;
        if (!id || selectedBoards.length >= maxBoards) return;
        const record = recordByID(id);
        if (!record || record.status !== "active") return;
        selectedBoards.push({
            catalogue_id: record.id,
            provider: record.provider,
            external_id: record.external_id,
            name: record.name,
            time_zone_offset_ms: record.time_zone_offset_ms,
            show_detailed_jumuah: true,
        });
        renderBoardConfiguration();
        scheduleBoardSave();
    }

    async function loadBoardConfiguration() {
        const [catalogue, selection] = await Promise.all([
            jsonRequest(catalogueEndpoint),
            jsonRequest(selectionEndpoint),
        ]);
        catalogueRecords = Array.isArray(catalogue.records) ? catalogue.records : [];
        selectedBoards = Array.isArray(selection.boards) ? selection.boards.slice(0, maxBoards).map(board => ({
            ...board,
            show_detailed_jumuah: board.show_detailed_jumuah !== false,
        })) : [];
        lastSavedBoards = selectedBoards.map(board => ({...board}));
        renderBoardConfiguration();
        setSaveStatus(boardSaveStatus, "Changes are saved automatically.");
    }

    async function saveBoardSelection() {
        if (selectedBoards.length < 1 || selectedBoards.length > maxBoards) {
            setSaveStatus(boardSaveStatus, "Select between one and three MasjidBoards", "unsaved");
            showBanner("Select between one and three MasjidBoards.", "error");
            return;
        }
        const inactive = selectedBoards.find(board => {
            const record = recordByID(board.catalogue_id);
            return !record || record.status !== "active";
        });
        if (inactive) {
            setSaveStatus(boardSaveStatus, "Remove or replace Masjids outside the current location", "unsaved");
            showBanner("Remove or replace MasjidBoards outside the current location catalogue before saving.", "error");
            return;
        }

        const candidate = selectedBoards.map(board => ({...board}));
        setSaveStatus(boardSaveStatus, "Saving…", "saving");
        try {
            const response = await jsonRequest(selectionEndpoint, {
                method: "PUT",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    catalogue_ids: candidate.map(board => board.catalogue_id),
                    detailed_jumuah: Object.fromEntries(candidate.map(board => [board.catalogue_id, board.show_detailed_jumuah !== false])),
                }),
            });
            lastSavedBoards = Array.isArray(response.boards) ? response.boards.map(board => ({...board})) : candidate;
            if (!boardSavePending) {
                selectedBoards = lastSavedBoards.map(board => ({...board}));
                renderBoardConfiguration();
            }
            await loadStatus();
            setSaveStatus(boardSaveStatus, boardSavePending ? "Saving…" : "Saved · display updated", boardSavePending ? "saving" : "saved");
        } catch (error) {
            boardSavePending = false;
            selectedBoards = lastSavedBoards.map(board => ({...board}));
            renderBoardConfiguration();
            setSaveStatus(boardSaveStatus, "Could not save changes", "error");
            showBanner(`Could not save selected Masjids: ${error.message}`, "error");
        }
    }

    async function drainBoardSaves() {
        if (boardSaving) return;
        boardSaving = true;
        while (boardSavePending) {
            boardSavePending = false;
            await saveBoardSelection();
        }
        boardSaving = false;
    }

    function scheduleBoardSave() {
        boardSavePending = true;
        setSaveStatus(boardSaveStatus, "Saving…", "saving");
        window.clearTimeout(boardSaveTimer);
        boardSaveTimer = window.setTimeout(() => { void drainBoardSaves(); }, 450);
    }

    async function refreshCatalogue() {
        refreshCatalogueButton.disabled = true;
        refreshCatalogueButton.textContent = "Refreshing…";
        try {
            await jsonRequest(catalogueRefreshEndpoint, {method: "POST"});
            await loadBoardConfiguration();
            const count = activeCatalogueRecords().length;
            showBanner(`Masjid list refreshed · ${count} Masjid${count === 1 ? "" : "s"} available`, "success");
        } catch (error) {
            showBanner(`Could not refresh MasjidBoards: ${error.message}`, "error");
        } finally {
            refreshCatalogueButton.disabled = false;
            refreshCatalogueButton.textContent = "Refresh Masjid List";
        }
    }

    async function refreshBoards() {
        refreshBoardsButton.disabled = true;
        refreshBoardsButton.textContent = "Refreshing…";
        try {
            await jsonRequest("/api/masjidboard/boards/refresh", {method: "POST"});
            await loadStatus();
            showBanner("Timetables refreshed.", "success");
        } catch (error) {
            showBanner(`Could not refresh timetables: ${error.message}`, "error");
        } finally {
            refreshBoardsButton.disabled = false;
            refreshBoardsButton.textContent = "Refresh Timetables";
        }
    }

    async function loadStatus() {
        try {
            const status = await jsonRequest(statusEndpoint);
            const boards = Array.isArray(status.boards) ? status.boards : [];
            masjidBoardStatus.replaceChildren();

            const summary = document.createElement("div");
            summary.className = "status-summary";
            summary.append(statusDetail("Configured", status.configured ? "Yes" : "No"));
            summary.append(statusDetail("Selected Masjids", String(boards.length)));
            masjidBoardStatus.append(summary);

            if (boards.length) {
                const list = document.createElement("div");
                list.className = "status-board-list status-board-list-detailed";
                for (const board of boards) {
                    const value = board.status || "unknown";
                    const item = document.createElement("article");
                    item.className = `status-board-card status-${value}`;
                    const heading = document.createElement("div");
                    heading.className = "status-board-heading";
                    heading.append(textElement("strong", "status-board-name", board.name || board.catalogue_id || "MasjidBoard"));
                    heading.append(textElement("span", `status-board-state ${value}`, statusLabel(value)));
                    item.append(heading);

                    const details = document.createElement("div");
                    details.className = "status-board-details";
                    if (board.using_cached_data || value === "stale") details.append(textElement("div", "status-board-note warning", "Using cached timetable data"));
                    if (value === "unavailable" && !board.board) details.append(textElement("div", "status-board-note error", "No timetable data is currently available"));
                    if (board.last_successful_update) details.append(statusDetail("Last successful update", formatTimestamp(board.last_successful_update)));
                    if (board.last_attempt) details.append(statusDetail("Last attempt", formatTimestamp(board.last_attempt)));
                    if (board.update_error) details.append(statusDetail("Update error", board.update_error, "error"));
                    if (board.persistence_error) details.append(statusDetail("Cache error", board.persistence_error, "error"));
                    item.append(details);
                    list.append(item);
                }
                masjidBoardStatus.append(list);
            }
        } catch (error) {
            masjidBoardStatus.innerHTML = `<p class="status-detail status-detail-error">Unable to load MasjidBoard status: ${error.message}</p>`;
        }
    }

    setupTheme();
    addLocationButton.addEventListener("click", () => { addLocation(); markLocationsUnsaved(); });
    saveLocationsButton.addEventListener("click", saveLocations);
    refreshHierarchyButton.addEventListener("click", refreshHierarchy);
    boardSearch.addEventListener("input", renderAvailableBoards);
    addBoardButton.addEventListener("click", addSelectedBoard);
    availableBoards.addEventListener("dblclick", addSelectedBoard);
    refreshCatalogueButton.addEventListener("click", refreshCatalogue);
    refreshBoardsButton.addEventListener("click", refreshBoards);

    loadConfiguration().catch(error => {
        locationConfiguration.innerHTML = `<p class="status-detail status-detail-error">Unable to load locations: ${error.message}</p>`;
        showBanner("MasjidBoard location configuration could not be loaded.", "error");
    });
    loadBoardConfiguration().catch(error => {
        boardSelection.innerHTML = `<p class="status-detail status-detail-error">Unable to load MasjidBoards: ${error.message}</p>`;
        availableBoards.replaceChildren();
        addBoardButton.disabled = true;
        setSaveStatus(boardSaveStatus, "Masjid selection could not be loaded", "error");
    });
    loadStatus();
})();
