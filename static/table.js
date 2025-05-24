// Set up column visibility toggles

let sortColumn = "name";
let sortDirection = "asc";
let data = [];
let columns = [];

function setupColumnToggles(columns) {
  const togglesContainer = document.querySelector(".column-toggles");
  columns.forEach((column) => {
    const label = document.createElement("label");
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.checked = column.visible;
    checkbox.dataset.column = column.id;
    checkbox.addEventListener("change", function () {
      columns.find((col) => col.id === column.id).visible = this.checked;
      renderTable();
    });

    label.appendChild(checkbox);
    label.appendChild(document.createTextNode(` ${column.label}`));
    togglesContainer.appendChild(label);
  });
}

// Sort function
function sortData(columnId) {
  if (sortColumn === columnId) {
    sortDirection = sortDirection === "asc" ? "desc" : "asc";
  } else {
    sortColumn = columnId;
    sortDirection = "asc";
  }
  renderTable();
}

function renderTable() {
  const tableHead = document.querySelector("#data-table thead tr");
  const tableBody = document.querySelector("#data-table tbody");
  const searchTerm = document
    .getElementById("search-input")
    .value.toLowerCase();

  // Clear existing content
  tableHead.innerHTML = "";
  tableBody.innerHTML = "";

  // Add table headers
  columns.forEach((column) => {
    if (column.visible) {
      const th = document.createElement("th");
      th.textContent = column.label;
      th.classList.add("sortable");
      if (sortColumn === column.id) {
        th.classList.add(sortDirection);
      }
      th.addEventListener("click", () => sortData(column.id));
      tableHead.appendChild(th);
    }
  });

  // Filter and sort data
  const filteredData = data.filter((d) => {
    return Object.values(d).some((value) =>
      String(value).toLowerCase().includes(searchTerm)
    );
  });

  const sortedData = [...filteredData].sort((a, b) => {
    const valueA = String(a[sortColumn] || "");
    const valueB = String(b[sortColumn] || "");

    return sortDirection === "asc"
      ? valueA.localeCompare(valueB)
      : valueB.localeCompare(valueA);
  });

  // Populate table rows
  sortedData.forEach((r) => {
    const row = document.createElement("tr");

    columns.forEach((column) => {
      if (column.visible) {
        const cell = document.createElement("td");
        cell.textContent = r[column.id] || "";
        row.appendChild(cell);
      }
    });

    tableBody.appendChild(row);
  });
}
