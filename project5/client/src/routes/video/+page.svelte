<script>
  import { onMount } from 'svelte';

  let people = [];

  onMount(async () => {
    try {
      const response = await fetch('/getPeople');
      if (response.ok) {
        people = await response.json();
      } else {
        console.error('Error fetching people:', response.status);
      }
    } catch (error) {
      console.error('Error fetching people:', error);
    }
  });
</script>

<h1>People</h1>
{#if people.length === 0}
  <p>Loading people...</p>
{:else}
  <ul>
    {#each people as person}
      <li>{person.name}</li>
    {/each}
  </ul>
{/if}
<h1>testing</h1>

<form action="/upload" method="post" enctype="multipart/form-data">
    <input type="file" name="file">
    <br><br>
    <input type="submit" value="Upload">
</form>
