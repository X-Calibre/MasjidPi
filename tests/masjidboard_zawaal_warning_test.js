"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

global.window = {};
require("../frontend/masjidboard-warning-utils.js");

const {isZawaalWarningActive} = window.MasjidBoardWarningUtils;
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
assert.equal(isZawaalWarningActive({astronomical: {istiwa: {hour: 12, minute: 6}}}, new Date()), false, "an incomplete range must not warn");
assert.equal(isZawaalWarningActive({astronomical: {istiwa: {hour: 12, minute: 6}, zawaal_end: {hour: 12, minute: 6}}}, new Date()), false, "a zero-length range must not warn");
assert.equal(isZawaalWarningActive({astronomical: {istiwa: {hour: 12, minute: 6}, zawaal_end: {hour: 12, minute: 5}}}, new Date()), false, "a reversed range must not warn");

const root = path.resolve(__dirname, "..");
const displayJS = fs.readFileSync(path.join(root, "frontend/masjidboard-display.js"), "utf8");
const displayCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-display.css"), "utf8");
const html = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");
assert.match(displayJS, /latestView\.boards\) \? latestView\.boards\[0\]/, "only the primary masjid must control the warning");
assert.match(displayJS, /classList\.toggle\("zawaal-warning-active"/, "the warning state must refresh with the display clock");
for (const selector of ["current-time", "detailed-gregorian-date", "detailed-islamic-date", "appliance-clock", "appliance-information time", "appliance-information div"]) {
    assert.match(displayCSS, new RegExp(selector.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")), `${selector} must receive the warning style`);
}
assert.match(displayCSS, /prefers-reduced-motion:reduce/, "the pulse must respect reduced-motion preferences");
assert.match(html, /masjidboard-warning-utils\.js\?v=20260905-zawaal-warning/, "the warning utility must be loaded with a fresh cache key");

console.log("MasjidBoard Zawaal warning tests passed");
