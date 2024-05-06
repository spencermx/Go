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

<main>
  <section class="gallery">
    <h2>Image Gallery</h2>
    <div class="image-grid">
      {#each images as image}
        <div class="image-item">
          <img src={image.url} alt={image.alt} />
        </div>
      {/each}
    </div>
  </section>

  <section class="upload">
    <h2>Upload Image</h2>
    <form action="/upload" method="post" enctype="multipart/form-data">
      <label for="file-input" class="file-label">Choose File</label>
      <input type="file" name="file" id="file-input" class="file-input" />
      <button type="submit" class="upload-button">Upload</button>
    </form>
  </section>
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  section {
    margin-bottom: 40px;
  }

  h2 {
    font-size: 24px;
    margin-bottom: 20px;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    margin-bottom: 10px;
  }

  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-gap: 20px;
  }

  .image-item {
    border: 1px solid #ccc;
    border-radius: 4px;
    overflow: hidden;
  }

  .image-item img {
    width: 100%;
    height: auto;
  }

  .file-label {
    display: inline-block;
    padding: 10px 20px;
    background-color: #007bff;
    color: #fff;
    border-radius: 4px;
    cursor: pointer;
  }

  .file-input {
    display: none;
  }

  .upload-button {
    display: block;
    width: 100%;
    padding: 10px;
    background-color: #28a745;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-top: 10px;
  }
</style>
