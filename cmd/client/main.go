package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseURL = "http://localhost:8080"

func main() {
	list := flag.Bool("list", false, "List all notes")
	get := flag.String("get", "", "Get a note by ID")
	add := flag.String("add", "", "Add a new note with the given content")
	del := flag.String("delete", "", "Delete a note by ID")
	deleteAll := flag.Bool("deleteAll", false, "Delete all notes")
	flag.Parse()

	switch {
	case *list:
		doList()
	case *get != "":
		doGet(*get)
	case *add != "":
		doAdd(*add)
	case *del != "":
		doDelete(*del)
	case *deleteAll:
		doDeleteAll()
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func doList() {
	resp, err := http.Get(baseURL + "/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: server returned %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var notes []models.Note
	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNotes(notes)
}

func doGet(id string) {
	resp, err := http.Get(baseURL + "/" + id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: server returned %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var note models.Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNote(note)
}

func doAdd(content string) {
	payload, err := json.Marshal(map[string]string{"note": content})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(baseURL+"/", "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		fmt.Fprintf(os.Stderr, "Error: server returned %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var note models.Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNote(note)
}

func doDelete(id string) {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/"+id, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: server returned %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var note models.Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNote(note)
}

func doDeleteAll() {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: server returned %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var notes []models.Note
	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNotes(notes)
}

func formatTime(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.Format("02 Jan 2006, 15:04")
}

func formatNote(n models.Note) string {
	return fmt.Sprintf("id: %d\nnote: %s\ncreated: %s", n.ID, n.Note, formatTime(n.CreatedAt))
}

func printNote(n models.Note) {
	fmt.Println(formatNote(n))
}

func printNotes(notes []models.Note) {
	parts := make([]string, len(notes))
	for i, n := range notes {
		parts[i] = formatNote(n)
	}

	fmt.Println(strings.Join(parts, "\n---\n"))
}
