package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	goenv "github.com/Netflix/go-env"

	"github.com/rs/xid"
	spotify "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type environment struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectPort string `env:"REDIRECT_PORT"`
}

var env environment

func authSpotify(authCallback func(string), authedCallback func(*oauth2.Token)) error {
	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	auth := spotify.NewAuthenticator(
		"http://localhost:"+env.RedirectPort+"/",
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserFollowRead,
		spotify.ScopeUserLibraryRead,
		spotify.ScopeUserTopRead,
	)

	auth.SetAuthInfo(env.ClientID, env.ClientSecret)

	// generate a unique session identifier
	sid := xid.New().String()

	// create server
	mux := http.NewServeMux()
	server := &http.Server{Addr: ":" + env.RedirectPort, Handler: mux}

	// create cancellable context for local server - used to handle auth redirect
	ctx, cancel := context.WithCancel(context.Background())

	// redirect handler
	redirectHandler := func(w http.ResponseWriter, r *http.Request) {
		// use the same state string here that you used to generate the URL
		token, err := auth.Token(sid, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}
		// shutdown server
		defer cancel()

		// the client can now be used to make authenticated requests
		log.Println("Authenticated")

		// callback
		authedCallback(token)
	}
	// get the user to this URL - how you do that is up to you
	// you should specify a unique state string to identify the session
	url := auth.AuthURL(sid)

	// set up local server to receive auth token
	mux.HandleFunc("/", redirectHandler)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// auth callback
	authCallback(url)

	select {
	case <-ctx.Done():
		// Shutdown the server when its context gets canceled
		server.Shutdown(ctx)
	}

	return nil
}

func main() {
	// parse environment
	_, err := goenv.UnmarshalFromEnviron(&env)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Environment", env)

	// spotify
	authCallback := func(url string) {
		// present authentication to user
		fmt.Println("To authenticate go to : ", url)
	}
	authedCallback := func(token *oauth2.Token) {
		fmt.Println("Access token:", token.AccessToken)
		fmt.Println("Refresh token:", token.RefreshToken)
	}

	// start spotify auth
	if err = authSpotify(authCallback, authedCallback); err != nil {
		log.Fatal(err)
	}
}
