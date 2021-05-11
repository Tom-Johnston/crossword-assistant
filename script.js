//worker is a web worker which loads the searcher.
var worker = new Worker('worker.js');

worker.onmessage = function (e) {
    if (e.data == "init") {
        //The web worker has just initialised and is ready for the first query.
        enableSearch()
        return
    }
    //The web worker is returning the results of a query.
    processSearchResults(e.data);
};

window.addEventListener('DOMContentLoaded', (event) => {
    document.getElementById("add-anagram").addEventListener("click", clickAddAnagram)
    document.getElementById("add-pattern").addEventListener("click", clickAddPattern)
    document.getElementById("go-button").addEventListener('click', goQuery)
    loadFromGet()
});

function clickAddAnagram() {
    addQuery("anagram", "")
}

function clickAddPattern() {
    addQuery("pattern", "")
}

function disableSearch() {
    var goButton = document.getElementById("go-button");
    goButton.disabled = true;
}

function enableSearch() {
    var goButton = document.getElementById("go-button");
    goButton.disabled = false;
}

//modifyInput makes the text uppercase and replaces a space with a ?
function modifyInput(e) {
    var end = e.target.selectionEnd;
    var curr = e.target.value;
    var res = curr.replace(/ /g, "?");
    res = res.toUpperCase();
    e.target.value = res
    e.target.setSelectionRange(end, end);
    e.preventDefault();
}

//switchType is used on the buttons next to the constraints so clicking the button changes the type.
function switchType(e) {
    var target = e.target;
    if (target.getAttribute("data-query-type") == "anagram") {
        target.setAttribute('data-query-type', "pattern")
        target.innerHTML = "pattern"
    } else if (target.getAttribute("data-query-type") == "pattern") {
        target.setAttribute('data-query-type', "anagram")
        target.innerHTML = "anagram"
    }
}

//addQuery adds a constraint of the given type and with the given value.
function addQuery(typeString, queryString) {
    var query = document.createElement("div")
    query.classList.add("query")

    var type = document.createElement("button")
    type.classList.add("query-type", typeString)
    type.setAttribute("type", "button")
    type.setAttribute("data-query-type", typeString)
    type.innerHTML = typeString
    type.addEventListener("click", switchType)
    query.appendChild(type)

    var content = document.createElement("input")
    content.value = queryString
    content.addEventListener("input", modifyInput)
    content.classList.add("query-content")
    content.setAttribute("autocomplete", "off")
    content.setAttribute("autocorrect", "off")
    content.setAttribute("autocapitalize", "off")
    content.setAttribute("spellcheck", "off")
    content.setAttribute("type", "text")
    query.appendChild(content)


    var remove = document.createElement("button")
    remove.classList.add("query-remove")
    remove.setAttribute("type", "button")
    remove.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" fill="rgba(256, 256, 256, 0.6)" viewBox="0 0 24 24"><title>Remove this constraint</title><path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM8 9h8v10H8V9zm7.5-5l-1-1h-5l-1 1H5v2h14V4z"/></svg>`;
    remove.addEventListener("click", removeQuery)
    query.appendChild(remove)

    document.getElementById("query-container").appendChild(query)
}

//removeQuery deletes the constraint.
function removeQuery(e) {
    var target = e.target;
    var query = target.closest(".query")
    query.parentElement.removeChild(query)
}

//queryStartTime will be set to the time that the query is started and is used to find the time the query takes.
var queryStartTime;

//goQuery is called when the user clicks the search button. It removes any old results, adds the search to the history and then sends the search to the web worker.
//This disables the search button so the user can't search again until the first search completes.
function goQuery() {
    //Store the start time.
    queryStartTime = performance.now();

    //Hide the about text
    var aboutText = document.getElementById("about-text")
    aboutText.classList.add("none")

    //Empty the results div and show it.
    var resultDiv = document.getElementById("results")
    while (resultDiv.lastChild) {
        resultDiv.removeChild(resultDiv.lastChild);
    }
    resultDiv.classList.remove("none")
    
    //Find the queries.
    var queriesHTML = document.getElementById("inputs").getElementsByClassName("query")
    var queries = [];
    for (var query of queriesHTML) {
        var type = query.getElementsByClassName('query-type')[0].getAttribute("data-query-type")
        var content = query.getElementsByClassName('query-content')[0].value
        content = content.replaceAll(".", "?")
        queries.push({ type: type, content: content })
    }

    //Send the queries to the web worker.
    worker.postMessage(queries)

    //Disable the search button so the user can't search until this one finishes.
    disableSearch()
    
    var statusText = document.getElementById("status-text")
    statusText.innerHTML = "Searching..."

    document.getElementById("limit-reached").style.display = "none"

    //Add the query to the history.
    if (saveToGet() != window.location.search.substr(0)) {
        history.pushState(null, "", saveToGet())
    }
}

//processSearchResults is called when the web worker returns the results and displays the results.
//This also enables the search button.
function processSearchResults(results) {
    //If the result is a single string, it is an error.
    if (typeof results == "string") {
        console.error(results)
        var t1 = performance.now()
        var statusText = document.getElementById("status-text")
        statusText.classList.add("error")
        statusText.innerHTML = `${results}. Took ${(t1 - queryStartTime).toFixed(1)} milliseconds.`
        //We can recover from the errors that we return so let the user try again.
        enableSearch()
        //Show the about text since there aren't any results to show.
        var aboutText = document.getElementById("about-text")
        aboutText.classList.remove("none")
        var resultDiv = document.getElementById("results")
        resultDiv.classList.add("none")
        return
    }

    var resultDiv = document.getElementById("results")

    //The results are returned in arrays split by the length of the words and we want to show the words in decreasing length.
    var keys = Object.entries(results).sort((a, b) => +b[0] - a[0])
    //We will count the total number of results so we know if we are likely to have missed other solutions.
    var numResults = 0
    for (const [key, value] of keys) {
        var lengthWrapper = document.createElement("div")
        lengthWrapper.classList.add("length")

        resultDiv.appendChild(lengthWrapper)

        var lengthTitle = document.createElement("h3")
        if (key == 1) {
            lengthTitle.innerHTML = "1 letter"
        } else {
            lengthTitle.innerHTML = key + " letters"
        }

        lengthWrapper.appendChild(lengthTitle)

        var resultWrapper = document.createElement("div")
        resultWrapper.classList.add("results")

        resultWrapper.style.gridTemplateColumns = "repeat(auto-fill, " + (+key + 2) + "em)"
        lengthWrapper.appendChild(resultWrapper)
        resultWrapper.style.gridTemplateColumns = "repeat(auto-fill, " + (+key + 2) + "em)"

        for (var result of value) {
            var node = document.createElement("span")
            node.classList.add("result")
            node.innerHTML = result
            resultWrapper.appendChild(node)
            numResults++
        }
    }

    //Update the status text with the time taken and the number of words found.
    var t1 = performance.now()
    var statusText = document.getElementById("status-text")
    statusText.classList.remove("error")
    if (numResults === 1) {
        statusText.innerHTML = `Found 1 word in ${(t1 - queryStartTime).toFixed(1)} milliseconds.`
    } else if (numResults >= 250) {
        statusText.innerHTML = `Found 250+ words in ${(t1 - queryStartTime).toFixed(1)} milliseconds. There may be other solutions.`
    } else {
        statusText.innerHTML = `Found ${numResults} words in ${(t1 - queryStartTime).toFixed(1)} milliseconds.`
    }

    //If we have at least 250 words, then the search may have been stopped by the LimitSearcher.
    if (numResults >= 250) {
        document.getElementById("limit-reached").style.display = "block"
    }

    //Enable the search button so the user can search again.
    enableSearch()
}

//loadFromGet sets the queries using the parameters given in the URL.
function loadFromGet() {
    //Remove all old queries
    var queryContainer = document.getElementById("query-container");
    while (queryContainer.lastChild) {
        queryContainer.removeChild(queryContainer.lastChild);
    }

    var s = window.location.search.substr(1);
    var searchParams = new URLSearchParams(s);
    for (let p of searchParams) {
        addQuery(p[0], p[1])
    }
}

//saveToGet encodes the current queries to use in the URL of a GET request.
function saveToGet() {
    var s = new URLSearchParams();
    var queriesHTML = document.getElementById("inputs").getElementsByClassName("query")
    for (var query of queriesHTML) {
        var type = query.getElementsByClassName('query-type')[0].getAttribute("data-query-type")
        var content = query.getElementsByClassName('query-content')[0].value
        s.append(type, content)
    }
    return "?" + s.toString()
}

//Load the queries from the URL on load.
window.addEventListener('DOMContentLoaded', () => {
    loadFromGet()
});

//Handle the user moving forwards/backwards in the browser.
window.addEventListener('popstate', () => {
    //Load the queries from the URL.
    loadFromGet()
    //Hide any results that are already there and show the about text.
    var aboutText = document.getElementById("about-text")
    aboutText.classList.remove("none")
    var resultDiv = document.getElementById("results")
    resultDiv.classList.add("none")
    var statusText = document.getElementById("status-text")
    statusText.innerHTML = ""
});

if ('serviceWorker' in navigator) {
    window.addEventListener('load', function () {
        navigator.serviceWorker.register('sw.js').then(function (registration) {
        }, function (err) {
            // Registration has failed
            console.log('ServiceWorker registration failed: ', err);
        });
    });
}

let deferredPrompt;

window.addEventListener('beforeinstallprompt', (e) => {
    // Stash the event so it can be triggered later.
    deferredPrompt = e;

    var installButton = document.getElementById("install-button")
    installButton.style.display = "block";
    installButton.addEventListener('click', (ev) => {
        deferredPrompt.prompt();
        installButton.style.display = 'none';
        deferredPrompt.userChoice
            .then((choiceResult) => {
                deferredPrompt = null;
            });
    });
});
