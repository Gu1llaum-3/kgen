package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var podName string
var podImage string

// podCmd represents the pod command
var podCmd = &cobra.Command{
    Use:   "pod",
    Short: "Create a Kubernetes pod file",
    Run: func(cmd *cobra.Command, args []string) {
        if podName == "" || podImage == "" {
            fmt.Println("Both --name and --image flags are required")
            return
        }

        content := fmt.Sprintf(`
apiVersion: v1
kind: Pod
metadata:
  name: %s
spec:
  containers:
  - name: %s
    image: %s
    ports:
    - containerPort: 80
`, podName, podName, podImage)
        
        filename := fmt.Sprintf("%s-pod.yaml", podName)
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

        fmt.Printf("Pod file '%s' created successfully.\n", filename)
    },
}

func init() {
    rootCmd.AddCommand(podCmd)

    // Add flags
    podCmd.Flags().StringVarP(&podName, "name", "n", "", "Name of the pod")
    podCmd.Flags().StringVarP(&podImage, "image", "i", "", "Docker image for the pod")
}
