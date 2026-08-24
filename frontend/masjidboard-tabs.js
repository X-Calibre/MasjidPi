(() => {
    "use strict";

    const tabs = Array.from(document.querySelectorAll("[data-board-tab]"));
    const panels = Array.from(document.querySelectorAll("[data-board-panel]"));
    const validTabs = new Set(tabs.map(tab => tab.dataset.boardTab));

    function activate(name, updateHash = true) {
        const selected = validTabs.has(name) ? name : "masjids";
        for (const tab of tabs) {
            const active = tab.dataset.boardTab === selected;
            tab.classList.toggle("active", active);
            tab.setAttribute("aria-selected", String(active));
        }
        for (const panel of panels) {
            const active = panel.dataset.boardPanel === selected;
            panel.hidden = !active;
            panel.classList.toggle("active", active);
        }
        if (updateHash) history.replaceState(null, "", "#" + selected);
    }

    for (const tab of tabs) {
        tab.addEventListener("click", () => activate(tab.dataset.boardTab));
    }

    const initial = location.hash.slice(1);
    activate(validTabs.has(initial) ? initial : "masjids", false);
})();