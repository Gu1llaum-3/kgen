package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var deploymentName string
var deploymentImage string

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
    Use:   "deployment",
    Short: "Create a Kubernetes deployment file",
    Run: func(cmd *cobra.Command, args []string) {
        if deploymentName == "" || deploymentImage == "" {
            fmt.Println("Both --name and --image flags are required")
            return
        }

        content := fmt.Sprintf(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
spec:
  replicas: 2
  selector:
    matchLabels:
      app: %s-deployment
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: %s
        image: %s
        ports:
        - containerPort: 80
`, deploymentName, deploymentName, deploymentName, deploymentName, deploymentImage)
        
        filename := fmt.Sprintf("%s-deployment.yaml", deploymentName)
        file, err := os.Create(filename)
        if err != nil {
            fmt.Println("Error creating file:", err)
            return
        }
        defer file.Close()

        _, err = file.WriteString(content)
        if err != nil {
            fmt.Println("Error writing to file:", err)
            return
        }

        fmt.Printf("Deployment file '%s' created successfully.\n", filename)
    },
}

func init() {
    rootCmd.AddCommand(deploymentCmd)

    // Add flags
    deploymentCmd.Flags().StringVarP(&deploymentName, "name", "n", "", "Name of the deployment")
    deploymentCmd.Flags().StringVarP(&deploymentImage, "image", "i", "", "Docker image for the deployment")
}
