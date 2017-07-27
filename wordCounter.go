// word-counter counts the words in a text file.
package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "io/ioutil"
    "os"
)

var wordRegExp = regexp.MustCompile(`\pL+('\pL+)*`)
var wordCnts = make(map[string]int)

func main() {
    lookup := flag.String("lookup", "", "Existing word count file to use for lookup")
    flag.Parse()
    if len(flag.Args()) < 1 {
        log.Fatal("Missing filename argument")
    }
    file := flag.Arg(0)

    // Choose program function.
    if *lookup != "" {
        lookUpCnts(file, *lookup)
    } else {
        cntWords(file)
    }

    f,err := os.Create("result.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)

    // Print the word counts.
    for k,v := range wordCnts {
        fmt.Printf("%s,%d\n", k, v)
        fmt.Fprintf(w, "%s,%d\n", k, v)
        w.Flush()
    }
}

// cntWords counts the words in a file.
func cntWords(file string) {
    hdl, err := http.Get(file)
    if err != nil {
        log.Fatal(err)
    }
    defer hdl.Body.Close()
    body, err := ioutil.ReadAll(hdl.Body)
    if err != nil {
        log.Fatal(err)
    }
    responseString :=  string(body)
    line := strings.ToLower(responseString)
    words := wordRegExp.FindAllString(line, -1)
    for _, word := range words {
        wordCnts[word]++
    }
}

// lookUpCnts looks up the word counts for a file in an existing word count file.
func lookUpCnts(file string, lookup string) {
    hdl, err := http.Get(file)
    if err != nil {
        log.Fatal(err)
    }
    defer hdl.Body.Close()
    body, err := ioutil.ReadAll(hdl.Body)
    if err != nil {
        log.Fatal(err)
    }

    responseStr := string(body)
    wordsToLookUp := make(map[string]bool)
    line := strings.ToLower(responseStr)
    words := wordRegExp.FindAllString(line, -1)
    for _, word := range words {
        wordsToLookUp[word] = true
    }

    lookupHdl, err := os.Open(lookup)
    if err != nil {
        log.Fatal(err)
    }
    defer lookupHdl.Close()
    scanner := bufio.NewScanner(lookupHdl)
    for scanner.Scan() {
        if err := scanner.Err(); err != nil {
            log.Fatal(err)
        }
        line := strings.ToLower(scanner.Text())
        fields := strings.Split(line, ",")
        if len(fields) != 2 {
            continue
        }
        word := fields[0]
        if wordsToLookUp[word] {
            cnt, err := strconv.Atoi(fields[1])
            if err == nil {
                wordCnts[word] = cnt
            }
        }
    }
}
