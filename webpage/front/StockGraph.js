/**
 * All the information to create a graph
 */
export class StockGraph {

    error
    id
    name
    numDays
    /**
     * Creates a stock graph
     * @param name Company name
     * @param symbol Stock ticker symbol
     */
    constructor(name, symbol, apiKey) {
        this.error = false
        this.id = symbol
        this.name = name

        // Make API call
        let stocks = this.stockApiCall(symbol, apiKey)
                         .then(this.json2stocks)

        // Scale sentiment data and draw
        stocks.then((stock => {

            // Warn about API key
            if(stock.length === 0) this.error = true

            this.sentimentApiCall(name)
                .then(this.json2sentiments)
                .then(sentiment => {
                    this.numDays = sentiment.length
                    this.draw(stock, sentiment, symbol)
                })
        }).bind(this))
    }

    /**
     * Makes a call to the sentiment API, and returns the json from it
     * @param symbol Stock trading symbol
     */

    sentimentApiCall(name) {
        const api = 'https://sentocks.sandyuraz.com/api/sentiments/market'
        const url = `${api}?company=${name}`

        return fetch(url).then(res => res.json())
    }

    /**
     * Makes a call to the stocks, and returns the json from it
     * @param symbol Stock trading symbol
     */
    stockApiCall(symbol, apiKey) {
        // const apiKey = 'VF5PZ9DBDNDKYYSN'
        const api = 'https://www.alphavantage.co/query?'
        const url = `${api}function=TIME_SERIES_INTRADAY&outputsize=full&symbol=${symbol}&interval=1min&apikey=${apiKey}`

        return fetch(url).then(res => res.json())
    }

    /**
     * Draws the graphs in their respective context
     * @param stock stock data
     * @param sentiment sentiment data
     * @param symbol The stock trading symbol
     */
    draw(stock, sentiment, symbol) {
        // console.log(symbol, 'stock', stock.slice(-this.numDays))
        // console.log(symbol, 'sentiment', sentiment.slice(-this.numDays))

        let chart = this.getChart(
            stock.slice(-this.numDays),
            sentiment.slice(-this.numDays)
        )

        //draw on canvas
        new Chart(document.getElementById(`${symbol}-preview`), chart)
        new Chart(document.getElementById(`${symbol}-full`), chart)
    }

    getChart(stock, sentiment) {
        // Generate colors
        const color = {
            r: Math.random()*100 + 50,
            g: Math.random()*100 + 50,
            b: Math.random()*100 + 50
        }

        // Generate a dataset
        const dataset = (title, data, color, side) => ({label: title, data: data, yAxisID: side, fill: false, borderColor: color, lineTension: 0, pointRadius: 0})

        return {
            type: 'line',
            data: {
                labels: this.getLabels(this.numDays),
                datasets: [
                    dataset('Stock Price', stock, `rgb(${color.r}, ${color.g}, ${color.b})`, 'left'),
                    dataset('Sentiment Rating', sentiment, `rgb(${color.r+100}, ${color.g+100}, ${color.b+100})`, 'right')
                ]
            },
            options: {
                responsive: true,
                hoverMode: 'index',
                stacked: false,
                scales: {
                    yAxes: [{
                        type: 'linear',
                        display: true,
                        position: 'left',
                        id: 'left'
                    }, {
                        type: 'linear',
                        display: true,
                        position: 'right',
                        id: 'right'
                    }],
                }
            }
        }
    }

    /**
     * Generates the list of labels for the last n days
     * @param n number of days
     */
    getLabels(n) {
        return new Array(n).fill(' ')
        // let date = new Date
        // date.setHours(9)
        // date.setMinutes(30)

        // return [...Array(n).keys()].map(n => {
        //     date = new Date(date.getTime() + 60000)
        //     let min = date.getMinutes()
        //     return `${date.getHours()}:${min<10?'0':''}${min}`
        // })
    }

    /**
     * converts the JSON object from the API into a list of sentiment data
     * @param json The JSON object to be converted
     */
    json2sentiments(json) {
        return json.map(e => e[1])
    }

    /**
     * Converts the JSON object from the API into a list of stock data
     * @param json The JSON object to be converted
     */
    json2stocks(json) {
        const minKey = 'Time Series (1min)'
        const dataKey = '4. close'

        const minutes = json[minKey]
        let data = []

        for(let min in minutes) data.push(minutes[min][dataKey])

        return data
    }
}
