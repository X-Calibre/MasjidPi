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

    window.MasjidBoardCommunityUtils = {
        collectCommunityItems,
        communityTypeLabel,
        fieldLabel,
        noticeTitle,
        orderedFields,
        plainText,
    };
})();
