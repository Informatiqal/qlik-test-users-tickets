<script>
  import { selectedVP } from "../store";

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
  <div class="title">Available virtual proxies</div>
  <div class="list">
    {#each virtualProxies as vp}
      <div
        class="vp"
        class:selected={$selectedVP == vp.prefix}
        on:click={() => selectUser(vp.prefix)}
      >
        {vp.description} ({vp.prefix ? vp.prefix : "/"})
      </div>
    {/each}
  </div>
</div>

<style>
  .vp {
    background-color: #3f3f46;
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
  }

  .list {
    min-height: 200px;
    background-color: #3f3f46;
  }
</style>
