<script>
  import { onMount } from 'svelte';
  import {
    ListWorkflows, ListSources, GetWorkflow,
    SaveWorkflow, DeleteWorkflow,
    GetWorkflowDiff, RunWorkflow
  } from '../../wailsjs/go/gui/app.js';

  // ── props ──────────────────────────────────────────────────────────────────
  // workflowId: currently open workflow (null = list view)
  // No external onBack needed; back is internal to this component.

  // ── state ──────────────────────────────────────────────────────────────────
  let workflows   = [];   // list
  let sources     = [];   // for dropdowns
  let selected    = null; // current workflow object (deep copy for editing)
  let mode        = 'list'; // list | view | edit | preview | applied
  let diffs       = [];   // StepDiff[]
  let runResult   = null; // WorkflowResult
  let loading     = false;
  let error       = '';

  // per-step toggle: show/hide unchanged tracks
  let showUnchanged = {};

  // ── lifecycle ──────────────────────────────────────────────────────────────
  async function loadList() {
    try { [workflows, sources] = await Promise.all([ListWorkflows(), ListSources()]); }
    catch (e) { error = String(e); }
    workflows = workflows ?? [];
    sources   = sources   ?? [];
  }

  onMount(loadList);

  // ── navigation ─────────────────────────────────────────────────────────────
  async function openWorkflow(w) {
    error = '';
    try { selected = JSON.parse(JSON.stringify(await GetWorkflow(w.id))); }
    catch (e) { error = String(e); return; }
    diffs = []; runResult = null; showUnchanged = {};
    mode = 'view';
  }

  function backToList() {
    selected = null; diffs = []; runResult = null; mode = 'list'; error = '';
  }

  // ── list actions ───────────────────────────────────────────────────────────
  async function newWorkflow() {
    selected = { id: '', name: 'New Workflow', steps: [] };
    diffs = []; runResult = null;
    mode = 'edit';
  }

  async function deleteWorkflow(w) {
    if (!confirm(`Delete "${w.name}"?`)) return;
    try { await DeleteWorkflow(w.id); await loadList(); }
    catch (e) { error = String(e); }
  }

  // ── edit helpers ───────────────────────────────────────────────────────────
  function startEdit() {
    selected = JSON.parse(JSON.stringify(selected)); // fresh deep copy
    mode = 'edit';
  }

  function cancelEdit() {
    // discard edits — reload from disk
    if (selected.id) {
      GetWorkflow(selected.id)
        .then(w => { selected = w; mode = 'view'; })
        .catch(e => { error = String(e); backToList(); });
    } else {
      backToList();
    }
  }

  async function saveWorkflow() {
    error = ''; loading = true;
    try {
      selected = await SaveWorkflow(selected);
      await loadList();
      mode = 'view';
    } catch (e) { error = String(e); }
    loading = false;
  }

  function addStep() {
    const sid = sources[0]?.id ?? '';
    selected.steps = [...selected.steps, {
      id: '', kind: 'sync',
      source:  { source_id: sid, resource: 'tracks',    query: '' },
      targets: [{ source_id: sid, resource: 'playlists', query: '' }],
      after: [], options: {}
    }];
  }

  function removeStep(i) {
    selected.steps = selected.steps.filter((_, j) => j !== i);
  }

  function addTarget(step) {
    step.targets = [...step.targets, { source_id: sources[0]?.id ?? '', resource: 'playlists', query: '' }];
    selected = selected; // trigger reactivity
  }

  function removeTarget(step, ti) {
    step.targets = step.targets.filter((_, j) => j !== ti);
    selected = selected;
  }

  // ── preview / apply ────────────────────────────────────────────────────────
  async function runPreview() {
    error = ''; loading = true; mode = 'preview';
    diffs = []; runResult = null; showUnchanged = {};
    try { diffs = await GetWorkflowDiff(selected.id); }
    catch (e) { error = String(e); mode = 'view'; }
    loading = false;
  }

  async function applyWorkflow() {
    if (!confirm('Apply all changes? This cannot be undone.')) return;
    error = ''; loading = true;
    try { runResult = await RunWorkflow(selected.id); mode = 'applied'; }
    catch (e) { error = String(e); }
    loading = false;
  }

  // ── helpers ────────────────────────────────────────────────────────────────
  function sourceName(id) {
    const s = sources.find(s => s.id === id);
    return s ? s.name : id?.slice(0, 8) + '…';
  }

  function endpointLabel(ep) {
    return `${sourceName(ep.source_id)} / ${ep.resource}${ep.query ? ' · ' + ep.query : ''}`;
  }

  // diffs keyed by step_id for quick lookup in preview/applied mode
  $: diffByStep = Object.fromEntries((diffs ?? []).map(d => [d.step_id, d]));

  // result keyed by step_id
  $: resultByStep = Object.fromEntries(
    (runResult?.steps ?? []).map(r => [r.step_id, r])
  );

  // toolbar button states
  $: canEdit    = mode === 'view' || mode === 'preview' || mode === 'applied';
  $: canPreview = (mode === 'view' || mode === 'applied') && !!selected?.id;
  $: canApply   = mode === 'preview' && !loading && diffs.length > 0;

  const KIND_COLOURS = { sync: '#2563eb', fix: '#d97706', edit: '#7c3aed' };
  const kindColour   = k => KIND_COLOURS[k?.toLowerCase()] ?? '#555';

  function statusColour(s) {
    return { success: '#22c55e', failed: '#ef4444', blocked: '#a855f7', running: '#f59e0b' }[s] ?? '#555';
  }

  const KINDS = ['sync', 'fix', 'edit'];
</script>

<!-- ══════════════════════════════════════════════════════════════════════ -->
<!--  LIST VIEW                                                            -->
<!-- ══════════════════════════════════════════════════════════════════════ -->
{#if mode === 'list'}
  <div class="list-screen">
    <div class="list-header">
      <h2>Workflows</h2>
      <button class="btn-primary" on:click={newWorkflow}>+ New Workflow</button>
    </div>
    {#if error}<p class="error">{error}</p>{/if}
    {#if workflows.length === 0}
      <p class="empty">No workflows yet. Create one to get started.</p>
    {:else}
      <div class="wf-cards">
        {#each workflows as w}
          <div class="wf-card" role="button" tabindex="0"
               on:click={() => openWorkflow(w)}
               on:keydown={e => e.key === 'Enter' && openWorkflow(w)}>
            <div class="wf-card-body">
              <span class="wf-name">{w.name}</span>
              <span class="wf-meta">{w.steps?.length ?? 0} step{w.steps?.length !== 1 ? 's' : ''}</span>
            </div>
            <button class="btn-delete" on:click|stopPropagation={() => deleteWorkflow(w)}>Delete</button>
          </div>
        {/each}
      </div>
    {/if}
  </div>

<!-- ══════════════════════════════════════════════════════════════════════ -->
<!--  DETAIL VIEW  (view / edit / preview / applied)                       -->
<!-- ══════════════════════════════════════════════════════════════════════ -->
{:else if selected}
  <div class="detail-screen">

    <!-- toolbar -->
    <div class="toolbar">
      <button class="btn-back" on:click={backToList}>← Workflows</button>

      {#if mode === 'edit'}
        <input class="wf-title-input" bind:value={selected.name} placeholder="Workflow name" />
      {:else}
        <h2 class="wf-title">{selected.name}</h2>
      {/if}

      <div class="spacer" />

      {#if error}<span class="error-inline">{error}</span>{/if}

      {#if mode === 'edit'}
        <button class="btn-ghost" on:click={cancelEdit} disabled={loading}>Cancel</button>
        <button class="btn-primary" on:click={saveWorkflow} disabled={loading}>
          {loading ? 'Saving…' : 'Save'}
        </button>
      {:else}
        <button class="btn-ghost"    on:click={startEdit}       disabled={!canEdit || loading}>Edit</button>
        <button class="btn-preview"  on:click={runPreview}      disabled={!canPreview || loading}>
          {loading && mode === 'preview' ? 'Loading…' : '🔍 Preview'}
        </button>
        <button class="btn-apply"    on:click={applyWorkflow}   disabled={!canApply || loading}>
          {loading && mode === 'applied' ? 'Applying…' : '▶ Apply'}
        </button>
      {/if}
    </div>

    <!-- step list -->
    <div class="steps-list">

      {#if selected.steps.length === 0}
        <p class="empty" style="padding: 2rem 0">
          {mode === 'edit' ? 'No steps yet — add one below.' : 'This workflow has no steps.'}
        </p>
      {/if}

      {#each selected.steps as step, i}
        {@const diff     = diffByStep[step.id]}
        {@const result   = resultByStep[step.id]}

        <div class="step-card"
             class:step-card--edit={mode === 'edit'}
             class:step-card--success={result?.status === 'success'}
             class:step-card--failed={result?.status === 'failed'}
             class:step-card--blocked={result?.status === 'blocked'}>

          <!-- ── step header ─────────────────────────────────────────── -->
          <div class="step-header">
            <span class="step-index">{i + 1}</span>

            {#if mode === 'edit'}
              <select class="kind-select" bind:value={step.kind}>
                {#each KINDS as k}
                  <option value={k}>{k.toUpperCase()}</option>
                {/each}
              </select>
            {:else}
              <span class="kind-badge" style="background: {kindColour(step.kind)}20; color: {kindColour(step.kind)}; border-color: {kindColour(step.kind)}40">
                {step.kind.toUpperCase()}
              </span>
            {/if}

            {#if result}
              <span class="status-pill" style="color:{statusColour(result.status)}">
                { result.status === 'success' ? '✓ done'
                : result.status === 'failed'  ? '✗ failed'
                : result.status === 'blocked' ? '⊘ blocked'
                : result.status }
              </span>
            {/if}

            <div class="step-header-spacer" />

            {#if mode === 'edit'}
              <button class="btn-remove-step" on:click={() => removeStep(i)}>Remove</button>
            {/if}
          </div>

          <!-- ── view / preview / applied: read-only endpoint display ── -->
          {#if mode !== 'edit'}
            <div class="step-body">
              <div class="endpoint-row">
                <span class="ep-label">FROM</span>
                <span class="ep-value">{endpointLabel(step.source)}</span>
              </div>

              {#each step.targets as tgt, ti}
                <div class="endpoint-row">
                  <span class="ep-label">TO{step.targets.length > 1 ? ' ' + (ti + 1) : ''}</span>
                  <span class="ep-value">{endpointLabel(tgt)}</span>
                </div>

                <!-- diff rows (preview + applied modes) -->
                {#if (mode === 'preview' || mode === 'applied') && diff}
                  <div class="diff-block">
                    {#if diff.added.length === 0 && diff.removed.length === 0}
                      <div class="diff-uptodate">✓ Already up to date</div>
                    {:else}
                      {#each diff.added as t}
                        <div class="diff-row diff-add">
                          <span class="diff-sign">+</span>
                          <span class="diff-title">{t.title || t.id}</span>
                          <span class="diff-artist">{t.artist}</span>
                          {#if t.bpm}<span class="diff-bpm">{t.bpm}</span>{/if}
                        </div>
                      {/each}
                      {#each diff.removed as t}
                        <div class="diff-row diff-remove">
                          <span class="diff-sign">−</span>
                          <span class="diff-title">{t.title || t.id}</span>
                          <span class="diff-artist">{t.artist}</span>
                          {#if t.bpm}<span class="diff-bpm">{t.bpm}</span>{/if}
                        </div>
                      {/each}
                      {#if diff.unchanged.length > 0}
                        <button class="btn-toggle-unchanged"
                                on:click={() => showUnchanged[diff.step_id + ti] = !showUnchanged[diff.step_id + ti]}>
                          {showUnchanged[diff.step_id + ti]
                            ? 'Hide unchanged'
                            : `Show ${diff.unchanged.length} unchanged`}
                        </button>
                        {#if showUnchanged[diff.step_id + ti]}
                          {#each diff.unchanged as t}
                            <div class="diff-row diff-unchanged">
                              <span class="diff-sign">·</span>
                              <span class="diff-title">{t.title || t.id}</span>
                              <span class="diff-artist">{t.artist}</span>
                              {#if t.bpm}<span class="diff-bpm">{t.bpm}</span>{/if}
                            </div>
                          {/each}
                        {/if}
                      {/if}
                    {/if}
                  </div>
                {/if}
              {/each}

              <!-- error / warning output from run -->
              {#if result?.error}
                <div class="result-error">✗ {result.error}</div>
              {/if}
              {#if result?.warnings?.length}
                {#each result.warnings as w}
                  <div class="result-warning">⚠ {w}</div>
                {/each}
              {/if}
              {#if result?.status === 'blocked'}
                <div class="result-blocked">⊘ Blocked — a dependency failed.</div>
              {/if}
            </div>

          <!-- ── edit mode: inline fields ──────────────────────────── -->
          {:else}
            <div class="step-body step-body--edit">

              <!-- source endpoint -->
              <div class="edit-section-label">Source</div>
              <div class="ep-edit-row">
                <select bind:value={step.source.source_id}>
                  {#each sources as s}<option value={s.id}>{s.name}</option>{/each}
                </select>
                <input bind:value={step.source.resource} placeholder="resource (e.g. tracks)" />
                <input bind:value={step.source.query} placeholder="query (optional)" class="grow" />
              </div>

              <!-- target endpoints -->
              <div class="edit-section-label" style="margin-top:0.75rem">
                Target{step.targets.length > 1 ? 's' : ''}
              </div>
              {#each step.targets as tgt, ti}
                <div class="ep-edit-row">
                  <select bind:value={tgt.source_id}>
                    {#each sources as s}<option value={s.id}>{s.name}</option>{/each}
                  </select>
                  <input bind:value={tgt.resource} placeholder="resource (e.g. playlists)" />
                  <input bind:value={tgt.query} placeholder="query (optional)" class="grow" />
                  {#if step.targets.length > 1}
                    <button class="btn-icon-sm" on:click={() => removeTarget(step, ti)}>✕</button>
                  {/if}
                </div>
              {/each}
              <button class="btn-link" on:click={() => addTarget(step)}>+ Add target</button>

              <!-- after dependencies -->
              {#if i > 0}
                <div class="edit-section-label" style="margin-top:0.75rem">Run after (step IDs)</div>
                <input
                  class="grow"
                  value={step.after?.join(', ') ?? ''}
                  placeholder="Leave blank to run in parallel"
                  on:input={e => { step.after = e.target.value.split(',').map(s => s.trim()).filter(Boolean); }}
                />
              {/if}
            </div>
          {/if}

        </div><!-- /step-card -->
      {/each}

      {#if mode === 'edit'}
        <button class="btn-add-step" on:click={addStep}>+ Add Step</button>
      {/if}

    </div><!-- /steps-list -->
  </div><!-- /detail-screen -->
{/if}

<style>
  /* ── layout ─────────────────────────────────────────────────────────── */
  .list-screen   { max-width: 760px; }
  .detail-screen { max-width: 800px; display: flex; flex-direction: column; gap: 0; }

  /* ── list ────────────────────────────────────────────────────────────── */
  .list-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; }
  h2 { font-size: 1.25rem; font-weight: 600; }
  .wf-cards { display: flex; flex-direction: column; gap: 0.5rem; }
  .wf-card {
    display: flex; align-items: center; justify-content: space-between;
    background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 10px;
    padding: 1rem 1.25rem; cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
  }
  .wf-card:hover { background: #222; border-color: #3a3a3a; }
  .wf-card-body { display: flex; flex-direction: column; gap: 0.2rem; }
  .wf-name { font-size: 0.95rem; font-weight: 500; }
  .wf-meta { font-size: 0.75rem; color: #666; }

  /* ── toolbar ─────────────────────────────────────────────────────────── */
  .toolbar {
    display: flex; align-items: center; gap: 0.6rem;
    padding: 0 0 1.25rem 0; flex-wrap: wrap;
  }
  .wf-title       { font-size: 1.1rem; font-weight: 600; margin: 0; }
  .wf-title-input {
    font-size: 1.1rem; font-weight: 600; background: transparent;
    border: none; border-bottom: 1px solid #444; color: #fff;
    padding: 0.15rem 0.25rem; flex: 1; min-width: 160px;
  }
  .spacer { flex: 1; }
  .error-inline { font-size: 0.8rem; color: #f87171; }

  /* ── step cards ──────────────────────────────────────────────────────── */
  .steps-list { display: flex; flex-direction: column; gap: 0.75rem; }

  .step-card {
    background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 10px;
    overflow: hidden; transition: border-color 0.2s;
  }
  .step-card--edit    { border-color: #3a3a3a; }
  .step-card--success { border-color: #166534; }
  .step-card--failed  { border-color: #7f1d1d; }
  .step-card--blocked { border-color: #581c87; }

  .step-header {
    display: flex; align-items: center; gap: 0.6rem;
    padding: 0.7rem 1rem; background: #111; border-bottom: 1px solid #222;
  }
  .step-index {
    width: 22px; height: 22px; border-radius: 50%;
    background: #2a2a2a; color: #888;
    font-size: 0.7rem; font-weight: 700;
    display: flex; align-items: center; justify-content: center; flex-shrink: 0;
  }
  .kind-badge {
    font-size: 0.65rem; font-weight: 700; letter-spacing: 0.06em;
    padding: 0.2rem 0.55rem; border-radius: 4px; border: 1px solid;
  }
  .kind-select {
    background: #0f0f0f; border: 1px solid #333; color: #e8e8e8;
    border-radius: 4px; padding: 0.2rem 0.4rem; font-size: 0.8rem;
  }
  .status-pill { font-size: 0.75rem; font-weight: 600; }
  .step-header-spacer { flex: 1; }

  /* ── step body ────────────────────────────────────────────────────────── */
  .step-body { padding: 0.9rem 1rem; display: flex; flex-direction: column; gap: 0.45rem; }
  .step-body--edit { gap: 0.5rem; }

  .endpoint-row { display: flex; align-items: baseline; gap: 0.6rem; }
  .ep-label {
    font-size: 0.6rem; font-weight: 700; letter-spacing: 0.08em;
    color: #555; text-transform: uppercase; min-width: 32px; flex-shrink: 0;
  }
  .ep-value { font-size: 0.8rem; color: #ccc; }

  /* ── diff block ──────────────────────────────────────────────────────── */
  .diff-block {
    margin: 0.35rem 0 0.1rem 38px;
    border-left: 2px solid #2a2a2a; padding-left: 0.75rem;
    display: flex; flex-direction: column; gap: 0.15rem;
  }
  .diff-row {
    display: flex; align-items: baseline; gap: 0.5rem;
    font-size: 0.78rem; padding: 0.15rem 0.4rem; border-radius: 3px;
  }
  .diff-add     { border-left: 3px solid #16a34a; background: #052e1680; }
  .diff-remove  { border-left: 3px solid #dc2626; background: #450a0a80; }
  .diff-unchanged { border-left: 3px solid #333; color: #666; }
  .diff-sign    { font-weight: 700; width: 12px; flex-shrink: 0; }
  .diff-add    .diff-sign { color: #4ade80; }
  .diff-remove .diff-sign { color: #f87171; }
  .diff-title   { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .diff-artist  { color: #777; font-size: 0.72rem; white-space: nowrap; }
  .diff-bpm     { color: #555; font-size: 0.7rem; white-space: nowrap; }
  .diff-uptodate {
    font-size: 0.78rem; color: #4ade80;
    padding: 0.35rem 0.5rem; background: #052e1640; border-radius: 4px;
  }
  .btn-toggle-unchanged {
    background: none; border: none; color: #60a5fa;
    font-size: 0.75rem; cursor: pointer; padding: 0.2rem 0; text-align: left;
  }

  /* ── edit fields ─────────────────────────────────────────────────────── */
  .edit-section-label {
    font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.08em; color: #555;
  }
  .ep-edit-row { display: flex; gap: 0.5rem; align-items: center; }
  .ep-edit-row select,
  .ep-edit-row input {
    background: #0f0f0f; border: 1px solid #2a2a2a; border-radius: 5px;
    padding: 0.35rem 0.6rem; color: #e8e8e8; font-size: 0.8rem;
  }
  .ep-edit-row input:focus, .ep-edit-row select:focus { outline: 1px solid #2563eb; }
  .grow { flex: 1; }

  /* ── result feedback ─────────────────────────────────────────────────── */
  .result-error   { font-size: 0.8rem; color: #f87171; padding: 0.4rem 0; }
  .result-warning { font-size: 0.8rem; color: #fbbf24; padding: 0.2rem 0; }
  .result-blocked { font-size: 0.8rem; color: #a855f7; padding: 0.4rem 0; font-style: italic; }

  /* ── add step ────────────────────────────────────────────────────────── */
  .btn-add-step {
    width: 100%; padding: 0.75rem; background: #1a1a1a;
    border: 1px dashed #333; border-radius: 10px;
    color: #555; cursor: pointer; font-size: 0.875rem;
    transition: border-color 0.15s, color 0.15s;
  }
  .btn-add-step:hover { border-color: #2563eb; color: #60a5fa; }

  /* ── buttons ─────────────────────────────────────────────────────────── */
  .btn-primary   { background: #2563eb; color: #fff; border: none; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-ghost     { background: transparent; color: #aaa; border: 1px solid #333; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-ghost:hover:not(:disabled) { background: #2a2a2a; color: #fff; }
  .btn-preview   { background: #1e3a5f; color: #60a5fa; border: 1px solid #2563eb; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-apply     { background: #166534; color: #4ade80; border: 1px solid #16a34a; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-back      { background: transparent; border: 1px solid #333; color: #999; border-radius: 6px; padding: 0.4rem 0.75rem; cursor: pointer; font-size: 0.875rem; }
  .btn-delete    { background: transparent; color: #f87171; border: 1px solid #7f1d1d; border-radius: 4px; padding: 0.25rem 0.75rem; cursor: pointer; font-size: 0.75rem; }
  .btn-remove-step { background: transparent; color: #f87171; border: 1px solid #7f1d1d; border-radius: 4px; padding: 0.2rem 0.6rem; cursor: pointer; font-size: 0.75rem; }
  .btn-icon-sm   { background: transparent; border: none; color: #666; cursor: pointer; font-size: 0.875rem; padding: 0; }
  .btn-icon-sm:hover { color: #f87171; }
  .btn-link      { background: none; border: none; color: #60a5fa; cursor: pointer; font-size: 0.8rem; padding: 0.2rem 0; }
  button:disabled { opacity: 0.35; cursor: not-allowed; }

  /* ── misc ────────────────────────────────────────────────────────────── */
  .empty { color: #555; font-style: italic; }
  .error { color: #f87171; font-size: 0.875rem; margin-bottom: 1rem; }
</style>
