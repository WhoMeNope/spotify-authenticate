# spotify-authenticate

Simple utility to obtain a refreshable user token for Spotify.

## Usage
1. Clone this repo
2. [Create a Spotify app](https://developer.spotify.com/dashboard/)
3. Add a redirect url pointing to localhost (e.g. `http://localhost:3000`)
4. Inside the cloned repo, create a `.env` file with `CLIENT_ID`, `CLIENT_SECRET` and `REDIRECT_PORT` set to your app's values.
5. `go run main.go` and follow the instructions there
