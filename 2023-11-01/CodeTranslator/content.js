// content.js
console.log("Content script loaded!");

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === "translateCode") {
        // Handle the code translation here or send it to the popup
        console.log("Selected Code:", message.code);
        sendResponse({ resultCode: "Translation Result Here" });
    }
});
