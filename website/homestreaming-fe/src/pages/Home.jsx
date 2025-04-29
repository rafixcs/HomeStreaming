import React, { useEffect, useState } from 'react';
import { fetchVideos } from '../api/videos';
import VideoGrid from '../components/VideoGrid/VideoGrid';
import VideoPlayerModal from '../components/VidePlayerModal/VideoPlayerModal';

function Home() {
  const [videos, setVideos] = useState([]);
  const [selectedVideo, setSelectedVideo] = useState(null);

  const mockVideos = [
    {
      id: "1",
      title: "Big Buck Bunny",
      thumbnailUrl: "https://peach.blender.org/wp-content/uploads/title_anouncement.jpg?x11217",
      videoUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4"
    },
    {
      id: "2",
      title: "Sintel",
      thumbnailUrl: "https://mango.blender.org/wp-content/uploads/2013/05/sintel_poster.jpg",
      videoUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4"
    },
    {
      id: "3",
      title: "Tears of Steel",
      thumbnailUrl: "https://mango.blender.org/wp-content/uploads/2013/05/tears_of_steel_poster.jpg",
      videoUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4"
    },
    {
      id: "4",
      title: "Elephants Dream",
      thumbnailUrl: "https://orange.blender.org/wp-content/themes/orange/images/common/ed_head.jpg",
      videoUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4"
    },
    {
      id: "5",
      title: "For Bigger Blazes",
      thumbnailUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerBlazes.jpg",
      videoUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4"
    },
    {
      id: "6",
      title: "For Bigger Escape",
      thumbnailUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerEscape.jpg",
      videoUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscape.mp4"
    }
  ];

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
