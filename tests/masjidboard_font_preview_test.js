"use strict";

const assert = require("node:assert/strict");

function loadTheme(search) {
    const listeners = {};
    global.document = {
        body: {dataset: {}},
        documentElement: {dataset: {}},
    };
    global.window = {
        location: {search, href: `http://localhost/masjidboard.html${search}`},
        addEventListener(name, handler) { listeners[name] = handler; },
    };
    delete require.cache[require.resolve("../frontend/masjidboard-theme.js")];
    require("../frontend/masjidboard-theme.js");
    return document.documentElement.dataset.boardFont;
}

assert.equal(loadTheme("?font=ibm-plex"), "ibm-plex");
assert.equal(loadTheme("?font=source-sans"), "source-sans");
assert.equal(loadTheme("?font=unknown"), undefined);
assert.equal(loadTheme(""), undefined);

console.log("MasjidBoard font preview tests passed");
