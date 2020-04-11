document.addEventListener("DOMContentLoaded",() => {
    let mygraph = document.querySelector("#mygraph");
    let req = fetch("http://167.172.114.123:5050/api/graphs?h=12");
    req.then((data) => data.blob()).then(function (image) { mygraph.data = image; });
});