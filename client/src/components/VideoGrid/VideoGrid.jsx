import React from 'react';
import VideoCard from '../VideoCard/VideoCard';
import './VideoGrid.css';

function VideoGrid({ videos, onVideoSelect }) {
  return (
    <div className="video-grid">
      {videos.map(video => (
        <VideoCard key={video.id} video={video} onClick={() => onVideoSelect(video)} />
      ))}
    </div>
  );
}

export default VideoGrid;
