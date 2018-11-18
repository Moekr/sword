function formatTimestamp(timestamp) {
    let date = new Date(timestamp * 1000);
    let year = date.getFullYear();
    let month = ("0" + (date.getMonth() + 1)).substr(-2);
    let day = ("0" + date.getDate()).substr(-2);
    let hour = ("0" + date.getHours()).substr(-2);
    let minute = ("0" + date.getMinutes()).substr(-2);
    return year + "-" + month + "-" + day + " " + hour + ":" + minute;
}

function paintAbbr(element, id, name, note, timeRange, dataSets) {
    let chart = echarts.init(document.getElementById(element));
    let option = {
        title: {
            text: name,
            link: "/detail.html?t=" + id + "&r=" + timeRange,
            target: "self",
            subtext: note,
        },
        grid: {
            top: note === "" ? 60 : 80,
        },
        tooltip: {
            trigger: "axis",
            axisPointer: {
                type: "cross",
                animation: false,
                label: {
                    backgroundColor: "#ccc",
                    borderColor: "#aaa",
                    borderWidth: 1,
                    shadowBlur: 0,
                    shadowOffsetX: 0,
                    shadowOffsetY: 0,
                    textStyle: {
                        color: "#222"
                    }
                }
            },
            formatter: function (params) {
                return params.map(param =>
                    param.seriesName + ": " + (param.value == null ? -1 : param.value) + "ms").join("<br>")
            }
        },
        legend: {
            top: "bottom",
            data: dataSets.map(dataSet => dataSet.observer.name)
        },
        xAxis: {
            data: dataSets[0].data.map(data => formatTimestamp(data[0]))
        },
        yAxis: {
            name: "延迟/ms"
        },
        series: dataSets.map(function (dataSet) {
            return {
                name: dataSet.observer.name,
                type: "line",
                data: dataSet.data.map(data => data[1] < 0 ? null : data[1]),
                showSymbol: false,
                showAllSymbol: true,
                sampling: "average",
            }
        })
    };
    chart.setOption(option);
    let pre = window.onresize;
    if (typeof pre !== "function") {
        window.onresize = function () {
            chart.resize();
        };
    } else {
        window.onresize = function () {
            chart.resize();
            pre();
        };
    }
}

function paintFull(element, title, dataSet) {
    let chart = echarts.init(document.getElementById(element));
    let option = {
        title: {
            text: title
        },
        tooltip: {
            trigger: "axis",
            axisPointer: {
                type: "cross",
                animation: false,
                label: {
                    backgroundColor: "#ccc",
                    borderColor: "#aaa",
                    borderWidth: 1,
                    shadowBlur: 0,
                    shadowOffsetX: 0,
                    shadowOffsetY: 0,
                    textStyle: {
                        color: "#222"
                    }
                }
            },
            formatter: function (params) {
                console.log(params);
                let avg = -1, max = -1, min = -1, lost = 0;
                params.forEach(param => {
                    let value = param.value == null ? -1 : param.value;
                    switch (param.seriesName) {
                        case "avg":
                            avg = value;
                            break;
                        case "max":
                            max = value;
                            break;
                        case "min":
                            min = value;
                            break;
                        case "lost":
                            lost = value;
                            break;
                    }
                });
                avg = Math.max(avg, -1);
                max = Math.max(max + min, -1);
                min = Math.max(min, -1);
                lost = Math.max(lost, 0);
                return "平均值：" + avg + "ms<br>最大值：" + max + "ms<br>最小值：" + min + "ms<br>丢包率：" + lost + "%"
            }
        },
        xAxis: {
            data: dataSet.data.map(data => formatTimestamp(data[0]))
        },
        yAxis: [
            {
                name: "延迟/ms"
            },
            {
                name: "丢包率",
                axisLabel: {
                    formatter: "{value}%"
                },
                splitLine: {
                    show: false,
                },
                max: 100,
                min: 0,
            }
        ],
        series: [
            {
                name: "avg",
                type: "line",
                data: dataSet.data.map(data => data[1] < 0 ? null : data[1]),
                showSymbol: false,
                showAllSymbol: true,
            },
            {
                name: "min",
                type: "line",
                data: dataSet.data.map(data => data[3] < 0 ? null : data[3]),
                symbol: "none",
                showAllSymbol: true,
                lineStyle: {
                    normal: {
                        opacity: 0
                    }
                },
                stack: "fluctuation",
            },
            {
                name: "max",
                type: "line",
                data: dataSet.data.map(data => data[2] < 0 ? null : data[2] - data[3]),
                symbol: "none",
                showAllSymbol: true,
                lineStyle: {
                    normal: {
                        opacity: 0
                    }
                },
                stack: "fluctuation",
                areaStyle: {
                    normal: {
                        color: "#ccc"
                    }
                },
            },
            {
                name: "lost",
                type: "bar",
                data: dataSet.data.map(data => data[4] < 0 ? null : data[4]),
                yAxisIndex: 1,
                barCategoryGap: "0%",
            }
        ]
    };
    chart.setOption(option);
    let pre = window.onresize;
    if (typeof pre !== "function") {
        window.onresize = function () {
            chart.resize();
        };
    } else {
        window.onresize = function () {
            chart.resize();
            pre();
        };
    }
}