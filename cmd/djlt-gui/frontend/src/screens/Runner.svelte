<script>
  import { PreviewWorkflow, RunWorkflow, GetWorkflow } from '../../wailsjs/go/gui/app.js';

  export let workflowId = null;
  export let onBack = () => {};

  let workflow = null;
  let previewResult = null;
  let runResult = null;
  let loading = false;
  let error = '';
  let phase = 'idle'; // idle | previewing | preview-done | running | done

  async function load() {
    if (!workflowId) return;
    try { workflow = await GetWorkflow(workflowId); }
    catch (e) { error = String(e); }
  }

  $: if (workflowId) { load(); previewResult = null; runResult = null; phase = 'idle'; error = ''; }

  async function preview() {
    error = ''; loading = true; phase = 'previewing';
    try {
      previewResult = await PreviewWorkflow(workflowId);
      phase = 'preview-done';
    } catch (e) { error = String(e); phase = 'idle'; }
    loading = false;
  }

  async function run() {
    if (!confirm('Apply all changes? This cannot be undone.')) return;
    error = ''; loading = true; phase = 'running';
    try {
      runResult = await RunWorkflow(workflowId);
      phase = 'done';
    } catch (e) { error = String(e); phase = 'preview-done'; }
    loading = false;
  }

  function statusColour(s) {
    return { pending: '#555', running: '#f59e0b', success: '#22c55e', failed: '#ef4444', blocked: '#a855f7' }[s] ?? '#555';
  }

  function statusIcon(s) {
    return { pending: '○', running: '◌', success: '✓', failed: '✗', blocked: '⊘' }[s] ?? '?';
  }

  const activeResult = () => runResult ?? previewResult;
</script>

<div class="screen">
  <div class="toolbar">
    <button class="btn-back" on:click={onBack}>← Back</button>
    {#if workflow}
      <h2>{workflow.name}</h2>
      <span class="step-count">{workflow.steps?.length ?? 0} steps</span>
    {/if}
    <div style="flex:1" />
    {#if phase === 'idle' || phase === 'preview-done'}
      <button class="btn-preview" on:click={preview} disabled={loading}>
        {loading && phase === 'previewing' ? 'Running preview…' : '🔍 Preview'}
      </button>
    {/if}
    {#if phase === 'preview-done'}
      <button class="btn-run" on:click={run} disabled={loading}>▶ Run (Apply)</button>
    {/if}
  </div>

  {#if error}<p class="error">{error}</p>{/if}

  {#if loading && (phase === 'previewing' || phase === 'running')}
    <p class="loading">{phase === 'previewing' ? 'Running preview…' : 'Applying changes…'}</p>
  {/if}

  {#if activeResult()}
    {@const result = activeResult()}
    <div class="mode-badge" class:preview={!runResult} class:run={!!runResult}>
      {runResult ? 'RUN RESULTS' : 'PREVIEW RESULTS'}
    </div>

    <div class="steps-grid">
      {#each result.steps as sr}
        {@const blocked = sr.status === 'blocked'}
        <div class="step-card" style="--status-color: {statusColour(sr.status)}">
          <div class="step-header">
            <span class="status-icon">{statusIcon(sr.status)}</span>
            <span class="step-id">{sr.step_id?.slice(0, 8)}…</span>
            <span class="status-label" style="color: {statusColour(sr.status)}">{sr.status}</span>
          </div>

          {#if blocked}
            <p class="blocked-reason">Blocked — a dependency failed.</p>
          {:else if sr.error}
            <div class="message-group error-group">
              <span class="msg-label">Error</span>
              <p class="msg-line error-msg">{sr.error}</p>
            </div>
          {/if}

          {#if sr.previews?.length}
            <div class="message-group">
              <span class="msg-label">Preview ({sr.previews.length})</span>
              {#each sr.previews as m}<p class="msg-line">{m}</p>{/each}
            </div>
          {/if}

          {#if sr.successes?.length}
            <div class="message-group success-group">
              <span class="msg-label">Applied ({sr.successes.length})</span>
              {#each sr.successes as m}<p class="msg-line">{m}</p>{/each}
            </div>
          {/if}

          {#if sr.warnings?.length}
            <div class="message-group warning-group">
              <span class="msg-label">Warnings ({sr.warnings.length})</span>
              {#each sr.warnings as m}<p class="msg-line">{m}</p>{/each}
            </div>
          {/if}

          {#if !sr.previews?.length && !sr.successes?.length && !sr.warnings?.length && !sr.error && !blocked}
            <p class="msg-line dim">No output.</p>
          {/if}
        </div>
      {/each}
    </div>
  {:else if !loading && workflow}
    <p class="hint">Press <strong>Preview</strong> to see what this workflow would change.</p>
  {/if}
</div>

<style>
  .screen { max-width: 1100px; }
  .toolbar { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; }
  h2 { font-size: 1.1rem; font-weight: 600; }
  .step-count { color: #666; font-size: 0.8rem; }
  .btn-back { background: transparent; border: 1px solid #333; color: #999; border-radius: 6px; padding: 0.4rem 0.75rem; cursor: pointer; font-size: 0.875rem; }
  .btn-preview { background: #1e3a5f; color: #60a5fa; border: 1px solid #2563eb; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  .btn-run { background: #16a34a; color: #fff; border: none; border-radius: 6px; padding: 0.4rem 1rem; cursor: pointer; font-size: 0.875rem; }
  button:disabled { opacity: 0.4; cursor: default; }
  .mode-badge { display: inline-block; padding: 0.2rem 0.75rem; border-radius: 4px; font-size: 0.7rem; font-weight: 700; letter-spacing: 0.08em; margin-bottom: 1rem; }
  .mode-badge.preview { background: #1e3a5f; color: #60a5fa; }
  .mode-badge.run { background: #14532d; color: #4ade80; }
  .steps-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
  .step-card { background: #1a1a1a; border: 1px solid var(--status-color, #2a2a2a); border-radius: 8px; overflow: hidden; }
  .step-header { display: flex; align-items: center; gap: 0.5rem; padding: 0.6rem 0.8rem; background: #111; border-bottom: 1px solid #222; }
  .status-icon { font-size: 0.9rem; color: var(--status-color); }
  .step-id { font-family: monospace; font-size: 0.7rem; color: #555; flex: 1; }
  .status-label { font-size: 0.75rem; font-weight: 600; text-transform: uppercase; }
  .message-group { padding: 0.6rem 0.8rem; border-bottom: 1px solid #1e1e1e; }
  .message-group:last-child { border-bottom: none; }
  .msg-label { display: block; font-size: 0.65rem; text-transform: uppercase; color: #555; letter-spacing: 0.05em; margin-bottom: 0.3rem; }
  .msg-line { font-size: 0.8rem; color: #ccc; padding: 0.1rem 0; white-space: pre-wrap; word-break: break-word; }
  .success-group .msg-line { color: #4ade80; }
  .warning-group .msg-line { color: #fbbf24; }
  .error-group .msg-line { color: #f87171; }
  .blocked-reason { font-size: 0.8rem; color: #a855f7; padding: 0.6rem 0.8rem; font-style: italic; }
  .dim { color: #555; }
  .error { color: #f87171; font-size: 0.875rem; margin-bottom: 1rem; }
  .loading { color: #f59e0b; font-size: 0.875rem; margin-bottom: 1rem; }
  .hint { color: #555; font-style: italic; margin-top: 2rem; }
</style>
