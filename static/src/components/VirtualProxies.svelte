<script>
  import { selectedVP, selectedProxy } from "../store";

  export let virtualProxies;

  function selectUser(prefix) {
    if (prefix == $selectedVP) {
      selectedVP.select(undefined);
      return;
    }

    selectedVP.select(prefix);
  }
</script>

<div class="content">
  <div class="title">Virtual proxies</div>
  <div class="list">
    {#if $selectedProxy}
      {#each virtualProxies as vp}
        <div
          class="vp"
          class:selected={$selectedVP == vp.prefix}
          on:click={() => selectUser(vp.prefix)}
        >
          {vp.description} ({vp.prefix ? vp.prefix : "/"})
        </div>
      {/each}
    {:else}
      <div class="no-proxy">
        <div>Please select proxy service first</div>
      </div>
    {/if}
  </div>
</div>

<style>
  .vp {
    cursor: pointer;
    padding: 0.5rem;
    text-align: left;
    border: 1px solid transparent;
  }

  .vp:hover {
    border: 1px solid #646cff;
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
    max-height: 300px;
    min-height: 300px;
    background-color: #3f3f46;
    overflow: auto;
  }

  .no-proxy {
    display: flex;
    max-height: 300px;
    min-height: 300px;
    align-items: center;
    justify-content: center;
  }
</style>
