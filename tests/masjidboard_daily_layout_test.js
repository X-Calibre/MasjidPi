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
    document:{}, communityCards:{clientHeight:300,clientWidth:400,querySelectorAll:selector=>selector === ".detailed-daily-card" ? [card] : []},
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
const makeDuaRow = (arabic, height, multiplier) => {
    const value = {style:{},clientHeight:height,clientWidth:400};
    const text = {get scrollHeight(){return parseInt(value.style.fontSize)*multiplier;},scrollWidth:350};
    return {value,classList:{contains:()=>arabic},querySelector:selector=>selector === '.detailed-dua-text' ? text : value};
};
const arabicRow = makeDuaRow(true, 400, 10);
const translationRow = makeDuaRow(false, 150, 6);
fitContext.communityCards.querySelectorAll = selector => selector === '.detailed-dua-arabic, .detailed-dua-translation' ? [arabicRow, translationRow] : [];
fitContext.communitySignature='dua';
fitContext.scheduleDailyFit(); callback();
assert.equal(arabicRow.value.style.fontSize,'40px');
assert.equal(translationRow.value.style.fontSize,'25px', 'Translation must fit independently of Arabic');
console.log('MasjidBoard independent Dua text sizing tests passed');
for (let count=0; count<=5; count++) {
    const notices=Array.from({length:count},()=>({type:'announcement',body:'Short notice'}));
    const banking=context.packCommunityPages([...notices,{type:'contribution',fields:{bank:'FNB',account_number:'62599032978'}}]);
    const last=banking[banking.length-1];
    assert.equal(last.layout,'single');
    assert.equal(last.entries.length,1);
    assert.equal(last.entries[0].item.type,'contribution','Banking details must receive the full sidebar height');
    assert.equal(banking.flatMap(page=>page.entries).length,count+1,'All preceding notices must be preserved');
}
console.log('MasjidBoard landscape contributions pagination tests passed');
for(const type of ['funeral','nikah','eid']) {
    const notices=context.packCommunityPages([{type:'announcement',body:'Short'},{type,fields:{}},{type:'well_wishes',body:'Short'}]);
    assert.equal(notices.length,3);
    assert.equal(notices[1].layout,'single');
    assert.equal(notices[1].entries[0].item.type,type);
}
const paired=context.packCommunityPages([{type:'programme',body:'Thursday after Esha'},{type:'well_wishes',body:'Please make dua'}]);
assert.equal(paired.length,1);
assert.equal(paired[0].layout,'halves','Short notices may share two readable cards');
console.log('Landscape notice full-height and paired pagination tests passed');
const noticeBody={style:{},clientHeight:150,clientWidth:400};
const noticeText={get scrollHeight(){return parseInt(noticeBody.style.fontSize)*6;},scrollWidth:350};
const noticeCard={style:{setProperty:()=>{}},clientHeight:300,scrollHeight:300,querySelector:selector=>selector === '.detailed-community-body' ? noticeBody : selector === '.detailed-notice-text' ? noticeText : null};
fitContext.communityCards.querySelectorAll=selector=>selector === '.detailed-refined-notice' ? [noticeCard] : [];
fitContext.communitySignature='notice';
fitContext.scheduleDailyFit();callback();
assert.equal(noticeBody.style.fontSize,'25px','Landscape announcement text must use the largest fitting size');
context.orderedFields=item=>item.testFields;
const groups=context.landscapeNoticeFields({type:'dawah',testFields:[{label:'Gasht Out',value:'Thursday'},{label:'Out Time',value:'After Asr'},{label:'Additional detail',value:'Preserve me'}]});
assert.equal(groups[0].value,'Thursday\nAfter Asr');
assert.equal(groups[1].value,'Preserve me');
console.log('Landscape notice text sizing and field preservation tests passed');
