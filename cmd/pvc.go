package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// pvcCmd represents the pvc command
var pvcCmd = &cobra.Command{
    Use:   "pvc",
    Short: "Create a Kubernetes PersistentVolumeClaim (PVC) file",
    Run: func(cmd *cobra.Command, args []string) {
        content := `
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
        file, err := os.Create("pvc.yaml")
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

        fmt.Println("PersistentVolumeClaim file created successfully.")
    },
}

func init() {
    rootCmd.AddCommand(pvcCmd)
}
