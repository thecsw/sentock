document.addEventListener('DOMContentLoaded', function() {
    var elems = document.querySelectorAll('.modal');
    var instances = M.Modal.init(elems, {});
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

// var myChart = new Chart(examplePreviewCTX, {
//     type: 'bar',
//     data: {
//         labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
//         datasets: [{
//             label: 'Favorite Color',
//             data: [
//                 Math.random(),
//                 Math.random(),
//                 Math.random(),
//                 Math.random(),
//                 Math.random(),
//                 Math.random()
//             ],
//             backgroundColor: [
//                 'rgba(255, 99, 132, 0.4)',
//                 'rgba(54, 162, 235, 0.4)',
//                 'rgba(255, 206, 86, 0.4)',
//                 'rgba(75, 192, 192, 0.4)',
//                 'rgba(153, 102, 255, 0.4)',
//                 'rgba(255, 159, 64, 0.4)'
//             ],
//             borderColor: [
//                 'rgba(255, 99, 132, 1)',
//                 'rgba(54, 162, 235, 1)',
//                 'rgba(255, 206, 86, 1)',
//                 'rgba(75, 192, 192, 1)',
//                 'rgba(153, 102, 255, 1)',
//                 'rgba(255, 159, 64, 1)'
//             ],
//             borderWidth: 1

//         }]
//     },
//     options: {
//         scales: {
//             yAxes: [{
//                 ticks: {
//                     beginAtZero: true
//                 }
//             }]
//         }
//     }
// })

// var myChart = new Chart(exampleFullCTX, {
//     type: 'bar',
//     data: {
//         labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
//         datasets: [{
//             label: 'Favorite Color',
//             data: [
//                 Math.random(),
//                 Math.random(),
//                 Math.random(),
//                 Math.random(),
//                 Math.random(),
//                 Math.random()
//             ],
//             backgroundColor: [
//                 'rgba(255, 99, 132, 0.4)',
//                 'rgba(54, 162, 235, 0.4)',
//                 'rgba(255, 206, 86, 0.4)',
//                 'rgba(75, 192, 192, 0.4)',
//                 'rgba(153, 102, 255, 0.4)',
//                 'rgba(255, 159, 64, 0.4)'
//             ],
//             borderColor: [
//                 'rgba(255, 99, 132, 1)',
//                 'rgba(54, 162, 235, 1)',
//                 'rgba(255, 206, 86, 1)',
//                 'rgba(75, 192, 192, 1)',
//                 'rgba(153, 102, 255, 1)',
//                 'rgba(255, 159, 64, 1)'
//             ],
//             borderWidth: 1

//         }]
//     },
//     options: {
//         scales: {
//             yAxes: [{
//                 ticks: {
//                     beginAtZero: true
//                 }
//             }]
//         }
//     }
// })