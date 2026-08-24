(() => {
    "use strict";

    function plainText(value) {
        const source = String(value || "").replace(/<br\s*\/?\s*>/gi, "\n");
        if (!source) return "";
        const parsed = new DOMParser().parseFromString(source, "text/html");
        return (parsed.body.textContent || "").replace(/\r/g, "").replace(/\n{3,}/g, "\n\n").trim();
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
            silver: "Silver / g", minimum_mahr: "Minimum Mahr", mahr_faatimi: "Mahr Faatimi",
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
            economic: ["nisaab", "krugerrand", "gold_24", "silver", "minimum_mahr", "mahr_faatimi"],
        }[item.type] || [];
        const titleFields = new Set(item.type === "funeral" ? ["name"]
            : item.type === "nikah" ? ["name_one", "name_two", "bride"]
                : item.type === "salaah_change" ? ["prayer"] : []);
        const names = [...preferred, ...Object.keys(fields).sort()].filter((name, index, all) =>
            !titleFields.has(name) && all.indexOf(name) === index && plainText(fields[name])
        );
        return names.slice(0, 6).map((name) => ({label: fieldLabel(name), value: plainText(fields[name])}));
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
            announcement: "Announcement", eid: "Eid Notice", funeral: "Funeral Notice",
            nikah: "Nikah Notice", well_wishes: "Well Wishes", salaah_change: "Salaah Time Change",
            programme: "Programme", new_moon: "New Moon", dawah: "Dawah / Gasht",
            three_day_jamaat: "Three-Day Jamaat", contribution: "Contributions",
            economic: "Islamic Economic Indicators",
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
                add({
                    type: "announcement",
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
        communityTypeLabel,
        fieldLabel,
        fixtureCommunityItems,
        noticeTitle,
        orderedFields,
        plainText,
    };
})();
