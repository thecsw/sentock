// document.addEventListener("DOMContentLoaded",() => {
//     let mygraph = document.querySelector("#mygraph");
//     let req = fetch("http://167.172.114.123:5050/api/graphs?h=12");
//     req.then((data) => data.blob()).then(function (image) { mygraph.data = image; });
// });
import {StockGraph} from './StockGraph.js'

// Initialize modals
document.addEventListener('DOMContentLoaded',
    () => M.Modal.init(document.querySelectorAll('.modal'), {})
)

// Set up stock graphs
new Vue({
    el: '#graphs',
    methods: {
        getError: function() {
            if(this.graphs.some(g => g.error))
                return '<h1 class="red-text">The API can only be used 5 times per minute, please wait a minute and try again.</h1>'
        }
    },
    data: {
        graphs: [
            new StockGraph('McDonalds', 'MCD', 'VF5PZ9DBDNDKYYSN')
            // new StockGraph('Chipotle', 'CMG', 'ZDDTI9CUZ6SEF5IZ'),
            // new StockGraph('Microsoft', 'MSFT', '9O7U8ZRHOF4WU0FN'),
            // new StockGraph('Disney', 'DIS', 'N0U4CZ8XN5YCEK9U'),
            // new StockGraph('FedEx', 'FDX', 'YTU231UOF5L0T7WB')
        ]
    }
})
