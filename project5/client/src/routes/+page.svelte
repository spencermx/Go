<script lang="ts">
  import { onMount } from 'svelte';
  import lunr from 'lunr';

  /*************** STRUCTURES **************/
  interface VTTData {
   video: Video;
   cues: VTTCue[];
  }

  interface Video {
    videoId: string;
    videoName: string;
    videoUrl: string;
    videoCaptionsUrl: string;
    videoThumbnailUrl: string;
  }

  interface VTTCue {
   startTime: number;
   endTime: number;
   text: string;
   timeRange: string;
  }

  interface SearchResultItem {
   video: Video;
   vttCue: VTTCue;
  }
  /*************** STRUCTURES **************/

  /******** TITLE SEARCH VIDEOS ***********/
  let json = `[ { "videoId": "56dca040-60d2-4582-8184-30235584dd73", "videoName": "Richard Feynman-The Character of Physical Law Part6: Probability and Uncertainty", "videoUrl": "https://d271tjczb1hjof.cloudfront.net/56dca040-60d2-4582-8184-30235584dd73-Richard Feynman-The Character of Physical Law Part6: Probability and Uncertainty.mp4", "videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/56dca040-60d2-4582-8184-30235584dd73-captions.vtt", "videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/56dca040-60d2-4582-8184-30235584dd73-thumbnail.jpg" }, { "videoId": "9a77d400-3ebc-4543-9868-084cb9c46ef6", "videoName": "Richard Feynman-The Character of Physical Law Part2: The Relation of Mathematics to Physics", "videoUrl": "https://d271tjczb1hjof.cloudfront.net/9a77d400-3ebc-4543-9868-084cb9c46ef6-Richard Feynman-The Character of Physical Law Part2: The Relation of Mathematics to Physics.mp4", "videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/9a77d400-3ebc-4543-9868-084cb9c46ef6-captions.vtt", "videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/9a77d400-3ebc-4543-9868-084cb9c46ef6-thumbnail.jpg" }, { "videoId": "9e1e2dd4-c836-43af-ba21-090b9a1032d3", "videoName": "Richard Feynman Messenger Lectures at Cornell   The Character of Physical Law   Part 1 The Law of Gravitation", "videoUrl": "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-Richard Feynman Messenger Lectures at Cornell   The Character of Physical Law   Part 1 The Law of Gravitation.mp4", "videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-captions.vtt", "videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-thumbnail.jpg" }, { "videoId": "d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5", "videoName": "Clip Richard Feynman: Mathematicians versus Physicists", "videoUrl": "https://d271tjczb1hjof.cloudfront.net/d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5-Clip Richard Feynman: Mathematicians versus Physicists.mp4", "videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5-captions.vtt", "videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5-thumbnail.jpg" } ]`
  let filteredVideos = [];
  let selectedFile = null;
  let searchQuery = '';
  /******** TITLE SEARCH VIDEOS ***********/

  /******** TEXT SEARCH VIDEOS ***********/
  let textSearchQuery = ""
  let searchIndex: lunr.Index;
  let vttData: VTTData[] = [];
  let searchResults: SearchResultItem[] = [];
  /******** TEXT SEARCH VIDEOS ***********/
  let videos: Video[] = [];

  function buildSearchIndex() {
    searchIndex = lunr(function (this: any) {
      this.ref('id');
      this.field('text');
      // Add the test document to the search index
      this.add({
        id: 'test_document',
        text: 'there is some planetary reason for all of these things where once we all built our houses on the land later through evoloution we found a new planetary calling where up castles were constructed'
      });
      // Add the test document to the search index
      this.add({
        id: 'test_document1',
        text: 'we some reason we live here'
      });

      for (let i = 0; i < vttData.length; i++) {
        const item = vttData[i];
        for (let j = 0; j < item.cues.length; j++) {
          const uniqueId = `${item.video.videoId}_${j}`;
          const cue = item.cues[j];
          this.add({
            id: uniqueId,
            text: cue.text
          });
        }
      }
    });

    let results0 = searchIndex.search("a");
    let results1 = searchIndex.search("along with the earth");
    let results2 = searchIndex.search("there is some reason")
    let results3 = searchIndex.search("some reason");

    let results4 = searchIndex.search("planetary");
    let results5 = searchIndex.search("planetary reason");
    let results6 = searchIndex.search("planetary calling");
  }  

  function performSearch() {
    if (searchIndex) {
      if (textSearchQuery.trim() === '') {
        searchResults = [];
      } else {
        searchResults = searchVideos(textSearchQuery, searchIndex, vttData);
        let x = 10
      }
    }
  }

  function searchVideos(query: string, index: lunr.Index, vttData: VTTData[]): SearchResultItem[] {
    const results = index.search(query);
    const searchResults: SearchResultItem[] = [];
  
    results.forEach((result: lunr.Index.Result) => {
      const [videoId, cueIndex] = result.ref.split('_');
      const vttDataItem = vttData.find(data => data.video.videoId === videoId);

      if (vttDataItem) {
        const vttCue = vttDataItem.cues[parseInt(cueIndex)];
        const searchResultItem: SearchResultItem = {
          video: vttDataItem.video,
          vttCue: vttCue
        };
        searchResults.push(searchResultItem);
      }
    });
  
    return searchResults;
  }

  async function loadVTTFiles(): Promise<VTTData[]> {
    vttData = []; 
 
    for (const video of videos) {
      if (video.videoCaptionsUrl) {
        try {
          const response = await fetch(video.videoCaptionsUrl);
          const vttContent = await response.text();
          const cues = parseVTTFile(vttContent);
 
          vttData.push({
            video: video,
            cues: cues
          });
        } catch (error) {
          console.error(`Error fetching captions for video ${video.videoId}:`, error);
        }
      }
    }
 
    return vttData;
  } 

  function parseVTTFile(vttContent: string): VTTCue[] {
    const lines = vttContent.trim().split('\n');
    const cues: VTTCue[] = [];
  
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i].trim();
  
      if (line.startsWith('NOTE') || line.startsWith('WEBVTT')) {
        continue;
      }
  
      if (line.includes('-->')) {
        const [startTime, endTime] = line.split('-->').map(parseTimestamp);
        const text = lines[++i].trim();
        cues.push({
          startTime: startTime, 
          endTime: endTime, 
          text: text,
          timeRange: line});
      }
    }
  
    return cues;
  }
   
  function parseTimestamp(timestamp: string): number {
    const [hours, minutes, seconds] = timestamp.trim().split(':');
    return (
      parseFloat(hours) * 3600 +
      parseFloat(minutes) * 60 +
      parseFloat(seconds.replace(',', '.'))
    );
  }

  function createAnchorTagForCard(video) {
    const videoUrl = encodeURIComponent(video.videoUrl);
    const videoName = encodeURIComponent(video.videoName);
    const videoCaptionsUrl = encodeURIComponent(video.videoCaptionsUrl);

    let nameTag = ``
    let result = `<a class="card-link" href="/video?url=${videoUrl}&videoName=${videoName}&captions=${videoCaptionsUrl}">
                <div class="card">
                  <div class="card-img-container">
                    <img src=${video.videoThumbnailUrl} videoName=${video.videoName} class="card-img" />`

    if (video.videoName) {

        result += `<div class="card-body">
                      <h6 class="card-title">${video.videoName}</h6>
                   </div>`;

    }

    result += `</div></div></a>`;

    return result;
  }
  
  function createAnchorTagForTimeStamp(result) {
    const videoUrl = encodeURIComponent(result.video.videoUrl);
    const videoName = encodeURIComponent(result.video.videoName);
    const videoCaptionsUrl = encodeURIComponent(result.video.videoCaptionsUrl);
    const timestamp = result.vttCue.timeRange.split('-->')[0].trim();
    const displayText = result.vttCue.text.slice(0, 100) 
    return `<a href="/video?url=${videoUrl}&videoName=${videoName}&captions=${videoCaptionsUrl}&starttime=${result.vttCue.startTime}">${displayText}</a>`;
  }
  onMount(async () => {
    videos = JSON.parse(json);
    // videos = await fetchVideos();
    await loadVTTFiles();

    buildSearchIndex();

  })

  $: {
    filteredVideos = videos.filter(video => {
      const videoNameText = video.videoName ? video.videoName.toLowerCase() : '';
      return videoNameText.includes(searchQuery.toLowerCase());
    });
  }


</script>

<body>
  <div class="container">
    <main>
      <section class="search my-4">
        <input type="text" class="form-control" placeholder="Search videos by title..." bind:value={searchQuery} />
      </section>

      <section class="search my-4">
        <input type="text" class="form-control" placeholder="Search videos by text content..." bind:value={textSearchQuery} on:input={performSearch} />
      </section>

      <section class="search-results my-4">
        <div class="search-results-container">
        {#if searchResults.length > 0}
          <ul>
            {#each searchResults as result, i}
              {#if i === 0 || result.video.videoId !== searchResults[i - 1].video.videoId}
                <li>
                  <strong>{result.video.videoName}</strong>
                  <ul>
                    <li>
                      {result.vttCue.timeRange.split('-->')[0].trim()} {@html createAnchorTagForTimeStamp(result)}
                    </li>
                  </ul>
                </li>
              {:else}
                <li>
                  <ul>
                    <li>
                      {result.vttCue.timeRange.split('-->')[0].trim()} {@html createAnchorTagForTimeStamp(result)}
                    </li>
                  </ul>
                </li>
              {/if}
            {/each}
          </ul>
        {:else if searchQuery}
          <p>No videos found.</p>
        {/if}
        </div>
      </section>

      <section class="gallery">
        <div class="row row-cols-2 row-cols-md-4 g-4">
          {#each filteredVideos as file (file.videoUrl)}
            <div class="col">
              {@html createAnchorTagForCard(file)}
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

  .search-results {
    margin-top: 1rem;
  }

  .search-results ul {
    list-style-type: none;
    padding: 0;
  }

  .search-results li {
  }

  .video-title {
    font-size: 1.2rem;
    font-weight: bold;
    margin-bottom: 0.5rem;
  }

  .match-item {
    background-color: #f5f5f5;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 0.5rem;
  }

  .match-text {
    margin-bottom: 0.25rem;
  }

  .match-timestamp {
    font-size: 0.9rem;
    color: #666;
  }
 .search-results-container {
    max-height: 300px;
    overflow-y: auto;
  }
</style>
