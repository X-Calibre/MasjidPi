(() => {
    "use strict";

    function plainText(value) {
        const source = String(value || "").replace(/<br\s*\/?\s*>/gi, "\n");
        if (!source) return "";
        const parsed = new DOMParser().parseFromString(source, "text/html");
        return (parsed.body.textContent || "").replace(/\r/g, "").replace(/\n{3,}/g, "\n\n").trim();
    }

    // South African Rand values use a fixed presentation contract throughout
    // MasjidBoard: R123,456.00. Do not delegate the currency symbol or
    // separators to the browser locale, which varies across WebKit builds.
    function formatRand(value, decimals = 2) {
        const amount = Number(value);
        if (!Number.isFinite(amount)) return "";
        return "R" + amount.toLocaleString("en-US", {
            minimumFractionDigits: decimals,
            maximumFractionDigits: decimals,
        });
    }

    function formatNoticeDate(value) {
        const raw = plainText(value);
        const months = {
            jan: 0, feb: 1, mar: 2, apr: 3, may: 4, jun: 5,
            jul: 6, aug: 7, sep: 8, oct: 9, nov: 10, dec: 11,
        };
        const textMatch = raw.match(/^(\d{1,2})\s+([A-Za-z]+)(?:\s+(\d{4}))?(?:\s+\d{1,2}:\d{2}(?::\d{2})?)?$/);
        const isoMatch = raw.match(/^(\d{4})-(\d{1,2})-(\d{1,2})(?:[T\s].*)?$/);
        if (!textMatch && !isoMatch) return raw;

        const now = new Date();
        const hasYear = Boolean(isoMatch || textMatch[3]);
        let year = isoMatch ? Number(isoMatch[1]) : textMatch[3] ? Number(textMatch[3]) : now.getFullYear();
        const month = isoMatch ? Number(isoMatch[2]) - 1 : months[textMatch[2].slice(0, 3).toLowerCase()];
        const day = Number(isoMatch ? isoMatch[3] : textMatch[1]);
        if (month === undefined || month < 0 || month > 11) return raw;
        let date = new Date(year, month, day, 12);
        if (!hasYear && date.getTime() < now.getTime() - 180 * 24 * 60 * 60 * 1000) {
            date = new Date(year + 1, month, day, 12);
        }
        if (Number.isNaN(date.getTime()) || date.getDate() !== day || date.getMonth() !== month) return raw;
        return date.toLocaleDateString("en-ZA", {
            weekday: "long", day: "numeric", month: "long",
        });
    }

    function formatUpdatedAt(value) {
        const updatedAt = new Date(value);
        if (Number.isNaN(updatedAt.getTime())) return String(value || "");
        return updatedAt.toLocaleString("en-ZA", {
            day: "numeric", month: "short", year: "numeric",
            hour: "2-digit", minute: "2-digit", hour12: false,
        });
    }

    function formatClock(time) {
        if (!time || !Number.isInteger(time.hour) || !Number.isInteger(time.minute)) return "";
        return `${String(time.hour).padStart(2, "0")}:${String(time.minute).padStart(2, "0")}`;
    }

    const duaAfterAdhanArabic = "اللَّهُمَّ رَبَّ هَذِهِ الدَّعْوَةِ التَّامَّةِ، وَالصَّلَاةِ الْقَائِمَةِ، آتِ مُحَمَّدًا الْوَسِيلَةَ وَالْفَضِيلَةَ، وَابْعَثْهُ مَقَامًا مَحْمُودًا الَّذِي وَعَدْتَهُ";
    const duaAfterAdhanTranslation = "O Allah, Lord of this perfect call and established prayer, grant Muhammad the Wasilah and virtue, and raise him to the praised position You have promised him.";
    const duaAfterAdhanWindowMinutes = 5;

    function boardLocalMinutes(board, now) {
        const timezone = String(board && board.time_zone || "").trim();
        const fixedOffset = timezone.match(/^(?:GMT|UTC)(?:([+-])(\d{1,2})(?::?(\d{2}))?)?$/i);
        if (fixedOffset) {
            const sign = fixedOffset[1] === "-" ? -1 : 1;
            const offset = fixedOffset[1]
                ? sign * (Number(fixedOffset[2]) * 60 + Number(fixedOffset[3] || 0))
                : 0;
            return (now.getUTCHours() * 60 + now.getUTCMinutes() + offset + 24 * 60) % (24 * 60);
        }
        try {
            const parts = new Intl.DateTimeFormat("en-GB", {
                timeZone: timezone || undefined,
                hour: "2-digit", minute: "2-digit", hourCycle: "h23",
            }).formatToParts(now);
            const values = Object.fromEntries(parts.map((part) => [part.type, part.value]));
            return Number(values.hour) * 60 + Number(values.minute);
        } catch (_) {
            return now.getHours() * 60 + now.getMinutes();
        }
    }

    function specialDhuhrItem(board) {
        const special = board && board.special_dhuhr;
        const time = special && special.time;
        const label = plainText(special && special.label);
        if (!time || !Number.isInteger(time.hour) || !Number.isInteger(time.minute) || !label) return null;
        return {label: `Dhuhr ${label}`, time};
    }

    function duaAfterAdhanContent() {
        return {
            type: "dua_after_adhan",
            title: "Dua after Adhan",
            body: "",
            fields: {arabic: duaAfterAdhanArabic, translation: duaAfterAdhanTranslation, note: "Wasilah is a rank in Jannah."},
            source: "MasjidPi",
        };
    }

    function duaAfterAdhanItem(boards, now, enabled, windowMinutes = duaAfterAdhanWindowMinutes, force = false) {
        if (!enabled || !now || typeof now.getTime !== "function" || Number.isNaN(now.getTime())) return null;
        if (force) return duaAfterAdhanContent();
        for (const board of Array.isArray(boards) ? boards : []) {
            const currentMinutes = boardLocalMinutes(board, now);
            for (const prayer of Array.isArray(board && board.prayers) ? board.prayers : []) {
                const adhan = prayer && prayer.adhan;
                if (!adhan || !Number.isInteger(adhan.hour) || !Number.isInteger(adhan.minute)) continue;
                let elapsed = currentMinutes - (adhan.hour * 60 + adhan.minute);
                if (elapsed < 0) elapsed += 24 * 60;
                if (elapsed >= 0 && elapsed < windowMinutes) {
                    return duaAfterAdhanContent();
                }
            }
        }
        return null;
    }

    function detailedJumuahItems(boards, now, isFriday) {
        const items = [];
        for (const board of Array.isArray(boards) ? boards : []) {
            if (!board || board.show_detailed_jumuah === false || !isFriday(board, now)) continue;
            const service = Array.isArray(board.jumuah) ? board.jumuah[0] : null;
            if (!service) continue;

            const schedule = [];
            const semanticName = (value) => plainText(value).normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
            const semantic = new Set();
            const representedTimes = new Set();
            for (const event of Array.isArray(service.events) ? service.events : []) {
                const heading = plainText(event && event.heading);
                const time = formatClock(event && event.time);
                if (!heading || !time) continue;
                schedule.push({heading, time});
                semantic.add(String(event.code || "").trim() === "0" ? "adhan" : semanticName(heading));
                representedTimes.add(time);
            }
            const addFallback = (semanticName, heading, time) => {
                const formatted = formatClock(time);
                if (formatted && !semantic.has(semanticName) && !representedTimes.has(formatted)) {
                    schedule.push({heading, time: formatted});
                    representedTimes.add(formatted);
                }
            };
            addFallback("adhan", "Adhan", service.adhan);
            if (!semantic.has("jamaah") && !semantic.has("salaah")) {
                addFallback("jamaah", "Salaah", service.jamaah || service.effective_salaah);
            }
            if (schedule.length === 0) continue;

            items.push({
                type: "jumuah_schedule",
                title: "Jumu’ah Schedule",
                body: plainText(service.khateeb),
                fields: {},
                schedule,
                source: plainText(board.name),
            });
        }
        return items;
    }

    function boardLocalDate(board, now) {
        const timezone = String(board && board.time_zone || "").trim();
        try {
            const parts = new Intl.DateTimeFormat("en-CA", {
                timeZone: timezone || undefined,
                year: "numeric", month: "2-digit", day: "2-digit",
            }).formatToParts(now);
            const values = Object.fromEntries(parts.map((part) => [part.type, part.value]));
            return new Date(Date.UTC(Number(values.year), Number(values.month) - 1, Number(values.day)));
        } catch (_) {
            return new Date(Date.UTC(now.getFullYear(), now.getMonth(), now.getDate()));
        }
    }

    function noticeEffectiveDate(value, today) {
        const raw = plainText(value).replace(/^(?:mon|tue|wed|thu|fri|sat|sun)(?:day)?\s*,?\s*/i, "").trim();
        const months = {
            jan: 0, feb: 1, mar: 2, apr: 3, may: 4, jun: 5,
            jul: 6, aug: 7, sep: 8, oct: 9, nov: 10, dec: 11,
        };
        const iso = raw.match(/^(\d{4})-(\d{1,2})-(\d{1,2})/);
        const dayFirst = raw.match(/^(\d{1,2})\s+([A-Za-z]+)(?:\s*,?\s*(\d{4}))?/);
        const monthFirst = raw.match(/^([A-Za-z]+)\s+(\d{1,2})(?:\s*,?\s*(\d{4}))?/);
        let year;
        let month;
        let day;
        let explicitYear = false;
        if (iso) {
            year = Number(iso[1]); month = Number(iso[2]) - 1; day = Number(iso[3]); explicitYear = true;
        } else if (dayFirst) {
            year = dayFirst[3] ? Number(dayFirst[3]) : today.getUTCFullYear();
            month = months[dayFirst[2].slice(0, 3).toLowerCase()];
            day = Number(dayFirst[1]);
            explicitYear = Boolean(dayFirst[3]);
        } else if (monthFirst) {
            year = monthFirst[3] ? Number(monthFirst[3]) : today.getUTCFullYear();
            month = months[monthFirst[1].slice(0, 3).toLowerCase()];
            day = Number(monthFirst[2]);
            explicitYear = Boolean(monthFirst[3]);
        } else {
            return null;
        }
        if (month === undefined) return null;
        let result = new Date(Date.UTC(year, month, day));
        if (!explicitYear && result < today) result = new Date(Date.UTC(year + 1, month, day));
        if (result.getUTCMonth() !== month || result.getUTCDate() !== day) return null;
        return result;
    }

    function salaahChangeEffectiveDate(item, today) {
        const fields = item && item.fields && typeof item.fields === "object" ? item.fields : {};
        const structured = noticeEffectiveDate(fields.effective_date, today);
        if (structured) return structured;
        const text = `${plainText(item && item.title)}\n${plainText(item && item.body)}`;
        const match = text.match(/\b(?:effective\s+(?:from\s+)?|as\s+of\s+)?((?:(?:mon|tue|wed|thu|fri|sat|sun)(?:day)?\s*,?\s*)?(?:\d{1,2}\s+[A-Za-z]+|[A-Za-z]+\s+\d{1,2})(?:\s*,?\s*\d{4})?)/i);
        return match ? noticeEffectiveDate(match[1], today) : null;
    }

    function isVisibleSalaahChange(item, board, now) {
        if (item.type !== "salaah_change" && item.type !== "salaah_change_announcement") return true;
        const today = boardLocalDate(board, now);
        const effective = salaahChangeEffectiveDate(item, today);
        if (!effective) return false;
        const daysAway = Math.round((effective.getTime() - today.getTime()) / 86_400_000);
        return daysAway >= 0 && daysAway <= 7;
    }

    function communityRank(item) {
        const ranks = {
            funeral: 10,
            urgent_announcement: 20,
            salaah_change: 21,
            salaah_change_announcement: 21,
            eid: 22,
            announcement: 30,
            general: 30,
            jumuah_schedule: 31,
            nikah: 32,
            well_wishes: 33,
            programme: 40,
            weekly_programme: 40,
            ramadan_programme: 40,
            class_time_change: 41,
            dawah: 42,
            three_day_jamaat: 43,
            new_moon: 50,
            contribution: 51,
        };
        return ranks[item.type] || 30;
    }

    function communityPriorityGroups(items) {
        const groups = new Map();
        for (const entry of items.map((item, index) => ({item, index, rank: communityRank(item)}))
            .sort((left, right) => left.rank - right.rank || left.index - right.index)) {
            const tier = Math.floor(entry.rank / 10);
            if (!groups.has(tier)) groups.set(tier, []);
            groups.get(tier).push(entry.item);
        }
        return [...groups.entries()]
            .sort((left, right) => left[0] - right[0])
            .map(([, entries]) => entries);
    }

    function orderedCommunityItemGroups(board, now, isFriday) {
        if (!board) return [];
        const items = collectCommunityItems([board]);
        items.push(...detailedJumuahItems([board], now, isFriday));
        return communityPriorityGroups(items.filter((item) => isVisibleSalaahChange(item, board, now)));
    }

    function orderedCommunityItems(board, now, isFriday) {
        return orderedCommunityItemGroups(board, now, isFriday).flat();
    }

    function fieldLabel(name) {
        const labels = {
            address: "Address", bride: "Bride", cemetery: "Cemetery", date: "Date",
            groom_relation: "Groom", lecture: "Lecture", name: "Name", name_one: "Name",
            name_two: "Name", pickup: "Pickup", relation: "Relation", relation_one: "Family",
            relation_two: "Family", salaah: "Salaah", salaah_time: "Janazah", salaah_venue: "Venue",
            time: "Time", venue: "Venue", account_name: "Account Name", account_number: "Account Number",
            bank: "Bank", branch_code: "Branch Code", bsb: "BSB", masjid_taleem: "Masjid Taleem",
            gasht_out_day: "Gasht Out", gasht_out_time: "Out Time", gasht_in_day: "Gasht In",
            gasht_in_time: "In Time", first_location: "First Jamaat", first_date: "First Date",
            second_location: "Second Jamaat", second_date: "Second Date",
            nisaab: "Nisaab", krugerrand: "Krugerrand", gold_24: "Gold 24 ct / g",
            rand_dollar: "Rand/Dollar", gold_22: "Gold 22 ct / g", gold_18: "Gold 18 ct / g",
            gold_14: "Gold 14 ct / g", gold_9: "Gold 9 ct / g", silver: "Silver / g",
            minimum_mahr: "Minimum Mahr", mahr_faatimi: "Mahr Faatimi",
            retrieved_at: "Retrieved at", ayah_number: "Ayah", reference: "Reference",
            arabic: "Arabic", translation: "Translation", note: "Note",
        };
        return labels[name] || name.replace(/_/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
    }

    function orderedFields(item) {
        const fields = item && item.fields && typeof item.fields === "object" ? item.fields : {};
        const preferred = {
            funeral: ["relation", "salaah_time", "salaah_venue", "pickup", "cemetery", "address"],
            nikah: ["groom_relation", "relation_one", "relation_two", "date", "time"],
            eid: ["date", "venue", "address", "lecture", "salaah"],
            salaah_change: ["effective_date", "new_time"],
            new_moon: ["birth_date", "birth", "best_visibility", "visibility_date", "first_moonset", "first_age"],
            dawah: ["masjid_taleem", "gasht_out_day", "gasht_out_time", "gasht_in_day", "gasht_in_time"],
            three_day_jamaat: ["first_location", "first_date", "second_location", "second_date"],
            contribution: ["bank", "account_name", "branch_code", "account_number", "bsb"],
            economic: ["rand_dollar", "nisaab", "krugerrand", "gold_24", "gold_22", "gold_18", "gold_14", "gold_9", "silver", "minimum_mahr", "mahr_faatimi", "retrieved_at"],
			dua_after_adhan: ["arabic", "translation", "note"],
			daily_ayah: ["ayah_number"], daily_hadith: ["reference"], daily_sunnah: ["reference"],
        }[item.type] || [];
        const titleFields = new Set(item.type === "funeral" ? ["name"]
            : item.type === "nikah" ? ["name_one", "name_two", "bride"]
                : item.type === "salaah_change" ? ["prayer"] : []);
        const names = [...preferred, ...Object.keys(fields).sort()].filter((name, index, all) =>
            !titleFields.has(name) && all.indexOf(name) === index && plainText(fields[name])
        );
        const fieldLimit = item.type === "economic" ? 13 : 6;
        return names.slice(0, fieldLimit).map((name) => ({label: fieldLabel(name), value: plainText(fields[name])}));
    }

    function noticeTitle(notice) {
        const fields = notice && notice.fields && typeof notice.fields === "object" ? notice.fields : {};
        if (notice.type === "funeral") return plainText(fields.name) || plainText(notice.title) || "Funeral Notice";
        if (notice.type === "nikah") {
            const names = [fields.name_one, fields.bride || fields.name_two].map(plainText).filter(Boolean);
            if (names.length > 0) return names.join(" & ");
        }
        return plainText(notice.title) || fieldLabel(notice.type || "general") + " Notice";
    }

    function communityTypeLabel(type) {
        const labels = {
            announcement: "Announcement", urgent_announcement: "Urgent Announcement", eid: "Eid Notice", funeral: "Funeral Notice",
            nikah: "Nikah Notice", well_wishes: "Well Wishes", salaah_change: "Salaah Time Change",
            programme: "Programme", new_moon: "New Moon", dawah: "Dawah / Gasht",
            three_day_jamaat: "Three-Day Jamaat", contribution: "Contributions",
            economic: "Islamic Economic Indicators",
            salaah_change_announcement: "Salaah Time Change", class_time_change: "Class Time Change",
            weekly_programme: "Weekly Programme", ramadan_programme: "Ramadan Programme",
            dua_after_adhan: "Dua after Adhan",
        };
        return labels[type] || fieldLabel(type || "general") + " Notice";
    }

    function collectCommunityItems(boards) {
        const items = [];
        const seen = new Set();
        function add(item) {
            const key = JSON.stringify([item.type, item.title, item.body, item.fields]);
            if (!item.title && !item.body && Object.keys(item.fields || {}).length === 0) return;
            if (seen.has(key)) return;
            seen.add(key);
            items.push(item);
        }

        for (const board of boards) {
            for (const notice of Array.isArray(board.notices) ? board.notices : []) {
                const fields = notice && notice.fields && typeof notice.fields === "object" ? notice.fields : {};
                add({
                    type: plainText(notice.type).toLowerCase() || "general",
                    title: noticeTitle(notice),
                    body: Object.keys(fields).length === 0 ? plainText(notice.content) : "",
                    fields,
                    source: board.name,
                });
            }
            for (const announcement of Array.isArray(board.announcements) ? board.announcements : []) {
                const categories = new Set(["announcement", "urgent_announcement", "salaah_change_announcement", "class_time_change", "weekly_programme", "ramadan_programme"]);
                const category = plainText(announcement.category).toLowerCase();
                add({
                    type: categories.has(category) ? category : "announcement",
                    typeLabel: categories.has(category) && category !== "announcement" ? communityTypeLabel(category) : "",
                    title: plainText(announcement.title) || "Masjid Announcement",
                    body: plainText(announcement.content),
                    fields: {},
                    source: board.name,
                });
            }
            for (const programme of Array.isArray(board.programmes) ? board.programmes : []) {
                add({
                    type: "programme",
                    title: plainText(programme.title) || "Masjid Programme",
                    body: plainText(programme.content),
                    fields: {},
                    source: board.name,
                });
            }
            if (board.new_moon && board.new_moon.fields && typeof board.new_moon.fields === "object") {
                add({type: "new_moon", title: "New Moon Information", body: "", fields: board.new_moon.fields, source: board.name});
            }
            if (board.banking && board.banking.fields && typeof board.banking.fields === "object") {
                add({
                    type: "contribution",
                    title: plainText(board.banking.title) || "Masjid Contributions",
                    body: "",
                    fields: board.banking.fields,
                    source: board.name,
                });
            }
        }
        return items;
    }

    function dailyIslamicItems(content) {
        if (!content || typeof content !== "object") return [];
        const source = plainText(content.source) || "MasjidBoard Live";
        const items = [];
        if (content.ayah && plainText(content.ayah.text)) {
            items.push({
                type: "daily_ayah",
                title: plainText(content.ayah.surah) || "Daily Ayah",
                body: plainText(content.ayah.text),
                fields: {ayah_number: plainText(content.ayah.ayah_number)},
                source,
            });
        }
        if (content.hadith && plainText(content.hadith.text)) {
            items.push({
                type: "daily_hadith",
                title: plainText(content.hadith.heading) || "Hadith",
                body: plainText(content.hadith.text),
                fields: {reference: plainText(content.hadith.reference)},
                source,
            });
        }
        if (content.sunnah && plainText(content.sunnah.text)) {
            items.push({
                type: "daily_sunnah",
                title: plainText(content.sunnah.heading) || "Sunnah",
                body: plainText(content.sunnah.text),
                fields: {reference: plainText(content.sunnah.reference)},
                source,
            });
        }
        return items;
    }

    function fixtureCommunityItems(mode) {
        const source = "Layout test fixture — not live";
        const items = [
            {
                type: "funeral",
                title: "Marhoom Abdullah Ismail",
                body: "",
                fields: {
                    relation: "Father of Yusuf Ismail",
                    salaah_time: "After Zuhr · 13:15",
                    salaah_venue: "Central Masjid",
                    pickup: "12:45 from the family residence",
                    cemetery: "Central Cemetery",
                    address: "10 Example Road",
                },
                source,
            },
            {
                type: "nikah",
                title: "Muhammad & Ayesha",
                body: "",
                fields: {
                    groom_relation: "Son of Ahmad & Fatima",
                    relation_two: "Daughter of Ismail & Maryam",
                    date: "Saturday, 29 August 2026",
                    time: "After Asr · 16:45",
                    venue: "Masjid Hall",
                },
                source,
            },
            {
                type: "eid",
                title: "Eid Salaah Notice",
                body: "",
                fields: {
                    date: "Monday, 25 May 2026",
                    venue: "Community Sports Ground",
                    address: "1 Example Field Road",
                    lecture: "Lecture · 07:00",
                    salaah: "Salaah · 07:30",
                },
                source,
            },
            {
                type: "announcement",
                title: "Important Access Notice",
                body: "Please use the northern entrance while maintenance work is under way. The main parking area will be closed after Maghrib. Elderly worshippers and families may use the reserved drop-off area near the hall entrance. Please follow the directions of the volunteers on duty.",
                fields: {},
                source,
            },
			{
				type: "dawah",
				title: "Dawah and Gasht",
				body: "",
				fields: {masjid_taleem: "Daily after Esha Salaah", gasht_out_day: "Thursday", gasht_out_time: "After Asr", gasht_in_day: "Monday", gasht_in_time: "After Maghrib"},
				source,
			},
			{
				type: "three_day_jamaat",
				title: "Three-Day Jamaat",
				body: "",
				fields: {first_location: "Hartbeespoort area", first_date: "4–6 September", second_location: "Pretoria West", second_date: "11–13 September"},
				source,
			},
			{
				type: "contribution",
				title: "Masjid Contributions — Lillah Only",
				body: "",
				fields: {bank: "Example Bank", account_name: "Example Masjid Trust", branch_code: "123456", account_number: "000 123 456", bsb: ""},
				source,
			},
            {
                type: "salaah_change",
                title: "Esha Time Change",
                body: "",
                fields: {prayer: "Esha", effective_date: "1 September", new_time: "19:45"},
                source,
            },
            {
                type: "programme",
                title: "Taleem Programme",
                body: "Wednesday 11:15–12:15\nResident's home",
                fields: {},
                source,
            },
            {
                type: "new_moon",
                title: "New Moon Information",
                body: "",
                fields: {birth_date: "23 August", birth: "05:27", visibility_date: "24 August", best_visibility: "18:37"},
                source,
            },
            {
                type: "announcement",
                title: "Masjid Announcement",
                body: "تذكير: سيكون البرنامج بعد صلاة العشاء بإذن الله. نرجو من الجميع الحضور في الوقت المحدد.",
                fields: {},
                source,
            },
            {
                type: "well_wishes",
                title: "Du'a Requested",
                body: "The community is requested to make du'a for those who are unwell and for their families.",
                fields: {},
                source,
            },
            {
                type: "announcement",
                title: "Weekly Programme",
                body: "The weekly community programme will take place after Esha on Thursday evening.",
                fields: {},
                source,
            },
        ];
        if (mode === "new") {
			const newTypes = new Set(["dawah", "three_day_jamaat", "contribution"]);
            return items.filter((item) => newTypes.has(item.type));
        }
        return items;
    }

    window.MasjidBoardCommunityUtils = {
        collectCommunityItems,
		dailyIslamicItems,
        detailedJumuahItems,
		duaAfterAdhanItem,
        duaAfterAdhanWindowMinutes,
        communityTypeLabel,
        fieldLabel,
        fixtureCommunityItems,
        formatNoticeDate,
        formatRand,
        formatUpdatedAt,
        noticeTitle,
        communityPriorityGroups,
        orderedCommunityItemGroups,
        orderedCommunityItems,
        orderedFields,
        plainText,
        specialDhuhrItem,
    };
})();
