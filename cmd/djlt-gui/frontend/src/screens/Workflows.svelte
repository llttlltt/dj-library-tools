<script>
  import { onMount } from 'svelte';
  import { ListWorkflows, ListSources, SaveWorkflow, DeleteWorkflow } from '../../wailsjs/go/gui/app.js';

  export let onRun = () => {};

  let workflows = [];
  let sources = [];
  let selected = null;
  let error = '';
  let dirty = false;

  async function load() {
    try {
      [workflows, sources] = await Promise.all([ListWorkflows(), ListSources()]);
      workflows = workflows ?? [];
      sources = sources ?? [];
    } catch (e) { error = String(e); }
  }

  onMount(load);

  function selectWorkflow(w) {
    selected = JSON.parse(JSON.stringify(w)); // deep copy
    dirty = false;
  }

  function newWorkflow() {
    selected = { id: '', name: 'New Workflow', steps: [] };
    dirty = true;
  }

  function sourceName(id) {
    const s = sources.find(s => s.id === id);
    return s ? `${s.name} (${s.provider})` : id?.slice(0, 8) + '…';
  }

  function addStep() {
    const sid = sources[0]?.id ?? '';
    selected.steps = [...selected.steps, {
      id: '', kind: 'sync',
      source: { source_id: sid, resource: 'tracks', query: '' },
      targets: [{ source_id: sid, resource: 'playlists', query: '' }],
      after: [], options: {}
    }];
    dirty = true;
  }

  function removeStep(i) {
    selected.steps = selected.steps.filter((_, j) => j !== i);
    dirty = true;
  }

  async function save() {
    error = '';
    try {
      const saved = await SaveWorkflow(selected);
      selected = saved;
      dirty = false;
      await load();
    } catch (e) { error = String(e); }
  }

  async function remove(id) {
    if (!confirm('Delete this workflow?')) return;
    try {
      await DeleteWorkflow(id);
      selected = null;
      await load();
    } catch (e) { error = String(e); }
  }

  const KINDS = ['sync', 'fix', 'edit'];
</script>

<div class="layout">
  <!-- Left: list -->
  <aside>
    <div class="aside-header">
      <span class="section-label">Workflows</span>
      <button class="btn-icon" on:click={newWorkflow} title="New Workflow">+</button>
    </div>
    {#if workflows.length === 0}
      <p class="empty">No workflows yet.</p>
    {:else}
      <ul>
        {#each workflows as w}
          <li class:active={selected?.id === w.id}>
            <button class="list-item" on:click={() => selectWorkflow(w)}>
              <span class="wf-name">{w.name}</span>
              <span class="wf-steps">{w.steps?.length ?? 0} steps</span>
            </button>
          </li>
        {/each}
      </ul>
    {/if}
  </aside>

  <!-- Right: editor -->
  <div class="editor">
    {#if error}<p class="error">{error}</p>{/if}

    {#if !selected}
      <p class="empty" style="margin-top:3rem">Select a workflow or create one.</p>
    {:else}
      <div class="editor-header">
        <input class="wf-title" bind:value={selected.name} on:input={() => dirty = true} />
        <div class="editor-actions">
          {#if selected.id}
            <button class="btn-run" on:click={() => onRun(selected.id)}>▶ Run</button>
            <button class="btn-danger-sm" on:click={() => remove(selected.id)}>Delete</button>
          {/if}
          <button class="btn-primary" disabled={!dirty} on:click={save}>Save</button>
        </div>
      </div>

      <div class="steps">
        {#each selected.steps as step, i}
          <div class="step-card">
            <div class="step-header">
              <span class="step-num">Step {i + 1}</span>
              <select bind:value={step.kind} on:change={() => dirty = true} class="kind-sel">
                {#each KINDS as k}<option>{k}</option>{/each}
              </select>
              <button class="btn-remove" on:click={() => removeStep(i)}>✕</button>
            </div>

            <div class="step-body">
              <div class="ep-group">
                <label>Source</label>
                <div class="ep-row">
                  <select bind:value={step.source.source_id} on:change={() => dirty = true}>
                    {#each sources as s}<option value={s.id}>{sourceName(s.id)}</option>{/each}
                  </select>
                  <input bind:value={step.source.resource} placeholder="resource" on:input={() => dirty = true} />
                  <input bind:value={step.source.query} placeholder="query (optional)" on:input={() => dirty = true} class="grow" />
                </div>
              </div>

              {#if step.kind === 'sync'}
                <div class="ep-group">
                  <label>Targets</label>
                  {#each step.targets as tgt, ti}
                    <div class="ep-row">
                      <select bind:value={tgt.source_id} on:change={() => dirty = true}>
                        {#each sources as s}<option value={s.id}>{sourceName(s.id)}</option>{/each}
                      </select>
                      <input bind:value={tgt.resource} placeholder="resource" on:input={() => dirty = true} />
                      <input bind:value={tgt.query} placeholder="query (optional)" on:input={() => dirty = true} class="grow" />
                      {#if step.targets.length > 1}
                        <button class="btn-remove" on:click={() => { step.targets = step.targets.filter((_,j)=>j!==ti); dirty=true; }}>✕</button>
                      {/if}
                    </div>
                  {/each}
                  <button class="btn-link" on:click={() => { step.targets = [...step.targets, { source_id: sources[0]?.id ?? '', resource: 'playlists', query: '' }]; dirty = true; }}>
                    + Add target
                  </button>
                </div>
              {/if}

              {#if i > 0}
                <div class="ep-group">
                  <label>Run after (step IDs)</label>
                  <input
                    value={step.after?.join(', ') ?? ''}
                    placeholder="Leave blank to run in parallel"
                    on:input={(e) => { step.after = e.target.value.split(',').map(s=>s.trim()).filter(Boolean); dirty=true; }}
                    class="grow"
                  />
                </div>
              {/if}
            </div>
          </div>
        {/each}
        <button class="btn-add-step" on:click={addStep}>+ Add Step</button>
      </div>
    {/if}
  </div>
</div>

<style>
  .layout { display: grid; grid-template-columns: 220px 1fr; gap: 1.5rem; height: calc(100vh - 8rem); }
  aside { background: #1a1a1a; border-radius: 8px; border: 1px solid #2a2a2a; overflow-y: auto; }
  .aside-header { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1rem; border-bottom: 1px solid #2a2a2a; }
  .section-label { font-size: 0.75rem; text-transform: uppercase; color: #666; letter-spacing: 0.05em; }
  ul { list-style: none; }
  li.active .list-item { background: #1e3a5f; }
  .list-item { width: 100%; display: flex; flex-direction: column; align-items: flex-start; gap: 0.1rem; padding: 0.65rem 1rem; background: transparent; border: none; color: inherit; cursor: pointer; text-align: left; }
  .list-item:hover { background: #222; }
  .wf-name { font-size: 0.875rem; font-weight: 500; }
  .wf-steps { font-size: 0.7rem; color: #666; }
  .editor { overflow-y: auto; }
  .editor-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1.5rem; }
  .wf-title { flex: 1; background: transparent; border: none; border-bottom: 1px solid #333; color: #fff; font-size: 1.25rem; font-weight: 600; padding: 0.25rem 0; }
  .editor-actions { display: flex; gap: 0.5rem; }
  .steps { display: flex; flex-direction: column; gap: 1rem; }
  .step-card { background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 8px; }
  .step-header { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; border-bottom: 1px solid #222; }
  .step-num { font-size: 0.75rem; color: #666; font-weight: 600; text-transform: uppercase; }
  .kind-sel { background: #0f0f0f; border: 1px solid #333; color: #e8e8e8; border-radius: 4px; padding: 0.2rem 0.5rem; font-size: 0.8rem; }
  .step-body { padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; }
  .ep-group { display: flex; flex-direction: column; gap: 0.4rem; }
  .ep-row { display: flex; gap: 0.5rem; align-items: center; }
  label { font-size: 0.7rem; color: #666; text-transform: uppercase; letter-spacing: 0.05em; }
  input, select { background: #0f0f0f; border: 1px solid #333; border-radius: 4px; padding: 0.35rem 0.6rem; color: #e8e8e8; font-size: 0.8rem; }
  .grow { flex: 1; }
  .btn-primary { background: #2563eb; color: #fff; border: none; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-primary:disabled { opacity: 0.4; cursor: default; }
  .btn-run { background: #16a34a; color: #fff; border: none; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-danger-sm { background: transparent; color: #f87171; border: 1px solid #7f1d1d; border-radius: 4px; padding: 0.35rem 0.75rem; cursor: pointer; font-size: 0.75rem; }
  .btn-remove { background: transparent; color: #555; border: none; cursor: pointer; font-size: 0.875rem; padding: 0 0.25rem; }
  .btn-remove:hover { color: #f87171; }
  .btn-icon { background: transparent; color: #999; border: 1px solid #333; border-radius: 4px; width: 24px; height: 24px; cursor: pointer; font-size: 1rem; display: flex; align-items: center; justify-content: center; }
  .btn-link { background: none; border: none; color: #60a5fa; cursor: pointer; font-size: 0.8rem; padding: 0.25rem 0; }
  .btn-add-step { width: 100%; padding: 0.65rem; background: #1a1a1a; border: 1px dashed #333; border-radius: 8px; color: #555; cursor: pointer; font-size: 0.875rem; }
  .btn-add-step:hover { border-color: #2563eb; color: #60a5fa; }
  .empty { color: #555; font-style: italic; }
  .error { color: #f87171; font-size: 0.875rem; margin-bottom: 1rem; }
</style>
