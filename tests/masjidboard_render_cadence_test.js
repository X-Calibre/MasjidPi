"use strict";

const assert = require("node:assert/strict");

// Load the display utilities with only the small DOM surface used at startup.
const inertElement = () => ({
    classList: {add() {}, remove() {}, toggle() {}, contains() { return false; }},
    style: {setProperty() {}},
    dataset: {},
    children: [],
    querySelectorAll() { return []; },
    replaceChildren() {},
});

global.document = {
    hidden: false,
    getElementById() { return inertElement(); },
};
global.CustomEvent = function CustomEvent() {};
global.fetch = () => new Promise(() => {});
global.window = {
    location: {search: ""},
    MasjidBoardDate: {isIslamicFriday() { return false; }},
    dispatchEvent() {},
    setInterval() {},
    setTimeout() {},
};

require("../frontend/masjidboard-display.js");

const {prayerMinuteKey, viewSignature} = window.MasjidBoardDisplayUtils;
const boards = [{catalogue_id: "test"}];

assert.equal(
    prayerMinuteKey(boards, new Date(2026, 7, 29, 7, 15, 1)),
    prayerMinuteKey(boards, new Date(2026, 7, 29, 7, 15, 59)),
    "second ticks must not trigger a prayer-grid rebuild",
);
assert.notEqual(
    prayerMinuteKey(boards, new Date(2026, 7, 29, 7, 15, 59)),
    prayerMinuteKey(boards, new Date(2026, 7, 29, 7, 16, 0)),
    "a new minute must refresh countdown placement and the upcoming event",
);

const view = {configured: true, boards: [{catalogue_id: "test"}]};
assert.equal(viewSignature(view), viewSignature(structuredClone(view)), "unchanged API data must not trigger a structural rebuild");
assert.notEqual(viewSignature(view), viewSignature({...view, configured: false}), "changed API data must trigger a structural rebuild");

console.log("MasjidBoard render cadence tests passed");
