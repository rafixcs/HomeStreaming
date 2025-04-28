import React from 'react';
import './VideoCard.css';

function VideoCard({ video, onClick }) {
  return (
    <div className="video-card" onClick={onClick}>
      <img src={video.thumbnailUrl} alt={video.title} className="video-thumb" />
      <div className="video-title">{video.title}</div>
    </div>
  );
}

export default VideoCard;
