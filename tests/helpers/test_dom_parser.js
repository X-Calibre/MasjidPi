"use strict";

function textContent(source) {
    let result = "";
    let insideTag = false;

    for (const character of String(source)) {
        if (character === "<") {
            insideTag = true;
        } else if (character === ">") {
            insideTag = false;
        } else if (!insideTag) {
            result += character;
        }
    }

    return result;
}

class TestDOMParser {
    parseFromString(source) {
        return {body: {textContent: textContent(source)}};
    }
}

module.exports = TestDOMParser;
