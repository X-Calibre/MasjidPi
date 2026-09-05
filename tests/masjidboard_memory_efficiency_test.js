"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const display = fs.readFileSync(path.join(root, "frontend/masjidboard-display.js"), "utf8");
const detailed = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");
const html = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");

assert.match(display, /const refreshIntervalMs = 30_000/);
assert.match(display, /"If-None-Match": displayETag/);
assert.match(display, /response\.status === 304/);
assert.match(display, /response\.headers\.get\("ETag"\)/);

assert.match(detailed, /let communityPageNodes = \[\]/);
assert.match(detailed, /function buildCommunityPageNodes\(\)/);
assert.match(detailed, /communityCards\.replaceChildren\(\.\.\.communityPageNodes\[pageIndex\]\)/);
assert.equal(
    (detailed.match(/renderCommunityCard\(entry\.item, entry\.span\)/g) || []).length,
    1,
    "notice cards must be built only when the notice-page cache changes",
);

assert.match(html, /masjidboard-display\.js\?v=20260905-zawaal-warning/);
assert.match(html, /masjidboard-detailed\.js\?v=20260905-daily-times-fit/);

console.log("MasjidBoard memory efficiency tests passed");
