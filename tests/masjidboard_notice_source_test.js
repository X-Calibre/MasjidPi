"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const landscape = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");
const appliance = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");

assert.match(landscape, /`Source: \$\{item\.source\}`/);
assert.match(appliance, /"Source: " \+ item\.source/);
assert.match(landscape, /item\.type === "salaah_change"[\s\S]*?`Source: \$\{item\.source\}`/);
assert.match(appliance, /item\.type === "salaah_change"[\s\S]*?"Source: " \+ item\.source/);
assert.match(landscape, /`From \$\{item\.source\}`/);
assert.match(appliance, /`From \$\{indicators\.source\}`/);
assert.doesNotMatch(appliance, /portrait/);

console.log("MasjidBoard notice source-label tests passed");
