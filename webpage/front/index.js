import {StockGraph} from './StockGraph.js'

document.addEventListener('DOMContentLoaded', function() {
    M.Modal.init(document.querySelectorAll('.modal'), {})
})

let app = new Vue({
    el: '#app',
    data: {
        graphs: [
            new StockGraph('Company Name', 'id1',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id2',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id3',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id4',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id5',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id6',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id7',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id8',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id9',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            ),
            new StockGraph('Company Name', 'id10',
                new Array(12).fill(0).map(()=>Math.random()*100),
                new Array(12).fill(0).map(()=>Math.random()*100)
            )
        ]
    }
})

//draw all graphs on canvas
for(let graph of app.graphs) graph.draw()
