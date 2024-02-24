function fetchData() {
    fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('data');
            tableBody.innerHTML = ''; // Clear existing table data
            data.forEach(calc => {
                // Main row for the spot symbol
                const mainRow = document.createElement('tr');
                mainRow.innerHTML = `
                    <td>${calc.spotSymbol}</td>
                    <td>${calc.spotPrice}</td>
                    <td colspan="6"></td> <!-- Empty cells for futures data -->
                `;
                tableBody.appendChild(mainRow);

                // Sub-rows for each future
                calc.futures.forEach((future, index) => {
                    const futureRow = document.createElement('tr');
                    futureRow.className = index % 2 === 0 ? 'even' : 'odd'; // Alternating row styles
                    futureRow.innerHTML = `
                        <td></td> <!-- Empty cell for spot symbol indentation -->
                        <td></td> <!-- Empty cell for spot price indentation -->
                        <td>${future.futureSymbol}</td>
                        <td>${future.futurePrice}</td>
                        <td>${future.settlementDate}</td>
                        <td>${future.apr.toFixed(2)}%</td>
                        <td>${future.apy.toFixed(2)}%</td>
                        <td>${future.yield.toFixed(2)}%</td>
                    `;
                    tableBody.appendChild(futureRow);
                });
            });
        })
        .catch(error => console.error('Error fetching data:', error));
}

// Fetch data on initial load
fetchData();
