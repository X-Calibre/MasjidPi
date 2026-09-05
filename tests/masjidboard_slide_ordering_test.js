"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

const root = path.resolve(__dirname, "..");
const utilitySource = fs.readFileSync(path.join(root, "frontend/masjidboard-community-utils.js"), "utf8");
const appliance = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");
const landscape = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");

class TestDOMParser {
    parseFromString(value) {
        return {body: {textContent: String(value).replace(/<[^>]*>/g, "")}};
    }
}

const context = {window: {}, DOMParser: TestDOMParser, Intl, Date};
vm.runInNewContext(utilitySource, context);
const {orderedCommunityItemGroups, orderedCommunityItems} = context.window.MasjidBoardCommunityUtils;

const board = {
    name: "Primary Masjid",
    time_zone: "Africa/Johannesburg",
    notices: [
        {type: "nikah", title: "Nikah", fields: {date: "Saturday"}},
        {type: "funeral", title: "Funeral", fields: {name: "Marhoom"}},
        {type: "eid", title: "Eid", fields: {date: "Eid day"}},
        {type: "salaah_change", title: "Fajr Time Change", fields: {effective_date: "5 September 2026", new_time: "05:15"}},
        {type: "salaah_change", title: "Asr Time Change", fields: {effective_date: "12 September 2026", new_time: "16:30"}},
        {type: "salaah_change", title: "Late Change", fields: {effective_date: "13 September 2026", new_time: "16:45"}},
        {type: "salaah_change", title: "Past Change", fields: {effective_date: "4 September 2026", new_time: "16:00"}},
        {type: "salaah_change", title: "Undated Change", fields: {new_time: "17:00"}},
        {type: "well_wishes", title: "Du'a Requested", content: "Please make du'a"},
        {type: "dawah", title: "Dawah and Gasht", fields: {gasht_out_day: "Friday"}},
        {type: "three_day_jamaat", title: "Three-Day Jamaat", fields: {first_location: "Pretoria"}},
    ],
    announcements: [
        {category: "announcement", title: "General Announcement", content: "General"},
        {category: "urgent_announcement", title: "Urgent Access Notice", content: "Use side entrance"},
        {category: "salaah_change_announcement", title: "Salaah Times Change", content: "As of Friday, September 12, 2026"},
        {category: "salaah_change_announcement", title: "Old Salaah Times Change", content: "As of September 1, 2026"},
        {category: "class_time_change", title: "Class Time Change", content: "Starts earlier"},
    ],
    programmes: [{title: "Taleem", content: "After Esha"}],
    new_moon: {fields: {birth_date: "12 September"}},
    banking: {title: "Contributions", fields: {bank: "Example Bank"}},
    jumuah: [{events: [{code: "0", heading: "Adhan", time: {hour: 12, minute: 30}}], jamaah: {hour: 13, minute: 0}}],
};

const now = new Date("2026-09-05T10:00:00Z");
const items = orderedCommunityItems(board, now, () => true);
const types = Array.from(items, (item) => item.type);
assert.deepEqual(Array.from(orderedCommunityItemGroups(board, now, () => true), (group) => group.length), [1, 5, 4, 4, 2]);
assert.deepEqual(types, [
    "funeral",
    "urgent_announcement", "salaah_change", "salaah_change", "salaah_change_announcement", "eid",
    "announcement", "jumuah_schedule", "nikah", "well_wishes",
    "programme", "class_time_change", "dawah", "three_day_jamaat",
    "new_moon", "contribution",
]);
assert.ok(!items.some((item) => item.title === "Late Change"));
assert.ok(!items.some((item) => item.title === "Past Change"));
assert.ok(!items.some((item) => item.title === "Undated Change"));
assert.ok(!items.some((item) => item.title === "Old Salaah Times Change"));

assert.match(appliance, /boards\.forEach\(\(board, boardIndex\)/, "Appliance must build each masjid section in selection order");
assert.match(appliance, /boardIndex === 0 && dailyItems\(board\)/, "Daily Times must belong only to the primary masjid");
assert.match(appliance, /boards\[0\] \? \[boards\[0\]\] : \[\]/, "Dua must inspect only the primary masjid");
assert.ok(appliance.indexOf("dailyIslamicItems(dailyContent)") > appliance.indexOf("boards.forEach((board, boardIndex)"), "shared daily content must follow all masjid sections");
assert.match(landscape, /boards\.flatMap\(\(board, boardIndex\)/, "Landscape must preserve per-masjid priority groups");
assert.match(landscape, /itemGroups\.flatMap\(\(group\) => packCommunityPages\(group\)\)/, "Landscape must not pack cards from different masjids together");
assert.match(landscape, /boards\[0\] \? \[boards\[0\]\] : \[\]/, "Landscape Dua must inspect only the primary masjid");

console.log("MasjidBoard slide-ordering tests passed");
