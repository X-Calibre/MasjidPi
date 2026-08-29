"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const {remainingText, sourceState} = require("../frontend/masjidboard-listen-toast.js");

const html = fs.readFileSync(path.resolve(__dirname, "../frontend/masjidboard.html"), "utf8");
assert.equal((html.match(/data-board-source-toast/g) || []).length, 2);
assert.match(html, /masjidboard-listen-toast\.js\?v=20260829-source-transition-toast/);

assert.deepEqual(sourceState(null), {kind: "none", key: "none"});
assert.deepEqual(sourceState({listening: false}), {kind: "none", key: "none"});

const masjid = sourceState({
    listening: true,
    active_source: "masjid",
    active_stream_id: "masjid-a",
    active_stream_name: "Masjid us Salaam",
});
assert.equal(masjid.kind, "masjid");
assert.equal(masjid.name, "Masjid us Salaam");

const radio = sourceState({
    listening: true,
    active_source: "radio",
    active_stream_id: "radio-a",
    radio_name: "Radio Islam International",
});
assert.equal(radio.kind, "radio");
assert.equal(radio.name, "Radio Islam International");

const waiting = sourceState({
    listening: true,
    active_source: "none",
    radio_id: "radio-a",
    radio_name: "Radio Islam International",
    radio_resume_pending: true,
    radio_resume_at: "2026-08-29T10:05:00Z",
});
assert.equal(waiting.kind, "waiting");
assert.equal(waiting.name, "Radio Islam International");
assert.equal(remainingText("2026-08-29T10:05:00Z", Date.parse("2026-08-29T10:00:04Z")), "4:56");

console.log("MasjidBoard Listen toast tests passed");
