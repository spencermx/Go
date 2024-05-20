<script lang="ts">
  import { onMount } from 'svelte';
  import lunr from 'lunr';

  /******** TITLE SEARCH VIDEOS ***********/
  let files =[
  	{
  		"videoId": "56dca040-60d2-4582-8184-30235584dd73",
  		"videoName": "Richard Feynman-The Character of Physical Law Part6: Probability and Uncertainty",
  		"videoUrl": "https://d271tjczb1hjof.cloudfront.net/56dca040-60d2-4582-8184-30235584dd73-Richard Feynman-The Character of Physical Law Part6: Probability and Uncertainty.mp4",
  		"videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/56dca040-60d2-4582-8184-30235584dd73-captions.vtt",
  		"videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/56dca040-60d2-4582-8184-30235584dd73-thumbnail.jpg"
  	},
  	{
  		"videoId": "9a77d400-3ebc-4543-9868-084cb9c46ef6",
  		"videoName": "Richard Feynman-The Character of Physical Law Part2: The Relation of Mathematics to Physics",
  		"videoUrl": "https://d271tjczb1hjof.cloudfront.net/9a77d400-3ebc-4543-9868-084cb9c46ef6-Richard Feynman-The Character of Physical Law Part2: The Relation of Mathematics to Physics.mp4",
  		"videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/9a77d400-3ebc-4543-9868-084cb9c46ef6-captions.vtt",
  		"videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/9a77d400-3ebc-4543-9868-084cb9c46ef6-thumbnail.jpg"
  	},
  	{
  		"videoId": "9e1e2dd4-c836-43af-ba21-090b9a1032d3",
  		"videoName": "Richard Feynman Messenger Lectures at Cornell   The Character of Physical Law   Part 1 The Law of Gravitation",
  		"videoUrl": "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-Richard Feynman Messenger Lectures at Cornell   The Character of Physical Law   Part 1 The Law of Gravitation.mp4",
  		"videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-captions.vtt",
  		"videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-thumbnail.jpg"
  	},
  	{
  		"videoId": "d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5",
  		"videoName": "Clip Richard Feynman: Mathematicians versus Physicists",
  		"videoUrl": "https://d271tjczb1hjof.cloudfront.net/d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5-Clip Richard Feynman: Mathematicians versus Physicists.mp4",
  		"videoCaptionsUrl": "https://d271tjczb1hjof.cloudfront.net/d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5-captions.vtt",
  		"videoThumbnailUrl": "https://d271tjczb1hjof.cloudfront.net/d5ef5f9a-8c76-4a72-858c-d2f0c9c70da5-thumbnail.jpg"
  	}
  ];

  let filteredFiles = [];
  let selectedFile = null;
  let searchQuery = '';
  /******** TITLE SEARCH VIDEOS ***********/

  /******** TEXT SEARCH VIDEOS ***********/
  interface VTTCue {
   startTime: number;
   endTime: number;
   text: string;
  }

  interface VTTData {
   video: string;
   cues: VTTCue[];
  }

  interface SearchResultItem {
   videoId: string;
   videoName: string;   
  // when search text content display the following
  // the video name 
  // the video time stamp of the match AND the caption/doc it was found in 
  }

  let vttUrls = [ "https://d271tjczb1hjof.cloudfront.net/9e1e2dd4-c836-43af-ba21-090b9a1032d3-captions.vtt" ]
  let textSearchQuery = ""
  let searchIndex: lunr.Index;
  let vttData: VTTData[] = [];
  let searchResults: VTTData[] = [];
  /******** TEXT SEARCH VIDEOS ***********/

  onMount(async () => {
    vttData = await loadVTTFiles();
    searchIndex = buildSearchIndex(vttData);

    $: {
      filteredFiles = files.filter(file => {
        const videoNameText = file.videoName ? file.videoName.toLowerCase() : '';
        return videoNameText.includes(searchQuery.toLowerCase());
      });
    }
  })
//  //  try {
//      const response = await fetch('/getVideos');
//      if (response.ok) {
//        files = await response.json();
//      } else {
//        console.error('Error fetching files:', response.status);
//      }
//    } catch (error) {
//      console.error('Error fetching files:', error);
//    }
//  });
//

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
         cues.push({ startTime, endTime, text });
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
   
   async function loadVTTFiles(): Promise<VTTData[]> {
     const vttData: VTTData[] = [];
   
     for (const vttUrl of vttUrls) {
       const response = await fetch(vttUrl);
       const vttContent = await response.text();
       const cues = parseVTTFile(vttContent);
   
       vttData.push({
         video: vttUrl,
         cues: cues
       });
     }
   
     return vttData;
   }
 function buildSearchIndex(vttData: VTTData[]) {
  const index = lunr(function (this: any) {
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
      const doc = vttData[i];
      for (let j = 0; j < doc.cues.length; j++) {
        const cue = doc.cues[j];
        const uniqueId = `${doc.video}_${j}`;
        this.add({
          id: uniqueId,
          text: cue.text
        });
      }
    }

  });

  let results0 = index.search("a");
  let results1 = index.search("along with the earth");
  let results2 = index.search("there is some reason")
  let results3 = index.search("some reason");

  let results4 = index.search("planetary");
  let results5 = index.search("planetary reason");
  let results6 = index.search("planetary calling");
  return index;
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



function searchVideos(query: string, index: lunr.Index, vttData: VTTData[]): VTTData[] {
  const results = index.search(query);
  const matchingVideos: VTTData[] = [];

  results.forEach((result: lunr.Index.Result) => {
    const [videoUrl, cueIndex] = result.ref.split('_');
    const video = vttData.find(data => data.video === videoUrl);
    if (video) {
      matchingVideos.push(video);
    }
  });

  return matchingVideos;
}
</script>

<body>
  <div class="container">
    <main>
      <section class="search my-4">
        <input type="text" class="form-control" placeholder="Search videos..." bind:value={searchQuery} />
      </section>
      <section class="search my-4">
        <input type="text" class="form-control" placeholder="Search videos..." bind:value={textSearchQuery} on:input={performSearch} />
      </section>

      <section class="search-results my-4">
        {#if searchResults.length > 0}
            <ul>
                {#each searchResults as result}
                    <li>
                    <strong>{result.videoName}</strong> - Matched: {result.match}
                    </li>
                {/each}
            </ul>
        {:else if searchQuery}
            <p>No videos found.</p>
        {/if}
      </section>



      <section class="gallery">
        <div class="row row-cols-2 row-cols-md-4 g-4">
          {#each filteredFiles as file (file.videoUrl)}
            <div class="col">
              <a href="/video?url={encodeURIComponent(file.videoUrl)}&videoName={encodeURIComponent(file.videoName)}&captions={encodeURIComponent(file.videoCaptionsUrl)}" class="card-link">
                <div class="card">
                  <div class="card-img-container">
                    <img src={file.videoThumbnailUrl} videoName={file.videoName} class="card-img" />
                  </div>
                  {#if file.videoName}
                    <div class="card-body">
                      <h6 class="card-title">{file.videoName}</h6>
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
