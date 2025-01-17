import {StockGraph} from './StockGraph.js'

// Initialize modals
document.addEventListener('DOMContentLoaded', () =>
    M.Modal.init(document.querySelectorAll('.modal'), {})
)

//display date on navbar
const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
const days = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat']
const date = new Date
new Vue({
    el: '#nav',
    data: {
        date: days[date.getDay()] + ' ' + months[date.getMonth()] + ' ' + date.getDate() + ' ' + date.getFullYear()
    }
})

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
            new StockGraph('McDonalds', 'MCD', 'VF5PZ9DBDNDKYYSN'),
            new StockGraph('Google', 'GOOG', 'ZDDTI9CUZ6SEF5IZ'),
            new StockGraph('Microsoft', 'MSFT', '9O7U8ZRHOF4WU0FN'),
            new StockGraph('Disney', 'DIS', 'N0U4CZ8XN5YCEK9U'),
            new StockGraph('Fedex', 'FDX', 'YTU231UOF5L0T7WB')
        ]
    }
})
