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

const context = {window: {}, DOMParser: TestDOMParser};
vm.runInNewContext(utilitySource, context);
const {specialDhuhrItem} = context.window.MasjidBoardCommunityUtils;

function board(label, timeZone = "Africa/Johannesburg") {
    return {
        time_zone: timeZone,
        prayers: [{key: "dhuhr", adhan: {hour: 12, minute: 30}, jamaah: {hour: 13, minute: 30}}],
        special_dhuhr: {time: {hour: 13, minute: 0}, label},
    };
}

assert.deepEqual(
    JSON.parse(JSON.stringify(specialDhuhrItem(board("(Sundays & Public Holidays)")))),
    {label: "Dhuhr (Sundays & Public Holidays)", time: {hour: 13, minute: 0}},
);
assert.equal(specialDhuhrItem(board("(Everyday)")).time.hour, 13);
assert.equal(specialDhuhrItem(board("(Public Holidays)")).time.hour, 13, "the label explains applicability without date inference");
assert.equal(specialDhuhrItem({special_dhuhr: {label: "(Everyday)"}}), null);
assert.equal(specialDhuhrItem({special_dhuhr: {time: {hour: 13, minute: 0}}}), null);
const duplicateJamaah = board("(Everyday)");
duplicateJamaah.prayers[0].jamaah = {hour: 13, minute: 0};
assert.equal(specialDhuhrItem(duplicateJamaah), null, "a normal Jamaah duplicate must be hidden");
const duplicateAdhan = board("(Everyday)");
duplicateAdhan.prayers[0].adhan = {hour: 13, minute: 0};
assert.equal(specialDhuhrItem(duplicateAdhan), null, "a normal Adhan duplicate must be hidden");

for (const source of [applianceSource, landscapeSource]) {
    assert.match(source, /specialDhuhrItem\(board\)/);
    assert.match(source, /add\([^\n]*specialDhuhr\.label, specialDhuhr\.time\)/);
}

console.log("MasjidBoard special Dhuhr tests passed");
