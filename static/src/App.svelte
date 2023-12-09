<script>
  import { onMount } from "svelte";
  import { SvelteToast, toast } from "@zerodevx/svelte-toast";
  import Select from "svelte-select";

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
  let clusters = [];
  let selectedCluster = "";
  let loaded = false;
  let qmcLink = "";
  let hubLink = "";
  let generateButtonEnabled = false;
  let generateButtonTitle = "";
  let attributesString = "";
  let attributesPlaceholderValues = [
    "Additional attributes (valid JSON) to be associated with the ticket.\n\n",
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

  $: if ($selectedProxy) {
    virtualProxies = proxies.filter((p) => p.id == $selectedProxy)[0]
      .virtualProxies;

    if (virtualProxies.length == 1) selectedVP.select(virtualProxies[0].prefix);
  } else {
    virtualProxies = [];
  }

  $: if ($selectedUser && $selectedVP != undefined) {
    generateButtonEnabled = true;
  } else {
    generateButtonEnabled = false;
  }

  $: if (!generateButtonEnabled) {
    generateButtonTitle = "Select User, Proxy and Virtual Proxy values";
  } else {
    generateButtonTitle = "Generate";
  }

  async function getUsers(cluster) {
    await fetch(`https://localhost:8081/api/users/${cluster}`, {
      method: "GET",
    })
      .then((r) => r.json())
      .then((r) => {
        users = r;

        // pre-select the user if there is only one available
        if (users.length == 1) selectedUser.select(users[0].userId);
      });
  }

  async function getProxies(cluster) {
    await fetch(`https://localhost:8081/api/proxies/${cluster}`, {
      method: "GET",
    })
      .then((r) => r.json())
      .then((r) => {
        proxies = r;

        // pre-select the proxy if there is only one available
        if (proxies.length == 1) selectedProxy.select(proxies[0].id);
      });
  }

  async function getClusters() {
    await fetch("https://localhost:8081/api/clusters", {
      method: "GET",
    })
      .then((r) => r.json())
      .then((r) => {
        clusters = r;
        // clusters = [r[0]];

        // if there is only one cluster returned then directly
        // get the proxies associated with it and dont wait
        // for the user to select the only available option
        if (clusters.length == 1)
          return getProxies(clusters[0]).then(() => {
            selectedCluster = clusters[0];
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
          cluster: selectedCluster["value"],
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
      new Promise((resolve) => setTimeout(resolve, 500)),
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

  async function selectCluster() {
    loaded = false;
    selectedProxy.reset();
    selectedVP.reset();
    selectedUser.reset();

    Promise.all([
      getProxies(selectedCluster["value"]),
      getUsers(selectedCluster["value"]),
      new Promise((resolve) => setTimeout(resolve, 500)),
    ])
      .then(() => (loaded = true))
      .catch((e) => {
        toast.push(e.message, {
          theme: { ...toastErrorTheme },
        });
        selectedCluster = "";
        loaded = true;
      });
  }

  onMount(async () => {
    selectedProxy.reset();
    selectedVP.reset();
    selectedUser.reset();

    await Promise.all([
      getClusters(),
      new Promise((resolve) => setTimeout(resolve, 500)),
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
    {#if !selectedCluster["value"]}
      <div class="overlay"></div>
      <div class="overlay-text">
        <div>Please select cluster</div>
      </div>
    {/if}
    <content>
      <clusters>
        <Select
          id="clusters"
          items={clusters}
          bind:value={selectedCluster}
          showChevron
          clearable={false}
          searchable={false}
          placeholder="QLIK CLUSTERS"
          --item-padding="0px"
          --z-index="101"
          on:change={selectCluster}
        >
          <div slot="item" let:item class="list-item">
            &nbsp;&nbsp; {item.label}
          </div>
        </Select>
      </clusters>
      <selections>
        <users><Users {users} /></users>
        <proxies><Proxies {proxies} /></proxies>
        <virtual-proxies><VirtualProxies {virtualProxies} /></virtual-proxies>
        <attributes>
          <span class="title">Attributes</span>
          <textarea
            id="attributes"
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
      </selections>
    </content>
  {/if}
</main>

<style>
  .overlay {
    margin-top: 50px;
    position: absolute;
    width: 100%;
    height: 100%;
    background-color: blueviolet;
    opacity: 65%;
    z-index: 100;
    right: 0px;
    -webkit-mask-image: linear-gradient(
      0deg,
      transparent,
      150px,
      blueviolet 300px
    );
    mask-image: linear-gradient(0deg, transparent, 150px, blueviolet 300px);
  }
  .overlay-text {
    position: absolute;
    width: 100%;
    height: 65%;
    z-index: 101;
    right: 0px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    font-size: 3em;
    text-transform: uppercase;
    color: white;
    letter-spacing: 5px;
    font-weight: 500;
  }

  clusters {
    z-index: 103;
  }

  .list-item {
    /* background-color: #242424; */
    background-color: darkgray;
    color: black;
    text-align: left;
    cursor: pointer;
    z-index: 99;
    transition: background-color 0.1s ease-in-out;
  }

  .list-item:hover {
    background-color: blueviolet;
    color: white;
    transition: background-color 0.2s ease-in-out;
  }

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
    position: relative;
    display: flex;
    flex-direction: column;
    grid-template-rows: 50px auto;
    grid-template-columns: auto;
    padding-left: 1rem;
    padding-right: 1rem;
    padding-top: 2rem;
  }

  selections {
    /* grid-column: 1; */
    display: grid;
    grid-template-columns: 25% 25% 25% 25%;
    grid-template-rows: 350px 50px 50px;
    padding-left: 1rem;
    padding-right: 4rem;
    padding-top: 2rem;
    gap: 1rem;
    flex: 1;
  }

  .loader {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  clusters {
    grid-column: 1 / span 2;
    grid-row: 1;
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
    height: 300px;
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
    align-items: center;
    justify-content: center;
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
    padding-left: 5px;
    padding-top: 5px;
  }

  textarea:focus {
    border: 1px solid #646cff;
    outline: none;
  }

  .title {
    letter-spacing: 3px;
  }
</style>
