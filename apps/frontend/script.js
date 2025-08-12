async function fetchOrder() {
    const orderId = document.getElementById('orderId').value.trim();
    const resultDiv = document.getElementById('result');

    if (!orderId) {
        resultDiv.innerHTML = '<p>Введите ID заказа</p>';
        return;
    }

    resultDiv.innerHTML = '<p>Загрузка...</p>'

    try {
        const response = await fetch(`${window.API_URL || `http://localhost:8080`}/orders/${orderId}`);

        if (!response.ok) {
            throw new Error(`Заказ не найден: ${response.status}`);
        }

        const rawOrder = await response.json();

        const adaptedOrder = {
            ...rawOrder.Order,
            Delivery: rawOrder.Delivery,
            Payment: rawOrder.Payment,
            Items: rawOrder.Items
        }
        resultDiv.innerHTML = generateOrderTable(adaptedOrder);
    } catch (error) {
        resultDiv.innerHTML = `<p style="color: red;">Ошибка: ${error.message}</p>`;
    }
}

function objectToTable(obj, title = "") {
    let rows = "";
    for (const [key, value] of Object.entries(obj)) {
        let displayValue = value;

        if (value && typeof value === 'object' && !Array.isArray(value)) {
            displayValue = objectToTable(value)
        }
        else if (Array.isArray(value)) {
            displayValue = arrayToTable(value);
        }
        else if (key.toLowerCase().includes('date') && value) {
            const date = new Date(value)
            displayValue = isNaN(date.getTime()) ? value : date.toLocaleString('ru-RU')
        }

        rows += `<tr><td><strong>${key}:</strong></td><td>${displayValue}</td></tr>`;
    }

    return`
    <table class="data-table">
      ${title ? `<caption><strong>${title}</strong></caption>` : ""}
      ${rows}
    </table>
    `;
}

function arrayToTable(array) {
    if (array.length === 0) return ("(пусто)");

    const headers = Object.keys(array[0] || {});
    let headerRow = headers.map(h => `<th style="text-align: left; border-bottom: 1px solid #ddd;">${h}</th>`).join("");
    let dataRows = array.map(item => {
        let cells = headers.map(h => {
            let val = item[h];
            if (val instanceof Date) val = val.toLocaleString();
            if (typeof val === 'object') val = JSON.stringify(val, null, 2);
            return `<td style="padding: 4px 8px; border-bottom: 1px solid #eee;">${val ?? '-'}</td>`;
        }).join("");
    return `<tr>${cells}</tr>`;
    }).join("");

    return `
    <table class="data-table">
      <thead>
        <tr>${headerRow}</tr>
      </thead>
      <tbody>
        ${dataRows}
      </tbody>
    </table>
    `;
}

function generateOrderTable(order) {
    if (!order) {
        return '<p>Данные заказа отсутствуют</p>'
    }

    const { Delivery, Payment, Items, ...mainData } = order;

    let mainTable = objectToTable(mainData, "order")
    let deliveryTable = Delivery ? objectToTable(Delivery, "delivery") : "";
    let paymentTable = Payment ? objectToTable(Payment, "payment") : "";
    let itemsTable = Items ? arrayToTable(Items, "items") : "";

    return `
    <div style="font-family: Arial, sans-serif; max-width: 1000px;">
      ${mainTable}
      ${deliveryTable}
      ${paymentTable}
      ${itemsTable}
    </div>
  `;
}