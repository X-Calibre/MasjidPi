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
assert.match(detailedJS, /document\.documentElement\.classList\.add\("landscape-layout"\)/);
assert.match(detailedCSS, /grid-template-columns:\s*clamp\(120px,9vw,10\.3125rem\) minmax\(0,1fr\)/);
assert.match(detailedCSS, /@media \(min-width:1101px\)/);
assert.doesNotMatch(detailedCSS, /@media \(min-width:1101px\) and \(max-width:2000px\)/);
assert.match(configHTML, /Landscape \(responsive\)/);
assert.match(configJS, /Responsive landscape/);
assert.doesNotMatch(configHTML, /Landscape \(1920 × 1080\)/);

console.log("MasjidBoard responsive landscape tests passed");
