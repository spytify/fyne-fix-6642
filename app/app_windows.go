//go:build !ci && !android && !ios && !wasm && !test_web_driver && !tinygo

package app

import (
	"errors"
	"net/url"
	"runtime"
	"time"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	internalapp "fyne.io/fyne/v2/internal/app"
)

// OpenURL hands the link to the shell, as url.dll's FileProtocolHandler
// does, without starting rundll32 to get there. Spytify's fork: security
// software watches rundll32 when a GUI app starts it.
func (a *fyneApp) OpenURL(u *url.URL) error {
	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	target, err := windows.UTF16PtrFromString(u.String())
	if err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		runtime.LockOSThread()
		if windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED|windows.COINIT_DISABLE_OLE1DDE) == nil {
			defer windows.CoUninitialize()
		}
		done <- windows.ShellExecute(0, verb, target, nil, nil, windows.SW_SHOWNORMAL)
	}()
	return <-done
}

// Spytify's fork: upstream shows notifications by writing a PowerShell
// script to %TEMP% and running it with -ExecutionPolicy Bypass and
// Invoke-Expression, which security software flags. The fork never starts
// PowerShell; Spytify shows its notifications itself, through the Windows
// Runtime.
var errNoNotifications = errors.New("this build of Fyne shows no notifications on Windows")

func (a *fyneApp) SendNotification(*fyne.Notification) {
	fyne.LogError("Notification not shown", errNoNotifications)
}

func (a *fyneApp) ScheduleNotification(*fyne.Notification, time.Time) (*fyne.ScheduledNotification, error) {
	return nil, errNoNotifications
}

func (a *fyneApp) CancelScheduledNotification(string) error {
	return errNoNotifications
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default, this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	a.Driver().(systrayDriver).SetSystemTrayMenu(menu)
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

// SetSystemTrayWindow assigns a window to be shown with the system tray menu is tapped.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayWindow(w fyne.Window) {
	a.Driver().(systrayDriver).SetSystemTrayWindow(w)
}

func watchTheme(s *settings) {
	go internalapp.WatchTheme(func() {
		fyne.Do(s.setupTheme)
	})
}

func (a *fyneApp) registerRepositories() {
	// no-op
}
