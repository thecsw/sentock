export class StockGraph {

    constructor(name, stockData, PositivityData) {

        const color = {
            r: Math.random()*100+50,
            g: Math.random()*100+50,
            b: Math.random()*100+50
        }

        const chart = {
            type: 'line',
            data: {
                labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
                datasets: [
                    {
                        label: 'Stock Price',
                        data: stockData,
                        fill: false,
                        borderColor: `rgb(${color.r}, ${color.g}, ${color.b})`,
                        lineTension: 0.2
                    },
                    {
                        label: 'Positivity Rating',
                        data: PositivityData,
                        fill: false,
                        borderColor: `rgb(${color.r+100}, ${color.g+100}, ${color.b+100})`,
                        lineTension: 0.2
                    }
                ]
            },
            options: {}
        }
        new Chart(document.getElementById(`${name}-preview`), chart)
        new Chart(document.getElementById(`${name}-full`), chart)

    }
}
