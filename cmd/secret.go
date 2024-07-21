package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var secretName string
var secretNamespace string
var secretTemplate bool

// Template for Kubernetes Secret
const secretTemplateContent = `
apiVersion: v1
kind: Secret
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
type: Opaque
data:
  # Example of encoded data:
  # username: bXl1c2Vy
  # password: bXlwYXNzd29yZA==
  key1: value1
  key2: value2
`

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
    Use:   "secret",
    Short: "Create a Kubernetes Secret file",
    Run: func(cmd *cobra.Command, args []string) {
        if secretTemplate {
            displaySecretTemplate("secret", "my-secret", "default")
            return
        }

        if secretName == "" {
            fmt.Println("The --name flag is required")
            return
        }

        createSecretFile(secretName, secretNamespace)
    },
}

func createSecretFile(name, namespace string) {
    tmpl, err := template.New("secret").Parse(secretTemplateContent)
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

    outputFilename := fmt.Sprintf("%s-secret.yaml", name)
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

    fmt.Printf("Secret file '%s' created successfully.\n", outputFilename)
}

func displaySecretTemplate(templateType, defaultName, defaultNamespace string) {
    tmpl, err := template.New("display").Parse(secretTemplateContent)
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
    rootCmd.AddCommand(secretCmd)

    // Add flags
    secretCmd.Flags().StringVarP(&secretName, "name", "n", "", "Name of the Secret")
    secretCmd.Flags().StringVarP(&secretNamespace, "namespace", "N", "default", "Namespace for the Secret")
    secretCmd.Flags().BoolVarP(&secretTemplate, "template", "t", false, "Display Secret template")
}
