<script>
  import { selectedUser } from "../store";
  export let users;

  function selectUser(userId) {
    if (userId == $selectedUser) {
      selectedUser.select(undefined);
      return;
    }

    selectedUser.select(userId);
  }
</script>

<div class="content">
  <div class="title">Test users</div>
  <div class="list">
    {#if users.length == 0}
      <div class="no-users">
        <div>No test users were found</div>
      </div>
    {:else}
      {#each users as user}
        <div
          class="user"
          class:selected={$selectedUser == user.userId}
          on:click={() => selectUser(user.userId)}
        >
          {user.userDirectory}\{user.userId}
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .user {
    cursor: pointer;
    padding: 0.5rem;
    text-align: left;
    border: 1px solid transparent;
    transition: border 0.2s ease-in-out;
  }

  .user:hover {
    border: 1px solid #646cff;
    border-top-left-radius: 8px;
    transition: border 0.2s ease-in-out;
  }

  .selected {
    background-color: blueviolet;
  }

  .content {
    position: relative;
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
    border-top-left-radius: 8px;
    overflow: auto;
    border: 1px solid transparent;
  }

  .no-users {
    display: flex;
    max-height: 300px;
    min-height: 300px;
    align-items: center;
    justify-content: center;
  }
</style>
