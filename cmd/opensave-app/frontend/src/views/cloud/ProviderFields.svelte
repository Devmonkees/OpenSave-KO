<script>
  import { t } from '../../lib/i18n.js';
  // Where a provider without a sign-in sends its backups.
  import { native } from '../../lib/api.js';

  export let config;

  async function pickFolder() {
    const dir = await native.selectDirectory('Select backup destination folder');
    if (dir) config.url = dir;
  }
</script>

{#if config.provider === 'local'}
  <div class="field">
    <label for="cb-folder">{$t('ui.destinationFolderEGANas')}</label>
    <div class="path-row">
      <input id="cb-folder" bind:value={config.url} placeholder="D:\Backups\OpenSave" />
      <button class="btn" on:click={pickFolder}>{$t('settings.storage.browse')}</button>
    </div>
  </div>
{:else if config.provider === 'webdav'}
  <div class="field">
    <label for="cb-url">WebDAV URL</label>
    <input id="cb-url" bind:value={config.url} placeholder="https://nas.local/dav/opensave/" />
  </div>
  <div class="two">
    <div class="field">
      <label for="cb-user">{$t('ui.username')}</label>
      <input id="cb-user" bind:value={config.username} />
    </div>
    <div class="field">
      <label for="cb-pass">{$t('ui.password')}</label>
      <input id="cb-pass" type="password" bind:value={config.password} />
    </div>
  </div>
{:else if config.provider === 'webhook'}
  <div class="field">
    <label for="cb-hook">{$t('ui.webhookUrlReceivesMultipartPost')}</label>
    <input id="cb-hook" bind:value={config.url} />
  </div>
  <div class="field">
    <label for="cb-headers">{$t('ui.customHeadersJson')}</label>
    <input id="cb-headers" bind:value={config.headers} placeholder={'{"Authorization": "Bearer …"}'} />
  </div>
{/if}

<style>
  .two {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
  }
</style>
