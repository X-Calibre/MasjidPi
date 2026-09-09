"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const read = relative => fs.readFileSync(path.join(root, relative), "utf8");
const launcher = read("frontend/appliance-720-demo.html");
const boardHTML = read("frontend/masjidboard.html");
const appliance = read("frontend/masjidboard-appliance.js");
const controls = read("frontend/masjidboard-touch-controls.js");
const toast = read("frontend/masjidboard-listen-toast.js");
const setupHTML = read("frontend/setup.html");
const setup = read("frontend/setup.js");

assert.match(launcher, /width:720px/);
assert.match(launcher, /height:1280px/);
assert.match(launcher, /visual-demo=cards/);
assert.match(launcher, /dua-fixture=1/);
assert.match(launcher, /jumuah-fixture=khateeb/);
for (const tab of ["masjid", "radio", "theme", "network", "display"]) {
    assert.ok(launcher.includes("visual-demo-controls=" + tab));
}
for (const state of ["masjid", "radio", "waiting", "cycle"]) {
    assert.ok(launcher.includes("visual-demo-toast=" + state));
}
for (const step of ["network", "password", "hidden", "success", "location", "picker", "masjid"]) {
    assert.ok(launcher.includes("visual-demo-step=" + step));
}

assert.match(appliance, /params\.get\("visual-demo"\) === "cards"/);
assert.match(appliance, /communityFixtureMode = visualDemoCards \? "1"/);
assert.match(appliance, /\|\| visualDemoCards\) return/);
assert.match(controls, /params\.get\("visual-demo-controls"\)/);
assert.match(controls, /function visualDemoResponse/);
assert.match(controls, /if \(visualDemoTab\) return/);
assert.match(controls, /\/setup\.html\?visual-demo-step=network/);
assert.match(toast, /get\("visual-demo-toast"\)/);
assert.match(toast, /visualDemoStates/);
assert.match(toast, /setInterval/);
assert.match(setup, /get\("visual-demo-step"\)/);
assert.match(setup, /showVisualDemoStep/);
assert.match(setup, /event\.stopImmediatePropagation\(\)/);
assert.match(setup, /document\.addEventListener\("click", event =>/);
assert.match(boardHTML, /masjidboard-appliance\.js\?v=20260909-visual-demo/);
assert.match(boardHTML, /masjidboard-listen-toast\.js\?v=20260909-visual-demo/);
assert.match(boardHTML, /masjidboard-touch-controls\.js\?v=20260909-visual-demo/);
assert.match(setupHTML, /setup\.js\?v=20260909-visual-demo/);

console.log("Touch Display 2 visual-demo tests passed");
