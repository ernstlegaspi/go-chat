const loginButton = document.getElementById("login-button")
const registerButton = document.getElementById("register-button")
const loginForm = document.getElementById("login-form")
const registerForm = document.getElementById("register-form")

const activeClass = ["bg-[#1d1d1d]", "text-white"]
const notActiveClass = ["bg-white", "text-[#1d1d1d]", "hover:bg-[#1d1d1d]", "hover:text-white"]

registerForm.classList.add("hidden")

loginButton.addEventListener("click", () => {
	registerButton.classList.remove(...activeClass)
	registerButton.classList.add(...notActiveClass)

	loginButton.classList.add(...activeClass)
	loginButton.classList.remove(...notActiveClass)
})

registerButton.addEventListener("click", () => {
	loginButton.classList.remove(...activeClass)
	loginButton.classList.add(...notActiveClass)

	registerButton.classList.add(...activeClass)
	registerButton.classList.remove(...notActiveClass)
})
