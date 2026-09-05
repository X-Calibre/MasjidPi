"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

const root = path.resolve(__dirname, "..");
const utilitySource = fs.readFileSync(path.join(root, "frontend/masjidboard-community-utils.js"), "utf8");
const applianceJS = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");
const landscapeJS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");
const applianceCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.css"), "utf8");
const landscapeCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.css"), "utf8");
const displayJS = fs.readFileSync(path.join(root, "frontend/masjidboard-display.js"), "utf8");
const html = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");

class TestDOMParser {
    parseFromString(value) {
        return {body: {textContent: String(value).replace(/<[^>]*>/g, "")}};
    }
}

const context = {window: {}, DOMParser: TestDOMParser};
vm.runInNewContext(utilitySource, context);
const {collectCommunityItems, duaAfterAdhanItem, duaAfterAdhanWindowMinutes, orderedFields} = context.window.MasjidBoardCommunityUtils;

const announcements = collectCommunityItems([{
    name: "Test Masjid",
    announcements: [
        {category: "weekly_programme", title: "Weekly Programs", content: "Tafseer after Esha"},
        {category: "ramadan_programme", title: "Taraweeh 2026", content: "One juz nightly"},
        {category: "unknown", title: "Custom notice", content: "Keep this generic"},
        {category: "announcement", title: "تهنئة", content: "نبارك للإمام الجديد"},
    ],
}]);
assert.deepEqual(Array.from(announcements, item => item.type), [
    "weekly_programme", "ramadan_programme", "announcement", "announcement",
]);
assert.equal(announcements[0].typeLabel, "Weekly Programme");
assert.equal(announcements[2].typeLabel, "");
assert.equal(announcements[3].body, "نبارك للإمام الجديد");

const board = {
    time_zone: "Africa/Johannesburg",
    prayers: [{key: "fajr", adhan: {hour: 5, minute: 30}}],
};
assert.equal(duaAfterAdhanItem([board], new Date("2026-09-05T03:35:00Z"), false), null);
assert.equal(duaAfterAdhanWindowMinutes, 5);
const dua = duaAfterAdhanItem([board], new Date("2026-09-05T03:34:00Z"), true);
assert.equal(dua.type, "dua_after_adhan");
assert.equal(dua.source, "MasjidPi");
assert.match(dua.fields.arabic, /اللَّهُمَّ/);
assert.match(dua.fields.translation, /^O Allah/);
assert.deepEqual(Array.from(orderedFields(dua), field => field.label), ["Arabic", "Translation", "Note"]);
assert.equal(duaAfterAdhanItem([board], new Date("2026-09-05T03:35:00Z"), true), null, "five-minute display window must be exclusive at its end");
const midnightBoard = {
    time_zone: "UTC",
    prayers: [{key: "esha", adhan: {hour: 23, minute: 58}}],
};
assert.equal(duaAfterAdhanItem([midnightBoard], new Date("2026-09-06T00:02:00Z"), true).type, "dua_after_adhan", "the window must continue through local midnight");
const fixedOffsetBoard = {
    time_zone: "GMT+10:00",
    prayers: [{key: "fajr", adhan: {hour: 5, minute: 30}}],
};
assert.equal(duaAfterAdhanItem([fixedOffsetBoard], new Date("2026-09-05T19:34:00Z"), true).type, "dua_after_adhan", "provider GMT offsets must not fall back to the browser timezone");
assert.equal(duaAfterAdhanItem([], new Date("2026-09-05T03:40:00Z"), true, 10, true).type, "dua_after_adhan");

for (const renderer of [applianceJS, landscapeJS]) {
    assert.match(renderer, /duaAfterAdhanItem/);
    assert.match(renderer, /show_dua_after_adhan/);
    assert.match(renderer, /params\.get\("dua-fixture"\) === "1"/);
    assert.match(renderer, /duaVisibleNow !== duaAfterAdhanVisible/, "the card must react to the display clock without waiting for new API data");
    assert.match(renderer, /\.dir = "auto"/);
}
assert.match(applianceCSS, /\[dir="rtl"\]/);
assert.match(landscapeCSS, /\[dir="rtl"\]/);
assert.match(landscapeJS, /item\.type === "dua_after_adhan"/, "the Landscape Dua card must reserve a detailed-content slot");
assert.match(landscapeJS, /items = \[duaItem\]/, "the Landscape Dua card must override all other cards");
assert.match(applianceJS, /slides = \[communitySlide\(\[duaItem\]\)\]/, "the Appliance Dua card must suspend all other slides");
assert.match(landscapeJS, /item\.type !== "dua_after_adhan"/, "Landscape must omit Dua source attribution");
assert.match(applianceJS, /item\.type !== "dua_after_adhan"/, "Appliance must omit Dua source attribution");
assert.match(displayJS, /window\.MasjidBoardCurrentView = view/, "the base renderer must retain the latest view for late layout modules");
assert.match(landscapeJS, /if \(window\.MasjidBoardCurrentView\) refresh\(window\.MasjidBoardCurrentView\)/, "Landscape must replay an already-fetched view");
assert.match(html, /masjidboard-community-utils\.js\?v=20260905-dua-priority/, "Dua priority assets must use a new cache key");

console.log("MasjidBoard Section 10 community-content tests passed");
