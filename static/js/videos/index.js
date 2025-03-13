
function getSample() {
    const path = window.location.pathname;
    const videoName = path.split('/').pop()

    if(videoName) {
        const videoElement = document.getElementById('videoPlayer')
        videoElement.src = '/video/' + videoName
    }
}

document.addEventListener("DOMContentLoaded", function() {
    getSample()
});
