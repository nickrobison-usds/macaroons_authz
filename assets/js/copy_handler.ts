$(document).ready(handleCopyClick)


function handleCopyClick($: JQueryStatic): void {
    console.log("Attaching copy handler");
    const copyable = $(".copyable");
    copyable.click(e => clickHandler(e));
    copyable.hover(hoverInHandler, hoverOutHandler);
}

function clickHandler(value: JQuery.Event<HTMLElement, null>): void {
    const el = value.currentTarget;
    const range = document.createRange();
    range.selectNodeContents(el);
    const sel = window.getSelection();
    sel.removeAllRanges();
    sel.addRange(range);
    document.execCommand("copy");
    sel.removeAllRanges();
    console.log(`Copied ${sel} to clipboard.`);
}

function hoverInHandler(value: JQuery.Event<HTMLElement, null>): void {
    console.log("Hovering in.");
    $(value.currentTarget).siblings("i:first").show();
}

function hoverOutHandler(value: JQuery.Event<HTMLElement, null>): void {
    console.log("Hovering out.");
    $(value.currentTarget).siblings("i:first").hide();
}
