# MasjidBoard Domain Model

**Status:** Implemented in MasjidPi v1.5.2

The provider translates positional MasjidBoard Live data into semantic Go types. The normalized model does not expose upstream row numbers to the rest of the application.

## Board

A normalized Board contains:

- identity and timezone;
- date context;
- five daily prayers;
- optional Jumu'ah services;
- optional special Dhuhr;
- optional astronomical times;
- announcements and programmes;
- structured notices;
- media metadata;
- contribution/banking data; and
- new-moon data.

Identity and the daily timetable form the core. Other content is optional and may be absent without invalidating the board.

## Identity and date context

Board identity stores the provider ID, display name, alternate name and timezone.

Date context stores the Gregorian source date plus Islamic-date adjustment information. This lets a cached board follow the masjid's sunset-based Islamic-date rollover without requiring a provider refresh at sunset.

## Prayer times

The five fixed prayers are:

- Fajr
- Dhuhr
- Asr
- Maghrib
- Esha

Each prayer may contain an Adhan and Jamaah local wall-clock time. Missing values remain absent rather than being guessed.

Special Dhuhr stores a provider time and its applicability label, such as “Sundays & Public Holidays”. Presentation suppresses it when it duplicates normal Dhuhr.

## Jumu'ah

Each JumuahService may contain:

- dedicated Adhan and Jamaah;
- provider-configured heading/time events;
- optional Khateeb or explanatory text; and
- Islamic-time representations retained separately from civil clock times.

The three source heading/time pairs describe stages of a service, not three separate services. A Khutbah time is not reclassified as Salaah when no explicit Jamaah/Salaah value exists.

Deprecated alternate-time fields remain only for decoding older caches.

## Astronomical times

Optional astronomical values include:

- Suhur/Sehri end;
- Fajr start;
- Sunrise;
- Ishraaq;
- Duha/Chaasht;
- Istiwaa caution;
- Istiwaa;
- Zawaal end;
- Asr Shafi'i;
- Asr Hanafi;
- Sunset; and
- Esha start.

Presentation uses the Istiwaa caution-to-Zawaal-end interval for the warning state, falling back to Istiwaa when the caution value is absent.

## Community content

Announcements contain a conservative semantic category, title, content and optional visibility window.

Programmes contain title/content plus optional schedule and visibility values.

Structured notices use a finite notice type and a field map for provider-specific details. Current types include general, Nikah, funeral, well-wishes, Eid, Salaah change, Dawah and three-day Jamaat. Taleem and contribution data have their own normalized representations.

Unknown or unconfirmed source semantics are not promoted into invented fields.

## Media, banking and new moon

Media metadata provides source identity, retrieval/local paths, type, hash, visibility and display duration. Media presentation is not yet enabled.

Banking and new-moon structures preserve known provider fields without making them prerequisites for timetable use.

## Shared content

Daily Ayah, Hadith and Sunnah and Islamic Economic Indicators are shared content sources, not properties of an individual selected masjid. They use separate models/caches and are joined only in the display view.

## Display model

The internal model is not returned directly to browsers. The display package creates a smaller JSON contract with:

- local clock values as hour/minute objects;
- stable prayer keys and labels;
- only presentation-relevant optional content;
- per-board status; and
- selected saved display preferences.

Provider diagnostics, raw rows and cache implementation details do not cross the display boundary.

## Optionality and cache rules

- Missing optional content does not invalidate core prayer data.
- Invalid core data is not cached as a successful board.
- One board's failure does not invalidate another board.
- Older compatible cache fields are migrated during presentation where required.
- No layer invents times or content that the source did not provide.

The implemented Go structures are authoritative when this overview and source code diverge.
