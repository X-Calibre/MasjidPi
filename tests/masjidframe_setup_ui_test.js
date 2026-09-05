"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const startup = fs.readFileSync(path.join(root, "frontend/masjidboard-startup.html"), "utf8");
const html = fs.readFileSync(path.join(root, "frontend/setup.html"), "utf8");
const css = fs.readFileSync(path.join(root, "frontend/setup.css"), "utf8");
const js = fs.readFileSync(path.join(root, "frontend/setup.js"), "utf8");

assert.match(startup, /http:\/\/127\.0\.0\.1:8080\/appliance/);
assert.match(html, /Set up MasjidFrame/);
assert.match(html, /Choose your 2\.4 GHz Wi-Fi network/);
assert.match(html, /id="wifiPassword"[^>]*readonly/);
assert.match(html, /id="keyboard"[^>]*aria-label="On-screen keyboard"/);
assert.match(css, /@media \(orientation: landscape\)/);
assert.match(css, /\.key\s*\{[^}]*height: 54px/s);
assert.match(js, /\/api\/setup\/wifi\/networks/);
assert.match(js, /\/api\/setup\/wifi\/connect/);
assert.match(js, /const symbolRows/);
assert.match(js, /const letterRows = \[\s*\["1", "2", "3", "4", "5", "6", "7", "8", "9", "0"\]/s);
assert.match(js, /symbols: "#\+="/);
assert.match(js, /password\.value = Array\.from\(password\.value\)\.slice\(0, -1\)/);
assert.doesNotMatch(js, /innerHTML\s*=.*network\.ssid/);

console.log("MasjidFrame first-run setup UI tests passed");
