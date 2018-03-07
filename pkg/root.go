// Copyright © 2017 Calum Gardner <calum@chronojam.co.uk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chronojam/dexy/pkg/providers"
	"github.com/coreos/go-oidc"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
	"github.com/pressly/chi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dexy",
	Short: "",
	Long:  `Dexy is a simple application used to grab an oauth2 token from a provider`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(viper.GetString("token_file"))
		if err == nil {
			// Check if expiry is what we expect.
			var tok returnToken
			err = json.Unmarshal(b, &tok)
			if err != nil {
				log.Fatalf("error while unmarshalling token from file %v", err)
			}

			if tok.ExpiryTime.After(time.Now()) {
				// We've not expired, so just return the token from the file.
				fmt.Println(string(b))
				return
			}
		}

		// Our token has either expired or doesn't exist.
		provider, err := oidc.NewProvider(context.Background(), viper.GetString("identity_host"))
		if err != nil {
			log.Fatalf("error while creating new oidc provider %v", err)
		}
		oauth2Config := oauth2.Config{
			ClientID:     C.ClientID,
			ClientSecret: C.ClientSecret,
			RedirectURL:  C.CallbackURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       append([]string{oidc.ScopeOpenID, "email", "profile"}, viper.GetStringSlice("auth.scopes")...),
		}

		tokenChan := make(chan *returnToken)
		w := &web{
			verifier:  provider.Verifier(&oidc.Config{ClientID: C.ClientID}),
			cfg:       oauth2Config,
			tokenChan: tokenChan,
		}

		var authCodeURL string = oauth2Config.AuthCodeURL("")
		if prov, ok := C.Provider.(providers.IProvider); ok {
			requestParams, err := prov.BuildRequestParameters()
			if err != nil {
				log.Println("Error calling BuildRequestParameters %v", err)
			}
			log.Println(requestParams)
			authCodeURL += requestParams

		}
		err = browser.OpenURL(authCodeURL)
		if err != nil {
			log.Fatalf("error while opening new web browser %v", err)
		}
		go w.Serve()
		tok := <-tokenChan

		b, err = json.Marshal(tok)
		if err != nil {
			log.Fatalf("error while marshalling token from provider %v", err)
		}

		err = ioutil.WriteFile(viper.GetString("token_file"), b, 0755)
		if err != nil {
			log.Fatalf("error while attempting to write token to file %v", err)
		}

		fmt.Println(string(b))
	},
}

type returnToken struct {
	AccessToken string    `json:"access_token"`
	ExpiryTime  time.Time `json:"expiry_time"`
}

type web struct {
	verifier  *oidc.IDTokenVerifier
	cfg       oauth2.Config
	tokenChan chan *returnToken
	provider  *string
}

func (s *web) Serve() {
	r := chi.NewRouter()

	oauth := chi.NewRouter()
	oauth.Get("/callback", s.oauth2Callback)
	r.Mount("/oauth2", oauth)

	http.ListenAndServe(":10111", r)
}

func (s *web) oauth2Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oauth2Token, err := s.cfg.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		log.Fatalf("error while attempting initial token exchange %v", err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
	}

	// Parse and verify ID Token payload.
	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		fmt.Println(err.Error())
		// handle error
	}

	// C.Finalise()

	ret := &returnToken{
		AccessToken: rawIDToken,
		ExpiryTime:  idToken.Expiry,
	}
	s.tokenChan <- ret

	fmt.Fprintf(w, "Done, you can now close this window")
	// }
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	browser.Stdout = ioutil.Discard
	browser.Stderr = ioutil.Discard

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dexy.yaml)")
}

var C providers.Config

var destinations = []providers.IProvider{
	providers.BaseProvider{},
	providers.GoogleApps{},
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".dexy" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/dexy")
		viper.AddConfigPath(".")
		viper.SetConfigName(".dexy")
	}
	viper.SetEnvPrefix("dexy")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetDefault("token_file", home+"/.dexy-token.yaml")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
	}

	if viper.IsSet("callback_url") {
		viper.Set("callback_url", fmt.Sprintf("http://%s:%d/oauth2/callback", viper.GetString("callback_host"), viper.GetInt("callback_port")))
	}

	for _, provider := range destinations {
		if ok, err := provider.IsSetCorrectly(); err == nil && ok {
			if err := viper.Unmarshal(&C); err != nil {
				log.Printf("Could not unmarshal config %+v", err)
			}

			log.Printf("%+v", C)
		}

	}
}
