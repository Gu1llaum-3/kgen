package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// pvCmd represents the pv command
var pvCmd = &cobra.Command{
    Use:   "pv",
    Short: "Create a Kubernetes PersistentVolume (PV) file",
    Run: func(cmd *cobra.Command, args []string) {
        content := `
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
        file, err := os.Create("pv.yaml")
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

        fmt.Println("PersistentVolume file created successfully.")
    },
}

func init() {
    rootCmd.AddCommand(pvCmd)
}
