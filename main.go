package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Node struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Children []Node `yaml:"children"`
}

func main() {
	// Looks for config here:
	//  * $XDG_CONFIG_HOME/dynhome/dynhome.yaml
	//  * $HOME/.config/dynhome/dynhome.yaml
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir = os.Getenv("HOME") + "/.config"
	}
	path := ""
	if dir != "" {
		path = fmt.Sprintf("%s/dynhome/dynhome.yaml", dir)
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(501)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.Host)
		nodes := parseConfig(path)

		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, "<DOCTYPE html><html><body><ul>")
		for _, node := range *nodes {
			writeNode(w, node)
		}
		fmt.Fprintln(w, "</ul></body></html>")
	})

	http.ListenAndServe(":7777", nil)
}

func writeNode(w http.ResponseWriter, node Node) {
	if node.Url != "" {
		fmt.Fprintf(w, "<li><a href=\"%s\">%s</a>", node.Url, node.Name)
	} else {
		fmt.Fprintf(w, "<li><strong>%s</strong>", node.Name)
	}

	if len(node.Children) > 0 {
		fmt.Fprint(w, "<ul>")
		for _, child := range node.Children {
			writeNode(w, child)
		}
		fmt.Fprint(w, "</ul>")
	}
	fmt.Fprint(w, "</li>")
}

func parseConfig(path string) *[]Node {
	data, _ := os.ReadFile(path)
	nodes := []Node{}
	yaml.Unmarshal(data, &nodes)
	return &nodes
}

