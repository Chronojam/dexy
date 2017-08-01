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
	"os"
	"time"

	"encoding/json"
	"github.com/coreos/go-oidc"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pressly/chi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dexy",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(viper.GetString("token_file"))
		if err == nil {
			tok := &oauth2.Token{}
			// Check if expiry is what we expect.
			err = json.Unmarshal(b, tok)
			if err != nil {
				panic(err)
			}

			if !tok.Expiry.Before(time.Now()) {
				fmt.Println(string(b))
				return
			}
		}

		// Our token has either expired or doesnt exist.
		provider, err := oidc.NewProvider(context.Background(), viper.GetString("auth.dex_host"))
		if err != nil {
			panic(err)
		}
		oauth2Config := oauth2.Config{
			ClientID:     viper.GetString("auth.client_id"),
			ClientSecret: viper.GetString("auth.client_secret"),
			RedirectURL:  viper.GetString("auth.callback_url"),
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "groups", "email"},
		}
		tokenChan := make(chan oauth2.Token)
		w := &web{
			cfg:       oauth2Config,
			tokenChan: tokenChan,
		}

		fmt.Println(oauth2Config.AuthCodeURL(""))
		go w.Serve()
		tok := <-tokenChan

		b, err = json.Marshal(tok)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(viper.GetString("token_file"), b, 0755)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))
	},
}

type web struct {
	cfg       oauth2.Config
	tokenChan chan oauth2.Token
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
		panic(err)
	}
	b, err := json.Marshal(oauth2Token)
	if err != nil {
		panic(err)
	}
	// If we get this far, just write our token out to our file.
	err = ioutil.WriteFile(viper.GetString("token_file"), b, 0755)
	if err != nil {
		panic(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
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
	viper.SetDefault("auth", map[string]interface{}{
		"dex_host":      "http://localhost:9999",
		"callback_host": "localhost",
		"callback_port": 10111,
		"client_id":     "dexy",
		"client_secret": "dexy",
	})
	viper.SetDefault("token_file", home+"/.dexy-token.yaml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.Set("auth.callback_url", fmt.Sprintf("http://%s:%d/oauth2/callback",
		viper.GetString("auth.callback_host"),
		viper.GetInt("auth.callback_port")))
}
