<script>
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
  <h1>INFO</h1>
  <div>
    <Tabs {items} />
  </div>
</about>

<style>
  about > h1 {
    letter-spacing: 10px;
  }

  about > div {
    padding-left: 2rem;
    padding-right: 2rem;
  }
</style>
