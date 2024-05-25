const
    v1 = "#bd93f9",
    v2 = "#61A151",
    v3 = "#6A8759",
    v4 = "#A9B7C6",
    v5 = "#FFC66D",
    v6 = "#CC7832",
    v7 = "#ff79c6",
    v8 = "#ffb86c",
    stone = "#7d8799", // Brightened compared to original to increase contrast
    darkBackground = "#21252b",
    highlightBackground = "#2c313a",
    background = "#2B2B2B",
    tooltipBackground = "#353a42",
    selection = "rgba(255, 255, 255, 0.1)",
    cursor = "#FFFFFF",
    selectionMatch = 'rgba(255, 255, 255, 0.2)',
    gutterBackground = 'rgba(255, 255, 255, 0.1)',
    gutterForeground = '#999',
    gutterBorder = 'transparent',
    lineHighlight = 'rgba(255, 255, 255, 0.1)'

const darculaTheme = {
    "&": {
        color: v4,
        backgroundColor: background
    },

    // comment
    ".ͼm": {
        color: v2,
    },

    // function
    ".ͼ6c": {
        color: v6,
    },
    ".ͼb": {
        color: v6,
    },

    //function args
    ".ͼ69": {
        color: v4,
    },
    ".ͼg": {
        color: v4,
    },

    //string
    ".ͼe": {
        color: v3,
    },

    // numbers
    ".ͼd": {
        color: v1
    },

    // object keys
    ".ͼl": {
        color: v5
    },
    ".ͼc": {
        color: v4
    },

    ".cm-content": {
        caretColor: cursor
    },

    ".cm-cursor, .cm-dropCursor": {borderLeftColor: cursor},
    "&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection": {backgroundColor: selection},

    ".cm-panels": {backgroundColor: darkBackground, color: v4},
    ".cm-panels.cm-panels-top": {borderBottom: "2px solid black"},
    ".cm-panels.cm-panels-bottom": {borderTop: "2px solid black"},

    ".cm-searchMatch": {
        backgroundColor: "#72a1ff59",
        outline: "1px solid #457dff"
    },
    ".cm-searchMatch.cm-searchMatch-selected": {
        backgroundColor: "#6199ff2f"
    },

    ".cm-activeLine": {backgroundColor: "rgba(255, 255, 255, 0.1)"},
    ".cm-selectionMatch": {backgroundColor: "rgba(170,254,102,0.41)"},

    "&.cm-focused .cm-matchingBracket, &.cm-focused .cm-nonmatchingBracket": {
        backgroundColor: "#bad0f847"
    },

    ".cm-gutters": {
        backgroundColor: background,
        color: stone,
        border: "none"
    },

    ".cm-activeLineGutter": {
        backgroundColor: highlightBackground
    },

    ".cm-foldPlaceholder": {
        backgroundColor: "transparent",
        border: "none",
        color: "#ddd"
    },

    ".cm-tooltip": {
        border: "none",
        backgroundColor: tooltipBackground
    },
    ".cm-tooltip .cm-tooltip-arrow:before": {
        borderTopColor: "transparent",
        borderBottomColor: "transparent"
    },
    ".cm-tooltip .cm-tooltip-arrow:after": {
        borderTopColor: tooltipBackground,
        borderBottomColor: tooltipBackground
    },
    ".cm-tooltip-autocomplete": {
        "& > ul > li[aria-selected]": {
            backgroundColor: highlightBackground,
            color: v4
        }
    }
}

export {darculaTheme}
