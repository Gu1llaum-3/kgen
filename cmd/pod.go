package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var podName string
var podImage string
var podNamespace string
var podTemplate bool

// Template for Kubernetes Pod
const podTemplateContent = `
apiVersion: v1
kind: Pod
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  containers:
  - name: {{.Name}}
    image: {{.Image}}
    ports:
    - containerPort: 80
`

// podCmd represents the pod command
var podCmd = &cobra.Command{
    Use:   "pod",
    Short: "Create a Kubernetes pod file",
    Run: func(cmd *cobra.Command, args []string) {
        //name := getEnv("KGEN_NAME", podName)
        //image := getEnv("KGEN_IMAGE", podImage)
        //namespace := getEnv("KGEN_NAMESPACE", podNamespace)
        // Prioritize flags over environment variables
        name := podName
        if name == "" {
            name = os.Getenv("KGEN_NAME")
        }

        image := podImage
        if image == "" {
            image = os.Getenv("KGEN_IMAGE")
        }

        namespace := podNamespace
        if namespace == "" {
            namespace = os.Getenv("KGEN_NAMESPACE")
        }

        if podTemplate {
            displayTemplate("pod", "nginx", "nginx:latest", "default")
            return
        }

        if name == "" || image == "" {
            fmt.Println("Both --name and --image flags are required")
            return
        }

        createPodFile(name, image, namespace)
    },
}

func createPodFile(name, image, namespace string) {
    tmpl, err := template.New("pod").Parse(podTemplateContent)
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

    outputFilename := fmt.Sprintf("%s-pod.yaml", name)
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

    fmt.Printf("Pod file '%s' created successfully.\n", outputFilename)
}

func displayTemplate(templateType, defaultName, defaultImage, defaultNamespace string) {
    tmpl, err := template.New("display").Parse(podTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name      string
        Image     string
        Namespace string
    }{
        Name:      defaultName,
        Image:     defaultImage,
        Namespace: defaultNamespace,
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
    rootCmd.AddCommand(podCmd)

    // Add flags
    podCmd.Flags().StringVarP(&podName, "name", "n", "", "Name of the pod")
    podCmd.Flags().StringVarP(&podImage, "image", "i", "", "Docker image for the pod")
    podCmd.Flags().StringVarP(&podNamespace, "namespace", "N", "", "Namespace for the pod")
    podCmd.Flags().BoolVarP(&podTemplate, "template", "t", false, "Display pod template")
}
