"use strict";

const assert = require("node:assert/strict");

global.window = {};
require("../frontend/masjidboard-date-utils.js");

const dates = window.MasjidBoardDate;
const board = {astronomical: {sunset: {hour: 18, minute: 0}}};

assert.equal(
    dates.isIslamicFriday(board, new Date(2026, 7, 27, 18, 3, 5)),
    false,
    "Thursday remains Thursday at the exact Islamic-date rollover boundary",
);
assert.equal(
    dates.isIslamicFriday(board, new Date(2026, 7, 27, 18, 3, 6)),
    true,
    "Jumu'ah begins immediately after Thursday's Islamic-date rollover",
);
assert.equal(
    dates.isIslamicFriday(board, new Date(2026, 7, 28, 18, 3, 5)),
    true,
    "Jumu'ah remains active at Friday's exact rollover boundary",
);
assert.equal(
    dates.isIslamicFriday(board, new Date(2026, 7, 28, 18, 3, 6)),
    false,
    "Dhuhr returns immediately after Friday's Islamic-date rollover",
);

assert.equal(
    dates.islamicWeekday(board, new Date(2026, 7, 27, 18, 3, 6)),
    "Al-Jumu'ah",
);
assert.equal(
    dates.islamicWeekday(board, new Date(2026, 7, 28, 18, 3, 6)),
    "As-Sabt",
);

console.log("MasjidBoard Islamic date utility tests passed");
