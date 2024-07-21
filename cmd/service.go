package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var serviceName string
var serviceNodePort int
var serviceLoadBalancer bool
var serviceNamespace string
var serviceTemplate bool

// Template for Kubernetes Service
const serviceTemplateContent = `
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}-service
  namespace: {{.Namespace}}
spec:
  selector:
    app: {{.Name}}
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
    {{- if eq .ServiceType "NodePort" }}
    nodePort: {{.NodePort}}
    {{- end }}
  type: {{.ServiceType}}
`

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
    Use:   "service",
    Short: "Create a Kubernetes service file",
    Run: func(cmd *cobra.Command, args []string) {

        //name := getEnv("KGEN_NAME", serviceName)
        //namespace := getEnv("KGEN_NAMESPACE", serviceNamespace)

        // Prioritize flags over environment variables
        name := serviceName
        if name == "" {
            name = os.Getenv("KGEN_NAME")
        }

        namespace := serviceNamespace
        if namespace == "" {
            namespace = os.Getenv("KGEN_NAMESPACE")
        }

        if serviceTemplate {
            defaultType := "ClusterIP"
            defaultNodePort := 30000

            if serviceNodePort != 0 {
                defaultType = "NodePort"
                defaultNodePort = serviceNodePort
            } else if serviceLoadBalancer {
                defaultType = "LoadBalancer"
            }

            displayServiceTemplate("service", "my-service", defaultType, defaultNodePort, "default")
            return
        }

        if name == "" {
            fmt.Println("The --name flag is required")
            return
        }

        if serviceNodePort != 0 {
            createServiceFile(name, "NodePort", serviceNodePort, namespace)
        } else if serviceLoadBalancer {
            createServiceFile(name, "LoadBalancer", 0, namespace)
        } else {
            createServiceFile(name, "ClusterIP", 0, namespace)
        }
    },
}

func createServiceFile(name, serviceType string, nodePort int, namespace string) {
    tmpl, err := template.New("service").Parse(serviceTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name        string
        ServiceType string
        NodePort    int
        Namespace   string
    }{
        Name:        name,
        ServiceType: serviceType,
        NodePort:    nodePort,
        Namespace:   namespace,
    }

    outputFilename := fmt.Sprintf("%s-service.yaml", name)
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

    fmt.Printf("Service file '%s' created successfully.\n", outputFilename)
}

func displayServiceTemplate(templateType, defaultName, defaultServiceType string, defaultNodePort int, defaultNamespace string) {
    tmpl, err := template.New("display").Parse(serviceTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name        string
        ServiceType string
        NodePort    int
        Namespace   string
    }{
        Name:        defaultName,
        ServiceType: defaultServiceType,
        NodePort:    defaultNodePort,
        Namespace:   defaultNamespace,
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
    rootCmd.AddCommand(serviceCmd)

    // Add flags
    serviceCmd.Flags().StringVarP(&serviceName, "name", "n", "", "Name of the service")
    serviceCmd.Flags().IntVar(&serviceNodePort, "nodeport", 0, "NodePort value (required for NodePort type)")
    serviceCmd.Flags().BoolVar(&serviceLoadBalancer, "loadbalancer", false, "Set service type to LoadBalancer")
    serviceCmd.Flags().StringVarP(&serviceNamespace, "namespace", "N", "", "Namespace for the service")
    serviceCmd.Flags().BoolVarP(&serviceTemplate, "template", "t", false, "Display service template")
}
