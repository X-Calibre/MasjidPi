"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const frontend = path.join(__dirname, "..", "frontend");
const html = fs.readFileSync(path.join(frontend, "masjidboard-config.html"), "utf8");
const display = fs.readFileSync(path.join(frontend, "masjidboard-layout-config.js"), "utf8");
const boards = fs.readFileSync(path.join(frontend, "masjidboard-config.js"), "utf8");

assert.doesNotMatch(html, /id="saveDisplayLayoutButton"/, "HDMI settings should not require a save button");
assert.doesNotMatch(html, /id="saveBoardsButton"/, "selected masjids should not require a save button");
assert.match(html, /id="saveLocationsButton"/, "location scope retains its deliberate save boundary");
assert.match(html, /id="displaySaveStatus"[^>]*aria-live="polite"/, "display autosave status must be announced");
assert.match(html, /id="boardSaveStatus"[^>]*aria-live="polite"/, "masjid autosave status must be announced");
assert.match(html, /id="locationSaveStatus"[^>]*aria-live="polite"/, "unsaved location status must be announced");
for (const id of ["showDailyAyah", "showDailyHadith", "showDailySunnah"]) {
    assert.match(html, new RegExp(`id="${id}"[^>]*checked`), `${id} must be enabled by default`);
}

assert.match(display, /slideDuration\.addEventListener\("change", saveAutomatically\)/, "slide duration saves only when adjustment is committed");
assert.match(display, /showEconomicIndicators\.addEventListener\("change", saveAutomatically\)/, "indicator visibility saves immediately");
assert.match(display, /Object\.values\(dailyContentInputs\).*addEventListener\("change", saveAutomatically\)/, "daily content visibility saves immediately");
assert.match(display, /show_daily_ayah:\s*!state \|\| state\.show_daily_ayah !== false/, "missing Ayah preference must default to enabled");
assert.match(display, /while \(savePending\)/, "display saves must be serialized");

assert.match(boards, /window\.setTimeout\(\(\) => \{ void drainBoardSaves\(\); \}, 450\)/, "reordering saves must be debounced");
assert.match(boards, /selectedBoards = lastSavedBoards\.map/, "failed masjid saves must roll back");
assert.match(boards, /show_detailed_jumuah:\s*true/, "newly selected masjids must enable the detailed Jumu'ah slide by default");
assert.match(boards, /board\.show_detailed_jumuah !== false/, "missing detailed Jumu'ah preference must default to enabled");
assert.match(boards, /detailed_jumuah:\s*Object\.fromEntries/, "per-masjid Jumu'ah preferences must be saved");
assert.match(boards, /markLocationsUnsaved\(\)/, "location edits must expose their unsaved state");

console.log("MasjidBoard autosave tests passed");
