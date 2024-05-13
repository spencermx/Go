<script>
  import { onMount } from 'svelte';

  let files = [];
  let filteredFiles = [];
  let selectedFile = null;
  let searchQuery = '';

  onMount(async () => {
    try {
      const response = await fetch('/getVideos');
      if (response.ok) {
        files = await response.json();
      } else {
        console.error('Error fetching files:', response.status);
      }
    } catch (error) {
      console.error('Error fetching files:', error);
    }
  });

  $: {
    filteredFiles = files.filter(file => {
      const altText = file.alt ? file.alt.toLowerCase() : '';
      return altText.includes(searchQuery.toLowerCase());
    });
  }
</script>

<body>
  <div class="content-wrapper">
    <header>
      <h1>Gallery</h1>
      <nav>
        <ul>
          <li><a href="/">Home</a></li>
        </ul>
      </nav>
    </header>

    <main>
      <section class="search">
        <input type="text" placeholder="Search videos..." bind:value={searchQuery} />
      </section>

      <section class="gallery">
        <h2>Gallery</h2>
        <div class="file-grid">
          {#each filteredFiles as file (file.url)}
            <div class="file-item" key={file.url}>
              <video controls crossorigin="anonymous">
                <source src={file.url} type="video/mp4">
                {#if file.videocaptionsurl}
                  <track src={file.videocaptionsurl} kind="captions" srclang="en" label="English" default crossorigin="anonymous">
                {/if}
                Your browser does not support the video tag.
              </video>
              {#if file.alt}
                <div class="file-title" title={file.alt}>
                  {file.alt}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      </section>

      <section class="upload">
        <h2>Upload Video</h2>
        <form action="/uploadVideo" method="post" enctype="multipart/form-data">
          <label for="file-input" class="file-label">Choose File</label>
          <input type="file" name="file" id="file-input" class="file-input" bind:files={selectedFile} accept="video/*" required/>
          {#if selectedFile}
            <div class="file-preview">
              <video src={URL.createObjectURL(selectedFile[0])} controls></video>
            </div>
          {/if}
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

  .search {
    margin-bottom: 20px;
  }

  .search input {
    width: 100%;
    padding: 10px;
    font-size: 16px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  .file-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-gap: 20px;
  }

.file-item {
  border: 1px solid #ccc;
  padding: 10px;
  text-align: center;
  position: relative;
  overflow: hidden;
  width: 200px;
  margin: 0 auto;
}

.file-item:hover {
  height: auto;
}

.file-item video {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.file-title {
  margin-top: 10px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
  height: 20px; /* Initial height */
  transition: height 0.3s ease; /* Add smooth transition */
}

.file-item:hover .file-title {
  white-space: normal;
  overflow: visible;
  height: auto; /* Expand height on hover */
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

  .file-preview {
    margin-top: 20px;
    border: 1px solid #ccc;
    padding: 10px;
    background-color: #f5f5f5;
    text-align: center;
  }

  .file-preview video {
    max-width: 300px;
    max-height: 200px;
    object-fit: contain;
  }
</style>
