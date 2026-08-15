(function () {
    "use strict";

    var form = document.getElementById("calc-form");
    var resultBox = document.getElementById("result-box");
    var resultValue = document.getElementById("result");
    var errorBox = document.getElementById("error-box");
    var errorText = document.getElementById("error-text");
    var historySection = document.getElementById("history-section");
    var historyList = document.getElementById("history-list");
    var calcBtn = document.getElementById("calc-btn");

    var history = [];

    var SYMBOLS = {
        add: "+",
        subtract: "−",
        multiply: "×",
        divide: "÷",
    };

    function hideResult() {
        resultBox.classList.add("hidden");
    }

    function hideError() {
        errorBox.classList.add("hidden");
    }

    function showError(msg) {
        errorText.textContent = msg;
        errorBox.classList.remove("hidden");
        hideResult();
    }

    function showResult(value) {
        resultValue.textContent = value;
        resultBox.classList.remove("hidden");
        hideError();
    }

    function addHistory(a, op, b, result) {
        var sym = SYMBOLS[op] || op;
        var entry = a + " " + sym + " " + b + " = ";
        var li = document.createElement("li");
        li.innerHTML = entry + '<span class="hist-result">' + result + "</span>";
        historyList.insertBefore(li, historyList.firstChild);
        history.push({ a: a, op: op, b: b, result: result });
        if (history.length > 10) {
            history.shift();
            if (historyList.lastChild) {
                historyList.removeChild(historyList.lastChild);
            }
        }
        historySection.classList.remove("hidden");
    }

    form.addEventListener("submit", function (e) {
        e.preventDefault();
        hideError();
        hideResult();

        var a = document.getElementById("a").value;
        var b = document.getElementById("b").value;
        var op = document.getElementById("operation").value;

        if (a === "" || b === "") {
            showError("Please enter both numbers.");
            return;
        }

        calcBtn.disabled = true;
        calcBtn.textContent = "Calculating…";

        fetch("/api/calculate", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ operation: op, a: parseFloat(a), b: parseFloat(b) }),
        })
            .then(function (res) {
                return res.json();
            })
            .then(function (data) {
                if (data.error) {
                    showError(data.error);
                } else {
                    var formatted = Number(data.result).toString();
                    showResult(formatted);
                    addHistory(a, op, b, formatted);
                }
            })
            .catch(function () {
                showError("Network error — unable to reach the server.");
            })
            .finally(function () {
                calcBtn.disabled = false;
                calcBtn.textContent = "Calculate";
            });
    });
})();