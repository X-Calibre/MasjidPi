"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const html = fs.readFileSync(path.join(root, "frontend/masjidboard.html"), "utf8");
const css = fs.readFileSync(path.join(root, "frontend/masjidboard-portrait.css"), "utf8");
const themes = fs.readFileSync(path.join(root, "frontend/masjidboard-themes.css"), "utf8");
const controller = fs.readFileSync(path.join(root, "frontend/masjidboard-touch-controls.js"), "utf8");
const portrait = fs.readFileSync(path.join(root, "frontend/masjidboard-portrait.js"), "utf8");

assert.match(html, /id="portraitListenPanel"/);
assert.match(html, /data-touch-tab="masjid"/);
assert.match(html, /data-touch-tab="radio"/);
assert.match(html, /data-touch-tab="theme"/);
assert.match(html, /Scheduled Play/);
assert.match(html, /Play Now/);
assert.match(html, /Stop Radio/);
assert.match(html, /Master Volume/);
assert.match(html, /Masjid Volume/);
assert.match(html, /Radio Volume/);
assert.match(controller, /dy < -70/);
assert.match(controller, /dy > 70/);
assert.match(controller, /\/api\/favourites/);
assert.match(controller, /\/api\/listen\/radio-mode/);
assert.match(controller, /\/api\/masjidboard\/layout/);
assert.match(controller, /button\.dataset\.sourceId = item\.id/);
assert.match(controller, /selectSource\("masjid", button\.dataset\.sourceId\)/);
assert.match(controller, /selectSource\("radio", button\.dataset\.sourceId\)/);
assert.doesNotMatch(controller, /dataset\.sourceID/);
assert.match(controller, /control\.addEventListener\("input", \(\) => scheduleVolumeSave\(name\)\)/);
assert.match(controller, /window\.setTimeout\(\(\) => saveVolume\(name, value\), 120\)/);
assert.match(controller, /pendingVolumes\[name\] \?\? values\[name\]/);
assert.match(controller, /\{volume:value, persist:true\}/);
assert.match(controller, /\["black-white", "Black & White"/);
assert.match(portrait, /masjidpi:portrait-listen-panel/);
assert.match(css, /var\(--portrait-panel\)/);
assert.match(css, /\.portrait-listen-sheet\s*\{[^}]*height:724px/s);
assert.match(css, /\.portrait-listen-content\s*\{[^}]*flex:1 1 auto;[^}]*min-height:0/s);
assert.match(themes, /--portrait-danger:var\(--danger\)/);

console.log("MasjidBoard touch-control tests passed");
