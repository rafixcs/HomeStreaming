export async function fetchVideos() {
    const res = await fetch('http://localhost:8080/videos');
    console.log(res)
    return res.json();
}
