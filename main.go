package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

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

func main() {
	fortnitePath, _ := os.Getwd()
	binariesPath := filepath.Join(fortnitePath, "FortniteGame\\Binaries\\Win64\\")
	launchArgs := "-epicapp=Fortnite -epicenv=Prod -epiclocale=en-us -epicportal -skippatchcheck -NOSSLPINNING -nobe -HTTP=WinINet"
	//launcherExe := filepath.Join(binariesPath, "FortniteLauncher.exe")
	//eacExe := filepath.Join(binariesPath, "FortniteClient-Win64-Shipping_EAC.exe")
	shippingExe := filepath.Join(binariesPath, "FortniteClient-Win64-Shipping.exe")

	cobaltDllPath := filepath.Join(fortnitePath, "Cobalt.dll")

	//eacCmd := exec.Command(eacExe, launchArgs)
	//launcherCmd := exec.Command(launcherExe, launchArgs)
	shippingCmd := exec.Command(shippingExe, launchArgs)

	shippingCmd.Start()
	err := injectDll(uint32(shippingCmd.Process.Pid), cobaltDllPath)
	if err != nil {
		print(err)
	}
}
