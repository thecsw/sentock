import {StockGraph} from './StockGraph.js'

document.addEventListener('DOMContentLoaded', function() {
    M.Modal.init(document.querySelectorAll('.modal'), {})
})

let app = new Vue({
    el: '#app',
    data: {}
})

new StockGraph('example',  [
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100
],  [
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100,
    Math.random()*100
])
