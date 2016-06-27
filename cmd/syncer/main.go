package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/andrewhsu/tiafm"
)

func download(filename string, url string) (written int64, err error) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	written, err = io.Copy(f, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return written, err
}

func main() {

	var local = flag.Bool("local", true, "use local file")
	flag.Parse()

	if len(os.Args) != 2 {
		log.Fatal("usage: syncer <base_url>")
	}

	base_url := os.Args[1]

	var err error
	var events []byte
	if *local {
		events, err = ioutil.ReadFile("tree.json")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		res, err := http.Get(base_url + "/results/tree.json")
		if err != nil {
			log.Fatal(err)
		}
		events, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	var rt tiafm.Results_Tree
	err = json.Unmarshal(events, &rt)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range rt.Seasons {
		fmt.Printf("%d %s\n", s.Year, s.Id)
		for _, e := range s.Events {
			a := strings.TrimSuffix(e.Code, "/afm")
			c := url.QueryEscape(e.Code)
			l := map[string]string{
				"round_event-%s.json": base_url + "/round/round_event.json?event=%s",
				"races-%s.json":       base_url + "/race/races.json?event=%s",
				"details-%s.json":     base_url + "/event/details.json?event=%s",
			}

			for f, u := range l {
				filename := fmt.Sprintf(f, a)
				url := fmt.Sprintf(u, c)

				if _, err := os.Stat(filename); err == nil {
					log.Printf("skipping %s %s\n", filename, url)
				} else {
					log.Printf("download %s %s", filename, url)
					written, err := download(filename, url)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf(" %d bytes\n", written)
				}
			}

			races, err := ioutil.ReadFile("races-" + a + ".json")
			if err != nil {
				log.Fatal(err)
			}

			var rr tiafm.Race_Races
			err = json.Unmarshal(races, &rr)
			if err != nil {
				log.Fatal(err)
			}

			for _, d := range rr.Days {
				for _, ses := range d.Sessions {
					for _, sr := range ses.Races {
						filename := fmt.Sprintf("results-%s.json", sr.Raid)
						url := fmt.Sprintf("%s/results/results.json?raid=%s", base_url, sr.Raid)

						if _, err := os.Stat(filename); err == nil {
							log.Printf("skipping %s %s\n", filename, url)
						} else {
							log.Printf("download %s %s", filename, url)
							written, err := download(filename, url)
							if err != nil {
								log.Fatal(err)
							}
							log.Printf(" %d bytes\n", written)
						}

						filename = fmt.Sprintf("laptimes-%s.json", sr.Raid)
						url = fmt.Sprintf("%s/laptimes/laptimes.json?raid=%s", base_url, sr.Raid)

						if _, err := os.Stat(filename); err == nil {
							log.Printf("skipping %s %s\n", filename, url)
						} else {
							log.Printf("download %s %s", filename, url)
							written, err := download(filename, url)
							if err != nil {
								log.Fatal(err)
							}
							log.Printf(" %d bytes\n", written)
						}
					}
				}
			}
		}
	}
}
