"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

const root = path.resolve(__dirname, "..");
const utilitySource = fs.readFileSync(path.join(root, "frontend/masjidboard-community-utils.js"), "utf8");
const applianceSource = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");
const landscapeSource = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");

class TestDOMParser {
    parseFromString(value) {
        return {body: {textContent: String(value).replace(/<[^>]*>/g, "")}};
    }
}

const context = {window: {}, DOMParser: TestDOMParser, Intl, Date};
vm.runInNewContext(utilitySource, context);
const {specialDhuhrItem} = context.window.MasjidBoardCommunityUtils;

function board(label, timeZone = "Africa/Johannesburg") {
    return {time_zone: timeZone, special_dhuhr: {time: {hour: 13, minute: 0}, label}};
}

const sunday = new Date("2026-09-06T10:00:00Z");
const saturday = new Date("2026-09-05T10:00:00Z");
assert.deepEqual(
    JSON.parse(JSON.stringify(specialDhuhrItem(board("(Sundays & Public Holidays)"), sunday))),
    {label: "Dhuhr (Sundays & Public Holidays)", time: {hour: 13, minute: 0}},
);
assert.equal(specialDhuhrItem(board("(Sundays & Public Holidays)"), saturday), null);
assert.equal(specialDhuhrItem(board("(Everyday)"), saturday).time.hour, 13);
assert.equal(specialDhuhrItem(board("Daily"), saturday).time.hour, 13);
assert.equal(specialDhuhrItem(board("(Public Holidays)"), sunday), null, "public holidays require a future calendar source");
assert.equal(specialDhuhrItem({special_dhuhr: {label: "(Everyday)"}}, saturday), null);
assert.equal(specialDhuhrItem({special_dhuhr: {time: {hour: 13, minute: 0}}}, saturday), null);

const offsetSunday = new Date("2026-09-05T14:30:00Z");
assert.equal(specialDhuhrItem(board("Sunday", "GMT+10:00"), offsetSunday).time.hour, 13, "weekday must use the board offset");

for (const source of [applianceSource, landscapeSource]) {
    assert.match(source, /specialDhuhrItem\(board,/);
    assert.match(source, /add\([^\n]*specialDhuhr\.label, specialDhuhr\.time\)/);
}

console.log("MasjidBoard special Dhuhr tests passed");
