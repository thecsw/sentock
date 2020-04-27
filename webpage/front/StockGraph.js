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
     * @param {string} name Company name
     * @param {string} symbol Stock ticker symbol
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
     * @param {string} symbol Stock trading symbol
     */

    sentimentApiCall(name) {
        let d = new Date
        d.setUTCSeconds(0)
        let before, after

        const day = d.getDay(),
              hour = d.getUTCHours() - 4, //get EDT
              minute = d.getUTCMinutes(),

              before930 = hour < 9 || (hour === 9 && minute < 30)

        if(day === 0 || (day === 1 && before930)) {
            //last friday
            d.setUTCDate(new Date().getUTCDate() + (6 - new Date().getUTCDay() - 1) - 7)

            d.setUTCHours(9+4)
            d.setUTCMinutes(30)
            before = Math.round(d.getTime()/1000)

            d.setUTCHours(12+4+4)
            d.setUTCMinutes(0)
            after = Math.round(d.getTime()/1000)

        } else if(day === 6) {
            // yesterday
            d.setUTCDay(5)

            d.setUTCHours(9+4)
            d.setUTCMinutes(30)
            before = Math.round(d.getTime()/1000)

            d.setUTCHours(12+4+4)
            d.setUTCMinutes(0)
            after = Math.round(d.getTime()/1000)
        } else {
            //weekday
            if(before930) {
                //before open
                d.setUTCDay(d.getUTCDay() - 1) // get data from yesterday

                d.setUTCHours(9+4)
                d.setUTCMinutes(30)
                before = Math.round(d.getTime()/1000)

                d.setUTCHours(12+4+4)
                d.setUTCMinutes(0)
                after = Math.round(d.getTime()/1000)

            } else if(hour > 4) {
                //after close
                d.setUTCHours(9+4)
                d.setUTCMinutes(30)
                before = Math.round(d.getTime()/1000)

                d.setUTCHours(12+4+4)
                d.setUTCMinutes(0)
                after = Math.round(d.getTime()/1000)

            } else {
                //during day
                after = Math.round(d.getTime()/1000)

                d.setUTCHours(9+4)
                d.setUTCMinutes(30)
                before = Math.round(d.getTime()/1000)
            }
        }

        const api = 'https://sentocks.sandyuraz.com/api/sentiments'
        const url = `${api}?company=${name}&before=${after}&after=${before}`

        return fetch(url).then(res => res.json())
    }

    /**
     * Makes a call to the stocks, and returns the json from it
     * @param {string} symbol Stock trading symbol
     */
    stockApiCall(symbol, apiKey) {
        // const apiKey = 'VF5PZ9DBDNDKYYSN'
        const api = 'https://www.alphavantage.co/query?'
        const url = `${api}function=TIME_SERIES_INTRADAY&outputsize=full&symbol=${symbol}&interval=1min&apikey=${apiKey}`

        return fetch(url).then(res => res.json())
    }

    /**
     * Draws the graphs in their respective context
     * @param {[number]} stock stock data
     * @param {[number]} sentiment sentiment data
     * @param {string} symbol The stock trading symbol
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
     * @param {number} n number of days
     */
    getLabels(n) {
        let date = new Date
        date.setHours(9)
        date.setMinutes(30)

        return [...Array(n).keys()].map(n => {
            date = new Date(date.getTime() + 60000)
            let min = date.getMinutes()
            return `${date.getHours()}:${min<10?'0':''}${min}`
        })
    }

    /**
     * converts the JSON object from the API into a list of sentiment data
     * @param {object} json The JSON object to be converted
     */
    json2sentiments(json) {
        return json.map(e => e[1])
    }

    /**
     * Converts the JSON object from the API into a list of stock data
     * @param {object} json The JSON object to be converted
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
