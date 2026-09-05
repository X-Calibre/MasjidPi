(() => {
    "use strict";

    const islamicWeekdays = [
        "Al-Ahad", "Al-Ithnayn", "Ath-Thulatha", "Al-Arbi'a",
        "Al-Khamis", "Al-Jumu'ah", "As-Sabt",
    ];

    function validTime(time) {
        return time && Number.isInteger(time.hour) && Number.isInteger(time.minute);
    }

    function boardGregorianDate(board) {
        const value = board && board.date ? board.date.gregorian : "";
        if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return new Date();
        const [year, month, day] = value.split("-").map(Number);
        return new Date(year, month - 1, day);
    }

    function formatGregorianDate(board) {
        return boardGregorianDate(board).toLocaleDateString("en-ZA", {
            weekday: "long", day: "numeric", month: "long", year: "numeric",
        });
    }

    function islamicDayDate(board, now = new Date()) {
        const date = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const sunset = board && board.astronomical ? board.astronomical.sunset : null;
        if (validTime(sunset)) {
            const rolloverSeconds = sunset.hour * 3600 + sunset.minute * 60 + 185;
            const nowSeconds = now.getHours() * 3600 + now.getMinutes() * 60 + now.getSeconds();
            if (nowSeconds > rolloverSeconds) date.setDate(date.getDate() + 1);
        }
        return date;
    }

    function islamicWeekday(board, now = new Date()) {
        const date = islamicDayDate(board, now);
        return islamicWeekdays[date.getDay()];
    }

    function isIslamicFriday(board, now = new Date()) {
        return islamicDayDate(board, now).getDay() === 5;
    }

    window.MasjidBoardDate = {boardGregorianDate, formatGregorianDate, islamicWeekday, isIslamicFriday};
})();
