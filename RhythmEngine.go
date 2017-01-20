package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"io"
	"strings"
	"sort"
)

type VnDictionary struct {
	// Size of the dictionary, which indicates the number of words.
	Size int
	// The list contains Vietnamese words.
	WordList []string
	// An index of all the words in the dictionary.
	Indexing map[string][]int
}

// The find method of VnDictionary looks for the entry that has
// the matched word.
// Input: the keyword. Output: list of matched entries.
// To avoid copying the VnDictionary object, we would better use pointer.
func (dict *VnDictionary) FindAll(keyword string) []string{
	// Get the list of entries for the given keyword.
	entriesList, found := dict.Indexing[keyword]
	// Check the existence of the keyword in the index.
	if found {
		sort.Ints(entriesList)
		entriesList = duplicatesEliminate(entriesList)
	}
	// Retrieves entries from the list of words.
	results := make([]string, 0)
	for _,v := range entriesList {
		results = append(results, dict.WordList[v])
	}
	return results
}

func duplicatesEliminate (list []int) []int {
	newList := make([]int, 0)
	currentValue := -1
	for _, value := range list {
		if currentValue != value {
			newList = append(newList, value)
			currentValue = value
		}
	}
	return newList
}


func main() {
	vnDict, _ := buildListAndIndex("./vn_dict.txt")
	writeFile("index.txt", vnDict.Indexing)

	// Keyword to search
	keyword := "Mưa"

	keyword = strings.ToLower(keyword)
	results := vnDict.FindAll(keyword)
	fmt.Printf("\nWords that rhythms with '%v' are:", keyword)
	for _, entry := range results{
		fmt.Printf("\n%v", entry)
	}
	fmt.Println("\n=== END ===")
}

func writeFile(fileName string, content interface{}) {
	file, _ := os.Create("./" + fileName)
	defer file.Close()
	io.WriteString(file, fmt.Sprintf("%v", content))
}

func buildListAndIndex(sourceFile string) (VnDictionary, error){
	fileName := sourceFile

	// Read the Vietnamese dictionary.
	content, err := ioutil.ReadFile(fileName)

	var vnDict VnDictionary
	vnDict.Indexing = make(map[string][]int)

	// Check for reading file error.
	if err == nil {
		// If it does not have any error, we print out the content.
		vn_dict := string(content)
		cutOff := 0
		entryIndex := 0
		// Read through content that we just parsed from the file.
		for i:=0; i < len(vn_dict); i++ {
			// Read single entry from the dictionary.
			if vn_dict[i] == '\n' {
				entry := string(vn_dict[cutOff:i])
				// Append the word into the word list.
				vnDict.WordList = append(vnDict.WordList, entry)
				// Index the word.
				// Get all the words within this entry.
				wordCutOff := 0
				for j := 0; j < len(entry); j++ {
					// If we have parsed a single word, add it into the index.
					if entry[j] == ' ' {
						// Read a single word from the entry.
						singleWord := strings.ToLower(entry[wordCutOff:j])
						// Append the position of the word into the index.
						vnDict.Indexing[singleWord] = append(vnDict.Indexing[singleWord], entryIndex)
						// Increase the next word cut off.
						wordCutOff = j + 1
					}
				}
				// Increase entry index.
				entryIndex++
				cutOff = i + 1
			}
		}
	}
	return vnDict, err
}