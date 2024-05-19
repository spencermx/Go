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
  <div class="container">
    <main>
      <section class="search my-4">
        <input type="text" class="form-control" placeholder="Search videos..." bind:value={searchQuery} />
      </section>

      <section class="gallery">
        <div class="row row-cols-2 row-cols-md-4 g-4">
          {#each filteredFiles as file (file.url)}
            <div class="col">
              <a href="/video?url={encodeURIComponent(file.url)}&alt={encodeURIComponent(file.alt)}&captions={encodeURIComponent(file.videocaptionsurl)}" class="card-link">
                <div class="card">
                  <div class="card-img-container">
                    <img src={file.videothumbnailurl} alt={file.alt} class="card-img" />
                  </div>
                  {#if file.alt}
                    <div class="card-body">
                      <h6 class="card-title">{file.alt}</h6>
                    </div>
                  {/if}
                </div>
              </a>
            </div>
          {:else}
            <p>No videos found.</p>
          {/each}
        </div>
      </section>
      
      <section class="upload my-4">
        <h2>Upload Video</h2>
        <form action="/uploadVideo" method="post" enctype="multipart/form-data">
          <div class="mb-3">
            <label for="file-input" class="form-label">Choose File</label>
            <input type="file" name="file" id="file-input" class="form-control" bind:files={selectedFile} accept="video/*" required />
          </div>
          {#if selectedFile}
            <div class="mb-3">
              <div class="video-preview-container">
                <video src={URL.createObjectURL(selectedFile[0])} controls class="video-preview"></video>
              </div>
            </div>
          {/if}
          <button type="submit" class="btn btn-primary">Upload</button>
        </form>
      </section>

    </main>
  </div>
</body>

<style>
  .card-link {
    text-decoration: none;
    color: inherit;
  }

  .card {
    border: none;
    border-radius: 0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.3s ease;
  }

  .card:hover {
    transform: translateY(-5px);
  }

  .card-img-container {
    width: 100%;
    height: 0;
    padding-bottom: 66.67%; /* 3:2 aspect ratio */
    position: relative;
    overflow: hidden;
  }

  .card-img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .card-body {
    padding: 0.5rem;
  }

  .card-title {
    margin: 0;
    font-size: 0.9rem;
    font-weight: bold;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .video-preview-container {
    width: 100%;
    max-width: 300px;
    margin: 0 auto;
  }

  .video-preview {
    width: 100%;
    height: auto;
    object-fit: contain;
  }

  @media (max-width: 576px) {
    .card-img-container {
      padding-bottom: 100%; /* 1:1 aspect ratio for mobile */
    }

    .card-title {
      font-size: 0.8rem;
    }

    .video-preview-container {
      max-width: 200px;
    }
  }
</style>
