export async function fetchVideos() {
    const res = await fetch('http://localhost:8080/api/videos');
    return res.json();
}
