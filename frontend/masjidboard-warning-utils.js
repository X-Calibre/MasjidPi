(() => {
    "use strict";

    function validTime(time) {
        return time && Number.isInteger(time.hour) && Number.isInteger(time.minute)
            && time.hour >= 0 && time.hour <= 23
            && time.minute >= 0 && time.minute <= 59;
    }

    function minutes(time) {
        return time.hour * 60 + time.minute;
    }

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

    function isZawaalWarningActive(board, now) {
        if (!board || !now || typeof now.getTime !== "function" || Number.isNaN(now.getTime())) return false;
        const astronomical = board.astronomical || {};
        const start = validTime(astronomical.istiwa_caution)
            ? astronomical.istiwa_caution
            : astronomical.istiwa;
        const end = astronomical.zawaal_end;
        if (!validTime(start) || !validTime(end)) return false;

        const startMinutes = minutes(start);
        const endMinutes = minutes(end);
        if (endMinutes <= startMinutes) return false;
        const currentMinutes = boardLocalMinutes(board, now);
        return currentMinutes >= startMinutes && currentMinutes < endMinutes;
    }

    window.MasjidBoardWarningUtils = {isZawaalWarningActive};
})();
