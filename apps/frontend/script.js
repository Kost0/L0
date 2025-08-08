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

        const order = await response.json();
        resultDiv.innerHTML = generateOrderTable(order);
    } catch (error) {
        resultDiv.innerHTML = '<p style="color: red;">Ошибка: ${error.message}</p>';
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
            displayValue = isNan(date) ? value : date.toLocalString('ru-RU')
        }

        rows += `<tr><td><strong>${key}:</strong></td><td>${displayValue}</td></tr>`;
    }

    return`
    <table style="width: 100%; border-collapse: collapse; margin-bottom: 10px;">
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
            if (val instanceof Date) val = val.toLocalString();
            if (typeof val === 'object') val = JSON.stringify(val, null, 2);
            return `<td style="padding: 4px 8px; border-bottom: 1px solid #eee;">${val ?? '-'}</td>`;
        }).join("");
    return `<tr>${cells}</tr>`;
    }).join("");

    return `
    <table style="width: 100%; border-collapse: collapse; margin: 5px 0; font-size: 0.9em;">
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
    const { delivery, payment, items, ...mainData } = order;

    const mainTable = objectToTable(mainData, "Основная информация");

    let deliveryTable = delivery ? objectToTable(delivery, "delivery") : "";
    let paymentTable = payment ? objectToTable(payment, "payment") : "";
    let itemsTable = items ? arrayToTable(items, "items") : "";

    return `
    <div style="font-family: Arial, sans-serif; max-width: 1000px;">
      <h2>Заказ #${order.id}</h2>
      ${mainTable}
      ${deliveryTable}
      ${paymentTable}
      ${itemsTable}
    </div>
  `;
}