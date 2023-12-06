<script>
  import { selectedProxy } from "../store";

  export let proxies;

  function selectUser(prefix) {
    if (prefix == $selectedProxy) {
      selectedProxy.select(undefined);
      return;
    }

    selectedProxy.select(prefix);
  }
</script>

<div class="content">
  <div class="title">Proxies</div>
  <div class="list">
    {#each proxies as p}
      <div
        class="proxy"
        class:selected={$selectedProxy == p.id}
        on:click={() => selectUser(p.id)}
      >
        <span class="wrap" title={p.name}>{p.name}</span>
        <span class="host-name wrap" title={p.hostName}>{p.hostName}</span>
      </div>
    {/each}
  </div>
</div>

<style>
  .proxy {
    cursor: pointer;
    padding: 0.5rem;
    text-align: left;
    border: 1px solid transparent;
    display: flex;
    flex-direction: row;
    align-items: center;
    transition: border 0.2s ease-in-out;
  }

  .proxy:hover {
    border: 1px solid #646cff;
    transition: border 0.2s ease-in-out;
  }

  .proxy > span {
    flex: 1;
  }

  .selected {
    background-color: blueviolet;
  }

  .content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .title {
    text-transform: uppercase;
    font-size: 18px;
    letter-spacing: 3px;
  }

  .list {
    height: 307px;
    background-color: #3f3f46;
    overflow: auto;
    border: 1px solid transparent;
  }

  .host-name {
    font-size: 12px;
    text-align: end;
  }

  .wrap {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
