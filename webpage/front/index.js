import {StockGraph} from './StockGraph.js'

document.addEventListener('DOMContentLoaded', function() {
    M.Modal.init(document.querySelectorAll('.modal'), {})
})

let app = new Vue({
    el: '#app',
    data: {
        graphs: [
            new StockGraph('McDonald\'s', 'MCD'),
            new StockGraph('Chipotle', 'CMG'),
            new StockGraph('Microsoft', 'MSFT'),
            new StockGraph('Disney', 'DIS'),
            new StockGraph('FedEx', 'FDX')
        ]
    }
})
