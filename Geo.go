package linkedin

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"strings"
)

type BatchGetGeoResponse struct {
	Statuses json.RawMessage `json:"statuses"`
	Errors   json.RawMessage `json:"errors"`
	Results  map[string]Geo  `json:"results"`
}

type Geo struct {
	DefaultLocalizedName struct {
		Locale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"locale"`
		Value string `json:"value"`
	} `json:"defaultLocalizedName"`
	Id int `json:"id"`
}

func (service *Service) BatchGetGeo(ids []string) (map[string]Geo, *errortools.Error) {
	var batchSize = 100
	var geos = make(map[string]Geo)

	for {
		var ids_ = ids
		if len(ids_) > batchSize {
			ids_ = ids_[:batchSize]
		}

		var batchGetGeoResponse BatchGetGeoResponse

		var header = http.Header{}
		header.Set("X-Restli-Protocol-Version", "2.0.0")

		url := service.urlV2(fmt.Sprintf("geo?ids=List(%s)", strings.Join(ids_, ",")))
		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               url,
			ResponseModel:     &batchGetGeoResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		for id, geo := range batchGetGeoResponse.Results {
			geos[id] = geo
		}

		if len(ids) <= batchSize {
			break
		}

		ids = ids[batchSize:]
	}

	return geos, nil
}
