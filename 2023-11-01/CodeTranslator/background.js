// background.js

chrome.runtime.onInstalled.addListener(() => {
    // Create a context menu item for code translation
    chrome.contextMenus.create({
        id: "translateCode",
        title: "Translate Code",
        contexts: ["selection"]
    });
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
    if (info.menuItemId === "translateCode" && info.selectionText) {
        // Trigger code translation via context menu
        // Assuming that you have set a default target language for simplicity
        const defaultTargetLanguage = "python";
        const message = {
            action: 'translateCode',
            code: info.selectionText,
            targetLanguage: defaultTargetLanguage
        };

        chrome.tabs.sendMessage(tab.id, message, (response) => {
            if (chrome.runtime.lastError) {
                console.error(chrome.runtime.lastError.message);
                return;
            }
            if (response && response.resultCode) {
                // Do something with the translated code
                console.log('Translated Code:', response.resultCode);
            } else {
                console.error('Translation failed:', response.error);
            }
        });
    }
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'translateCode') {
        chrome.storage.sync.get('openAiApiKey', async ({ openAiApiKey }) => {
            const promptMessage = `Translate the following code to ${message.targetLanguage} without including any additional text or formatting:\n\n${message.code}`;

            const data = {
                model: "gpt-4",
                messages: [
                    { role: "user", content: promptMessage }
                ],
                temperature: 0.7
            };

            try {
                const response = await fetch('https://api.openai.com/v1/chat/completions', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${openAiApiKey}`
                    },
                    body: JSON.stringify(data)
                });

                const result = await response.json();
                if (result.choices && result.choices.length > 0) {
                    const translatedCode = result.choices[0].message.content.trim().replace(/`+/g, '');
                    sendResponse({ resultCode: translatedCode });
                } else {
                    console.error('No choices returned from OpenAI:', result);
                    sendResponse({ resultCode: null, error: 'No choices returned from OpenAI' });
                }
            } catch (error) {
                console.error('Error during translation:', error);
                sendResponse({ resultCode: null, error: 'Error during translation' });
            }
        });
        return true; // Keep the message channel open for async sendResponse
    }
});
