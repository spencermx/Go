<!-- src/routes/video/+page.svelte -->
<script>
  import { page } from '$app/stores';

  let videoUrl = $page.url.searchParams.get('url');
  let videoAlt = $page.url.searchParams.get('alt');
  let videoCaptionsUrl = $page.url.searchParams.get('captions');
</script>
<div class="video-container">
  <div class="video-wrapper">
    <video controls crossorigin="anonymous">
      <source src={videoUrl} type="video/mp4">
      {#if videoCaptionsUrl}
        <track src={videoCaptionsUrl} kind="captions" srclang="en" label="English" default crossorigin="anonymous">
      {/if}
      Your browser does not support the video tag.
    </video>
  </div>
  <h2>{videoAlt}</h2>
</div>
<style>
  .video-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: calc(100vh - 60px);
    padding: 20px;
    box-sizing: border-box;
  }

  .video-wrapper {
    width: 100%;
    max-width: 800px;
    position: relative;
    padding-bottom: 56.25%; /* 16:9 aspect ratio */
    height: 0;
  }

  video {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  h2 {
    margin-top: 10px;
    text-align: center;
  }
</style>
