<script> 
  import { onMount } from 'svelte';

  let files = [];
  let selectedFile = null;

  onMount(async () => {
    try {
      const response = await fetch('/getVideos');
      if (response.ok) {
        files = await response.json();
        files.forEach(file => {
          const fileUrl = file.url;
          const fileExtension = fileUrl.toLowerCase().substr(fileUrl.lastIndexOf('.'));
          const isImageFile = isImage(fileUrl);
          console.log('File:', fileUrl);
          console.log('Extension:', fileExtension);
          console.log('Is Image:', isImageFile);
        });
      } else {
        console.error('Error fetching files:', response.status);
      }
    } catch (error) {
      console.error('Error fetching files:', error);
    }
  });

  function isImage(fileUrl) {
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif'];
    const fileExtension = fileUrl.toLowerCase().substr(fileUrl.lastIndexOf('.'));
    const isImageFile = imageExtensions.includes(fileExtension);
    console.log('File:', fileUrl);
    console.log('Extension:', fileExtension);
    console.log('Is Image:', isImageFile);
    return isImageFile;
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
     <section class="gallery">
        <h2>Gallery</h2>
        <div class="file-grid">
          {#each files as file}
            <div class="file-item">
                <video controls>
                  <source src={file.url} type="video/mp4">
                  Your browser does not support the video tag.
                </video>
              <div class="file-overlay">
                {#if file.alt}
                  <h3>{file.alt}</h3>
                {/if}
                {#if file.alt}
                  <p>{file.alt}</p>
                {/if}
              </div>
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
  .file-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-gap: 20px;
  }

  .file-item {
    border: 1px solid #ccc;
    padding: 10px;
    text-align: center;
  }

  .file-item img,
  .file-item video {
    max-width: 100%;
    max-height: 200px;
    object-fit: cover;
  }
  .file-preview {
    margin-top: 20px;
    border: 1px solid #ccc;
    padding: 10px;
    background-color: #f5f5f5;
    text-align: center;
  }

  .file-preview video {
    max-width: 300px; /* Set the maximum width */
    max-height: 200px; /* Set the maximum height */
    object-fit: contain; /* Maintain aspect ratio */
  }
</style>
