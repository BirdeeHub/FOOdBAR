document.addEventListener('htmx:configRequest', (event) => {
  event.detail.parameters['responseType'] = 'json';
})
