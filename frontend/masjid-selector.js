// Keep the favourites section independent from catalogue search.
// This file contains the UI-specific behaviour for the Masjid selector.

const originalRenderFavourites = renderFavourites;

renderFavourites = function () {
    originalRenderFavourites();
};

// Re-render favourites after the initial app.js startup, so the section is
// independent of whatever search query was present when the catalogue loaded.
renderFavourites();

// Selecting a favourite should work even while the catalogue is filtered.
// Clear the catalogue search so the selected masjid is visible in the catalogue.
favourites.addEventListener("click", event => {
    if (event.target.closest(".favourite-remove")) return;

    const item = event.target.closest(".favourite-item");
    if (!item) return;

    const id = item.dataset.id;
    if (!catalogue.some(stream => stream.id === id)) return;

    if (streamSearch.value) {
        streamSearch.value = "";
        renderStreams(id);
    }
});
