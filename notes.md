tabbed view gpt example

```html
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Tabbed View</title>
<style>
  /* Style the tab buttons */
  .tab {
    overflow: hidden;
    border: 1px solid #ccc;
    background-color: #f1f1f1;
  }

  /* Style the buttons inside the tab */
  .tab button {
    background-color: inherit;
    float: left;
    border: none;
    outline: none;
    cursor: pointer;
    padding: 14px 16px;
    transition: 0.3s;
  }

  /* Change background color of buttons on hover */
  .tab button:hover {
    background-color: #ddd;
  }

  /* Create an active/current tablink class */
  .tab button.active {
    background-color: #ccc;
  }

  /* Style the tab content */
  .tabcontent {
    display: none;
    padding: 6px 12px;
    border: 1px solid #ccc;
    border-top: none;
  }
</style>
</head>
<body>

<div class="tab">
  <button class="tablinks" onclick="openTab(event, 'Tab1')" id="defaultOpen">Tab 1</button>
  <button class="tablinks" onclick="openTab(event, 'Tab2')">Tab 2</button>
  <button class="tablinks" onclick="openTab(event, 'Tab3')">Tab 3</button>
</div>

<div id="Tab1" class="tabcontent">
  <h3>Tab 1 Content</h3>
  <p>This is the content of tab 1.</p>
</div>

<div id="Tab2" class="tabcontent">
  <h3>Tab 2 Content</h3>
  <p>This is the content of tab 2.</p>
</div>

<div id="Tab3" class="tabcontent">
  <h3>Tab 3 Content</h3>
  <p>This is the content of tab 3.</p>
</div>

<script>
function openTab(evt, tabName) {
  // Get all elements with class="tabcontent" and hide them
  var tabcontent = document.getElementsByClassName("tabcontent");
  for (var i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }

  // Get all elements with class="tablinks" and remove the class "active"
  var tablinks = document.getElementsByClassName("tablinks");
  for (var i = 0; i < tablinks.length; i++) {
    tablinks[i].className = tablinks[i].className.replace(" active", "");
  }

  // Show the current tab, and add an "active" class to the button that opened the tab
  document.getElementById(tabName).style.display = "block";
  evt.currentTarget.className += " active";
}

// Get the element with id="defaultOpen" and click on it
document.getElementById("defaultOpen").click();
</script>

</body>
</html>
```

To get the result of a media query for `prefers-color-scheme` in JavaScript and post it to an endpoint, you can follow these steps:

1. Use `window.matchMedia()` to check if the user's preferred color scheme matches a specific media query.
2. Extract the result of the media query.
3. Send the result to an endpoint using an HTTP request. You can use `fetch()` to perform this.

Here's a basic example:

```javascript
// Check if the user prefers dark mode
const darkModeQuery = window.matchMedia('(prefers-color-scheme: dark)');
const isDarkMode = darkModeQuery.matches;

// Define the data to be sent
const data = {
  isDarkMode: isDarkMode
};

// Define the endpoint URL
const endpoint = 'https://example.com/endpoint';

// Define the fetch options
const options = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(data)
};

// Send the data to the endpoint
fetch(endpoint, options)
  .then(response => {
    if (response.ok) {
      console.log('Data successfully sent to the endpoint');
    } else {
      console.error('Failed to send data to the endpoint');
    }
  })
  .catch(error => {
    console.error('Error:', error);
  });
```

In this example:
- We use `window.matchMedia('(prefers-color-scheme: dark)')` to check if the user prefers dark mode.
- The `matches` property of the media query result indicates whether the condition is met.
- We define the data to be sent as an object with a key `isDarkMode` containing the result of the media query.
- We define the endpoint URL where we want to send the data.
- We use `fetch()` to make a POST request to the endpoint with the data.
- We handle the response from the endpoint to log success or failure.

Ensure to replace `'https://example.com/endpoint'` with the actual endpoint URL where you want to send the data. Also, handle errors and responses appropriately based on your requirements.
