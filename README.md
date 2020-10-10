# spotify-authenticate

Simple utility to obtain a refreshable user token for Spotify.

## Usage

1. [Create a Spotify app](https://developer.spotify.com/dashboard/)
2. In the Spotify app settings - add `http://localhost:3000/` to
   `Redirect URIs`
3. Clone this repo
4. Add `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET` to your
   environment
5. `go run main.go`

