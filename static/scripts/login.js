function showPassword() {
    const input = document.getElementById("password");
    const icon = document.getElementById("eye-icon");
  
    input.type = "text";
    icon.src = "/static/svgs/eye-open.svg";
  }
  
  function hidePassword() {
    const input = document.getElementById("password");
    const icon = document.getElementById("eye-icon");
  
    input.type = "password";
    icon.src = "/static/svgs/eye-closed.svg";
  }