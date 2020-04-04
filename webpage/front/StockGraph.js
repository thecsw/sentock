export class StockGraph {

    constructor(name, stockData, PositivityData) {

        //generate random color
        let color1 = {
            r: Math.random()*255,
            g: Math.random()*255,
            b: Math.random()*255
        }

        //make sure luminance is high to see
        while(0.2126*color1.r + 0.7152*color1.g + 0.0722*color1.b < 64) {
            if(0.2126*color1.r < 0.7152*color1.g && 0.2126*color1.r < 0.0722*color1.b) color1.r++
            if(0.7152*color1.g < 0.2126*color1.r && 0.7152*color1.g < 0.0722*color1.b) color1.g++
            if(0.0722*color1.b < 0.7152*color1.g && 0.0722*color1.b < 0.2126*color1.r) color1.b++
        }

        //generate random color
        let color2 = {
            r: Math.random()*255,
            g: Math.random()*255,
            b: Math.random()*255
        }

        //make sure luminance is high to see
        while(0.2126*color2.r + 0.7152*color2.g + 0.0722*color2.b < 64) {
            if(0.2126*color2.r < 0.7152*color2.g && 0.2126*color2.r < 0.0722*color2.b) color2.r++
            if(0.7152*color2.g < 0.2126*color2.r && 0.7152*color2.g < 0.0722*color2.b) color2.g++
            if(0.0722*color2.b < 0.7152*color2.g && 0.0722*color2.b < 0.2126*color2.r) color2.b++
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
                        borderColor: `rgb(${color1.r}, ${color1.g}, ${color1.b})`,
                        lineTension: 0.2
                    },
                    {
                        label: 'Positivity Rating',
                        data: PositivityData,
                        fill: false,
                        borderColor: `rgb(${color2.r}, ${color2.g}, ${color2.b})`,
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
