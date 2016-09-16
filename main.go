package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BenchR267/goalfred"
)

func main() {
	queryTerms, _ := goalfred.NormalizedArguments()
	defer goalfred.Print()

	pods := getPods(strings.Join(queryTerms, "%20"))

	for _, pod := range pods {
		goalfred.Add(pod)
	}
}

type pod struct {
	Link      string
	Name      string `json:"id"`
	Summary   string
	Version   string
	Platforms []string
}

func (p pod) Documentation() string {
	return fmt.Sprintf("http://cocoadocs.org/docsets/%s/%s", p.Name, p.Version)
}

func (p pod) Item() goalfred.Item {
	title := fmt.Sprintf("%s (%s)", p.Name, p.Version)
	instruction := fmt.Sprintf("pod '%s', '%s'", p.Name, p.Version)
	i := goalfred.Item{
		Title:    title,
		Subtitle: p.Summary,
		Arg:      p.Link,
	}
	i.Mod.Cmd = &goalfred.ModContent{
		Arg:      p.Documentation(),
		Subtitle: "Open documentation!",
	}
	i.Mod.Alt = &goalfred.ModContent{
		Arg:      instruction,
		Subtitle: "Copy pod install instructions",
	}

	return i
}

type response struct {
	Allocations [][]json.RawMessage `json:"allocations"`
}

func getPods(searchTerm string) []pod {

	url := fmt.Sprintf("https://search.cocoapods.org/api/v1/pods.picky.hash.json?query=%v&ids=20&offset=0&sort=quality", searchTerm)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var res response
	json.Unmarshal(body, &res)

	rawJSON := res.Allocations[0][5]

	var pods []pod
	json.Unmarshal(rawJSON, &pods)

	return pods
}
