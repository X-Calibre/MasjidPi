"use strict";
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const source = fs.readFileSync(require('node:path').join(__dirname, '../frontend/masjidboard-detailed.js'), 'utf8');
const start = source.indexOf('    function isDetailedCommunityItem(');
const end = source.indexOf('    function renderCommunityCard(', start);
const context = {plainText: value => String(value || ''), orderedFields: item => Object.values(item.fields || {})};
vm.createContext(context);
vm.runInContext(source.slice(start, end), context);
const daily = ['daily_ayah', 'daily_hadith', 'daily_sunnah'].map((type, i) => ({type, body: 'Text '.repeat(i === 0 ? 100 : 10)}));
const pages = context.packCommunityPages([...daily, {type:'economic'}]);
assert.equal(pages.length, 2, 'Long Ayah must remain with Hadith and Sunnah');
assert.equal(pages[0].layout, 'daily');
assert.deepEqual(Array.from(pages[0].entries, entry => entry.item.type), daily.map(item => item.type));
assert.equal(pages[1].entries[0].item.type, 'economic', 'Economic indicators retain their own slide');
assert.equal(context.packCommunityPages([daily[1]])[0].entries.length, 1, 'Missing daily items create no empty cards');
const mixed = context.packCommunityPages([{type:'announcement', body:'Notice'}, ...daily]);
assert.equal(mixed[0].entries[0].item.type, 'announcement');
assert.equal(mixed[1].entries.length, 3, 'Daily content is grouped separately from notices');
console.log('MasjidBoard daily landscape grouping tests passed');
// Exercise the actual sizing callback against a card with measurable limits.
let callback;
let measurements = 0;
const body = {style:{}, clientHeight:100, clientWidth:400};
const text = {get scrollHeight(){ measurements++; return Number.parseInt(body.style.fontSize) * 4; }, scrollWidth:300};
const card = {clientHeight:200, scrollHeight:200, querySelector: selector => selector === '.detailed-daily-text' ? text : body};
const fitContext = {
    window:{requestAnimationFrame: fn => {callback = fn;}, addEventListener:()=>{}},
    document:{}, communityCards:{clientHeight:300,clientWidth:400,querySelectorAll:()=>[card]},
    communityPage:0,communitySignature:'content',communityPages:[{layout:'daily',entries:[{item:daily[0]}]}],
};
vm.createContext(fitContext);
const fitStart = source.indexOf('    let dailyFitPending = false;');
const fitEnd = source.indexOf('    function buildCommunityPageNodes(', fitStart);
vm.runInContext(source.slice(fitStart, fitEnd), fitContext);
fitContext.scheduleDailyFit(); callback();
assert.equal(body.style.fontSize, '25px', 'Choose the largest font that fits the measured height');
const measured = measurements;
fitContext.scheduleDailyFit(); callback();
assert.equal(measurements, measured, 'Unchanged clock refreshes must not measure text again');
body.clientHeight=128;
fitContext.communityCards.clientHeight=400;
fitContext.scheduleDailyFit(); callback();
assert.equal(body.style.fontSize, '32px', 'Resize must allow larger text');
console.log('MasjidBoard daily text sizing and measurement caching tests passed');
