import React from 'react';
import './VideoPlayerModal.css';

function VideoPlayerModal({ video, onClose }) { 
  console.log(video)

  if (!video) return null;
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <video controls autoPlay className="video-player">
          <source src={video.videoUrl} type="video/mp4" />
          Your browser does not support the video tag.
        </video>
        <button className="close-btn" onClick={onClose}>Close</button>
      </div>
    </div>
  );
}

export default VideoPlayerModal;
