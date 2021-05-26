package testenv

import (
	"fmt"
	"strings"

	enterprisev1 "github.com/splunk/splunk-operator/pkg/apis/enterprise/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// GenerateAppSourceSpec return AppSourceSpec struct with given values
func GenerateAppSourceSpec(appSourceName string, appSourceLocation string, appSourceDefaultSpec enterprisev1.AppSourceDefaultSpec) enterprisev1.AppSourceSpec {
	return enterprisev1.AppSourceSpec{
		Name:                 appSourceName,
		Location:             appSourceLocation,
		AppSourceDefaultSpec: appSourceDefaultSpec,
	}
}

// GetPodAppStatus Get the app install status and version number
func GetPodAppStatus(deployment *Deployment, podName string, ns string, appname string) (string, string, error) {
	output := GetPodAppInstallStatus(deployment, podName, ns, appname)
	status := strings.Fields(output)[1]
	path := "etc/apps"
	if strings.Contains(podName, "cluster-master") {
		path = "etc/master-apps/_cluster"
	}
	if strings.Contains(podName, "-deployer-") {
		path = "etc/shcluster/apps"
	}
	filePath := fmt.Sprintf("/opt/splunk/%s/%s/default/app.conf", path, appname)

	confline, err := GetConfLineFromPod(podName, filePath, ns, "Version", "", false)
	if err != nil {
		logf.Log.Error(err, "Failed to get Version from pod", "Pod Name", podName)
		return "", "", err
	}
	version := strings.TrimSpace(strings.Split(confline, " ")[2])
	return status, version, err

}

// GetPodAppInstallStatus Get the app install status
func GetPodAppInstallStatus(deployment *Deployment, podName string, ns string, appname string) string {
	stdin := fmt.Sprintf("/opt/splunk/bin/splunk display app '%s'", appname)
	command := []string{"/bin/sh"}
	stdout, stderr, err := deployment.PodExecCommand(podName, command, stdin, false)
	if err != nil {
		logf.Log.Error(err, "Failed to execute command on pod", "pod", podName, "command", command, "stdin", stdin)
		return "Failed"
	}
	logf.Log.Info("Command executed on pod", "pod", podName, "command", command, "stdin", stdin, "stdout", stdout, "stderr", stderr)

	return strings.TrimSuffix(stdout, "\n")
}
