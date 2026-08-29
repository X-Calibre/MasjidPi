"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const css = fs.readFileSync(path.join(root, "frontend/listen.css"), "utf8");

assert.match(css, /\.module-power-control\s*\{[^}]*border-radius:\s*999px/s);
assert.match(css, /\.module-power-control input \+ \.power-toggle::after\s*\{[^}]*background:\s*#fff/s);
assert.match(css, /\.module-power-control:has\(input:focus-visible\)/);
assert.doesNotMatch(css, /--card-bg/);

console.log("Module power UI tests passed");
