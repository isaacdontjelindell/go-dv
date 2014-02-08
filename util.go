package main

import (
	"bufio"
	"fmt"
	"os"
    "net"
    "time"
    "encoding/json"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type Node struct {
	Name      string
	Route     string
	TotalCost int // int cost + (cost to route)
}

func (n Node) String() string {
	return fmt.Sprintf("Name: %s, Route: %s, Cost: %d", n.Name, n.Route, n.TotalCost)
}

type Update struct {
	RoutingTable []Node
	From         string
}

func (u Update) String() string {
	var s string
	s = fmt.Sprintf("From: %s\nRouting Table:\n", u.From)
	for _, node := range u.RoutingTable {
		s += fmt.Sprint("  ")
		s += fmt.Sprintln(node)
	}
	return s
}


func testClient() {
	conn, err := net.Dial("udp", "127.0.0.1:1337")
	if err != nil {
		fmt.Println("[testClient] Error dialing connection.", err.Error())
	}

	for {
		time.Sleep(time.Second * 2)

		// build a test update struct
		testRoutingTable := []Node{
			Node{"t1", "yoda", 3},
			Node{"t2", "yoda", 5},
			Node{"t3", "yoda", 7},
		}
		update := Update{testRoutingTable, "yoda"}

		u, err := json.Marshal(update) // u is []byte
		if err != nil {
			fmt.Println("[testClient] error marshaling update to JSON", err.Error())
		}

		conn.Write(u)
	}
}
