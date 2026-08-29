(() => {
    "use strict";

    const pollIntervalMs = 2_000;
    const transientDurationMs = 10_000;

    function sourceState(status) {
        if (!status || !status.listening) return {kind: "none", key: "none"};
        if (status.radio_resume_pending && status.radio_resume_at) {
            const name = status.radio_name || "Selected radio station";
            return {
                kind: "waiting",
                key: `waiting:${status.radio_id || name}:${status.radio_resume_at}`,
                name,
                resumeAt: status.radio_resume_at,
            };
        }
        if (status.active_source === "masjid") {
            const name = status.active_stream_name || status.masjid_name || "Selected masjid";
            return {kind: "masjid", key: `masjid:${status.active_stream_id || name}`, name};
        }
        if (status.active_source === "radio") {
            const name = status.active_stream_name || status.radio_name || "Selected radio station";
            return {kind: "radio", key: `radio:${status.active_stream_id || name}`, name};
        }
        return {kind: "none", key: "none"};
    }

    function remainingText(resumeAt, now = Date.now()) {
        const remainingMs = new Date(resumeAt).getTime() - now;
        const totalSeconds = Math.max(0, Math.ceil(remainingMs / 1000));
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${minutes}:${String(seconds).padStart(2, "0")}`;
    }

    const exported = {sourceState, remainingText};
    if (typeof module !== "undefined" && module.exports) module.exports = exported;
    if (typeof document === "undefined") return;

    const toasts = Array.from(document.querySelectorAll("[data-board-source-toast]"));
    if (toasts.length === 0) return;

    let initialised = false;
    let lastKey = "";
    let currentState = {kind: "none", key: "none"};
    let hideTimer = 0;

    function hide() {
        clearTimeout(hideTimer);
        hideTimer = 0;
        for (const toast of toasts) {
            toast.classList.add("hidden");
            delete toast.dataset.sourceKind;
        }
    }

    function content(state) {
        if (state.kind === "waiting") {
            return {label: "Waiting for Radio", message: `${state.name} resumes in ${remainingText(state.resumeAt)}`};
        }
        if (state.kind === "masjid") return {label: "Masjid now playing", message: state.name};
        if (state.kind === "radio") return {label: "Radio now playing", message: state.name};
        return null;
    }

    function show(state, persistent) {
        clearTimeout(hideTimer);
        hideTimer = 0;
        const value = content(state);
        if (!value) {
            hide();
            return;
        }
        for (const toast of toasts) {
            toast.dataset.sourceKind = state.kind;
            toast.querySelector(".board-source-toast-label").textContent = value.label;
            toast.querySelector(".board-source-toast-message").textContent = value.message;
            toast.classList.remove("hidden");
        }
        if (!persistent) hideTimer = setTimeout(hide, transientDurationMs);
    }

    function applyStatus(status) {
        const next = sourceState(status);
        currentState = next;

        if (next.kind === "waiting") {
            show(next, true);
        } else if (!initialised) {
            hide();
        } else if (next.key !== lastKey) {
            if (next.kind === "masjid" || next.kind === "radio") show(next, false);
            else hide();
        }

        initialised = true;
        lastKey = next.key;
    }

    async function poll() {
        try {
            const response = await fetch("/api/listen/status", {cache: "no-store"});
            if (!response.ok) return;
            applyStatus(await response.json());
        } catch (_) {
            // Board-only profiles and transient network loss need no display error.
        }
    }

    setInterval(() => {
        if (currentState.kind === "waiting") show(currentState, true);
    }, 1_000);
    setInterval(poll, pollIntervalMs);
    poll();
})();
