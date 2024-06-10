const getID = id => document.getElementById(id)

const loginButton = document.getElementById("login-button")
const registerButton = document.getElementById("register-button")
const loginForm = document.getElementById("login-form")
const registerForm = document.getElementById("register-form")

const activeClass = ["bg-[#1d1d1d]", "text-white"]
const notActiveClass = ["bg-white", "text-[#1d1d1d]", "hover:bg-[#1d1d1d]", "hover:text-white"]

loginButton.addEventListener("click", () => {
	registerButton.classList.remove(...activeClass)
	registerButton.classList.add(...notActiveClass)

	loginButton.classList.add(...activeClass)
	loginButton.classList.remove(...notActiveClass)

	registerForm.classList.remove("block")
	registerForm.classList.add("hidden")

	loginForm.classList.remove("hidden")
	loginForm.classList.add("block")
})

registerButton.addEventListener("click", () => {
	loginButton.classList.remove(...activeClass)
	loginButton.classList.add(...notActiveClass)

	registerButton.classList.add(...activeClass)
	registerButton.classList.remove(...notActiveClass)

	loginForm.classList.remove("block")
	loginForm.classList.add("hidden")

	registerForm.classList.remove("hidden")
	registerForm.classList.add("block")
})

registerForm.addEventListener("htmx:configRequest", async e => {
	const name = getID("name")
	const email = getID("email")
	const password = getID("password")
	const confirmPassword = getID("confirm-password")

	if(!name || !email || !password || !confirmPassword) {
		e.preventDefault()
		alert("All fields are required.")
		console.log(name)
		console.log(email)
		console.log(password)
		console.log(confirmPassword)
		return
	}

	cons

	if(password !== confirmPassword) {
		e.preventDefault()
		alert("Passwords should be equal")
		return
	}

	if(password.length < 7 || confirmPassword.length < 7) {
		e.preventDefault()
		alert("Password shoud be 7 characters.")
		return
	}

	if(!/^[a-zA-Z ]+$/.test(name)) {
		e.preventDefault()
		alert("Name should be letters and spaces only.")
		return
	}

	if(name.length < 4) {
		e.preventDefault()
		alert("Name should be 8 characters.")
		return
	}

	if(email.length < 8) {
		e.preventDefault()
		alert("Email should be 8 characters.")
		return
	}
})

document.addEventListener("animationend", e => {
	if(e.animationName === "hide_alert") {
		const alert = document.getElementById("alert")

		alert.parentNode.removeChild(alert)
	}
})

document.addEventListener("htmx:responseError", e => {
	if(e.detail.xhr.status === 400) document.getElementById("error-message").innerHTML = e.detail.xhr.responseText
})

