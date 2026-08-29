"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const html = fs.readFileSync(path.join(root, "frontend/index.html"), "utf8");
const app = fs.readFileSync(path.join(root, "frontend/app.js"), "utf8");

const nowPlayingStart = html.indexOf('<section class="card now-playing-card">');
const nowPlayingEnd = html.indexOf("</section>", nowPlayingStart);
const countdown = html.indexOf('id="radioResumeDelayStatus"');
const delayControl = html.indexOf('id="radioResumeDelaySlider"');

assert.ok(nowPlayingStart >= 0 && countdown > nowPlayingStart && countdown < nowPlayingEnd);
assert.ok(countdown < delayControl);
assert.match(app, /if \(status\.radio_resume_pending\)/);
assert.match(app, /Radio will resume after the configured delay\./);

console.log("Radio resume status UI tests passed");
