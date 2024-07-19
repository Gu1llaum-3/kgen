package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
    Use:   "service",
    Short: "Create a Kubernetes service file",
    Run: func(cmd *cobra.Command, args []string) {
        content := `
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
        file, err := os.Create("service.yaml")
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

        fmt.Println("Service file created successfully.")
    },
}

func init() {
    rootCmd.AddCommand(serviceCmd)
}
