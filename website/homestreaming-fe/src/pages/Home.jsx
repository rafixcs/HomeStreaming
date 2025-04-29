import React, { useEffect, useState } from 'react';
import { fetchVideos } from '../api/videos';
import VideoGrid from '../components/VideoGrid/VideoGrid';
import VideoPlayerModal from '../components/VidePlayerModal/VideoPlayerModal';

function Home() {
  const [videos, setVideos] = useState([]);
  const [selectedVideo, setSelectedVideo] = useState(null);

  useEffect(() => {
    fetchVideos().then((value) => {
      console.log(value)
      setVideos(value)
    })
  }, []);

  return (
    <div>
      <h1>My Home Streaming</h1>
      <VideoGrid videos={videos} onVideoSelect={setSelectedVideo} />
      <VideoPlayerModal video={selectedVideo} onClose={() => setSelectedVideo(null)} />
    </div>
  );
}

export default Home;
