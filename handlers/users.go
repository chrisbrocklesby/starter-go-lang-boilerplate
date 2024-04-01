package handlers

import (
	"fmt"
	"net/http"
	"os"
	"starter-go-app/database"
	"starter-go-app/helpers"
	"starter-go-app/templates"

	"github.com/go-http-utils/cookie"
	"github.com/google/uuid"
)

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			templates.Tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{
				"Title": "Register",
				"Error": "Email and password are required.",
			})
			return
		}

		var user database.User
		exists := database.DB.First(&user, "email = ?", email)

		if exists.RowsAffected == 0 {
			hashPassword, err := helpers.HashPassword(password)
			if err != nil {
				templates.Tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{
					"Title": "Register",
					"Error": "Failed to encrypt password.",
				})
				return
			}

			newUser := database.User{
				Email:    email,
				Password: hashPassword,
				Code:     uuid.New().String(),
			}

			database.DB.Create(&newUser)

			helpers.SendEmail(email, "Verify your email", "Click here to verify your email: http://localhost:3000/users/verify?code="+newUser.Code+"&email="+newUser.Email)

			templates.Tmpl.ExecuteTemplate(w, "verify.html", map[string]interface{}{
				"Title":   "Register",
				"Message": "User created successfully.",
			})
			return
		} else {
			templates.Tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{
				"Title": "Register",
				"Error": "User already exists.",
			})
		}
	} else {
		templates.Tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{
			"Title": "Register",
		})
	}
}

func UserVerifyHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	email := r.FormValue("email")

	if code == "" || email == "" {
		templates.Tmpl.ExecuteTemplate(w, "verify.html", map[string]interface{}{
			"Title": "Verify",
			"Error": "Code and email are required.",
		})
		return
	}

	var user database.User
	exists := database.DB.First(&user, "email = ? AND code = ?", email, code)

	if exists.RowsAffected == 0 {
		templates.Tmpl.ExecuteTemplate(w, "verify.html", map[string]interface{}{
			"Title": "Verify",
			"Error": "Invalid code or email.",
		})
		return
	}

	user.Verified = true
	user.Code = ""
	database.DB.Save(&user)

	templates.Tmpl.ExecuteTemplate(w, "verify.html", map[string]interface{}{
		"Title":   "Verify",
		"Message": "Email verified successfully.",
	})
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			templates.Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Title": "Login",
				"Error": "Email and password are required.",
			})
			return
		}

		var user database.User
		exists := database.DB.First(&user, "email = ?", email)

		if !user.Verified {
			templates.Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Title": "Login",
				"Error": "Email not verified.",
			})
			return
		}

		if exists.RowsAffected == 0 {
			templates.Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Title": "Login",
				"Error": "Invalid email or password.",
			})
			return
		}

		if helpers.ComparePassword(user.Password, password) != nil {
			templates.Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Title": "Login",
				"Error": "Password is incorrect.",
			})
			return
		}

		cookies := cookie.New(w, r, os.Getenv("COOKIE_KEY"))

		cookies.Set("session", user.ID, &cookie.Options{
			Signed:   true,
			HTTPOnly: true,
			Path:     "/",
		})

		fmt.Println("Redirect: ", redirect)
		// if redirect then redirect else home
		if redirect != "" {
			fmt.Println("Redirecting to: > ", redirect)
			http.Redirect(w, r, redirect, http.StatusSeeOther)
		} else {
			fmt.Println("Redirecting to: / > ", redirect)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		templates.Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Title": "Login",
		})
	}
}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookies := cookie.New(w, r, os.Getenv("COOKIE_KEY"))

	cookies.Remove("session", &cookie.Options{
		Signed:   true,
		HTTPOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func UserForgotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")

		if email == "" {
			templates.Tmpl.ExecuteTemplate(w, "forgot.html", map[string]interface{}{
				"Title": "Forgot Password",
				"Error": "Email is required.",
			})
			return
		}

		var user database.User
		exists := database.DB.First(&user, "email = ?", email)

		if exists.RowsAffected == 0 {
			templates.Tmpl.ExecuteTemplate(w, "forgot.html", map[string]interface{}{
				"Title": "Forgot Password",
				"Error": "Invalid email.",
			})
			return
		}

		user.Code = uuid.New().String()
		database.DB.Save(&user)

		helpers.SendEmail(email, "Reset your password", "Click here to reset your password: http://localhost:3000/users/reset?code="+user.Code+"&email="+user.Email)

		templates.Tmpl.ExecuteTemplate(w, "forgot.html", map[string]interface{}{
			"Title":   "Forgot Password",
			"Message": "Password reset email sent.",
		})
	} else {
		templates.Tmpl.ExecuteTemplate(w, "forgot.html", map[string]interface{}{
			"Title": "Forgot Password",
		})
	}
}

func UserResetHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	code := r.URL.Query().Get("code")

	if email == "" || code == "" {
		templates.Tmpl.ExecuteTemplate(w, "reset.html", map[string]interface{}{
			"Title": "Reset Password",
			"Error": "Email and code are required.",
		})
		return
	}

	var user database.User
	exists := database.DB.First(&user, "email = ? AND code = ?", email, code)

	if exists.RowsAffected == 0 {
		templates.Tmpl.ExecuteTemplate(w, "reset.html", map[string]interface{}{
			"Title": "Reset Password",
			"Error": "Invalid email or code.",
		})
		return
	}

	if r.Method == "POST" {
		password := r.FormValue("password")

		if password == "" {
			templates.Tmpl.ExecuteTemplate(w, "reset.html", map[string]interface{}{
				"Title": "Reset Password",
				"Error": "Password is required.",
			})
			return
		}

		hashPassword, err := helpers.HashPassword(password)
		if err != nil {
			templates.Tmpl.ExecuteTemplate(w, "reset.html", map[string]interface{}{
				"Title": "Reset Password",
				"Error": "Failed to encrypt password.",
			})
			return
		}

		user.Password = hashPassword
		user.Code = ""
		database.DB.Save(&user)

		templates.Tmpl.ExecuteTemplate(w, "reset.html", map[string]interface{}{
			"Title":   "Reset Password",
			"Message": "Password reset successfully.",
		})
	} else {
		templates.Tmpl.ExecuteTemplate(w, "reset.html", map[string]interface{}{
			"Title": "Reset Password",
		})
	}
}
