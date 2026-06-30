<script>
  import Sources  from './screens/Sources.svelte';
  import Workflow from './screens/Workflow.svelte';

  let activeTab = 'sources'; // sources | workflows
</script>

<main>
  <header>
    <h1>DJ Library Tools</h1>
    <nav>
      <button class:active={activeTab === 'sources'}   on:click={() => activeTab = 'sources'}>Sources</button>
      <button class:active={activeTab === 'workflows'} on:click={() => activeTab = 'workflows'}>Workflows</button>
    </nav>
  </header>

  <section>
    {#if activeTab === 'sources'}
      <Sources />
    {:else if activeTab === 'workflows'}
      <Workflow />
    {/if}
  </section>
</main>

<style>
  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
  :global(body) { font-family: system-ui, sans-serif; background: #0f0f0f; color: #e8e8e8; }

  main { display: flex; flex-direction: column; height: 100vh; }

  header {
    display: flex; align-items: center; gap: 2rem;
    padding: 0.75rem 1.5rem;
    background: #1a1a1a; border-bottom: 1px solid #2a2a2a;
  }

  h1 { font-size: 1rem; font-weight: 600; color: #fff; white-space: nowrap; }

  nav { display: flex; gap: 0.25rem; }

  nav button {
    padding: 0.4rem 1rem; border: none; border-radius: 6px;
    background: transparent; color: #999; cursor: pointer; font-size: 0.875rem;
    transition: background 0.15s, color 0.15s;
  }
  nav button:hover { background: #2a2a2a; color: #fff; }
  nav button.active { background: #2563eb; color: #fff; }

  section { flex: 1; overflow: auto; padding: 1.5rem; }
</style>
