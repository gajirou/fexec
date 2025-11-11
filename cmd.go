package cmd

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gajirou/fexec/pkg/awshelper"
	"github.com/gajirou/fexec/pkg/utils"
)

func Run() error {
	profile := flag.String("p", "default", "利用プロファイル名")
	flag.Parse()

	ssmPlugin := "session-manager-plugin"
	_, err := exec.LookPath(ssmPlugin)
	if err != nil {
		return err
	}

	configService := awshelper.NewConfigService()
	awsConfig, err := configService.FindAWSCredential(*profile)
	if err != nil {
		utils.PrintMessage("ERR002")
		return err
	}
	if awsConfig.Region == "" {
		utils.PrintMessage("INF001")
		return nil
	}

	ecsService := awshelper.EcsService{}
	ecsService.SetEcsClient(awsConfig)
	clusters, err := ecsService.GetClusters()
	if err != nil {
		utils.PrintMessage("ERR003")
		return err
	}
	cluster, err := utils.ScreenDraw(clusters, "clusters")
	if err != nil {
		utils.PrintMessage("ERR999")
		return err
	}
	if cluster == "" {
		utils.PrintMessage("INF002")
		return nil
	}

	services, err := ecsService.GetServices(cluster)
	if err != nil {
		utils.PrintMessage("ERR004")
		return err
	}
	if services == nil {
		utils.PrintMessage("INF003")
		return nil
	}
	service, err := utils.ScreenDraw(services, "services")
	if err != nil {
		utils.PrintMessage("ERR999")
		return err
	}
	if service == "" {
		utils.PrintMessage("INF004")
		return nil
	}

	tasks, err := ecsService.GetTasks(cluster, service)
	if err != nil {
		utils.PrintMessage("ERR005")
		return err
	}
	if tasks == nil {
		utils.PrintMessage("INF005")
		return nil
	}
	task, err := utils.ScreenDraw(tasks, "tasks")
	if err != nil {
		utils.PrintMessage("ERR999")
		return err
	}
	if task == "" {
		utils.PrintMessage("INF006")
		return nil
	}

	containers, err := ecsService.GetContainers(cluster, task)
	if err != nil {
		utils.PrintMessage("ERR006")
		return err
	}
	if containers == nil {
		utils.PrintMessage("INF007")
		return nil
	}
	container, err := utils.ScreenDraw(containers, "containers")
	if err != nil {
		utils.PrintMessage("ERR999")
		return err
	}
	if container == "" {
		utils.PrintMessage("INF008")
		return nil
	}

	execCmd, err := ecsService.ExecuteContainer(cluster, task, container)
	if err != nil {
		utils.PrintMessage("INF009")
		return err
	}
	execSes, err := json.MarshalIndent(execCmd.Session, "", " ")
	if err != nil {
		utils.PrintMessage("ERR999")
		return err
	}
	cmd := exec.Command("session-manager-plugin", string(execSes), awsConfig.Region, "StartSession")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	signal.Ignore(os.Interrupt, syscall.SIGTERM)
	defer signal.Reset(os.Interrupt, syscall.SIGTERM)
	
	if err := cmd.Run(); err != nil {
		utils.PrintMessage("ERR999")
		return err
	}

	return nil
}
