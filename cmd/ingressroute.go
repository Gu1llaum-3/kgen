package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var ingressRouteName string
var ingressRouteNamespace string
var ingressRouteHost string
var ingressRouteServiceName string
var ingressRouteServicePort int
var ingressRouteTemplate bool

// Template for Traefik IngressRoute
const ingressRouteTemplateContent = `
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  entryPoints:
    - web
  routes:
  - match: Host('{{.Host}}')
    kind: Rule
    services:
    - name: {{.ServiceName}}
      port: {{.ServicePort}}
`

// ingressRouteCmd represents the ingressroute command
var ingressRouteCmd = &cobra.Command{
    Use:   "ingressroute",
    Short: "Create a Traefik IngressRoute file",
    Run: func(cmd *cobra.Command, args []string) {
        if ingressRouteTemplate {
            displayIngressRouteTemplate("ingressroute", "my-ingressroute", "default", "example.com", "my-ingressroute-service", 80)
            return
        }

        if ingressRouteName == "" || ingressRouteHost == "" || ingressRouteServicePort == 0 {
            fmt.Println("The --name, --host, and --service-port flags are required")
            return
        }

        if ingressRouteServiceName == "" {
            ingressRouteServiceName = fmt.Sprintf("%s-service", ingressRouteName)
        }

        createIngressRouteFile(ingressRouteName, ingressRouteNamespace, ingressRouteHost, ingressRouteServiceName, ingressRouteServicePort)
    },
}

func createIngressRouteFile(name, namespace, host, serviceName string, servicePort int) {
    tmpl, err := template.New("ingressroute").Parse(ingressRouteTemplateContent)
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

    outputFilename := fmt.Sprintf("%s-ingressroute.yaml", name)
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

    fmt.Printf("IngressRoute file '%s' created successfully.\n", outputFilename)
}

func displayIngressRouteTemplate(templateType, defaultName, defaultNamespace, defaultHost, defaultServiceName string, defaultServicePort int) {
    tmpl, err := template.New("display").Parse(ingressRouteTemplateContent)
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
    rootCmd.AddCommand(ingressRouteCmd)

    // Add flags
    ingressRouteCmd.Flags().StringVarP(&ingressRouteName, "name", "n", "", "Name of the IngressRoute")
    ingressRouteCmd.Flags().StringVarP(&ingressRouteNamespace, "namespace", "N", "default", "Namespace for the IngressRoute")
    ingressRouteCmd.Flags().StringVarP(&ingressRouteHost, "host", "H", "", "Host for the IngressRoute")
    ingressRouteCmd.Flags().StringVarP(&ingressRouteServiceName, "service-name", "s", "", "Service name for the IngressRoute backend")
    ingressRouteCmd.Flags().IntVarP(&ingressRouteServicePort, "service-port", "p", 80, "Service port for the IngressRoute backend")
    ingressRouteCmd.Flags().BoolVarP(&ingressRouteTemplate, "template", "t", false, "Display IngressRoute template")
}
