import {StockGraph} from './StockGraph.js'

// Initialize modals
document.addEventListener('DOMContentLoaded',
    () => M.Modal.init(document.querySelectorAll('.modal'), {})
)

// Set up stock graphs
new Vue({
    el: '#graphs',
    data: {
        graphs: [
            new StockGraph('McDonald\'s', 'MCD')
            // new StockGraph('Chipotle', 'CMG'),
            // new StockGraph('Microsoft', 'MSFT'),
            // new StockGraph('Disney', 'DIS'),
            // new StockGraph('FedEx', 'FDX')
        ]
    }
})
