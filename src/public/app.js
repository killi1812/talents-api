async function search() {
  const query = document.getElementById('searchInput').value;
  const apiKey = document.getElementById('apiKey').value;
  const shouldSave = document.getElementById('saveKey').checked;
  const resultsDiv = document.getElementById('results');
  const statusDiv = document.getElementById('status');

  if (!apiKey) {
    alert("Please enter API Key");
    return;
  }

  // Save/Remove key based on checkbox
  if (shouldSave) {
    localStorage.setItem('talent_search_key', apiKey);
  } else {
    localStorage.removeItem('talent_search_key');
  }

  statusDiv.innerHTML = "Searching...";
  resultsDiv.innerHTML = "";

  try {
    const response = await fetch(`/api/talents?q=${encodeURIComponent(query)}&limit=100`, {
      headers: { 'X-API-KEY': apiKey }
    });

    if (response.status === 401) {
      statusDiv.innerHTML = "<p style='color:red'>Unauthorized: Invalid API Key</p>";
      return;
    }

    const res = await response.json();
    const data = res.data;
    const total = res.total;

    statusDiv.innerHTML = `<p>Found ${total} talents.</p>`;

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
    statusDiv.innerHTML = "<p style='color:red'>Error fetching data</p>";
  }
}

// Initialization for Search Page
if (document.getElementById('apiKey')) {
  const saved = localStorage.getItem('talent_search_key');
  if (saved) {
    document.getElementById('apiKey').value = saved;
    document.getElementById('saveKey').checked = true;
  }
}

// Admin Functions
async function loadKeys() {
  const apiKey = document.getElementById('adminKey').value;
  const keyList = document.getElementById('keyList');
  const statusDiv = document.getElementById('status');

  if (!apiKey) return;

  try {
    const response = await fetch('/admin/keys', {
      headers: { 'X-API-KEY': apiKey }
    });

    if (response.status !== 200) {
      statusDiv.innerHTML = "<p style='color:red'>Failed to load keys. Check Admin Key.</p>";
      return;
    }

    const keys = await response.json();
    keyList.innerHTML = "";
    statusDiv.innerHTML = "";

    keys.forEach(k => {
      const row = document.createElement('tr');
      const lastUsed = k.lastUsed ? new Date(k.lastUsed).toLocaleString() : 'Never';
      
      row.innerHTML = `
        <td>${k.description || 'No description'}</td>
        <td><code>${k.key}</code></td>
        <td>${lastUsed}</td>
        <td>
          <button class="btn-small" onclick="copyToClipboard('${k.key}')">Copy</button>
          <button class="btn-small delete-btn" onclick="deleteKey('${k.key}')">Delete</button>
        </td>
      `;
      keyList.appendChild(row);
    });
  } catch (err) {
    statusDiv.innerHTML = "<p style='color:red'>Error connecting to server</p>";
  }
}

async function createKey() {
  const adminKey = document.getElementById('adminKey').value;
  const desc = document.getElementById('newKeyDesc').value;

  if (!desc) {
    alert("Please enter a description");
    return;
  }

  try {
    const response = await fetch('/admin/keys', {
      method: 'POST',
      headers: { 
        'X-API-KEY': adminKey,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ description: desc })
    });

    if (response.status === 201) {
      document.getElementById('newKeyDesc').value = "";
      loadKeys();
    } else {
      alert("Failed to create key");
    }
  } catch (err) {
    alert("Error creating key");
  }
}

async function deleteKey(key) {
  if (!confirm(`Are you sure you want to delete this key?`)) return;

  const adminKey = document.getElementById('adminKey').value;
  try {
    const response = await fetch(`/admin/keys/${encodeURIComponent(key)}`, {
      method: 'DELETE',
      headers: { 'X-API-KEY': adminKey }
    });

    if (response.status === 204) {
      loadKeys();
    } else {
      alert("Failed to delete key");
    }
  } catch (err) {
    alert("Error deleting key");
  }
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    alert("Key copied to clipboard");
  });
}
