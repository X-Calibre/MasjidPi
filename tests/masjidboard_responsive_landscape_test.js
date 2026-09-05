"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const detailedCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.css"), "utf8");
const detailedJS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");
const configHTML = fs.readFileSync(path.join(root, "frontend/masjidboard-config.html"), "utf8");
const configJS = fs.readFileSync(path.join(root, "frontend/masjidboard-layout-config.js"), "utf8");

assert.match(detailedCSS, /html\.landscape-layout\s*\{[^}]*font-size:\s*clamp\(11px,min\(\.833333vw,1\.481481vh\),32px\)/s);
assert.match(detailedJS, /params\.get\("profile"\) === "appliance"/);
assert.match(detailedJS, /document\.documentElement\.classList\.add\("landscape-layout"\)/);
assert.match(detailedCSS, /grid-template-columns:\s*clamp\(120px,9vw,10\.3125rem\) minmax\(0,1fr\)/);
assert.match(detailedCSS, /repeat\(var\(--daily-time-columns,10\),minmax\(0,1fr\)\)/, "Daily Times must use one dynamically sized column per item");
assert.match(detailedCSS, /\.additional-times-list\.is-crowded/, "crowded Daily Times rows must reduce spacing and typography");
assert.match(detailedJS, /--daily-time-columns/, "Landscape must expose the actual Daily Times item count to CSS");
assert.match(detailedJS, /items\.length > 10/, "the eleventh Daily Times item must activate compact sizing");
assert.match(detailedCSS, /@media \(min-width:1101px\)/);
assert.doesNotMatch(detailedCSS, /@media \(min-width:1101px\) and \(max-width:2000px\)/);
assert.match(detailedCSS, /@media \(max-width:1500px\)[^}]*\.landscape-layout \.time-value-stack:has\(\.event-countdown\) \{ gap:0; \}/s);
assert.match(configHTML, /Local display profile/);
assert.match(configHTML, /The 7-inch Waveshare appliance display is detected at startup/);
assert.match(configHTML, /masjidboard\.html\?profile=appliance/);
assert.doesNotMatch(configHTML, /id="displayLayout"/);
assert.doesNotMatch(configJS, /portrait/);
assert.doesNotMatch(configJS, /landscape/);

console.log("MasjidBoard responsive landscape tests passed");
