package cmd

import (
    "bytes"
    "fmt"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var pvcName string
var pvcNamespace string
var pvcTemplate bool
var pvcStorageSize string
var pvcAccessMode string
var pvcStorageClassName string

// Template for Kubernetes PersistentVolumeClaim (PVC)
const pvcTemplateContent = `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{.Name}}-pvc
  namespace: {{.Namespace}}
spec:
  accessModes:
    - {{.AccessMode}}
  resources:
    requests:
      storage: {{.StorageSize}}
  {{- if .StorageClassName }}
  storageClassName: {{.StorageClassName}}
  {{- end }}
`

// pvcCmd represents the pvc command
var pvcCmd = &cobra.Command{
    Use:   "pvc",
    Short: "Create a Kubernetes PersistentVolumeClaim (PVC) file",
    Run: func(cmd *cobra.Command, args []string) {
        if pvcTemplate {
            displayPVCTemplate("pvc", "my-pvc", "default", "1Gi", "ReadWriteOnce", "")
            return
        }

        if pvcName == "" {
            fmt.Println("The --name flag is required")
            return
        }

        createPVCFile(pvcName, pvcNamespace, pvcStorageSize, pvcAccessMode, pvcStorageClassName)
    },
}

func createPVCFile(name, namespace, storageSize, accessMode, storageClassName string) {
    tmpl, err := template.New("pvc").Parse(pvcTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name            string
        Namespace       string
        StorageSize     string
        AccessMode      string
        StorageClassName string
    }{
        Name:            name,
        Namespace:       namespace,
        StorageSize:     storageSize,
        AccessMode:      accessMode,
        StorageClassName: storageClassName,
    }

    outputFilename := fmt.Sprintf("%s-pvc.yaml", name)
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

    fmt.Printf("PersistentVolumeClaim file '%s' created successfully.\n", outputFilename)
}

func displayPVCTemplate(templateType, defaultName, defaultNamespace, defaultStorageSize, defaultAccessMode, defaultStorageClassName string) {
    tmpl, err := template.New("display").Parse(pvcTemplateContent)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        return
    }

    data := struct {
        Name            string
        Namespace       string
        StorageSize     string
        AccessMode      string
        StorageClassName string
    }{
        Name:            defaultName,
        Namespace:       defaultNamespace,
        StorageSize:     defaultStorageSize,
        AccessMode:      defaultAccessMode,
        StorageClassName: defaultStorageClassName,
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
    rootCmd.AddCommand(pvcCmd)

    // Add flags
    pvcCmd.Flags().StringVarP(&pvcName, "name", "n", "", "Name of the PersistentVolumeClaim")
    pvcCmd.Flags().StringVarP(&pvcNamespace, "namespace", "N", "default", "Namespace for the PersistentVolumeClaim")
    pvcCmd.Flags().BoolVarP(&pvcTemplate, "template", "t", false, "Display PersistentVolumeClaim template")
    pvcCmd.Flags().StringVarP(&pvcStorageSize, "storagesize", "s", "1Gi", "Storage size for the PersistentVolumeClaim")
    pvcCmd.Flags().StringVarP(&pvcAccessMode, "accessmode", "a", "ReadWriteOnce", "Access mode for the PersistentVolumeClaim")
    pvcCmd.Flags().StringVarP(&pvcStorageClassName, "storageclassname", "c", "", "Storage class name for the PersistentVolumeClaim")
}
