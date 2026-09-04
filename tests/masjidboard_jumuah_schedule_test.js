"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

const root = path.resolve(__dirname, "..");
const source = fs.readFileSync(path.join(root, "frontend/masjidboard-community-utils.js"), "utf8");
const applianceCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.css"), "utf8");
const landscapeCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.css"), "utf8");
const applianceJS = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");
const landscapeJS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");

class TestDOMParser {
    parseFromString(value) {
        return {body: {textContent: String(value).replace(/<[^>]*>/g, "")}};
    }
}

const context = {window: {}, DOMParser: TestDOMParser};
vm.runInNewContext(source, context);
const {detailedJumuahItems} = context.window.MasjidBoardCommunityUtils;

const board = {
    name: "Masjid us Salaam",
    show_detailed_jumuah: true,
    jumuah: [{
        adhan: {hour: 12, minute: 20},
        jamaah: {hour: 12, minute: 55},
        khateeb: "Sunnats immediately before khutbah",
        events: [
            {code: "0", heading: "Adhān", time: {hour: 12, minute: 20}},
            {code: "3", heading: "Sunan", time: {hour: 12, minute: 50}},
            {code: "6", heading: "Khutbah", time: {hour: 12, minute: 55}},
        ],
    }],
};

const fridayItems = detailedJumuahItems([board], new Date("2026-09-04T10:00:00Z"), () => true);
assert.equal(fridayItems.length, 1);
assert.equal(fridayItems[0].title, "Jumu’ah Schedule");
assert.equal(fridayItems[0].source, "Masjid us Salaam");
assert.equal(fridayItems[0].body, "Sunnats immediately before khutbah");
assert.deepEqual(Array.from(fridayItems[0].schedule, event => `${event.heading}:${event.time}`), [
    "Adhān:12:20", "Sunan:12:50", "Khutbah:12:55",
]);

assert.equal(detailedJumuahItems([board], new Date(), () => false).length, 0, "schedule must be Friday-only");
assert.equal(detailedJumuahItems([{...board, show_detailed_jumuah: false}], new Date(), () => true).length, 0, "per-masjid switch must hide the schedule");
assert.equal(detailedJumuahItems([{...board, show_detailed_jumuah: undefined}], new Date(), () => true).length, 1, "missing preference must default to enabled");

const fallback = detailedJumuahItems([{
    name: "Fallback Masjid",
    jumuah: [{adhan: {hour: 12, minute: 30}, jamaah: {hour: 13, minute: 0}}],
}], new Date(), () => true);
assert.deepEqual(Array.from(fallback[0].schedule, event => `${event.heading}:${event.time}`), ["Adhan:12:30", "Salaah:13:00"]);

assert.match(applianceCSS, /compact\.appliance-community-jumuah_schedule \.appliance-community-source\s*\{[^}]*font-size:26px/s);
assert.match(applianceCSS, /compact\.appliance-community-jumuah_schedule \.appliance-community-body\s*\{[^}]*font-size:22px/s);
assert.match(applianceCSS, /compact \.appliance-jumuah-schedule\s*\{[^}]*margin-top:48px/s);
assert.match(landscapeCSS, /\.detailed-community-jumuah_schedule \.detailed-community-source\s*\{[^}]*font-size:clamp\(1\.35rem,1\.55vw,1\.7rem\)/s);
for (const renderer of [applianceJS, landscapeJS]) {
    assert.match(renderer, /params\.get\("jumuah-fixture"\) === "khateeb"/);
    assert.match(renderer, /item\.body = "Khateeb: To be announced"/);
}

console.log("MasjidBoard detailed Jumu'ah schedule tests passed");
