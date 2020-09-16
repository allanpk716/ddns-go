package util

import (
	"context"
	"log"
	"os"
	"os/exec"
	"runtime"
)

type OneJob struct {
	Ctx     context.Context
	Cmder   *exec.Cmd
	Cancel  func()
	Err     error
	CmdLine string
	CmdArgs []string
	Running bool
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func InitFrpArgs(nowdir string, oneJob_s *OneJob, oneJob_c *OneJob) bool {
	nowdir = nowdir + string(os.PathSeparator) + "frpThings" + string(os.PathSeparator)
	absPath_Frps := nowdir + "frps"
	absPath_Frpc := nowdir + "frpc"
	absPath_Frps_ini := nowdir + "frps.ini"
	absPath_Frpc_ini := nowdir + "frpc.ini"

	switch runtime.GOOS {
	case "darwin":
		break
	case "linux":
		break
	case "windows":
		absPath_Frps += ".exe"
		absPath_Frpc += ".exe"
	}

	if FileExist(absPath_Frps) == false {
		log.Panicln(absPath_Frps + " not exist.")
		return false
	}
	if FileExist(absPath_Frpc) == false {
		log.Panicln(absPath_Frpc + " not exist.")
		return false
	}
	if FileExist(absPath_Frps_ini) == false {
		log.Panicln(absPath_Frps_ini + " not exist.")
		return false
	}
	if FileExist(absPath_Frpc_ini) == false {
		log.Panicln(absPath_Frpc_ini + " not exist.")
		return false
	}

	oneJob_s.CmdLine = absPath_Frps
	oneJob_s.CmdArgs = []string{
		"-c",
		absPath_Frps_ini,
	}
	oneJob_c.CmdLine = absPath_Frpc
	oneJob_c.CmdArgs = []string{
		"-c",
		absPath_Frpc_ini,
	}

	return true
}

func StartFrpThings(oneJob_s *OneJob, oneJob_c *OneJob) bool {

	if oneJob_s.Running == true && oneJob_c.Running == true {
		return false
	}

	log.Printf("Start frps ...")
	startFrp(oneJob_s)
	if oneJob_s.Err != nil {
		log.Printf("Start frps Error")
		log.Panicln(oneJob_s.Err.Error())
		return false
	}
	log.Printf("Start frps Done.")

	log.Printf("Start frpc ...")
	startFrp(oneJob_c)
	if oneJob_c.Err != nil {
		log.Printf("Start frpc Error")
		log.Panicln(oneJob_c.Err.Error())

		log.Printf("Close frps ...")
		CloseFrp(oneJob_s)
		log.Printf("Close frps Done.")
		return false
	}
	log.Printf("Start frpc Done.")

	oneJob_s.Running = true
	oneJob_c.Running = true

	return true
}

func CloseFrp(oneJob *OneJob) bool {
	oneJob.Cancel()
	_ = oneJob.Cmder.Wait()
	oneJob.Running = false
	return true
}

func startFrp(oneJob *OneJob) {
	oneJob.Ctx, oneJob.Cancel = context.WithCancel(context.Background())
	oneJob.Cmder = exec.CommandContext(oneJob.Ctx, oneJob.CmdLine, oneJob.CmdArgs...)
	// oneJob.Cmder.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	oneJob.Cmder.Stdout = os.Stdout
	err := oneJob.Cmder.Start()
	if err != nil {
		oneJob.Err = err
	}
}
