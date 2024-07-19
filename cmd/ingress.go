package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// ingressCmd represents the ingress command
var ingressCmd = &cobra.Command{
    Use:   "ingress",
    Short: "Create a Kubernetes ingress file",
    Run: func(cmd *cobra.Command, args []string) {
        content := `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
spec:
  rules:
  - host: my-app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: my-service
            port:
              number: 80
`
        file, err := os.Create("ingress.yaml")
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

        fmt.Println("Ingress file created successfully.")
    },
}

func init() {
    rootCmd.AddCommand(ingressCmd)
}
