function xAxis(timeRange) {
    let res = [];
    let ts = new Date().getTime() / 1000;
    let step = function (timeRange) {
        switch (timeRange) {
            case "2":
                return 7 * 60;
            case "3":
                return 30 * 60;
            case "4":
                return 6 * 60 * 60;
            default:
                return 60;
        }
    }(timeRange);
    for (let i = 1440; i > 0; i--) {
        res.push(formatTimestamp(ts - i * step));
    }
    return res;
}

function formatTimestamp(timestamp) {
    let date = new Date(timestamp * 1000);
    let year = date.getFullYear();
    let month = ("0" + (date.getMonth() + 1)).substr(-2);
    let day = ("0" + date.getDate()).substr(-2);
    let hour = ("0" + date.getHours()).substr(-2);
    let minute = ("0" + date.getMinutes()).substr(-2);
    return year + "-" + month + "-" + day + " " + hour + ":" + minute;
}

function autoResize(chart) {
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
            data: dataSets.map(dataSet => dataSet.client.name)
        },
        xAxis: {
            data: xAxis(timeRange),
        },
        yAxis: {
            name: "延迟/ms"
        },
        series: dataSets.map(function (dataSet) {
            return {
                name: dataSet.client.name,
                type: "line",
                data: dataSet.avg.map(i => i < 0 ? null : i),
                showSymbol: false,
                showAllSymbol: true,
                sampling: "average",
            }
        })
    };
    chart.setOption(option);
    autoResize(chart);
}

function paintFull(element, title, timeRange, dataSet) {
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
                let avg = -1, max = -1, min = -1, std = -1, los = 0;
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
                        case "std":
                            std = value;
                            break;
                        case "los":
                            los = value;
                            break;
                    }
                });
                avg = Math.max(avg, -1);
                max = Math.max(max + min, -1);
                min = Math.max(min, -1);
                std = Math.max(std, -1);
                los = Math.max(los, 0);
                return "平均：" + avg + "ms<br>最大：" + max + "ms<br>最小：" + min + "ms<br>抖动：" + std + "ms<br>丢包：" + los + "%"
            }
        },
        xAxis: {
            data: xAxis(timeRange),
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
                data: dataSet.avg.map(i => i < 0 ? null : i),
                showSymbol: false,
                showAllSymbol: true,
            },
            {
                name: "min",
                type: "line",
                data: dataSet.min.map(i => i < 0 ? null : i),
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
                data: dataSet.max.map((i, idx) => i < 0 ? null : i - dataSet.min[idx]),
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
                name: "los",
                type: "bar",
                data: dataSet.los.map(i => i < 0 ? null : i),
                yAxisIndex: 1,
                barCategoryGap: "0%",
            },
            {
                name: 'std',
                type: 'line',
                symbolSize: 0,
                showSymbol: false,
                lineStyle: {
                    width: 0,
                    color: 'rgba(0, 0, 0, 0)',
                },
                data: dataSet.std.map(i => i < 0 ? null : i),
            }
        ]
    };
    chart.setOption(option);
    autoResize(chart);
}