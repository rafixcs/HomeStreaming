export async function fetchVideos() {
    const res = await fetch('http://192.168.0.14:8080/videos');
    //const res = await fetch('http://localhost:8080/videos');
    console.log(res)
    return res.json();
}
