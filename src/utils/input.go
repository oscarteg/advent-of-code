package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ReadLines(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		number, _ := strconv.Atoi(scanner.Text())
		lines = append(lines, number)
	}
	return lines, scanner.Err()
}

func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}


func GetInput(url string) (body []byte) {
	response, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)

	fmt.Println(body)

	if err != nil {
		log.Fatalln(err)
	}

	return
}


func ParseResponse(bytes []byte) []string {

	stringArray := make([]string, len(bytes))

	for elem := range bytes {
		stringArray = append(stringArray, strconv.Itoa(elem))
	}

	return stringArray
}