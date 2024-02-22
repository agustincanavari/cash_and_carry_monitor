function fetchData() {
    fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            const dataDiv = document.getElementById('data');
            dataDiv.innerHTML = '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
        })
        .catch(error => console.error('Error fetching data:', error));
}

// Fetch data on initial load
fetchData();
