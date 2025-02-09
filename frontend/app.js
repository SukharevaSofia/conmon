"use strict";

let intervalId;

function id(x){
  return document.getElementById(x);
}

function create_table_html(data_array) {
  let result = ""
  result += '<table class="table">'
  result += `<thead>
    <tr>
      <th scope="col">IP</th>
      <th scope="col">Ping time</th>
      <th scope="col">Last successful ping</th>
    </tr>
  </thead>`;
  result += '<tbody>';


  data_array.forEach((x) => {
    result +=
      `<tr>
         <th>${x.ip}</th>
         <th>${x.ping_time}</th>
         <th>${x.last_success}</th>
       </tr>`;
  });

  result += '</tbody></table>';
  return result;
}

function render_table(data) {
  data = JSON.parse(data)
  console.log(data)

  if (data.length == 0) {
    id("containers").innerHTML = "No containers :("
  } else {
    id("containers").innerHTML = create_table_html(data);
  }
}

function act_render_table() {
  fetch("/backend/get-table").
    then(res => res.text()).
    then(res => { render_table(res); });
}

const value = id("value");
const input = id("update-interval");

act_render_table()
intervalId = setInterval(act_render_table, 1000)

value.textContent = input.value;
input.addEventListener("input", (event) => {
  // Update HTML
  value.textContent = event.target.value;

  // Reset interval function
  clearInterval(intervalId);
  intervalId = setInterval(act_render_table, 
    id("update-interval").value * 1000);
});
