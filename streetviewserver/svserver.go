// streetviewserver serves a streetview url that is injected with
// a script
package streetviewserver

import (
	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/urlbuilder"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func modifyMainPage(target string, w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(target)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	bodyAsString := string(body)
	replacedBody := strings.Replace(bodyAsString, "<head>", "<head> <script src=\"static/modify_frontend/modify.js\"/>", -1)
	w.Write([]byte(replacedBody))
}

func modifyInformation(target string, w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(target)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	body = filterStrings(body)

	w.Write(body)
}

// ServeLocation serves a specific location to the user.
func ServeLocation(l s2.LatLng, w http.ResponseWriter, r *http.Request) {
	mapsUrl := urlbuilder.BuildUrl(l)
	log.Println(mapsUrl)
	modifyMainPage(mapsUrl, w, r)
}

func ServeMaps(w http.ResponseWriter, r *http.Request) {
	fullURL := r.URL
	fullURL.Host = "www.google.com"
	fullURL.Scheme = "https"
	if strings.Contains(fullURL.String(), "photometa") {
		modifyInformation(fullURL.String(), w, r)
	} else {
		http.Redirect(w, r, fullURL.String(), 301)
	}
}