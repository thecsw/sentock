let app = new Vue({
    el: '#app',
    data: {
        message: 'blah blah'
    }
})

var ctx = document.getElementById('myChart').getContext('2d');

var myLineChart = new Chart(ctx, {
    type: 'line',
    data: [{
        x: 10,
        y: 20
    }, {
        x: 15,
        y: 10
    }],
    options: {
        scales: {
            yAxes: [{
                stacked: true
            }]
        }
    }
});