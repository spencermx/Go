<script>
  import { onMount } from 'svelte';

  let people = [];

  onMount(async () => {
    try {
      const response = await fetch('http://localhost:8080/getPeople');
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
