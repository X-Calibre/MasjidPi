"use strict";
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

function element(value = "") {
    const listeners = {};
    return {
        value, checked: false, disabled: false, textContent: "",
        classList: { toggle() {} },
        addEventListener(name, fn) { listeners[name] = fn; },
        trigger(name) { return listeners[name](); }
    };
}

const ids = Object.fromEntries([
    "radioScheduleEnabled", "radioScheduleStart", "radioScheduleStop",
    "radioScheduleTimes", "radioScheduleSave", "radioScheduleStatus"
].map(id => [id, element()]));
const events = {};
const requests = [];
vm.runInNewContext(fs.readFileSync(path.join(__dirname, "../frontend/radio-schedule.js"), "utf8"), {
    document: { getElementById: id => ids[id] },
    window: {
        addEventListener: (name, fn) => { events[name] = fn; },
        MasjidFrameRefreshListenStatus: async () => {},
        MasjidFrameUI: { notify() {} }
    },
    fetch: async (url, options) => {
        requests.push({ url, body: JSON.parse(options.body) });
        return { ok: true, json: async () => ({
            radio_schedule_enabled: true, radio_schedule_start: "08:00",
            radio_schedule_stop: "20:00", radio_schedule_allows_now: true
        }) };
    }
});

const baseline = { radio_schedule_enabled: true, radio_schedule_start: "06:00", radio_schedule_stop: "22:00", radio_schedule_allows_now: false };
events["masjidframe:listen-status"]({ detail: baseline });
ids.radioScheduleStart.value = "08:00";
ids.radioScheduleStart.trigger("change");
ids.radioScheduleStop.value = "20:00";
ids.radioScheduleStop.trigger("change");
events["masjidframe:listen-status"]({ detail: baseline });
assert.equal(ids.radioScheduleStart.value, "08:00");
assert.equal(ids.radioScheduleStop.value, "20:00");
assert.equal(requests.length, 0);
assert.equal(ids.radioScheduleSave.disabled, false);
ids.radioScheduleSave.trigger("click").then(() => {
    assert.equal(requests.length, 1);
    assert.deepEqual(requests[0].body, { enabled: true, start: "08:00", stop: "20:00" });
    assert.equal(ids.radioScheduleSave.disabled, true);
    console.log("Radio schedule UI tests passed");
}).catch(error => { console.error(error); process.exitCode = 1; });
