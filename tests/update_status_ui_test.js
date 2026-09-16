"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const script = fs.readFileSync(path.join(root, "frontend/update-status.js"), "utf8");
const index = fs.readFileSync(path.join(root, "frontend/index.html"), "utf8");
const board = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");

for (const html of [index, board]) {
    assert.match(html, /data-update-installation/);
    assert.match(html, /data-update-install/);
}

assert.match(script, /\/api\/update\/install/);
assert.match(script, /interrupt_playback:\s*true/);
assert.match(script, /window\.confirm/);
assert.match(script, /reboot_pending/);
assert.match(script, /probation/);
assert.match(script, /installed/);
assert.match(script, /rolled_back/);
assert.match(script, /download\?\.status === "verified"/);

console.log("Update installation UI tests passed");
