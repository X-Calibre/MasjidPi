"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

const root = path.resolve(__dirname, "..");
const utilitySource = fs.readFileSync(path.join(root, "frontend/masjidboard-community-utils.js"), "utf8");
const appliance = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.js"), "utf8");
const landscape = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.js"), "utf8");
const applianceCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-appliance.css"), "utf8");
const landscapeCSS = fs.readFileSync(path.join(root, "frontend/masjidboard-detailed.css"), "utf8");

class TestDOMParser {
    parseFromString(source) {
        return {body: {textContent: String(source).replace(/<[^>]*>/g, "")}};
    }
}

const context = {window: {}, DOMParser: TestDOMParser};
vm.runInNewContext(utilitySource, context);
const {dailyIslamicItems, formatNoticeDate} = context.window.MasjidBoardCommunityUtils;

assert.equal(formatNoticeDate("06 Oct 2026 00:00"), "Tuesday, 06 October");
assert.equal(formatNoticeDate("2026-10-06T00:00:00"), "Tuesday, 06 October");

const items = dailyIslamicItems({
    source: "MasjidBoard Live",
    ayah: {surah: "Surah 59 – Al-Hashr", ayah_number: "Ayah 7", text: "Ayah text"},
    hadith: {heading: "Hadith", text: "Hadith text", reference: "Bukhari"},
    sunnah: {heading: "Sunnah of Travelling", text: "Sunnah text", reference: "Muslim"},
});
assert.equal(items.length, 3);
assert.deepEqual(Array.from(items, (item) => item.type), ["daily_ayah", "daily_hadith", "daily_sunnah"]);
assert.ok(items.every((item) => item.source === "MasjidBoard Live"));
assert.ok(items.every((item) => item.typeLabel === undefined), "daily cards must not repeat their inferred content type");
assert.equal(items[0].fields.ayah_number, "Ayah 7");
assert.equal(items[1].fields.reference, "Bukhari");
assert.equal(dailyIslamicItems({source: "MasjidBoard Live", ayah: {text: "Only Ayah"}}).length, 1);

assert.match(appliance, /dailyIslamicItems\(dailyContent\)/);
assert.match(appliance, /slides\.push\(communitySlide\(\[item\]\)\)/, "each appliance item must receive a dedicated slide");
assert.match(appliance, /view\.daily_islamic_content/);
assert.match(landscape, /const sharedItems = dailyIslamicItems\(dailyContent\)/);
assert.match(landscape, /itemGroups\.push\(sharedItems\)/, "landscape shared content must form one final group");
assert.match(landscape, /view\.daily_islamic_content/);
assert.match(appliance, /"Source: " \+ item\.source/);
assert.match(landscape, /`Source: \$\{item\.source\}`/);
assert.match(appliance, /"appliance-daily-ayah-number"/);
assert.match(appliance, /item\.type !== "daily_ayah" \|\| field\.label !== "Ayah"/);
assert.match(landscape, /"detailed-daily-ayah-number"/);
assert.match(landscape, /item\.type !== "daily_ayah" \|\| field\.label !== "Ayah"/);
assert.match(appliance, /title\.dir = "auto"/);
assert.match(landscape, /title\.dir = "auto"/);
assert.match(applianceCSS, /\.appliance-community-daily_ayah/);
assert.match(applianceCSS, /\.appliance-daily-ayah-number/);
assert.match(applianceCSS, /\.content-very-long/);
assert.match(landscapeCSS, /\.detailed-community-daily_ayah/);
assert.match(landscapeCSS, /\.detailed-daily-ayah-number/);
assert.match(landscapeCSS, /\.content-very-long/);

console.log("MasjidBoard daily Islamic content tests passed");
