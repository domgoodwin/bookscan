package lookup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/domgoodwin/bookscan/items"
	"github.com/sirupsen/logrus"
)

const (
	discogsURL = "https://api.discogs.com"
)

type discogsDataStore struct{}

func (o discogsDataStore) Name() string {
	return "discogs"
}

func discogsAPIRequest(path string) (*http.Response, error) {
	url := fmt.Sprintf("%v/%v", discogsURL, path)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Bookscan/0.1 +https://dgood.win")
	req.Header.Set("Authorization", fmt.Sprintf("Discogs token=%v", os.Getenv("DISCOGS_ACCESS_TOKEN")))
	logrus.Debugf("discogs @%v req: %v", path, req)
	rsp, err := client.Do(req)
	return rsp, err
}

func (o discogsDataStore) LookupBarcode(barcode string) (*items.Record, error) {
	rsp, err := discogsAPIRequest(fmt.Sprintf("database/search?barcode=%v", barcode))
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		logrus.Error(rsp)
		return nil, fmt.Errorf("code %v", rsp.StatusCode)
	}
	defer rsp.Body.Close()
	response := &discogsDatabaseSearchResponse{}
	err = json.NewDecoder(rsp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("not found")
	}
	// not sure if we need this
	if len(response.Results) > 1 {
		logrus.Infof("found too many records: %v", response)
	}
	logrus.Debugf("found record in rsp: %v", response)
	discogsResult := response.Results[0]
	return discogsResult.Record(barcode)
}

func (d discogsDatabaseResult) Record(barcode string) (*items.Record, error) {
	release, err := d.Release()
	if err != nil {
		return nil, err
	}
	year, err := strconv.ParseInt(d.Year, 10, 0)
	if err != nil {
		return nil, err
	}
	return &items.Record{
		Title:    release.Title,
		Artists:  release.GetArtists(),
		Barcode:  barcode,
		Year:     int(year),
		Link:     d.Link(),
		CoverURL: "",
	}, nil
}

func (d discogsDatabaseResult) Release() (*discogsRelease, error) {
	rsp, err := discogsAPIRequest(fmt.Sprintf("releases/%v", d.ID))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	release := &discogsRelease{}
	err = json.NewDecoder(rsp.Body).Decode(release)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("release: %v", release)
	return release, nil
}

func (d discogsDatabaseResult) Link() string {
	return fmt.Sprintf("https://www.discogs.com%v", d.URI)
}

func (r discogsRelease) GetArtists() []string {
	var artists []string
	for _, artist := range r.Artists {
		artists = append(artists, artist.Name)
	}
	return artists
}

type discogsDatabaseSearchResponse struct {
	Pagination discogsPagnination      `json:"pagination"`
	Results    []discogsDatabaseResult `json:"results"`
}

type discogsPagnination struct {
	PerPage int                    `json:"per_page"`
	Pages   int                    `json:"pages"`
	Page    int                    `json:"page"`
	URLs    discogsPagninationURLs `json:"urls"`
	Items   int                    `json:"items"`
}

type discogsPagninationURLs struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type discogsDatabaseResult struct {
	Style       []string `json:"style"`
	Thumb       string   `json:"thumb"`
	Title       string   `json:"title"`
	Country     string   `json:"country"`
	Format      []string `json:"format"`
	URI         string   `json:"uri"`
	CatNo       string   `json:"catno"`
	Year        string   `json:"year"`
	Genre       []string `json:"genre"`
	ResourceURL string   `json:"resource_url"`
	Type        string   `json:"type"`
	ID          int      `json:"id"`
}

type discogsRelease struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	Year        int    `json:"year"`
	ResourceURL string `json:"resource_url"`
	URI         string `json:"uri"`
	Artists     []struct {
		Name        string `json:"name"`
		Anv         string `json:"anv"`
		Join        string `json:"join"`
		Role        string `json:"role"`
		Tracks      string `json:"tracks"`
		ID          int    `json:"id"`
		ResourceURL string `json:"resource_url"`
	} `json:"artists"`
	ArtistsSort string `json:"artists_sort"`
	Labels      []struct {
		Name           string `json:"name"`
		Catno          string `json:"catno"`
		EntityType     string `json:"entity_type"`
		EntityTypeName string `json:"entity_type_name"`
		ID             int    `json:"id"`
		ResourceURL    string `json:"resource_url"`
	} `json:"labels"`
	Series []struct {
		Name           string `json:"name"`
		Catno          string `json:"catno"`
		EntityType     string `json:"entity_type"`
		EntityTypeName string `json:"entity_type_name"`
		ID             int    `json:"id"`
		ResourceURL    string `json:"resource_url"`
	} `json:"series"`
	Companies []any `json:"companies"`
	Formats   []struct {
		Name         string   `json:"name"`
		Qty          string   `json:"qty"`
		Text         string   `json:"text"`
		Descriptions []string `json:"descriptions"`
	} `json:"formats"`
	DataQuality string `json:"data_quality"`
	Community   struct {
		Have   int `json:"have"`
		Want   int `json:"want"`
		Rating struct {
			Count   int     `json:"count"`
			Average float64 `json:"average"`
		} `json:"rating"`
		Submitter struct {
			Username    string `json:"username"`
			ResourceURL string `json:"resource_url"`
		} `json:"submitter"`
		Contributors []struct {
			Username    string `json:"username"`
			ResourceURL string `json:"resource_url"`
		} `json:"contributors"`
		DataQuality string `json:"data_quality"`
		Status      string `json:"status"`
	} `json:"community"`
	FormatQuantity    int      `json:"format_quantity"`
	DateAdded         string   `json:"date_added"`
	DateChanged       string   `json:"date_changed"`
	NumForSale        int      `json:"num_for_sale"`
	LowestPrice       any      `json:"lowest_price"`
	Title             string   `json:"title"`
	Country           string   `json:"country"`
	Released          string   `json:"released"`
	Notes             string   `json:"notes"`
	ReleasedFormatted string   `json:"released_formatted"`
	Identifiers       []any    `json:"identifiers"`
	Genres            []string `json:"genres"`
	Styles            []string `json:"styles"`
	Tracklist         []struct {
		Position string `json:"position"`
		Type     string `json:"type_"`
		Title    string `json:"title"`
		Duration string `json:"duration"`
	} `json:"tracklist"`
	Extraartists []any `json:"extraartists"`
	Images       []struct {
		Type        string `json:"type"`
		URI         string `json:"uri"`
		ResourceURL string `json:"resource_url"`
		URI150      string `json:"uri150"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	} `json:"images"`
	Thumb           string `json:"thumb"`
	EstimatedWeight int    `json:"estimated_weight"`
	BlockedFromSale bool   `json:"blocked_from_sale"`
}
