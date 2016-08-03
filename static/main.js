google.charts.load('current', {packages: ['corechart', 'line']});
google.charts.setOnLoadCallback(drawCharts);

var chartOptions = {
	hAxis: {
		format: 'MM/dd/yyyy',
		textStyle: {
			fontSize: 12
		}
	},
	vAxis: {
		minValue: 0,
		titleTextStyle: {
			italic: false
		},
		textStyle: {
			fontSize: 12
		}
	},
	legend: {
		position: 'top'
	},
	curveType: 'function',
	pointsVisible: true,
	interpolateNulls: true,
	pointShape: 'diamond',
	height: 300
};

function drawCharts() {
	$.get('api/v1/benchmarks', function(response) {
		var charts = {};

		for (var i = 0; i < response.length; i++) {
			var group = response[i].group;
			var metric = response[i].metric;
			var value = response[i].value;

			var timestamp = response[i].timestamp;
			var date = new Date(timestamp / 10 ** 6);

			if (charts[group] === undefined) {
				charts[group] = {};
			}
			if (charts[group][metric] === undefined) {
				charts[group][metric] = [];
			}

			charts[group][metric].push([date, value]);
		}

		Object.keys(charts).forEach(function(title, i) {
			chartOptions.title = title;

			var div = document.createElement('div');
			div.id = 'chart_div_' + i;
			$('#charts').append(div);

			var data;
			var indexes = [];
			var metrics = charts[title];

			Object.keys(metrics).forEach(function(metric, j) {
				var rows = metrics[metric];
				var mData = new google.visualization.DataTable();

				mData.addColumn('date', 'X');
				mData.addColumn('number', metric);
				mData.addRows(rows);

				if (j === 0) {
					data = mData;
				} else {
					indexes.push(j);
					data = google.visualization.data.join(data, mData, 'full', [[0, 0]], indexes, [1]);
				}
			});

			var chart = new google.visualization.LineChart(document.getElementById(div.id));
			chart.draw(data, chartOptions);
		})
	});
}
