package action

import (
	"log"
	"os"

	"github.com/lxn/walk"
	"golang.org/x/sys/windows/registry"
)

const REG_KEY = "ClipboardOnline"

var REG_VALUE = os.Args[0]

func NewAutoRunAction() (*walk.Action, error) {
	action := walk.NewAction()
	if err := action.SetText("开机启动"); err != nil {
		return nil, err
	}

	if err := action.SetCheckable(true); err != nil {
		return nil, err
	}

	isAutoRun, err := queryAutoRun()
	if err != nil {
		return nil, err
	}
	action.SetChecked(isAutoRun)

	action.Triggered().Attach(func() {
		if action.Checked() {
			if err := enableAutoRun(); err != nil {
				action.SetChecked(false)
				log.Println(err)
			}
		} else {
			if err := disableAutoRun(); err != nil {
				action.SetChecked(true)
				log.Println(err)
			}
		}
	})

	return action, nil
}

func queryAutoRun() (bool, error) {
	key, err := openAutoRunKey(registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer key.Close()
	val, _, err := key.GetStringValue(REG_KEY)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	if val == REG_VALUE {
		return true, nil
	}
	return false, nil
}

func openAutoRunKey(access uint32) (registry.Key, error) {
	autorunKey := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	key, err := registry.OpenKey(registry.CURRENT_USER, autorunKey, access)
	if err != nil {
		return 0, err
	}
	return key, nil
}

func enableAutoRun() error {
	key, err := openAutoRunKey(registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.SetStringValue(REG_KEY, REG_VALUE)
}

func disableAutoRun() error {
	key, err := openAutoRunKey(registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.DeleteValue(REG_KEY)
}
