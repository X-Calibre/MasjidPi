"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const read = file => fs.readFileSync(path.join(root, file), "utf8");
const html = read("frontend/masjidboard.html");
const css = read("frontend/masjidboard-appliance.css");
const appliance720 = read("frontend/masjidboard-appliance-720.css");
const controller = read("frontend/masjidboard-touch-controls.js");
const appliance = read("frontend/masjidboard-appliance.js");
const themes = read("frontend/masjidboard-themes.css");
const themeController = read("frontend/masjidboard-theme.js");
const config = read("frontend/masjidboard-config.html");
const layoutController = read("frontend/masjidboard-layout-config.js");
const setup = read("frontend/setup.html");

assert.match(html, /id="applianceQuickPanel"/);
assert.match(html, /id="applianceListenPanel"/);
assert.match(html, /id="applianceQuickConnection"/);
assert.match(html, /QUICK SETTINGS/);
assert.match(html, /Display Brightness/);
for (const label of ["Master Volume", "Masjid Volume", "Radio Volume"]) assert.match(html, new RegExp(label));
assert.equal((html.match(/data-radio-mode="schedule"/g) || []).length, 2);
assert.equal((html.match(/data-radio-mode="play_now"/g) || []).length, 2);
assert.equal((html.match(/data-radio-mode="stopped"/g) || []).length, 2);
for (const tab of ["masjid", "radio", "theme", "network"]) assert.match(html, new RegExp(`data-touch-tab="${tab}"`));
assert.doesNotMatch(html, /data-touch-tab="display"/);
assert.doesNotMatch(html, /applianceDisplayPanel|White balance|data-temperature|display-adjustments\.js/);
assert.doesNotMatch(setup, /display-adjustments\.js/);
assert.match(html, /masjidboard-touch-controls\.js\?v=20260909-split-panels/);
assert.match(html, /masjidboard-appliance-720\.css\?v=20260909-split-panels/);
assert.match(html, /masjidboard-appliance\.css\?v=20260909-split-panels/);

assert.match(controller, /profile !== "appliance-720"/);
assert.match(controller, /start\.y <= 120 && dy > 70/);
assert.match(controller, /start\.y >= start\.height - 120 && dy < -70/);
assert.match(controller, /setOpenPanel\("quick"\)/);
assert.match(controller, /setOpenPanel\("bottom"\)/);
assert.match(controller, /wireCloseGesture\(quickPanel,[^\n]*,-1\)/);
assert.match(controller, /wireCloseGesture\(bottomPanel,[^\n]*,1\)/);
assert.match(controller, /querySelectorAll\("\[data-radio-mode\]"\)/);
assert.match(controller, /\/api\/listen\/radio-mode/);
assert.match(controller, /\/api\/display\/settings/);
assert.doesNotMatch(controller, /color_temperature|data-temperature|MasjidPiDisplaySettings/);
assert.match(controller, /const inactivityTimeout = 60000/);
assert.match(controller, /setOpenPanel\(""\)/);
assert.match(controller, /return=board&profile=\$\{profile\}/);
assert.match(controller, /\/api\/favourites/);
assert.match(controller, /\/api\/masjidboard\/layout/);
assert.match(controller, /\/api\/setup\/device-access/);
assert.match(controller, /pendingVolumes\[name\] \?\? values\[name\]/);
assert.match(controller, /\{volume:value,persist:true\}/);
assert.match(appliance, /masjidpi:appliance-listen-panel/);

assert.match(css, /\.appliance-quick-sheet\s*\{[^}]*height:66\.667%/s);
assert.match(css, /\.appliance-quick-controls\s*\{[^}]*grid-template-columns:1fr 1fr/s);
assert.match(css, /\.appliance-listen-tabs\s*\{[^}]*grid-template-columns:repeat\(4,1fr\)/s);
assert.match(appliance720, /\.appliance-quick-sheet\s*\{[^}]*height:66\.667%/s);
assert.match(appliance720, /\.appliance-quick-controls\s*\{[^}]*grid-template-columns:1fr/s);
assert.match(appliance720, /\.appliance-listen-sheet\s*\{[^}]*height:790px/s);
assert.match(appliance720, /\.appliance-listen-tabs\s*\{[^}]*repeat\(4,1fr\)/s);
assert.doesNotMatch(appliance720, /appliance-temperature|appliance-display-settings/);

assert.match(appliance720, /body\.appliance-720-layout/);
assert.match(appliance720, /width:720px/);
assert.match(appliance720, /height:1280px/);
assert.match(appliance720, /\.appliance-clock\s*\{[^}]*font-size:104px/s);
assert.match(appliance720, /\.appliance-information strong\s*\{[^}]*font-size:40px;[^}]*white-space:normal/s);
assert.match(appliance720, /strong\.name-long\s*\{[^}]*font-size:34px/s);
assert.match(appliance720, /strong\.name-very-long\s*\{[^}]*font-size:29px/s);
assert.match(appliance720, /\.appliance-community-body\s*\{[^}]*font-size:34px/s);
assert.match(appliance720, /\.appliance-community-field\s*\{[^}]*font-size:28px/s);
assert.match(appliance720, /\.appliance-next-name\s*\{[^}]*font-size:32px/s);
assert.match(appliance720, /\.appliance-countdown\s*\{[^}]*font-size:30px/s);
assert.match(appliance720, /@keyframes appliance-720-toast-drop/);

for (const theme of ["ivory", "sage", "sky", "rose"]) {
    assert.match(themeController, new RegExp(`"${theme}"`));
    assert.match(layoutController, new RegExp(`"${theme}"`));
    assert.match(config, new RegExp(`name="boardTheme" value="${theme}"`));
    assert.match(themes, new RegExp(`body\\[data-board-theme="${theme}"\\]`));
}

console.log("MasjidBoard appliance touch-control tests passed");
