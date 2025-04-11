package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/windows"
)

// Taken from FNLauncher by DottoXD
func getPid(name string) uint32 {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0
	}
	processEntry := windows.ProcessEntry32{Size: 568}
	for {
		err := windows.Process32Next(handle, &processEntry)
		if err != nil {
			return 0
		}

		if windows.UTF16ToString(processEntry.ExeFile[:]) == name {
			return processEntry.ProcessID
		}
	}
}

// Taken from FNLauncher by DottoXD
func suspendProcess(pid uint32) error {
	handle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, pid)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	if r1, _, _ := windows.NewLazySystemDLL("ntdll.dll").NewProc("NtSuspendProcess").Call(uintptr(handle)); r1 != 0 {
		return fmt.Errorf("NtStatus='0x%.8X'", r1)
	}
	return nil
}

func injectDll(pid uint32, path string) error {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")
	fmt.Println("Got all procedures.")

	handle, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, pid)
	if err != nil {
		return err
	}
	print("Opened process!")

	r1, _, _ := VirtualAllocEx.Call(uintptr(handle), 0, uintptr(len(path)+1), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if r1 == 0 {
		return fmt.Errorf("VirtualAllocEx failed.")
	}
	print("VirtualAllocEx")

	bPtr, err := windows.BytePtrFromString(path)
	if err != nil {
		return err
	}
	print("BytePtrFromString")

	zero := uintptr(0)
	err = windows.WriteProcessMemory(handle, r1, bPtr, uintptr(len(path)+1), &zero)
	if err != nil {
		return err
	}
	print("WriteProcessMemory")

	LoadLibAddy, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
	if err != nil {
		return err
	}
	print("LoadLibAddy")

	tHandle, _, err := CreateRemoteThreadEx.Call(uintptr(handle), 0, 0, LoadLibAddy, r1, 0, 0)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(syscall.Handle(tHandle))
	print("CreateRemoteThreadEx")

	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func launchGame(fortnitePath string) {
	binariesPath := filepath.Join(fortnitePath, "FortniteGame\\Binaries\\Win64\\")
	launchArgs := "-epicapp=Fortnite -epicenv=Prod -epiclocale=en-us -epicportal -skippatchcheck -NOSSLPINNING -nobe"

	bServer := false
	if len(os.Args) > 1 {
		for i := 0; i < len(os.Args); i++ {
			slice := os.Args[i]
			if slice == "--server" {
				bServer = true
				os.Args = remove(os.Args, i)
			}
		}

		launchArgs += strings.Join(os.Args, " ")
	}

	cobaltDllPath := filepath.Join(fortnitePath, "Cobalt.dll")
	rebootDllPath := filepath.Join(fortnitePath, "Reboot.dll")

	shippingExe := filepath.Join(binariesPath, "FortniteClient-Win64-Shipping.exe")
	shippingCmd := exec.Command(shippingExe, launchArgs)

	shippingCmd.Start()
	err := injectDll(uint32(shippingCmd.Process.Pid), cobaltDllPath)
	if err != nil {
		fmt.Println(err)
	}

	if bServer {
		fmt.Println("Press Enter to inject Reboot!")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		err := injectDll(uint32(shippingCmd.Process.Pid), rebootDllPath)
		if err != nil {
			fmt.Println(err)
		}
	}

	for {
		time.Sleep(1 * time.Second)

		pid := getPid("FortniteClient-Win64-Shipping.exe")
		_, err := os.FindProcess(int(pid))

		if err != nil {
			break
		}
	}
}

func main() {
	localAppData := os.Getenv("LOCALAPPDATA")
	fortressAppData := filepath.Join(localAppData, ".FortressLauncher")
	if _, err := os.Stat(fortressAppData); os.IsNotExist(err) {
		os.Mkdir(fortressAppData, fs.FileMode(os.O_CREATE))
	}

	a := app.New()
	w := a.NewWindow("Fortress Launcher")

	mainContent := makePlayContent()
	sidebar := container.NewVBox(
		widget.NewLabelWithStyle("Fortress Launcher", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewButton("Play", func() { changePages(mainContent, makePlayContent()) }),
		widget.NewButton("Options", func() {}),
		widget.NewButton("Mods", func() {}),
		widget.NewButton("Exit", func() { changePages(mainContent, makeExitContent(w)) }),
	)
	split := container.NewHSplit(sidebar, mainContent)

	w.SetContent(split)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}
