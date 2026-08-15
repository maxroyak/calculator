(function () {
    "use strict";

    var display = document.getElementById("display");
    var buttons = document.querySelectorAll(".btn");

    // ---- State ----
    var currentInput = "0";
    var firstOperand = null;
    var operation = null;
    var waitingForOperand = false;

    // ---- Keyboard mapping ----
    var KEY_MAP = {
        "0": { action: "digit", value: "0" },
        "1": { action: "digit", value: "1" },
        "2": { action: "digit", value: "2" },
        "3": { action: "digit", value: "3" },
        "4": { action: "digit", value: "4" },
        "5": { action: "digit", value: "5" },
        "6": { action: "digit", value: "6" },
        "7": { action: "digit", value: "7" },
        "8": { action: "digit", value: "8" },
        "9": { action: "digit", value: "9" },
        ".": { action: "digit", value: "." },
        "+": { action: "op", value: "add" },
        "-": { action: "op", value: "subtract" },
        "*": { action: "op", value: "multiply" },
        "/": { action: "op", value: "divide" },
        "Enter": { action: "equals" },
        "=": { action: "equals" },
        "Escape": { action: "clear" },
        "Backspace": { action: "backspace" },
    };

    // ---- Display ----
    function updateDisplay() {
        var val = currentInput;
        // Truncate overly long strings
        if (val.length > 14) {
            val = String(parseFloat(val).toPrecision(12));
        }
        display.textContent = val;
        display.classList.remove("error");
    }

    function showError() {
        display.textContent = "Error";
        display.classList.add("error");
    }

    // ---- Digit handling ----
    function inputDigit(digit) {
        if (display.classList.contains("error")) {
            clearAll();
        }
        if (digit === ".") {
            if (waitingForOperand) {
                currentInput = "0.";
                waitingForOperand = false;
                updateDisplay();
                return;
            }
            if (currentInput.indexOf(".") !== -1) {
                return;
            }
            currentInput += ".";
            updateDisplay();
            return;
        }
        if (waitingForOperand) {
            currentInput = digit;
            waitingForOperand = false;
        } else {
            if (currentInput === "0") {
                currentInput = digit;
            } else {
                currentInput += digit;
            }
        }
        updateDisplay();
    }

    // ---- Operation handling ----
    function inputOp(op) {
        if (display.classList.contains("error")) {
            return;
        }
        var inputValue = parseFloat(currentInput);
        if (firstOperand === null) {
            firstOperand = inputValue;
        } else if (operation && !waitingForOperand) {
            // Chain: calculate intermediate result first
            var result = calculateSync(firstOperand, inputValue, operation);
            if (result === null) {
                showError();
                return;
            }
            firstOperand = result;
            currentInput = String(result);
            updateDisplay();
        }
        operation = op;
        waitingForOperand = true;
    }

    // ---- Equals handling ----
    function equals() {
        if (display.classList.contains("error")) {
            return;
        }
        if (operation === null || firstOperand === null) {
            return;
        }
        var inputValue = parseFloat(currentInput);
        var a = firstOperand;
        var b = inputValue;
        var op = operation;

        // Reset state immediately so UI is ready for next op
        firstOperand = null;
        operation = null;
        waitingForOperand = true;

        fetch("/api/calculate", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ operation: op, a: a, b: b }),
        })
            .then(function (res) {
                return res.json();
            })
            .then(function (data) {
                if (data.error) {
                    showError();
                } else {
                    currentInput = String(data.result);
                    // Keep firstOperand as result for chaining
                    firstOperand = parseFloat(currentInput);
                    waitingForOperand = true;
                    updateDisplay();
                }
            })
            .catch(function () {
                showError();
            });
    }

    // ---- Synchronous calc for intermediate chain ops ----
    function calculateSync(a, b, op) {
        switch (op) {
            case "add":
                return a + b;
            case "subtract":
                return a - b;
            case "multiply":
                return a * b;
            case "divide":
                if (b === 0) return null;
                return a / b;
            default:
                return null;
        }
    }

    // ---- Clear ----
    function clearAll() {
        currentInput = "0";
        firstOperand = null;
        operation = null;
        waitingForOperand = false;
        updateDisplay();
    }

    // ---- Backspace ----
    function backspace() {
        if (display.classList.contains("error")) {
            clearAll();
            return;
        }
        if (waitingForOperand) {
            return;
        }
        if (currentInput.length > 1) {
            currentInput = currentInput.slice(0, -1);
            if (currentInput === "-" || currentInput === "") {
                currentInput = "0";
            }
        } else {
            currentInput = "0";
        }
        updateDisplay();
    }

    // ---- Action dispatch ----
    function handleAction(action, value) {
        switch (action) {
            case "digit":
                inputDigit(value);
                break;
            case "op":
                inputOp(value);
                break;
            case "equals":
                equals();
                break;
            case "clear":
                clearAll();
                break;
            case "backspace":
                backspace();
                break;
        }
    }

    // ---- Button click handlers ----
    buttons.forEach(function (btn) {
        btn.addEventListener("click", function () {
            var action = btn.getAttribute("data-action");
            var value = btn.getAttribute("data-digit") || btn.getAttribute("data-op");
            handleAction(action, value);
        });
    });

    // ---- Keyboard handler ----
    document.addEventListener("keydown", function (e) {
        var mapped = KEY_MAP[e.key];
        if (!mapped) return;
        e.preventDefault();
        handleAction(mapped.action, mapped.value);
        // Visual feedback: find matching button and flash pressed state
        flashButton(mapped.action, mapped.value);
    });

    function flashButton(action, value) {
        var selector;
        if (action === "digit") {
            selector = '[data-action="digit"][data-digit="' + value + '"]';
        } else if (action === "op") {
            selector = '[data-action="op"][data-op="' + value + '"]';
        } else if (action === "equals") {
            selector = '[data-action="equals"]';
        } else if (action === "clear") {
            selector = '[data-action="clear"]';
        } else if (action === "backspace") {
            return; // no button for backspace
        } else {
            return;
        }
        var btn = document.querySelector(selector);
        if (btn) {
            btn.classList.add("pressed");
            setTimeout(function () {
                btn.classList.remove("pressed");
            }, 100);
        }
    }

    // ---- Initial display ----
    updateDisplay();
})();