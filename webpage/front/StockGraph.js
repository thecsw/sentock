/**
 * All the information to create a graph
 */
export class StockGraph {

    // Member variables
    id
    name
    sentimentData
    stockData

    /**
     * Creates a stock graph
     * @param {string} name Company name
     * @param {string} symbol Stock ticker symbol
     */
    constructor(name, id) {
        this.id = id
        this.name = name

        this.sentimentData = new Array(100).fill(0).map(()=>Math.random()*2-1)

        //load in data
        fetch(`https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=${id}&apikey=VF5PZ9DBDNDKYYSN`)
        .then(res => res.json())
        .then(function(json) {
            let data = []
            for(let day in json["Time Series (Daily)"]) {
                data.push(json["Time Series (Daily)"][day]["4. close"])
            }
            return data
        })
        .then(function(data) {
            this.stockData = data

            let stockNumbers = this.stockData.map(n => parseInt(n))
            let maxStock = stockNumbers.reduce((a, b) => a > b ? a : b)
            let minStock = stockNumbers.reduce((a, b) => a < b ? a : b)

            this.sentimentData = this.sentimentData.map(n => n * (maxStock - minStock) / 2 + (maxStock + minStock) / 2)

            this.draw()
        }.bind(this))
    }

    draw() {
        //generate colors
        const color = {
            r: Math.random()*100+50,
            g: Math.random()*100+50,
            b: Math.random()*100+50
        }

        //get current date
        let date = new Date
        date.setDate(date.getDate() - 99)

        //construct chart data
        const chart = {
            type: 'line',
            data: {
                labels: [...Array(100).keys()].map(function(n) {
                    date.setDate(date.getDate() + 1)
                    return `${date.getMonth()}/${date.getDate()}/${date.getFullYear()}`
                }),
                datasets: [
                    {
                        label: 'Stock Price',
                        data: this.stockData,
                        fill: false,
                        borderColor: `rgb(${color.r}, ${color.g}, ${color.b})`,
                        lineTension: 0.3,
                        pointRadius: 0,
                        fillColor: `rgb(${color.r}, ${color.g}, ${color.b})`
                    },
                    {
                        label: 'Sentiment Rating',
                        data: this.sentimentData,
                        fill: false,
                        borderColor: `rgb(${color.r+100}, ${color.g+100}, ${color.b+100})`,
                        lineTension: 0.3,
                        pointRadius: 0
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
