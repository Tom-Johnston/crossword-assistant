//Import the bridging code between Go and Javascript.
importScripts('wasm_exec.js');

const go = new Go();
WebAssembly.instantiateStreaming(fetch("crossword-assistant.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    //We send back the message init once we have finished initialising.
    postMessage("init")
});

onmessage = function(e) {
    //search is defined in the WASM module and returns the results of the query which we send back.
    var results = search(e.data);
    postMessage(results)
}
