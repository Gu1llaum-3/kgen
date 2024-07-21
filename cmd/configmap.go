package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var configMapName string
var configMapNamespace string
var configMapTemplate bool

// Template for Kubernetes ConfigMap
const configMapTemplateContent = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
data:
  key1: value1
  key2: value2
`

// configMapCmd represents the configmap command
var configMapCmd = &cobra.Command{
    Use:   "configmap",
    Short: "Create a Kubernetes ConfigMap file",
    Run: func(cmd *cobra.Command, args []string) {
        if configMapTemplate {
            displayConfigMapTemplate("configmap", "my-configmap", "default")
            return
        }

        if configMapName == "" {
            fmt.Println("The --name flag is required")
            return
        }

        createConfigMapFile(configMapName, configMapNamespace)
    },
}

func createConfigMapFile(name, namespace string) {
    tmpl, err := template.New("configmap").Parse(configMapTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name      string
        Namespace string
    }{
        Name:      name,
        Namespace: namespace,
    }

    outputFilename := fmt.Sprintf("%s-configmap.yaml", name)
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

    fmt.Printf("ConfigMap file '%s' created successfully.\n", outputFilename)
}

func displayConfigMapTemplate(templateType, defaultName, defaultNamespace string) {
    tmpl, err := template.New("display").Parse(configMapTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name      string
        Namespace string
    }{
        Name:      defaultName,
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
    rootCmd.AddCommand(configMapCmd)

    // Add flags
    configMapCmd.Flags().StringVarP(&configMapName, "name", "n", "", "Name of the ConfigMap")
    configMapCmd.Flags().StringVarP(&configMapNamespace, "namespace", "N", "default", "Namespace for the ConfigMap")
    configMapCmd.Flags().BoolVarP(&configMapTemplate, "template", "t", false, "Display ConfigMap template")
}
