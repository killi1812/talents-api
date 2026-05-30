async function search() {
  const query = document.getElementById('searchInput').value;
  const apiKey = document.getElementById('apiKey').value;
  const resultsDiv = document.getElementById('results');

  if (!apiKey) {
    alert("Please enter API Key");
    return;
  }

  resultsDiv.innerHTML = "Searching...";

  try {
    const response = await fetch(`/api/talents?q=${encodeURIComponent(query)}`, {
      headers: { 'X-API-KEY': apiKey }
    });

    if (response.status === 401) {
      resultsDiv.innerHTML = "<p style='color:red'>Unauthorized: Invalid API Key</p>";
      return;
    }

    const data = await response.json();
    resultsDiv.innerHTML = "";

    if (data.length === 0) {
      resultsDiv.innerHTML = "<p>No talents found.</p>";
      return;
    }

    data.forEach(talent => {
      const card = document.createElement('div');
      card.className = 'talent-card';
      card.innerHTML = `
                <h3>${talent.name}</h3>
                <p><strong>Version:</strong> ${talent.version} | <strong>Max Level:</strong> ${talent.maxLvl}</p>
                <p>${talent.description}</p>
            `;
      resultsDiv.appendChild(card);
    });
  } catch (error) {
    console.error(error)
    resultsDiv.innerHTML = "<p style='color:red'>Error fetching data</p>";
  }
}
