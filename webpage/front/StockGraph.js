const NUM_DAYS = 30

/**
 * All the information to create a graph
 */
export class StockGraph {

    error
    id
    name
    /**
     * Creates a stock graph
     * @param {string} name Company name
     * @param {string} symbol Stock ticker symbol
     */
    constructor(name, symbol, apiKey) {
        this.error = false
        this.id = symbol
        this.name = name

        // Make API call
        let stocks = this.apiCall(symbol, apiKey).then(this.json2stocks)

        // Scale sentiment data and draw
        stocks.then((stock => {

            //warn about API key
            if(stock.length === 0) this.error = true

            const stockNumbers = stock.slice(-NUM_DAYS).map(n => parseFloat(n))
            const maxStock = stockNumbers.reduce((a, b) => a > b ? a : b)
            const minStock = stockNumbers.reduce((a, b) => a < b ? a : b)
            const midStock = (maxStock + minStock) / 2
            const stockDiff = maxStock - minStock

            let sentiment = new Array(NUM_DAYS).fill(0).map(()=>Math.random()*2-1)
            sentiment = sentiment.map(n => n * stockDiff / 2 + midStock)

            this.draw(stock, sentiment, symbol)
        }).bind(this))
    }

    /**
     * Makes an API call, and returns the json from it
     * @param {string} symbol Stock trading symbol
     */
    apiCall(symbol, apiKey) {
        const api = 'https://www.alphavantage.co/query?'
        const url = `${api}function=TIME_SERIES_DAILY&symbol=${symbol}&apikey=${apiKey}`

        return fetch(url).then(res => res.json())
    }

    /**
     * Draws the graphs in their respective context
     * @param {[number]} stock stock data
     * @param {[number]} sentiment sentiment data
     * @param {string} symbol The stock trading symbol
     */
    draw(stock, sentiment, symbol) {

        let chart = this.getChart(
            stock.slice(-NUM_DAYS),
            sentiment.slice(-NUM_DAYS)
        )

        //draw on canvas
        new Chart(document.getElementById(`${symbol}-preview`), chart)
        new Chart(document.getElementById(`${symbol}-full`), chart)
    }

    getChart(stock, sentiment) {
        // Generate colors
        const color = {
            r: this.random(50, 150),
            g: this.random(50, 150),
            b: this.random(50, 150)
        }

        // Generate a dataset
        const dataset = (title, data, color) => ({label: title, data: data, fill: false, borderColor: color, lineTension: 0.3, pointRadius: 0})

        return {
            type: 'line',
            data: {
                labels: this.getLabels(NUM_DAYS),
                datasets: [
                    dataset('Stock Price', stock, `rgb(${color.r}, ${color.g}, ${color.b})`),
                    dataset('Sentiment Rating', sentiment, `rgb(${color.r+100}, ${color.g+100}, ${color.b+100})`)
                ]
            },
            options: {}
        }
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
