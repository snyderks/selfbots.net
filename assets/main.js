var creationHTML = `
<div>
  <label for="token-input">Discord User Token</label><br />
  <input type="text" placeholder="Please read disclaimer below." id="token-input"/>
</div>
<h2>
  Bot Features:
</h2>
<div class="selections" id="selections">

</div>
<div>
  <button onclick="submit()" id="submit-token">Submit</button>
</div>
`

var settingsHTML = `
<div>
  <h2>
    Bot Features:
  </h2>
  <div class="selections" id="selections">
  </div>
  <div>
    <button onclick="updateSettings()">Update Bot Settings</button>
  </div>
`

function userToken() {
  return localStorage.token !== undefined ? localStorage.token : null;
}

function insertUserInteraction() {
  var el = document.getElementById("bot-editor")
  if (userToken() === null) {
    el.innerHTML = creationHTML
  } else {
    el.innerHTML = settingsHTML
  }
}

function get(url) {
  return new Promise((resolve, reject) => {
    const req = new XMLHttpRequest()
    req.open('GET', url)
    req.onload = () => req.status === 200 ? resolve(req.response) : reject(Error(req.statusText))
    req.onerror = (e) => reject(Error(`Network Error: ${e}`))
    req.send()
  });
}
function sendTokenAndSelections(token, selections, url) {
  return new Promise((resolve, reject) => {
    const req = new XMLHttpRequest()
    req.open('POST', url)
    req.onload = () => req.status === 200 ? resolve(req.response) : reject(Error(req.statusText))
    req.onerror = (e) => reject(Error(`Network Error: ${e}`))
    req.send(JSON.stringify({
      "token": token,
      "selections": selections,
    }))
  });
}

function getUserSelections() {
  var chkBoxes = document.getElementsByClassName('selection-check');
  var selections = []
  var r = /(>)(\w+)/
  for (var i in chkBoxes) {
    if (chkBoxes[i].checked === true) {
      selections.push(r.exec(chkBoxes[i].parentElement.innerHTML)[2])
    }
  }
  return selections
}

function submit() {
  var token = document.getElementById("token-input").value
  token = token.replace(/\"/g, "")
  sendTokenAndSelections(token, getUserSelections(), "/sendToken")
  get("/api/discordLoginUrl")
    .then(function (data) {
      var url = JSON.parse(data)
      if (url !== undefined && url !== null && url.URL !== undefined) {
        window.location = url.URL
      }
  });
}

function updateSettings() {
  var token = localStorage.token
  if (token !== undefined) {
    sendTokenAndSelections(token, selections, "/sendToken")
  }
}

function toggleCheckbox(event) {
  event.stopPropagation()
  if (event.srcElement.checked) {
    event.srcElement.parentElement.classList.add("selection-checked");
  } else {
    event.srcElement.parentElement.classList.remove("selection-checked");
  }
}
function toggleCheckContainer(el) {
  var chk = el.srcElement.getElementsByTagName("input")[0]
  var val = chk.checked
  if (val) {
    chk.checked = false
    el.srcElement.classList.remove("selection-checked");
  } else {
    chk.checked = true
    el.srcElement.classList.add("selection-checked");
  }
  chk.value = !val
  console.log(chk)
}

function getSelections() {
  var path = "/assets/selections.json"
  var container = document.getElementById("selections")

  get(path).then(function (results) {
    if (results !== undefined && results.length > 0) {
      s = JSON.parse(results)
      for (var i in s) {
        var header = document.createElement("div")
        header.innerHTML = `<h3 class="selections-category-header">` + s[i].category + `</h3>`
        header.className = "selections-category"
        container.appendChild(header)

        for (var j in s[i].selections) {
          var el = document.createElement("div")
          el.className = "selection"
          el.id = "s" + i
          el.innerHTML = "<input type=\"checkbox\" class=\"selection-check\" />" +
          s[i].selections[j]
          el.onclick = toggleCheckContainer
          var chk = el.getElementsByTagName("input")[0]
          if (s[i].category === "Basic") {
            el.classList.add("selection-checked");
            chk.checked = true

          }
          chk.onclick = toggleCheckbox
          header.appendChild(el)
        }
      }
    }
  })
}

function init() {
  insertUserInteraction()
  getSelections()
}

init()
