<script>
  import { onMount } from "svelte";
  import { SvelteToast, toast } from "@zerodevx/svelte-toast";
  import {
    selectedProxy,
    selectedUser,
    selectedVP,
    showHideAbout,
  } from "./store";

  import CopySVG from "./assets/copy.svelte";
  import LoaderSVG from "./assets/loader.svelte";
  import GitHubSVG from "./assets/github.svelte";
  import InfoSVG from "./assets/info.svelte";

  import Users from "./components/Users.svelte";
  import VirtualProxies from "./components/VirtualProxies.svelte";
  import About from "./components/About.svelte";
  import Proxies from "./components/Proxies.svelte";

  let users = [];
  let virtualProxies = [];
  let proxies = [];
  let loaded = false;
  let qmcLink = "";
  let hubLink = "";
  let generateButtonEnabled = false;
  let generateButtonTitle = "";
  let attributesString = "";
  let attributesPlaceholderValues = [
    "Additional attributes to be associated with the ticket.\n\n",
    `[\n  { "group": "some group" },\n`,
    `  { "group": "another group" },\n`,
    `  { "otherProperty": "some value" },\n`,
    `  ...\n]`,
  ];
  let attributesPlaceholder = attributesPlaceholderValues.join("");
  const options = {};
  const toastErrorTheme = {
    "--toastColor": "white",
    "--toastBackground": "#ff6e64",
    "--toastBarBackground": "darkred",
  };

  $: if ($selectedUser && $selectedVP != undefined) {
    generateButtonEnabled = true;
  } else {
    generateButtonEnabled = false;
  }

  $: if ($selectedProxy) {
    virtualProxies = proxies.filter((p) => p.id == $selectedProxy)[0]
      .virtualProxies;
  } else {
    virtualProxies = [];
  }

  $: if (!generateButtonEnabled) {
    generateButtonTitle = "Select User, Proxy and Virtual Proxy values";
  } else {
    generateButtonTitle = "Generate";
  }

  async function getUsers() {
    await fetch("https://localhost:8081/api/users", {
      method: "GET",
    })
      .then((r) => r.json())
      .then((r) => {
        users = r;
      })
      .catch((e) => {
        toast.push(e.message, {
          theme: { ...toastErrorTheme },
        });
      });
  }

  async function getProxies() {
    await fetch("https://localhost:8081/api/proxies", {
      method: "GET",
    })
      .then((r) => r.json())
      .then((r) => {
        proxies = r;

        // pre-select the proxy if there is only one available
        // if (proxies.length == 1) selectedProxy.select(proxies[0].id);
      })
      .catch((e) => {
        toast.push(e.message, {
          theme: { ...toastErrorTheme },
        });
      });
  }

  async function generateTicket() {
    qmcLink = "";
    hubLink = "";
    let attributes = [];

    try {
      attributes = !attributesString ? [] : JSON.parse(attributesString);
    } catch (e) {
      toast.push("Unable to parse the attributes", {
        theme: { ...toastErrorTheme },
      });
      return;
    }

    loaded = false;

    await Promise.all([
      fetch("https://localhost:8081/api/ticket", {
        method: "POST",
        body: JSON.stringify({
          userId: $selectedUser,
          virtualProxyPrefix: $selectedVP,
          proxyId: $selectedProxy,
          attributes,
        }),
      })
        .then((a) => a.json())
        .then((ticketResponse) => {
          qmcLink = ticketResponse.links.qmc;
          hubLink = ticketResponse.links.hub;
        })
        .catch((e) => {
          toast.push(e.message, {
            theme: { ...toastErrorTheme },
          });
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
      getProxies(),
      new Promise((resolve) => setTimeout(resolve, 1000)),
    ]).then(() => (loaded = true));
  });
</script>

<main>
  <SvelteToast {options} />
  <header>
    <span>Test Users Ticket Generator</span>
    <div class="logo">
      <span title="About" on:click={() => showHideAbout.set(!$showHideAbout)}
        ><InfoSVG /></span
      >
      <span
        title="Source code"
        on:click={() =>
          window.open(
            "https://github.com/informatiqal/qlik-test-users-tickets",
            "_blank"
          )}><GitHubSVG /></span
      >
    </div>
  </header>
  {#if $showHideAbout}
    <About />
  {:else if !loaded}
    <div class="loader">
      <LoaderSVG />
    </div>
  {:else}
    <content>
      <users><Users {users} /></users>
      <proxies><Proxies {proxies} /></proxies>
      <virtual-proxies><VirtualProxies {virtualProxies} /></virtual-proxies>
      <attributes>
        <span class="title">Attributes</span>
        <textarea
          bind:value={attributesString}
          placeholder={attributesPlaceholder}
        />
      </attributes>
      <generate>
        <button
          class:button-disabled={!generateButtonEnabled}
          on:click={() => generateTicket()}
          title={generateButtonTitle}
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
    height: 50px;
    /* color: #ff6e64; */
  }

  .logo {
    position: absolute;
    right: 10px;
    cursor: pointer;
    display: flex;
    flex-direction: row;
    gap: 10px;
    justify-content: center;
    align-items: center;
  }

  .logo > span {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  content {
    display: grid;
    grid-template-columns: 25% 25% 25% 25%;
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

  virtual-proxies {
    grid-column: 3;
    grid-row: 1;
  }

  attributes {
    grid-column: 4;
    grid-row: 1;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  generate {
    grid-column: 1 / span 4;
    grid-row: 2;
  }

  links {
    grid-column: 1 / span 4;
    grid-row: 3;
  }

  attributes > span {
    text-transform: uppercase;
    font-size: 18px;
  }

  attributes > textarea {
    height: 100%;
    resize: none;
  }

  button {
    border: 1px solid transparent;
    padding: 0.6em 1.2em;
    font-size: 1em;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: border-color 0.25s;
    width: 100%;
    background-color: darkgreen;
    border-bottom-left-radius: 8px;
    border-bottom-right-radius: 8px;
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

  textarea {
    border: 1px solid gray;
    border-top-right-radius: 8px;
    padding: 5px;
  }

  textarea:focus {
    border: 1px solid #646cff;
    outline: none;
  }

  .title {
    letter-spacing: 3px;
  }
</style>
