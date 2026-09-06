"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

global.window = {};
require("../frontend/masjidboard-warning-utils.js");

const {isZawaalWarningActive, zawaalWindow} = window.MasjidBoardWarningUtils;
const board = {
    time_zone: "Africa/Johannesburg",
    astronomical: {
        istiwa_caution: {hour: 12, minute: 1},
        istiwa: {hour: 12, minute: 6},
        zawaal_end: {hour: 12, minute: 11},
    },
};

assert.equal(isZawaalWarningActive(board, new Date("2026-09-05T10:00:59Z")), false, "warning must be off before the interval");
assert.equal(isZawaalWarningActive(board, new Date("2026-09-05T10:01:00Z")), true, "warning must begin at Istiwa caution");
assert.equal(isZawaalWarningActive(board, new Date("2026-09-05T10:10:59Z")), true, "warning must remain active through the final minute");
assert.equal(isZawaalWarningActive(board, new Date("2026-09-05T10:11:00Z")), false, "warning end must be exclusive");

const istiwaFallback = {
    time_zone: "UTC",
    astronomical: {istiwa: {hour: 12, minute: 6}, zawaal_end: {hour: 12, minute: 11}},
};
assert.equal(isZawaalWarningActive(istiwaFallback, new Date("2026-09-05T12:06:00Z")), true, "Istiwa must be the start fallback");
assert.deepEqual(zawaalWindow({istiwa: {hour: 12, minute: 6}}), {
    start: {hour: 12, minute: 1},
    end: {hour: 12, minute: 11},
}, "Istiwa alone must derive a five-minute boundary on each side");
assert.deepEqual(zawaalWindow({
    istiwa_caution: {hour: 11, minute: 59},
    istiwa: {hour: 12, minute: 6},
}), {
    start: {hour: 11, minute: 59},
    end: {hour: 12, minute: 11},
}, "a supplied caution boundary must be retained while the missing end is derived");
assert.deepEqual(zawaalWindow({
    istiwa: {hour: 12, minute: 6},
    zawaal_end: {hour: 12, minute: 14},
}), {
    start: {hour: 12, minute: 1},
    end: {hour: 12, minute: 14},
}, "a missing caution boundary must be derived while the supplied end is retained");
const suppliedBoundaries = {
    istiwa_caution: {hour: 12, minute: 0},
    zawaal_end: {hour: 12, minute: 12},
};
assert.deepEqual(zawaalWindow(suppliedBoundaries), {
    start: suppliedBoundaries.istiwa_caution,
    end: suppliedBoundaries.zawaal_end,
}, "two supplied boundaries must work without Istiwa");
const onlyIstiwa = {time_zone: "UTC", astronomical: {istiwa: {hour: 12, minute: 6}}};
assert.equal(isZawaalWarningActive(onlyIstiwa, new Date("2026-09-05T12:00:00Z")), false);
assert.equal(isZawaalWarningActive(onlyIstiwa, new Date("2026-09-05T12:01:00Z")), true);
assert.equal(isZawaalWarningActive(onlyIstiwa, new Date("2026-09-05T12:11:00Z")), false);
assert.equal(isZawaalWarningActive({astronomical: {istiwa_caution: {hour: 12, minute: 6}, zawaal_end: {hour: 12, minute: 6}}}, new Date()), false, "a zero-length range must not warn");
assert.equal(isZawaalWarningActive({astronomical: {istiwa_caution: {hour: 12, minute: 6}, zawaal_end: {hour: 12, minute: 5}}}, new Date()), false, "a reversed range must not warn");

const root = path.resolve(__dirname, "..");
const displayJS = fs.readFileSync(path.join(root, "frontend/masjidboard-display.js"), "utf8");
const applianceJS = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");
const detailedJS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");
const displayCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-display.css"), "utf8");
const html = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");
assert.match(displayJS, /latestView\.boards\) \? latestView\.boards\[0\]/, "only the primary masjid must control the warning");
assert.match(displayJS, /classList\.toggle\("zawaal-warning-active"/, "the warning state must refresh with the display clock");
assert.match(applianceJS, /warningUtils\.zawaalWindow\(astronomical\)/, "Appliance must display the shared Zawaal window");
assert.match(detailedJS, /warningUtils\.zawaalWindow\(astronomical\)/, "Landscape must display the shared Zawaal window");
for (const selector of ["current-time", "detailed-gregorian-date", "detailed-islamic-date", "appliance-clock", "appliance-information time", "appliance-information div"]) {
    assert.match(displayCSS, new RegExp(selector.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")), `${selector} must receive the warning style`);
}
assert.match(displayCSS, /prefers-reduced-motion:reduce/, "the pulse must respect reduced-motion preferences");
assert.match(html, /masjidboard-warning-utils\.js\?v=20260905-zawaal-fallbacks/, "the warning utility must be loaded with a fresh cache key");

console.log("MasjidBoard Zawaal warning tests passed");
