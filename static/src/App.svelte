<script>
  import { onMount } from "svelte";
  import { SvelteToast, toast } from "@zerodevx/svelte-toast";
  import { selectedUser, selectedVP } from "./store";

  import CopySVG from "./assets/copy.svelte";
  import LoaderSVG from "./assets/loader.svelte";
  import GitHubSVG from "./assets/github.svelte";
  import InfoSVG from "./assets/info.svelte";

  import Users from "./components/Users.svelte";
  import VirtualProxies from "./components/VirtualProxies.svelte";
  import About from "./components/About.svelte";

  let users = [];
  let virtualProxies = [];
  let loaded = false;
  let qmcLink = "";
  let hubLink = "";
  let generateButtonEnabled = false;
  let isAboutSection = false;
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

  $: if ($selectedUser && $selectedVP) {
    generateButtonEnabled = true;
  } else {
    generateButtonEnabled = false;
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

  async function getVirtualProxies() {
    await fetch("https://localhost:8081/api/virtualproxies", {
      method: "GET",
    })
      .then((r) => r.json())
      .then((r) => {
        virtualProxies = r;
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
      getVirtualProxies(),
      new Promise((resolve) => setTimeout(resolve, 1000)),
    ]).then(() => (loaded = true));
  });
</script>

<main>
  <SvelteToast {options} />
  <header>
    <span>Test Users Ticket Generator</span>
    <div class="logo">
      <span title="About" on:click={() => (isAboutSection = !isAboutSection)}
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
  {#if isAboutSection}
    <About />
  {:else if !loaded}
    <div class="loader">
      <LoaderSVG />
    </div>
  {:else}
    <content>
      <users><Users {users} /></users>
      <proxies><VirtualProxies {virtualProxies} /></proxies>
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
    grid-template-columns: auto auto auto;
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

  attributes {
    grid-column: 3;
    grid-row: 1;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  generate {
    grid-column: 1 / span 3;
    grid-row: 2;
  }

  links {
    grid-column: 1 / span 3;
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
