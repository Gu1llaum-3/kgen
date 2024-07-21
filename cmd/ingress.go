package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var ingressName string
var ingressNamespace string
var ingressHost string
var ingressServiceName string
var ingressServicePort int
var ingressTemplate bool

// Template for Kubernetes Ingress
const ingressTemplateContent = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  rules:
  - host: {{.Host}}
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{.ServiceName}}
            port:
              number: {{.ServicePort}}
`

// ingressCmd represents the ingress command
var ingressCmd = &cobra.Command{
    Use:   "ingress",
    Short: "Create a Kubernetes Ingress file",
    Run: func(cmd *cobra.Command, args []string) {
        if ingressTemplate {
            displayIngressTemplate("ingress", "my-ingress", "default", "example.com", "my-service", 80)
            return
        }

        if ingressName == "" || ingressHost == "" || ingressServiceName == "" || ingressServicePort == 0 {
            fmt.Println("The --name, --host, --service-name, and --service-port flags are required")
            return
        }

        createIngressFile(ingressName, ingressNamespace, ingressHost, ingressServiceName, ingressServicePort)
    },
}

func createIngressFile(name, namespace, host, serviceName string, servicePort int) {
    tmpl, err := template.New("ingress").Parse(ingressTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name        string
        Namespace   string
        Host        string
        ServiceName string
        ServicePort int
    }{
        Name:        name,
        Namespace:   namespace,
        Host:        host,
        ServiceName: serviceName,
        ServicePort: servicePort,
    }

    outputFilename := fmt.Sprintf("%s-ingress.yaml", name)
    file, err := os.Create(outputFilename)
    if err != nil {
        fmt.Printf("Error creating file '%s': %v\n", outputFilename, err)
        return
    }
    defer file.Close()

    err = tmpl.Execute(file, data)
    if err != nil {
        fmt.Printf("Error executing template: %v\n", err)
        return
    }

    fmt.Printf("Ingress file '%s' created successfully.\n", outputFilename)
}

func displayIngressTemplate(templateType, defaultName, defaultNamespace, defaultHost, defaultServiceName string, defaultServicePort int) {
    tmpl, err := template.New("display").Parse(ingressTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name        string
        Namespace   string
        Host        string
        ServiceName string
        ServicePort int
    }{
        Name:        defaultName,
        Namespace:   defaultNamespace,
        Host:        defaultHost,
        ServiceName: defaultServiceName,
        ServicePort: defaultServicePort,
    }

    var result bytes.Buffer
    err = tmpl.Execute(&result, data)
    if err != nil {
        fmt.Printf("Error executing template: %v\n", err)
        return
    }

    fmt.Printf("Template for %s:\n\n%s\n", templateType, result.String())
}

func init() {
    rootCmd.AddCommand(ingressCmd)

    // Add flags
    ingressCmd.Flags().StringVarP(&ingressName, "name", "n", "", "Name of the Ingress")
    ingressCmd.Flags().StringVarP(&ingressNamespace, "namespace", "N", "default", "Namespace for the Ingress")
    ingressCmd.Flags().StringVarP(&ingressHost, "host", "H", "", "Host for the Ingress")
    ingressCmd.Flags().StringVarP(&ingressServiceName, "service-name", "s", "", "Service name for the Ingress backend")
    ingressCmd.Flags().IntVarP(&ingressServicePort, "service-port", "p", 80, "Service port for the Ingress backend")
    ingressCmd.Flags().BoolVarP(&ingressTemplate, "template", "t", false, "Display Ingress template")
}
