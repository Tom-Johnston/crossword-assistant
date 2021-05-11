![Crossword Assistant Logo](images/icon-circle.svg)

# Crossword Assistant

Crossword Assistant is a PWA to solve anagrams and find words with missing letters. 

## Features

- **Powerful:** There are 357 anagrams of the letters `OSSRD????` and 200 words of the form `??O?????D`, but only 12 words which are both. With this solver you can check for anagrams with checking letters, drastically reducing the number of options.
- **Cross-platform:** Works on any platform with a modern browser including Windows, Linux, Android and more. 
- **Offline:** Once the page has been loaded once, the solver is cached and works offline. 
- **Completely free:** No one should have to pay for something as simple as an anagram solver. This one is completely free to use with no ads or tracking!

## How it works

This app wouldn't be very useful if it didn't have a list of words to search, and we gratefully make use of the word lists provided by [SCOWL](http://wordlist.aspell.net/). The accents and punctuation are removed from each word and the normalised forms are converted into a series of [directed acyclic word graphs](https://en.wikipedia.org/wiki/) in a file `words.go` by running [compiler/compiler.go](compiler/compiler.go). Information about accents and punctuation is stored alongside the directed acyclic word graphs, but only for the necessary words. The code in [crossword-assistant.go](crossword-assistant.go) uses a depth first search to search the DAWGs created in the previous step . This is compiled (with `words.go`) to Web Assembly which is then run in a worker thread in [worker.js ](worker.js). The UI is fairly basic HTML, CSS and Javascript which is cached by the service worker in [sw.js](sw.js) to work offline.


## Building and serving the PWA

The website is built using the Github pages gem and can be built using `bundle exec jekyll build`. Unfortunately, the `crossword-assistant.wasm` file is served with the wrong MIME type when served by `bundle exec jekyll serve` and an alternative must be used. For example, the following snippet will move into the `_site` folder and serve the contents on `localhost:8080`. 

```
cd _site
goexec 'http.ListenAndServe(\":8080\", http.FileServer(http.Dir(\".\")))'
```

### Building the WASM solver

There are two steps to building the WASM solver. First, we must process the list of words and create the file `words.go` which contains the gob encoded DAWGs etc. This is done by piping in (UTF-8 encoded) words to the program generated from [compiler/compiler.go](compiler/compiler.go). If you have problems, with accents on the words not being handled correctly, check that the shell you are using is set to use UTF-8 encoding.

```powershell
# Don't forget that this first step needs to be run with the correct environment variables for your system and not for WASM. e.g.
# $env:GOOS = "windows"
# $env:GOARCH = "amd64"
cat word-lists/* | go run compiler/compiler.go

$env:GOOS = "js"
$env:GOARCH = "wasm"
go build -o crossword-assistant.wasm crossword-assistant.go words.go
# words.go is quite a large file with long lines which can cause problems with tools like go-pls so let's remove it. 
rm words.go
```

## Roadmap

- Store the DAWGs etc. as separate files which are passed into the WASM module.
- Add more conditions on types such as subwords, superwords and contains.
- Add a convenient way to look up definitions even if it is just a button to do a Google search.

### iOS/Safari support

As far as I know there is no reason this PWA can't be made to work on iOS, but I believe it currently doesn't and I can't debug this without an iPhone and a Mac.