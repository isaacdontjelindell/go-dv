package main

import (
    "bufio"
    "os"
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
    Name string
    Route string
    TotalCost int // int cost + (cost to route)
}

type Update struct {
    RoutingTable []Node
    From string
}

