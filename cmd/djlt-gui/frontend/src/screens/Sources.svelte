<script>
  import { onMount } from 'svelte';
  import { ListSources, CreateSource, DeleteSource } from '../../wailsjs/go/gui/app.js';

  let sources = [];
  let error = '';
  let adding = false;
  let form = { name: '', provider: 'rb', file_path: '', host: '', port: '', token: '' };

  async function load() {
    try { sources = await ListSources() ?? []; }
    catch (e) { error = e; }
  }

  onMount(load);

  async function add() {
    error = '';
    const cfg = {};
    if (form.provider === 'rb' || form.provider === 'm3u') {
      cfg.file_path = form.file_path;
    } else if (form.provider === 'plex') {
      cfg.host = form.host;
      cfg.port = form.port;
      cfg.token = form.token;
    }
    try {
      await CreateSource(form.name, form.provider, cfg);
      adding = false;
      form = { name: '', provider: 'rb', file_path: '', host: '', port: '', token: '' };
      await load();
    } catch (e) { error = String(e); }
  }

  async function remove(id) {
    if (!confirm('Delete this source?')) return;
    try { await DeleteSource(id); await load(); }
    catch (e) { error = String(e); }
  }
</script>

<div class="screen">
  <div class="toolbar">
    <h2>Sources</h2>
    <button class="btn-primary" on:click={() => adding = !adding}>
      {adding ? 'Cancel' : '+ Add Source'}
    </button>
  </div>

  {#if error}<p class="error">{error}</p>{/if}

  {#if adding}
    <form class="add-form" on:submit|preventDefault={add}>
      <div class="field">
        <label>Name</label>
        <input bind:value={form.name} placeholder="Main Library" required />
      </div>
      <div class="field">
        <label>Provider</label>
        <select bind:value={form.provider}>
          <option value="rb">Rekordbox</option>
          <option value="m3u">M3U</option>
          <option value="plex">Plex</option>
        </select>
      </div>
      {#if form.provider === 'rb' || form.provider === 'm3u'}
        <div class="field">
          <label>File Path</label>
          <input bind:value={form.file_path} placeholder="/path/to/library.xml" required />
        </div>
      {:else if form.provider === 'plex'}
        <div class="field"><label>Host</label><input bind:value={form.host} placeholder="localhost" /></div>
        <div class="field"><label>Port</label><input bind:value={form.port} placeholder="32400" /></div>
        <div class="field"><label>Token</label><input bind:value={form.token} placeholder="your-plex-token" /></div>
      {/if}
      <button type="submit" class="btn-primary">Save Source</button>
    </form>
  {/if}

  {#if sources.length === 0}
    <p class="empty">No sources configured. Add one above.</p>
  {:else}
    <table>
      <thead><tr><th>Name</th><th>Provider</th><th>ID</th><th></th></tr></thead>
      <tbody>
        {#each sources as s}
          <tr>
            <td class="bold">{s.name}</td>
            <td><span class="badge">{s.provider}</span></td>
            <td class="mono dim">{s.id}</td>
            <td><button class="btn-danger" on:click={() => remove(s.id)}>Delete</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .screen { max-width: 900px; }
  .toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
  h2 { font-size: 1.25rem; font-weight: 600; }
  .add-form { background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 8px; padding: 1rem; margin-bottom: 1.5rem; display: flex; flex-direction: column; gap: 0.75rem; max-width: 480px; }
  .field { display: flex; flex-direction: column; gap: 0.25rem; }
  label { font-size: 0.75rem; color: #999; text-transform: uppercase; letter-spacing: 0.05em; }
  input, select { background: #0f0f0f; border: 1px solid #333; border-radius: 6px; padding: 0.5rem 0.75rem; color: #e8e8e8; font-size: 0.875rem; }
  table { width: 100%; border-collapse: collapse; }
  th { text-align: left; font-size: 0.75rem; color: #666; text-transform: uppercase; letter-spacing: 0.05em; padding: 0.5rem 0.75rem; border-bottom: 1px solid #2a2a2a; }
  td { padding: 0.75rem; border-bottom: 1px solid #1e1e1e; font-size: 0.875rem; }
  .bold { font-weight: 500; }
  .mono { font-family: monospace; font-size: 0.75rem; }
  .dim { color: #666; }
  .badge { background: #1e3a5f; color: #60a5fa; padding: 0.2rem 0.5rem; border-radius: 4px; font-size: 0.75rem; font-weight: 500; }
  .btn-primary { background: #2563eb; color: #fff; border: none; border-radius: 6px; padding: 0.5rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-danger { background: transparent; color: #f87171; border: 1px solid #7f1d1d; border-radius: 4px; padding: 0.25rem 0.75rem; cursor: pointer; font-size: 0.75rem; }
  .empty { color: #555; font-style: italic; margin-top: 2rem; }
  .error { color: #f87171; margin-bottom: 1rem; font-size: 0.875rem; }
</style>
