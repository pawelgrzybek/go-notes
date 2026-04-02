package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("failed to connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewNoteServiceClient(conn)

	list := flag.Bool("list", false, "List all notes")
	watch := flag.Bool("watch", false, "Watch notes for real-time updates")
	get := flag.String("get", "", "Get a note by ID")
	add := flag.String("add", "", "Add a new note with the given content")
	del := flag.String("delete", "", "Delete a note by ID")
	deleteAll := flag.Bool("deleteAll", false, "Delete all notes")
	flag.Parse()

	switch {
	case *list:
		doList(client)
	case *watch:
		doWatch(client)
	case *get != "":
		doGet(client, *get)
	case *add != "":
		doAdd(client, *add)
	case *del != "":
		doDelete(client, *del)
	case *deleteAll:
		doDeleteAll(client)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func doWatch(client pb.NoteServiceClient) {
	stream, err := client.WatchNotes(context.Background(), &pb.WatchNotesRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Stream ended: %v\n", err)
			os.Exit(1)
		}
		fmt.Print("\033[H\033[2J")
		printNotes(resp.GetNotes())
	}
}

func doList(client pb.NoteServiceClient) {
	notes, err := client.ListNotes(context.Background(), &pb.ListNotesRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNotes(notes.Notes)
}

func doGet(client pb.NoteServiceClient, id string) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid ID %q\n", id)
		os.Exit(1)
	}

	note, err := client.GetNote(context.Background(), &pb.GetNoteRequest{Id: new(int32(idInt))})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printNote(note.GetNote())
}

func doAdd(client pb.NoteServiceClient, content string) {
	note, err := client.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: &content})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNote(note.GetNote())
}

func doDelete(client pb.NoteServiceClient, id string) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid ID %q\n", id)
		os.Exit(1)
	}

	note, err := client.DeleteNote(context.Background(), &pb.DeleteNoteRequest{Id: new(int32(idInt))})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printNote(note.GetNote())
}

func doDeleteAll(client pb.NoteServiceClient) {
	notes, err := client.DeleteAllNotes(context.Background(), &pb.DeleteAllNotesRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printNotes(notes.GetNotes())
}

func formatNote(n *pb.Note) string {
	return fmt.Sprintf("id: %d\nnote: %s\ncreated: %s", *n.Id, *n.Note, n.GetCreatedAt().AsTime().Format("02 Jan 2006, 15:04"))
}

func printNote(n *pb.Note) {
	fmt.Println(formatNote(n))
}

func printNotes(notes []*pb.Note) {
	parts := make([]string, len(notes))
	for i, n := range notes {
		parts[i] = formatNote(n)
	}

	fmt.Println(strings.Join(parts, "\n---\n"))
}
