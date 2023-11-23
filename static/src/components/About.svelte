<script>
  import { showHideAbout } from "../store";

  import Tab1 from "./About/Tab1.svelte";
  import Tab2 from "./About/Tab2.svelte";
  import Tab3 from "./About/Tab3.svelte";
  import Tabs from "./Tabs.svelte";

  let items = [
    { label: "General", value: 1, component: Tab1 },
    { label: "Clear cache", value: 2, component: Tab2 },
    { label: "Attributes", value: 3, component: Tab3 },
  ];
</script>

<about>
  <title>
    <h1>INFO</h1>
    <div
      class="close-button"
      title="Back"
      on:click={() => showHideAbout.set(false)}
    >
      X
    </div>
  </title>
  <tabs>
    <Tabs {items} />
  </tabs>
</about>

<style>
  h1 {
    letter-spacing: 10px;
    flex: 1;
  }

  title {
    display: flex;
    flex-direction: row;
  }

  .close-button {
    width: 50px;
    align-self: center;
    font-size: 20px;
    font-weight: bold;
    right: 2rem;
    position: relative;
    cursor: pointer;
    text-align: end;
  }
</style>
