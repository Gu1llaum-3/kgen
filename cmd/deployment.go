package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var deploymentName string
var deploymentImage string
var deploymentNamespace string
var deploymentTemplate bool

// Template for Kubernetes Deployment
const deploymentTemplateContent = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  replicas: 2
  selector:
    matchLabels:
      app: {{.Name}}-deployment
  template:
    metadata:
      labels:
        app: {{.Name}}
    spec:
      containers:
      - name: {{.Name}}
        image: {{.Image}}
        ports:
        - containerPort: 80
`

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
    Use:   "deployment",
    Short: "Create a Kubernetes deployment file",
    Run: func(cmd *cobra.Command, args []string) {
        //name := getEnv("KGEN_NAME", deploymentName)
        //image := getEnv("KGEN_IMAGE", deploymentImage)
        //namespace := getEnv("KGEN_NAMESPACE", deploymentNamespace)

        // Prioritize flags over environment variables
        name := deploymentName
        if name == "" {
            name = os.Getenv("KGEN_NAME")
        }

        image := deploymentImage
        if image == "" {
            image = os.Getenv("KGEN_IMAGE")
        }

        namespace := deploymentNamespace
        if namespace == "" {
            namespace = os.Getenv("KGEN_NAMESPACE")
        }


        if deploymentTemplate {
            displayDeploymentTemplate("deployment", "nginx", "nginx:latest", "default")
            return
        }

        if name == "" || image == "" {
            fmt.Println("Both --name and --image flags are required")
            return
        }

        createDeploymentFile(name, image, namespace)
    },
}

func createDeploymentFile(name, image, namespace string) {
    tmpl, err := template.New("deployment").Parse(deploymentTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name      string
        Image     string
        Namespace string
    }{
        Name:      name,
        Image:     image,
        Namespace: namespace,
    }

    outputFilename := fmt.Sprintf("%s-deployment.yaml", name)
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

    fmt.Printf("Deployment file '%s' created successfully.\n", outputFilename)
}

func displayDeploymentTemplate(templateType, name, image, namespace string) {
    tmpl, err := template.New("display").Parse(deploymentTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name      string
        Image     string
        Namespace string
    }{
        Name:      name,
        Image:     image,
        Namespace: namespace,
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
    rootCmd.AddCommand(deploymentCmd)

    // Add flags
    deploymentCmd.Flags().StringVarP(&deploymentName, "name", "n", "", "Name of the deployment")
    deploymentCmd.Flags().StringVarP(&deploymentImage, "image", "i", "", "Docker image for the deployment")
    deploymentCmd.Flags().StringVarP(&deploymentNamespace, "namespace", "N", "", "Namespace for the deployment")
    deploymentCmd.Flags().BoolVarP(&deploymentTemplate, "template", "t", false, "Display deployment template")
}
