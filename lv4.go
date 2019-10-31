package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, _ := os.Create("proverb.txt")
	defer f.Close()
	f.WriteString("don't communicate by sharing memory share memory by communicating")
	inputFile, _ := os.Open("proverb.txt")
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	inputString, _ := inputReader.ReadString('\n')
	fmt.Printf(" %s", inputString)

}
