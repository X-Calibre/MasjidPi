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

    function clockFromMinutes(value) {
        const normalized = (value + 24 * 60) % (24 * 60);
        return {hour: Math.floor(normalized / 60), minute: normalized % 60};
    }

    function zawaalWindow(astronomical) {
        const times = astronomical || {};
        const istiwa = validTime(times.istiwa) ? times.istiwa : null;
        const start = validTime(times.istiwa_caution)
            ? times.istiwa_caution
            : istiwa ? clockFromMinutes(minutes(istiwa) - 5) : null;
        const end = validTime(times.zawaal_end)
            ? times.zawaal_end
            : istiwa ? clockFromMinutes(minutes(istiwa) + 5) : null;
        return {start, end};
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
        const {start, end} = zawaalWindow(board.astronomical);
        if (!validTime(start) || !validTime(end)) return false;

        const startMinutes = minutes(start);
        const endMinutes = minutes(end);
        if (endMinutes <= startMinutes) return false;
        const currentMinutes = boardLocalMinutes(board, now);
        return currentMinutes >= startMinutes && currentMinutes < endMinutes;
    }

    window.MasjidBoardWarningUtils = {isZawaalWarningActive, zawaalWindow};
})();
