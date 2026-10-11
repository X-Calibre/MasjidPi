'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const path=require('node:path');
const vm=require('node:vm');
const source=fs.readFileSync(path.join(__dirname,'../frontend/masjidboard-appliance.js'),'utf8');
const context={orderedFields:item=>item.testFields||[]};
vm.createContext(context);
vm.runInContext(source.slice(source.indexOf('    function isRefinedNotice('),source.indexOf('    function communityCard(')),context);
for(const [type,labels] of Object.entries({dawah:['Masjid Taleem','Gasht Out','Out Time','Gasht In','In Time'],three_day_jamaat:['First Jamaat','First Date','Second Jamaat','Second Date'],new_moon:['Birth Date','Birth','Best Visibility','Visibility Date']})){
 const fields=labels.map((label,i)=>({label,value:`Value ${i}`}));
 fields.push({label:'Additional detail',value:'Preserved'});
 const grouped=context.portraitNoticeFields({type,testFields:fields});
 assert.equal(grouped.length,type==='dawah'?4:3);
 assert.equal(grouped.map(field=>field.value).join('\n'),fields.map(field=>field.value).join('\n'),'Every supplied detail must survive grouping');
 const missing=context.portraitNoticeFields({type,testFields:[fields[1]]});
 assert.equal(missing.length,1,'A missing partner field must not create empty groups');
 assert.equal(missing[0].value,fields[1].value);
}
context.isCompactCommunityItem=()=>true;
context.communitySlide=entries=>entries;
vm.runInContext(source.slice(source.indexOf('    function communitySlides('),source.indexOf('    function economicSlide(')),context);
const cards=['funeral','nikah','eid','programme','announcement','well_wishes','dawah','three_day_jamaat','new_moon'].map(type=>({type}));
const slides=context.communitySlides(cards);
assert.equal(slides.length,cards.length,'Every refined notice receives a full slide');
assert.ok(slides.every(slide=>slide.length===1));
console.log('Portrait notice grouping and full-height slide tests passed');
