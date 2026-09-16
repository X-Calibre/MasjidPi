(() => {
    "use strict";

    const widgets = [...document.querySelectorAll("[data-update-widget]")];
    if (widgets.length === 0) return;

    const formatDate = value => {
        if (!value) return "Not yet";
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return value;
        return new Intl.DateTimeFormat("en-ZA", {
            day: "numeric",
            month: "long",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit"
        }).format(date);
    };

    function setText(widget, selector, value) {
        const element = widget.querySelector(selector);
        if (element) element.textContent = value;
    }

    function render(widget, state) {
        const release = state?.available_release || null;
        const status = widget.querySelector("[data-update-state]");
        const link = widget.querySelector("[data-update-link]");

        setText(widget, "[data-update-current]", state?.current_version || "Unknown");
        setText(widget, "[data-update-last-checked]", formatDate(state?.last_checked_at));
        setText(widget, "[data-update-deadline]", release ? formatDate(state?.approval_deadline) : "Not applicable");

        if (release) {
            status.textContent = `${release.version} is available`;
            status.className = "update-state update-state-available";
            setText(widget, "[data-update-published]", formatDate(release.published_at));
            if (link) {
                link.href = release.page_url;
                link.classList.remove("hidden");
            }
        } else {
            status.textContent = "No stable update available";
            status.className = "update-state update-state-current";
            setText(widget, "[data-update-published]", "Not applicable");
            if (link) {
                link.removeAttribute("href");
                link.classList.add("hidden");
            }
        }

        const error = widget.querySelector("[data-update-error]");
        if (error) {
            error.textContent = state?.last_check_error || "";
            error.classList.toggle("hidden", !state?.last_check_error);
        }
    }

    function renderFailure(widget, message) {
        const status = widget.querySelector("[data-update-state]");
        status.textContent = "Update status unavailable";
        status.className = "update-state update-state-error";
        const error = widget.querySelector("[data-update-error]");
        if (error) {
            error.textContent = message;
            error.classList.remove("hidden");
        }
    }

    async function requestState(method) {
        const response = await fetch(
            method === "POST" ? "/api/update/check" : "/api/update/status",
            {method}
        );
        let state = null;
        try { state = await response.json(); } catch (_) {}
        if (!response.ok && !state?.schema_version) {
            throw new Error(state?.error || `Request failed (${response.status})`);
        }
        return {state, ok: response.ok};
    }

    async function refresh(widget, checkNow = false) {
        const button = widget.querySelector("[data-update-check]");
        if (button) {
            button.disabled = true;
            button.textContent = checkNow ? "Checking…" : "Loading…";
        }
        try {
            const result = await requestState(checkNow ? "POST" : "GET");
            render(widget, result.state);
        } catch (error) {
            renderFailure(widget, error.message);
        } finally {
            if (button) {
                button.disabled = false;
                button.textContent = "Check now";
            }
        }
    }

    for (const widget of widgets) {
        widget.querySelector("[data-update-check]")?.addEventListener("click", () => refresh(widget, true));
        refresh(widget);
    }
})();
