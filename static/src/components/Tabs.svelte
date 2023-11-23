<script>
  export let items = [];
  export let activeTabValue = 1;

  const handleClick = (tabValue) => () => (activeTabValue = tabValue);
</script>

<ul>
  {#each items as item}
    <li class={activeTabValue === item.value ? "active" : ""}>
      <span on:click={handleClick(item.value)}>{item.label}</span>
    </li>
  {/each}
</ul>
<div class="items">
  {#each items as item}
    {#if activeTabValue == item.value}
      <div class="box">
        <svelte:component this={item.component} />
      </div>
    {/if}
  {/each}
</div>

<style>
  .box {
    margin-bottom: 10px;
    padding: 30px;
    border: 1px solid blueviolet;
    border-radius: 0 0 0.5rem 0.5rem;
    border-top: 0;
  }

  ul {
    display: flex;
    flex-wrap: wrap;
    padding-left: 0;
    margin-bottom: 0;
    list-style: none;
    margin-left: 2rem;
    margin-right: 2rem;
    border-bottom: 1px solid blueviolet;
  }

  .items {
    padding-left: 2rem;
    padding-right: 2rem;
  }

  li {
    margin-bottom: -1px;
  }

  span {
    border: 1px solid transparent;
    border-top-left-radius: 0.25rem;
    border-top-right-radius: 0.25rem;
    display: block;
    padding: 0.5rem 1rem;
    cursor: pointer;
  }

  span:hover {
    border: 1px solid #646cff;
  }

  li.active > span {
    color: #495057;
    color: white;
    background-color: blueviolet;
  }
</style>
