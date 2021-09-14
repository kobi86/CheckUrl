package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ConfigurationFile struct {
	Httpslist  string   `json:"httpslist"`
	Outputfile string   `json:"outputfile"`
	Search     string   `json:"search"`
	ErrorCodes []string `json:"errorCodes"`
}

func main() {
	fromConfig := LoadConfiguration("urlCheckConfig.json")

	Httpslist := fromConfig.Httpslist
	Outputfile := fromConfig.Outputfile
	Search := fromConfig.Search
	ErrorCodes := fromConfig.ErrorCodes

	DeleteOutputFile(Outputfile)
	ListAllURIs(Httpslist, Search, Outputfile)
	fmt.Print("Starting to check URI's\n")
	RunUrlTest(Outputfile, ErrorCodes)
}

//Read json Config File
func LoadConfiguration(file string) ConfigurationFile {
	var config ConfigurationFile
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

//Checkes if the Message has error type that indicats that the site doesn't exists
func checkMessage(message string, errorCodes []string, site string) {
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	for _, errorCodes := range errorCodes {
		if strings.Contains(message, errorCodes) {
			fmt.Println(string(colorRed), "****** Not Exists ****** "+site, message, colorGreen)
			return
		}
	}
	fmt.Println("Exists "+site, message)
}

//Openes the URL list file that was created and sends each line to CheckMessage function
func RunUrlTest(URIList string, errorCodes []string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	data, err := os.Open(URIList)
	if err != nil {
		fmt.Print("URL file list cant be found")
		panic(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		site := scanner.Text()
		resp, err := http.Get(site)
		if err != nil {
			checkMessage(err.Error(), errorCodes, site)
		}
		if resp != nil {
			checkMessage(resp.Status, errorCodes, site)
		}
	}
}

//Deletes the OutPutFile in case it exists, Done only once
func DeleteOutputFile(fileToDelete string) {
	if _, err := os.Stat(fileToDelete); os.IsNotExist(err) {
		return
	}
	err := os.Remove(fileToDelete)
	if err != nil {
		panic(err)
	}
	fmt.Print("****** Old copy of Output file deleted ******\n")
}

//Add URI string from the Bookmarks.yaml list to output file
func UpdateOutputFile(lineToAdd string, outputfile string) {
	output, err := os.OpenFile(outputfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic((err))
	}
	defer output.Close()
	if _, err := output.WriteString(lineToAdd + "\n"); err != nil {
		panic(err)
	}
}

//Taks the list of URI's and runs line by line to search for uri, in case found send it to be written
func ListAllURIs(URIList string, search string, outputfile string) {
	data, err := os.Open(URIList)
	if err != nil {
		fmt.Print("Cant find config json file")
		panic(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		answer := strings.Contains(line, search)
		if answer == true {
			CleanLine := strings.Replace(line, "-", "", 1)
			FixedLine := strings.Replace(CleanLine, " ", "", 10)
			UpdateOutputFile(FixedLine, outputfile)
		}
	}
}
