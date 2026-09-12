"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");

const root = path.resolve(__dirname, "..");
const app = fs.readFileSync(path.join(root, "frontend/app.js"), "utf8");
const html = fs.readFileSync(path.join(root, "frontend/index.html"), "utf8");
const css = fs.readFileSync(path.join(root, "frontend/style.css"), "utf8");

assert.match(app, /function featuredMasjidRank/);
assert.match(app, /quraan recitation/);
assert.match(app, /quran recitation/);
assert.match(app, /takbeer/);
assert.match(app, /sautun noor/);
assert.match(app, /masjidCatalogue = sortMasjidCatalogue\(masjids\)/);
assert.match(app, /function orderedFavouriteStreams/);
assert.match(app, /\[\.\.\.favouriteIds\]\.map\(id => streamsByID\.get\(id\)\)/);
assert.match(app, /className = "favourite-move"/);
assert.match(app, /data\.direction|dataset\.direction/);
assert.match(app, /Favourite order saved/);
assert.match(app, /favouriteIds = new Set\(previousIds\)/);
assert.match(html, /app\.js\?v=20260911-stream-sorting/);
assert.match(html, /style\.css\?v=20260911-favourites-order/);
assert.match(css, /\.favourite-order-controls/);
assert.match(css, /\.favourite-move/);

console.log("Favourite ordering and Masjid catalogue sorting tests passed");
