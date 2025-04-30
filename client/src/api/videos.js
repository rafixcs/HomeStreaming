export async function fetchVideos() {
    const endpoint = process.env.REACT_APP_API_URL + '/videos' || 'http://192.168.0.14:8080/videos'

    const res = await fetch(endpoint);
    console.log(res)
    return res.json();
}
