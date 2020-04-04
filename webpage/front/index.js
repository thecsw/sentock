document.addEventListener('DOMContentLoaded', function() {
    M.Modal.init(document.querySelectorAll('.modal'), {})
})

let app = new Vue({
    el: '#app',
    data: {

    }
})

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

let exampleChart = {
    type: 'line',
    data: {
        labels: ['January','February','March','April','May','June','July'],
        datasets: [
            {
                label: 'Stock Price',
                data: [
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100
                ],
                fill: false,
                borderColor: `rgb(${color1.r}, ${color1.g}, ${color1.b})`,
                lineTension:0.1
            },
            {
                label: 'Positivity Rating',
                data: [
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100,
                    Math.random()*100
                ],
                fill: false,
                borderColor: `rgb(${color2.r}, ${color2.g}, ${color2.b})`,
                lineTension:0.1
            }
        ]
    },
    options: {}
}
let examplePreview = new Chart(document.getElementById('example-preview').getContext('2d'), exampleChart)
let exampleFull = new Chart(document.getElementById('example-full').getContext('2d'), exampleChart)
