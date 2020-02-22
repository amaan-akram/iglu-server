package routes

import (
	"fmt"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/tokens"

	"github.com/go-macaron/session"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

// AccountSettingsHandler handles the settings
func AccountSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Account Settings"
	ctx.Data["Accounts"] = models.GetUsers()
	ctx.HTML(200, "settings/accounts")
}

// PostAccountSettingsHandler handles the settings
func PostAccountSettingsHandler(ctx *macaron.Context, f *session.Flash) {
	switch ctx.Query("action") {
	case "get_invite":
		code := tokens.GenerateInviteKey()
		f.Info(fmt.Sprintf("Your new invitation code is: %s", code))
		break
	default:
		f.Error("Unknown action")
	}
	ctx.Redirect("/settings/accounts")
}

func DeleteAccountHandler(ctx *macaron.Context) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Params("username") {
			ctx.Data["DelUser"] = u
			ctx.HTML(200, "settings/del_account")
			return
		}
	}
	ctx.Redirect("/settings/accounts")
}

func PostDeleteAccountHandler(ctx *macaron.Context, f *session.Flash) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Query("username") {
			err := models.DeleteUser(ctx.Query("username"))
			if err != nil {
				f.Error("Failed to delete user!")
			} else {
				f.Success("User deleted!")
			}
		}
	}
	ctx.Redirect("/settings/accounts")
}

func EditAccountHandler(ctx *macaron.Context) {
	for _, u := range models.GetUsers() {
		if u.Username == ctx.Params("username") {
			ctx.Data["EditUser"] = u
			ctx.HTML(200, "settings/edit_account")
			return
		}
	}
	ctx.Redirect("/settings/accounts")
}

func PostEditAccountHandler(ctx *macaron.Context, f *session.Flash,
	form forms.EditAccountForm) {
	for _, u := range models.GetUsers() {
		if u.Username == form.Email {
			var updateCols []string
			user := models.User{Username: form.Email}

			// Update password
			if form.Password != "" {
				if form.Password != form.RePassword {
					f.Error("Passwords do not match!")
					ctx.Redirect(fmt.Sprintf("/settings/accounts/edit/%s", form.Email))
					return
				}
				pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
				if err != nil {
					panic(err)
				}
				updateCols = append(updateCols, "password")
				user.Password = string(pass)
			}

			// Update first name
			if form.FirstName != "" && form.FirstName != u.FirstName {
				updateCols = append(updateCols, "first_name")
				user.FirstName = form.FirstName
			}

			// Update first name
			if form.LastName != "" && form.LastName != u.LastName {
				updateCols = append(updateCols, "last_name")
				user.LastName = form.LastName
			}

			models.UpdateUserCols(&user, updateCols...)
			f.Success("User updated!")
		}
	}
	ctx.Redirect("/settings/accounts")
}

// AppearanceSettingsHandler handles the settings
func AppearanceSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Appearance Settings"
	ctx.HTML(200, "settings/appearance")
}

// PluginsSettingsHandler handles the settings
func PluginsSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Plugins"
	ctx.Data["Plugins"] = plugin.LoadedPlugins
	ctx.HTML(200, "settings/plugins")
}

// SettingsHandler handles the settings
func SettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "Settings"
	ctx.HTML(200, "settings")
}

// AboutSettingsHandler handles the about settings page
func AboutSettingsHandler(ctx *macaron.Context) {
	ctx.Data["NavTitle"] = "About"
	ctx.HTML(200, "settings/about")
}
