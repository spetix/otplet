async function loadVersions() {
  try {
    const res = await fetch('/versions/versions.json');   // file at gh-pages root
    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }

    const versions = await res.json();  // JSON → JS object [1](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Scripting/JSON)

    const container = document.getElementById('versions-list');
    if (!container) return;

    container.innerHTML = '';

    versions.forEach(v => {
      const li = document.createElement('li');

      const link = document.createElement('a');
      link.href = `/versions/${v.version}/`;
      link.textContent = v.title || v.version;

      if (v.aliases && v.aliases.includes('latest')) {
        link.textContent += ' (latest)';
      }

      li.appendChild(link);
      container.appendChild(li);
    });

  } catch (err) {
    console.error('Failed to load versions:', err);
  }
}

document.addEventListener('DOMContentLoaded', loadVersions);
``