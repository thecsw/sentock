/**
 * All the information to create a graph
 */
export class StockGraph {

    // Member variables
    id
    name
    numDays = 100

    /**
     * Creates a stock graph
     * @param {string} name Company name
     * @param {string} symbol Stock ticker symbol
     */
    constructor(name, id) {
        this.id = id
        this.name = name

        // Make API call
        let stocks = this.apiCall(id).then(this.json2stocks)

        // Scale sentiment data and draw
        stocks.then((stock => {
            const stockNumbers = stock.slice(-this.numDays).map(n => parseInt(n))
            const maxStock = stockNumbers.reduce((a, b) => a > b ? a : b)
            const minStock = stockNumbers.reduce((a, b) => a < b ? a : b)

            let sentiment = new Array(this.numDays).fill(0).map(()=>Math.random()*2-1)
            sentiment = sentiment.map(n => n * (maxStock - minStock) / 2 + (maxStock + minStock) / 2)

            this.draw(stock, sentiment)
        }).bind(this))
    }

    /**
     * Makes an API call, and returns the json from it
     * @param {string} symbol Stock trading symbol
     */
    apiCall(symbol) {
        const apiKey = 'VF5PZ9DBDNDKYYSN'
        const api = 'https://www.alphavantage.co/query?'
        const url = `${api}function=TIME_SERIES_DAILY&symbol=${symbol}&apikey=${apiKey}`

        return fetch(url).then(res => res.json())
    }

    /**
     * Draws the graphs in their respective context
     * @param {[number]} stock
     * @param {[number]} sentiment
     */
    draw(stock, sentiment) {
        //generate colors
        const color = {
            r: this.random(50, 150),
            g: this.random(50, 150),
            b: this.random(50, 150)
        }

        //construct chart data
        const chart = {
            type: 'line',
            data: {
                labels: this.getLabels(this.numDays),
                datasets: [
                    {
                        label: 'Stock Price',
                        data: stock.slice(-this.numDays),
                        fill: false,
                        borderColor: `rgb(${color.r}, ${color.g}, ${color.b})`,
                        lineTension: 0.3,
                        pointRadius: 0,
                        fillColor: `rgb(${color.r}, ${color.g}, ${color.b})`
                    },
                    {
                        label: 'Sentiment Rating',
                        data: sentiment.slice(-this.numDays),
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

    /**
     * Generates the list of labels for the last n days
     * @param {number} n number of days
     */
    getLabels(n) {
        let date = new Date
        date.setDate(date.getDate() - n)

        return [...Array(n).keys()].map(n => {
            date.setDate(date.getDate() + 1)
            return `${date.getMonth()}/${date.getDate()}/${date.getFullYear()}`
        })
    }

    /**
     * Converts the JSON object from the API into a list of stock data
     * @param {object} json The JSON object to be converted
     */
    json2stocks(json) {
        const daysKey = 'Time Series (Daily)'
        const dataKey = '4. close'

        const days = json[daysKey]
        let data = []

        for(let day in days) data.push(days[day][dataKey])

        return data
    }

    /**
     * Generate a random number from min to max
     * @param {number} min minimum possible number
     * @param {number} max maximum possible number
     */
    random(min, max) {
        return Math.random() * (max - min) + min
    }
}
