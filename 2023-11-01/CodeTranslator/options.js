document.addEventListener('DOMContentLoaded', () => {
    const saveButton = document.getElementById('save');
    const apiKeyInput = document.getElementById('api-key');
    const toggleVisibilityButton = document.getElementById('toggle-visibility');

    // Load the current API key from preferences
    chrome.storage.sync.get('openAiApiKey', ({ openAiApiKey }) => {
        if (openAiApiKey) {
            apiKeyInput.value = openAiApiKey;
        }
    });

    // Save the API key to preferences when the save button is clicked
    saveButton.addEventListener('click', () => {
        chrome.storage.sync.set({ openAiApiKey: apiKeyInput.value }, () => {
            console.log('API Key saved');
        });
    });

    // Toggle the visibility of the API key
    toggleVisibilityButton.addEventListener('click', () => {
        if (apiKeyInput.type === 'password') {
            apiKeyInput.type = 'text';
            toggleVisibilityButton.textContent = 'Hide';
        } else {
            apiKeyInput.type = 'password';
            toggleVisibilityButton.textContent = 'Show';
        }
    });
});
