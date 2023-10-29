const radios = document.querySelectorAll("#task-list input[type=radio]");
for(let j = 0; j < radios.length; j++) {
    const id = radios[j].id;
    const taskId = id.split("-")[0];
    console.log("Task ID: ", taskId);
}
//for(var i = 0, max = radios.length; i < max; i++) {
//    radios[i].onclick = function() {
//        alert(this.value);
//    }
//}

const setCookie = (taskId, exdays) => {
    console.log('Tried to set a cookie');
    //document.cookie = "username=John Doe; expires=Thu, 18 Dec 2013 12:00:00 UTC"; 
    const d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    let expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
} 

const getCookie = (cname) => {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');
  for(let i = 0; i <ca.length; i++) {
    let c = ca[i];
    console.log(c);
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

const checkCookie = () => {
  let user = getCookie("1");
  if (user != "") {
    alert("Welcome again " + user);
  } else {
    user = prompt("Please enter your name:", "");
    if (user != "" && user != null) {
      setCookie("username", user, 365);
    }
  }
}
