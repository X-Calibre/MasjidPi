"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const html = fs.readFileSync(path.join(root, "frontend/index.html"), "utf8");
const app = fs.readFileSync(path.join(root, "frontend/app.js"), "utf8");

assert.match(html, /id="refreshAudioDevices"/);
assert.match(app, /device\.unavailable \? " — unavailable" : ""/);
assert.match(app, /option\.disabled = Boolean\(device\.unavailable\)/);
assert.match(app, /refreshAudioDevicesButton\.addEventListener\("click"/);
assert.match(app, /await loadAudioDevices\(\{throwOnError: true\}\);\s*await refreshStatus\(\);/);
assert.match(app, /Audio devices refreshed\./);

console.log("Audio device refresh UI tests passed");
