export class StockGraph {

    html
    id
    name
    positivityData
    stockData

    constructor(name, id, stockData, positivityData) {
        this.id = id
        this.name = name
        this.positivityData = positivityData
        this.stockData = stockData

        //generate HTML
        this.html = `
        <div class="col s12 m6 xl4">
            <div class="hoverable card grey darken-4">
                <a class="modal-trigger white-text" href="#${id}-modal">
                    <div class="card-image">
                        <canvas id="${id}-preview"></canvas>
                    </div>
                    <div class="card-content">
                        <span class="card-title">${name}</span>
                    </div>
                </a>
            </div>
            <div id="${id}-modal" class="modal grey darken-4">
                <div class="modal-content white-text">
                    <h4>${name}</h4>
                    <canvas id="${id}-full"></canvas>
                </div>
            </div>
        </div>
        `
    }

    draw() {
        //generate colors
        const color = {
            r: Math.random()*100+50,
            g: Math.random()*100+50,
            b: Math.random()*100+50
        }

        //construct chart data
        const chart = {
            type: 'line',
            data: {
                labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
                datasets: [
                    {
                        label: 'Stock Price',
                        data: this.stockData,
                        fill: false,
                        borderColor: `rgb(${color.r}, ${color.g}, ${color.b})`,
                        lineTension: 0.2
                    },
                    {
                        label: 'Positivity Rating',
                        data: this.positivityData,
                        fill: false,
                        borderColor: `rgb(${color.r+100}, ${color.g+100}, ${color.b+100})`,
                        lineTension: 0.2
                    }
                ]
            },
            options: {}
        }

        //draw on canvas
        new Chart(document.getElementById(`${this.id}-preview`), chart)
        new Chart(document.getElementById(`${this.id}-full`), chart)
    }
}
