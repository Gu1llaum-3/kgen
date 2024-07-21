package cmd

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var namespaceName string
var namespaceTemplate bool

// namespaceCmd represents the namespace command
var namespaceCmd = &cobra.Command{
    Use:   "namespace",
    Short: "Create a Kubernetes Namespace file",
    Run: func(cmd *cobra.Command, args []string) {
        if namespaceTemplate {
            displayNamespaceTemplate("namespace", "my-namespace")
            return
        }

        if namespaceName == "" {
            fmt.Println("The --name flag is required")
            return
        }

        createNamespaceFile(namespaceName)
    },
}

func createNamespaceFile(name string) {
    filename := fmt.Sprintf("templates/namespace.yaml")
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file '%s': %v\n", filename, err)
        return
    }

    tmpl, err := template.New("namespace").Parse(string(content))
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name string
    }{
        Name: name,
    }

    outputFilename := fmt.Sprintf("%s-namespace.yaml", name)
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

    fmt.Printf("Namespace file '%s' created successfully.\n", outputFilename)
}

func displayNamespaceTemplate(templateType, defaultName string) {
    filename := fmt.Sprintf("templates/%s.yaml", templateType)
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file '%s': %v\n", filename, err)
        return
    }

    tmpl, err := template.New("display").Parse(string(content))
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name string
    }{
        Name: defaultName,
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
    rootCmd.AddCommand(namespaceCmd)

    // Add flags
    namespaceCmd.Flags().StringVarP(&namespaceName, "name", "n", "", "Name of the Namespace")
    namespaceCmd.Flags().BoolVarP(&namespaceTemplate, "template", "t", false, "Display Namespace template")
}
