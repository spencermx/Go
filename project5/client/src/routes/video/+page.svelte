<!-- src/routes/video/+page.svelte -->
<script>
  import { page } from '$app/stores';

  let videoUrl = $page.url.searchParams.get('url');
  let videoAlt = $page.url.searchParams.get('alt');
  let videoCaptionsUrl = $page.url.searchParams.get('captions');
</script>

<div class="video-player">
  <h2>{videoAlt}</h2>
  <video controls crossorigin="anonymous">
    <source src={videoUrl} type="video/mp4">
    {#if videoCaptionsUrl}
      <track src={videoCaptionsUrl} kind="captions" srclang="en" label="English" default crossorigin="anonymous">
    {/if}
    Your browser does not support the video tag.
  </video>
</div>

