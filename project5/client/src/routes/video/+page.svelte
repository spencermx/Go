<script>
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  let videoUrl = $page.url.searchParams.get('url');
  let videoAlt = $page.url.searchParams.get('alt');
  let videoCaptionsUrl = $page.url.searchParams.get('captions');
  let starttime = $page.url.searchParams.get('starttime');

  let videoElement;

  onMount(async () => {
    videoElement.addEventListener('loadedmetadata', () => {
      videoElement.currentTime = starttime;
    });
  });
</script>

<div class="video-container">
  <video bind:this={videoElement} controls crossorigin="anonymous">
    <source src={videoUrl} type="video/mp4" />
    {#if videoCaptionsUrl}
      <track
        src={videoCaptionsUrl}
        kind="captions"
        srclang="en"
        label="English"
        default
        crossorigin="anonymous"
      />
    {/if}
    Your browser does not support the video tag.
  </video>
</div>
<style>
  /* Add these styles to the root elements */
  html, body {
    height: 100%;
    margin: 0;
    padding: 0;
    overflow: hidden;
  }

  body {
    display: flex;
    flex-direction: column;
  }

  .video-container {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-grow: 1;
  }
  video {
    max-width: 100%; /* This ensures that the video never exceeds the width of its parent container */
    max-height: 100%; /* This ensures that the video never exceeds the height of its parent container */
  }
</style>
