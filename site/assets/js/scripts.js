
fetch('/versions.json')
  .then(r => r.json())
  .then(versions => {
    const select = document.getElementById('versions');

    versions.forEach(v => {
      const opt = document.createElement('option');
      opt.value = `/docs/${v.version}/`;
      opt.textContent = v.title;
      select.appendChild(opt);
    });

    select.addEventListener('change', () => {
      location.href = select.value;
    });
  })
  .catch(() => {
    console.error('Errore nel caricamento di versions.json');
  });