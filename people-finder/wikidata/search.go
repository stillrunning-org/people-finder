package wikidata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Person struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	BirthDate    string `json:"birthDate"`
	DeathDate    string `json:"deathDate"`
	Pic          string `json:"pic"`
	SiteLinksCnt int    `json:"siteLinksCnt"`
	Age          int    `json:"age"`
}

func FetchWikidataDeaths(date1 string, date2 string) ([]Person, error) {
	query := fmt.Sprintf(`SELECT ?item ?itemLabel ?birthDate ?deathDate ?pic ?sitelink
WHERE {
   ?item wdt:P570 ?deathDate .
   FILTER(?deathDate >= "%s"^^xsd:dateTime && ?deathDate < "%s"^^xsd:dateTime)
   OPTIONAL { ?item wdt:P569 ?birthDate. }
   OPTIONAL { ?item wdt:P18 ?pic. }
   OPTIONAL { ?item wikibase:sitelinks ?sitelink. }
   SERVICE wikibase:label {
     bd:serviceParam wikibase:language "en" .
     ?item rdfs:label ?itemLabel .
   }
} ORDER BY DESC(?sitelink)`, date1, date2)

	endpoint := "https://query.wikidata.org/sparql"
	encodedQuery := url.QueryEscape(query)
	req, err := http.NewRequest("GET", endpoint+"?query="+encodedQuery, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/sparql-results+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Results struct {
			Bindings []map[string]struct {
				Value string `json:"value"`
			} `json:"bindings"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	people := make([]Person, 0)

	cnt := 0
	for _, binding := range data.Results.Bindings {
		cnt++
		person := Person{
			Id:        binding["item"].Value,
			Name:      binding["itemLabel"].Value,
			BirthDate: binding["birthDate"].Value,
			DeathDate: binding["deathDate"].Value,
			Pic:       binding["pic"].Value,
		}
		if siteLinks, ok := binding["sitelink"]; ok {
			person.SiteLinksCnt, err = strconv.Atoi(siteLinks.Value)
			if err != nil {
				return nil, err
			}
		}
		if person.Pic != "" && person.SiteLinksCnt >= 1 {
			people = append(people, person)
		}
	}

	return people, nil
}
