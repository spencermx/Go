<script>
  const maintenanceMessage = "We are currently performing maintenance to address AWS S3 bucket vulnerabilities. We apologize for any inconvenience caused and appreciate your patience. Please check back later.";
</script>

<style>
  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    background-color: #f8f8f8;
    font-family: Arial, sans-serif;
    text-align: center;
  }

  h1 {
    font-size: 24px;
    margin-bottom: 20px;
  }

  p {
    font-size: 18px;
    max-width: 500px;
    margin-bottom: 30px;
  }
</style>

<div class="container">
  <h1>Down for Maintenance (May 7th 2024)</h1>
  <p>{maintenanceMessage}</p>
</div>

<!--<script> 
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
<body>
  <div class="content-wrapper">
    <header>
      <h1>Gallery</h1>
      <nav>
        <ul>
          <li><a href="/">Home</a></li>
          <li><a href="/contact">Contact</a></li>
        </ul>
      </nav>
    </header>

    <main>
      <section class="gallery">
        <h2>Image Gallery</h2>
        <div class="image-grid">
          {#each images as image}
            <div class="image-item">
              <img src={image.url} alt={image.alt} />
              <div class="image-overlay">
                {#if image.title}
                  <h3>{image.title}</h3>
                {/if}
                {#if image.description}
                  <p>{image.description}</p>
                {/if}
              </div>
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
  </div>
</body>

<style>
  body {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    margin: 0;
    font-family: 'Roboto', sans-serif;
    color: var(--text-color);
    background-color: var(--background-color);
  }

  .content-wrapper {
    flex: 1;
  }

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
    position: relative;
  }

  .image-item img {
    width: 100%;
    height: auto;
  }

  .image-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: rgba(0, 0, 0, 0.7);
    color: #fff;
    padding: 10px;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .image-item:hover .image-overlay {
    opacity: 1;
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

  header {
    background-color: #333;
    color: #fff;
    padding: 20px;
  }

  h1 {
    margin: 0;
    font-size: 28px;
  }

  nav ul {
    list-style-type: none;
    margin: 0;
    padding: 0;
  }

  nav ul li {
    display: inline-block;
    margin-right: 10px;
  }

  nav ul li a {
    color: #fff;
    text-decoration: none;
  }

  footer {
    background-color: #f5f5f5;
    padding: 20px;
    text-align: center;
    margin-top: 40px;
  }
</style>-->
