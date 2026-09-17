"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const script = fs.readFileSync(path.join(root, "frontend/update-status.js"), "utf8");
const index = fs.readFileSync(path.join(root, "frontend/index.html"), "utf8");
const board = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");
const boardConfig = fs.readFileSync(path.join(root, "frontend/masjidboard-config.html"), "utf8");

const touchPanel = board.match(
    /<section id="applianceUpdatesPanel"[\s\S]*?<\/section>/
)?.[0] || "";

for (const html of [index, board]) {
    assert.match(html, /data-update-install/);
}

assert.match(index, /data-update-installation/);
assert.match(index, /data-update-link/);
assert.match(boardConfig, /data-update-link/);
assert.match(touchPanel, /data-update-installed-version/);
assert.match(touchPanel, /data-update-available/);
assert.match(touchPanel, /data-update-install-date/);
assert.doesNotMatch(touchPanel, /data-update-last-checked/);
assert.doesNotMatch(touchPanel, /data-update-published/);
assert.doesNotMatch(touchPanel, /data-update-deadline/);
assert.doesNotMatch(touchPanel, /data-update-download/);
assert.doesNotMatch(touchPanel, /data-update-installation/);
assert.doesNotMatch(touchPanel, /data-update-link/);
assert.doesNotMatch(touchPanel, /View release notes/);

assert.match(script, /\/api\/update\/install/);
assert.match(script, /data-update-available/);
assert.match(script, /data-update-install-date/);
assert.match(script, /\["probation", "installed"\]/);
assert.match(script, /installation\?\.staged_at/);
assert.match(script, /interrupt_playback:\s*true/);
assert.match(script, /window\.confirm/);
assert.match(script, /reboot_pending/);
assert.match(script, /probation/);
assert.match(script, /installed/);
assert.match(script, /rolled_back/);
assert.match(script, /download\?\.status === "verified"/);

console.log("Update installation UI tests passed");
