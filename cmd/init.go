package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Create default Kubernetes files (deployment, service, pv, pvc)",
    Run: func(cmd *cobra.Command, args []string) {
        createFile("deployment.yaml", getDeploymentContent())
        createFile("service.yaml", getServiceContent())
        createFile("pv.yaml", getPVContent())
        createFile("pvc.yaml", getPVCContent())
        fmt.Println("Default Kubernetes files created successfully.")
    },
}

func createFile(filename, content string) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(content)
    if err != nil {
        fmt.Println("Error writing to file:", err)
    }
}

func getDeploymentContent() string {
    return `
apiVersion: apps/v1
kind: Deployment
metadata:
name: my-deployment
spec:
replicas: 2
selector:
 matchLabels:
   app: my-app
template:
 metadata:
   labels:
     app: my-app
 spec:
   containers:
   - name: my-container
     image: my-image
    `
}

func getServiceContent() string {
    return `
apiVersion: v1
kind: Service
metadata:
name: my-service
spec:
selector:
 app: my-app
ports:
- protocol: TCP
 port: 80
 targetPort: 80
    `
}

func getPVContent() string {
    return `
apiVersion: v1
kind: PersistentVolume
metadata:
name: my-pv
spec:
capacity:
 storage: 1Gi
accessModes:
 - ReadWriteOnce
hostPath:
 path: "/mnt/data"
    `
}

func getPVCContent() string {
    return `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
name: my-pvc
spec:
accessModes:
 - ReadWriteOnce
resources:
 requests:
   storage: 1Gi
    `
}

func init() {
    rootCmd.AddCommand(initCmd)
}
