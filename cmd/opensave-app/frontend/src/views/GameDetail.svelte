<script>
  import { t } from '../lib/i18n.js';
  import { slidingIndicator } from '../lib/motion.js';
  import { games } from '../lib/stores.js';
  import { createRunner } from '../lib/runner.js';
  import GameHeader from './game/GameHeader.svelte';
  import SnapshotsTab from './game/SnapshotsTab.svelte';
  import BranchesTab from './game/BranchesTab.svelte';
  import CloudTab from './game/CloudTab.svelte';
  import ConfigTab from './game/ConfigTab.svelte';
  import ManageTab from './game/ManageTab.svelte';
  import History from 'lucide-svelte/icons/history';
  import GitBranch from 'lucide-svelte/icons/git-branch';
  import Cloud from 'lucide-svelte/icons/cloud';
  import SlidersHorizontal from 'lucide-svelte/icons/sliders-horizontal';
  import Wrench from 'lucide-svelte/icons/wrench';

  export let params = {};

  $: game = $games[params.gameId];

  // One action at a time across the whole page — see lib/runner.js.
  const runner = createRunner();

  const tabs = [
    { id: 'snapshots', icon: History, component: SnapshotsTab },
    { id: 'branches', icon: GitBranch, component: BranchesTab },
    { id: 'cloud', icon: Cloud, component: CloudTab },
    { id: 'config', icon: SlidersHorizontal, component: ConfigTab },
    { id: 'danger', icon: Wrench, component: ManageTab }
  ];
  let tab = 'snapshots';

  // A tab is built the first time it is opened and then kept, hidden, while
  // another is shown: half-typed configuration or an open file list should
  // still be there after a look at another tab. Each tab loads what it needs
  // when it is built, so one never opened costs nothing.
  //
  // This page is reused, not rebuilt, when moving from one game to another,
  // so all of it starts over when the game changes — otherwise one game's
  // cloud list or unsaved edits would show up on the next.
  let opened = new Set([tab]);
  // Opened on a snapshot (from Ctrl+K): the list it is in.
  $: if (params.snapshot) tab = 'snapshots';

  let shownFor = params.gameId;
  $: if (params.gameId !== shownFor) {
    shownFor = params.gameId;
    tab = 'snapshots';
    opened = new Set([tab]);
  }
  $: if (!opened.has(tab)) opened = new Set(opened).add(tab);
</script>

{#if !game}
  <div class="empty"><h3>{$t('ui.gameNotFound')}</h3></div>
{:else}
  {#key params.gameId}
    <GameHeader {game} {runner} renameRequest={params.rename} />

    <div class="pill-tabs tabs" use:slidingIndicator={{ inset: 10 }}>
      {#each tabs as item (item.id)}
        <button class:active={tab === item.id} on:click={() => (tab = item.id)}><svelte:component this={item.icon} size={15} />{$t('gameTab.' + item.id)}</button>
      {/each}
    </div>

    {#each tabs as item (item.id)}
      {#if opened.has(item.id)}
        <div hidden={tab !== item.id}>
          <svelte:component this={item.component} {game} {runner} {...item.id === 'snapshots' ? { focus: params.snapshot } : {}} />
        </div>
      {/if}
    {/each}
  {/key}
{/if}

<style>
  .tabs {
    margin-bottom: 18px;
  }
</style>
