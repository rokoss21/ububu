package auth

import (
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type AuthDialog struct {
	window   fyne.Window
	callback func(bool)
	dialog   *dialog.CustomDialog
}

func NewAuthDialog(window fyne.Window, callback func(bool)) *AuthDialog {
	return &AuthDialog{
		window:   window,
		callback: callback,
	}
}

func (a *AuthDialog) Show() {
	// Создаем поля ввода
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter your password")
	
	rememberCheck := widget.NewCheck("Remember for this session", nil)
	
	// Информационный текст
	infoLabel := widget.NewRichTextFromMarkdown(`
🔒 **Administrator Authentication Required**

This operation requires root access to perform system modifications.

**Selected operations will:**
• Install system updates
• Modify system files  
• Access hardware information
`)
	
	// Кнопки
	cancelButton := widget.NewButton("Cancel", func() {
		a.dialog.Hide()
		a.callback(false)
	})
	
	authenticateButton := widget.NewButton("Authenticate", func() {
		password := passwordEntry.Text
		if password == "" {
			dialog.ShowError(fmt.Errorf("password cannot be empty"), a.window)
			return
		}
		
		// Проверяем пароль
		if a.validatePassword(password) {
			a.dialog.Hide()
			a.callback(true)
		} else {
			dialog.ShowError(fmt.Errorf("incorrect password"), a.window)
			passwordEntry.SetText("")
		}
	})
	
	// Делаем кнопку Authenticate основной
	authenticateButton.Importance = widget.HighImportance
	
	// Обработка Enter в поле пароля
	passwordEntry.OnSubmitted = func(text string) {
		authenticateButton.OnTapped()
	}
	
	// Компоновка
	content := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Password:", passwordEntry),
		),
		rememberCheck,
		widget.NewSeparator(),
		container.NewHBox(
			cancelButton,
			widget.NewSeparator(),
			authenticateButton,
		),
	)
	
	// Создаем диалог
	a.dialog = dialog.NewCustom("Administrator Authentication", "", content, a.window)
	a.dialog.Resize(fyne.NewSize(450, 350))
	
	// Фокус на поле пароля
	a.window.Canvas().Focus(passwordEntry)
	
	a.dialog.Show()
}

func (a *AuthDialog) validatePassword(password string) bool {
	// Проверяем пароль через sudo -v
	cmd := exec.Command("sudo", "-S", "-v")
	cmd.Stdin = strings.NewReader(password + "\n")
	
	err := cmd.Run()
	return err == nil
}
