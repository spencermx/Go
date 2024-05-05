<script>
  import { onMount } from 'svelte';

  let people = [];
  let images = [];

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

    try {
      const response = await fetch('/getImages');
      if (response.ok) {
        images = await response.json();
      } else {
        console.error('Error fetching images:', response.status);
      }
    } catch (error) {
      console.error('Error fetching images:', error);
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

<h1>Images</h1>
<div class="image-gallery">
  {#each images as image}
    <img src={image.url} alt={image.alt} />
  {/each}
</div>

<form action="/upload" method="post" enctype="multipart/form-data">
    <input type="file" name="file">
    <br><br>
    <input type="submit" value="Upload">
</form>
