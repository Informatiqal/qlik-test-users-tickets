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
  <div class="title">Available test users</div>
  <div class="list">
    {#each users as user}
      <div
        class="user"
        class:selected={$selectedUser == user.userId}
        on:click={() => selectUser(user.userId)}
      >
        {user.userDirectory}\{user.userId}
      </div>
    {/each}
  </div>
</div>

<style>
  .user {
    background-color: #3f3f46;
    cursor: pointer;
    padding: 0.5rem;
    text-align: left;
    border: 1px solid transparent;
  }

  .user:hover {
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
