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
    body: inertElement(),
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

require("../frontend/masjidboard-warning-utils.js");
require("../frontend/masjidboard-display.js");

const {prayerRenderKey, viewSignature} = window.MasjidBoardDisplayUtils;
const boards = [{
    catalogue_id: "test",
    prayers: [
        {key: "fajr", adhan: {hour: 5, minute: 30}, jamaah: {hour: 6, minute: 0}},
        {key: "dhuhr", adhan: {hour: 12, minute: 30}, jamaah: {hour: 13, minute: 0}},
    ],
}];

assert.equal(
    prayerRenderKey(boards, new Date(2026, 7, 29, 7, 15, 1)),
    prayerRenderKey(boards, new Date(2026, 7, 29, 7, 59, 59)),
    "second ticks must not trigger a prayer-grid rebuild",
);
assert.equal(
    prayerRenderKey(boards, new Date(2026, 7, 29, 7, 59, 59)),
    prayerRenderKey(boards, new Date(2026, 7, 29, 8, 0, 0)),
    "ordinary minute and hour boundaries must not rebuild an unchanged prayer grid",
);
assert.notEqual(
    prayerRenderKey(boards, new Date(2026, 7, 29, 5, 59, 59)),
    prayerRenderKey(boards, new Date(2026, 7, 29, 6, 1, 0)),
    "the grid must rebuild when the upcoming prayer event changes",
);
assert.notEqual(
    prayerRenderKey(boards, new Date(2026, 7, 29, 23, 59, 59)),
    prayerRenderKey(boards, new Date(2026, 7, 30, 0, 0, 0)),
    "the grid must rebuild when the local display date changes",
);

const view = {configured: true, boards: [{catalogue_id: "test"}]};
assert.equal(viewSignature(view), viewSignature(structuredClone(view)), "unchanged API data must not trigger a structural rebuild");
assert.notEqual(viewSignature(view), viewSignature({...view, configured: false}), "changed API data must trigger a structural rebuild");

console.log("MasjidBoard render cadence tests passed");
