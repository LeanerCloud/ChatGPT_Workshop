

// popup.js
document.addEventListener('DOMContentLoaded', () => {
    const sourceCodeInput = document.getElementById('sourceCode');
    const languageSelect = document.getElementById('languageSelect');
    const convertButton = document.getElementById('convertButton');
    const resultCodeInput = document.getElementById('resultCode');

    convertButton.addEventListener('click', () => {
        const code = sourceCodeInput.value;
        const targetLanguage = languageSelect.value;

        chrome.runtime.sendMessage({ action: 'translateCode', code, targetLanguage }, (response) => {
            resultCodeInput.value = response.resultCode;
        });
    });
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === "translateCode") {
        document.getElementById("sourceCode").value = message.code;
        // Trigger the translation here if needed
    }
});
