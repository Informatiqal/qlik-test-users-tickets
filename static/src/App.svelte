<script>
  import { onMount } from "svelte";
  import { selectedUser, selectedVP } from "./store";

  import CopySVG from "./assets/copy.svelte";
  import LoaderSVG from "./assets/loader.svelte";

  import Users from "./components/Users.svelte";
  import VirtualProxies from "./components/VirtualProxies.svelte";

  let users = [];
  let virtualProxies = [];
  let loaded = false;
  let qmcLink = "";
  let hubLink = "";
  let generateButtonEnabled = false;

  $: if ($selectedUser && $selectedVP) {
    generateButtonEnabled = true;
  } else {
    generateButtonEnabled = false;
  }

  async function getUsers() {
    const r = await fetch("https://localhost:8081/api/users", {
      method: "GET",
    });
    users = await r.json();
  }

  async function getVirtualProxies() {
    const r = await fetch("https://localhost:8081/api/virtualproxies", {
      method: "GET",
    });
    virtualProxies = await r.json();
  }

  async function generateTicket() {
    loaded = false;
    qmcLink = "";
    hubLink = "";

    await Promise.all([
      fetch("https://localhost:8081/api/ticket", {
        method: "POST",
        body: JSON.stringify({
          userId: $selectedUser,
          virtualProxyPrefix: $selectedVP,
        }),
      })
        .then((a) => a.json())
        .then((ticketResponse) => {
          qmcLink = ticketResponse.links.qmc;
          hubLink = ticketResponse.links.hub;
        }),
      new Promise((resolve) => setTimeout(resolve, 1000)),
    ]).then(() => {
      loaded = true;
    });
  }

  async function copyToClipBoard(text) {
    try {
      await navigator.clipboard.writeText(text);
      console.log("Content copied to clipboard");
    } catch (err) {
      console.error("Failed to copy: ", err);
    }
  }

  onMount(async () => {
    await Promise.all([
      getUsers(),
      getVirtualProxies(),
      new Promise((resolve) => setTimeout(resolve, 1000)),
    ]).then(() => (loaded = true));
  });
</script>

<main>
  <header>Test Users Ticket Generator</header>
  {#if !loaded}
    <div class="loader">
      <LoaderSVG />
    </div>
  {:else}
    <content>
      <users><Users {users} /></users>
      <proxies><VirtualProxies {virtualProxies} /></proxies>
      <generate>
        <button
          class:button-disabled={!generateButtonEnabled}
          on:click={() => generateTicket()}
          disabled={!generateButtonEnabled}>GENERATE TICKET</button
        >
      </generate>
      <links>
        {#if hubLink && qmcLink}
          <div class="links-content">
            <span>HUB</span>
            <span>QMC</span>
          </div>
          <div class="links-content left">
            <div>
              <span>{hubLink} </span>
              <span
                class="copy"
                title="Copy to clipboard"
                on:click={() => copyToClipBoard(hubLink)}><CopySVG /></span
              >
            </div>
            <div>
              <span>{qmcLink} </span>
              <span
                class="copy"
                title="Copy to clipboard"
                on:click={() => copyToClipBoard(qmcLink)}><CopySVG /></span
              >
            </div>
          </div>
        {/if}
      </links>
    </content>
  {/if}
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  header {
    background-color: blueviolet;
    height: 3rem;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    font-weight: bold;
    text-transform: uppercase;
    letter-spacing: 5px;
  }

  content {
    display: grid;
    grid-template-columns: auto auto;
    grid-template-rows: auto auto auto;
    padding-left: 4rem;
    padding-right: 4rem;
    padding-top: 2rem;
    gap: 1rem;
  }

  .loader {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  users {
    grid-column: 1;
    grid-row: 1;
  }

  proxies {
    grid-column: 2;
    grid-row: 1;
  }

  generate {
    grid-column: 1 / span 2;
    grid-row: 2;
  }

  links {
    grid-column: 1 / span 2;
    grid-row: 3;
  }

  button {
    /* border-radius: 8px; */
    border: 1px solid transparent;
    padding: 0.6em 1.2em;
    font-size: 1em;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: border-color 0.25s;
    width: 100%;
    background-color: darkgreen;
  }

  button:hover {
    border-color: #646cff;
  }

  button:focus,
  button:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }

  .button-disabled {
    background-color: gray;
    color: lightgray;
    cursor: not-allowed;
  }

  links {
    display: grid;
    grid-template-rows: auto auto;
  }

  .links-content {
    display: flex;
    flex-direction: row;
    gap: 1rem;
  }

  .links-content > span {
    flex: 1;
  }

  .links-content > div {
    flex: 1;
    display: flex;
    gap: 1rem;
  }

  .left {
    text-align: left;
  }

  .copy {
    cursor: pointer;
  }
</style>
