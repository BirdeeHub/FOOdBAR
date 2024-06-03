document.addEventListener("DOMContentLoaded", (event) => {
    document.body.addEventListener("htmx:beforeSwap", function(evt) {
        if (evt.detail.xhr.status === 422) {
            evt.detail.shouldSwap = true;
            evt.detail.isError = false;
        }
    });
})
document.addEventListener('htmx:configRequest', (event) => {
    function generateUniqueID() {
        var id = Date.now().toString(36) + Math.random().toString(36).substr(2, 5);
        sessionStorage.setItem('tab_id', id);
        return id;
    }
    var tabId = sessionStorage.getItem('tab_id') || generateUniqueID();
    event.detail.headers['tab_id'] = tabId;
});
document.addEventListener('htmx:afterRequest', (event) => {
    const tabIdHeader = event.detail.xhr.getResponseHeader('tab_id');
    if (tabIdHeader) {
        sessionStorage.setItem('tab_id', tabIdHeader);
    }
});
document.addEventListener('myExtraFieldEvent', (event) => {
    console.log("myExtraFieldEvent");
    console.log(event.detail.id);
    console.log(event.detail.value);
});
