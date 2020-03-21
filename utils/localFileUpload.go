package utils

import (
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/99designs/gqlgen/graphql"
    "github.com/h2non/filetype"
)

// UploadFile is used to upload files.
func UploadFile(upload *graphql.Upload, path string) (string, bool) {
    file := upload.File
    // if file not found, it means that user didn't upload a file(optional to upload)
    fileBytes, _ := ioutil.ReadAll(file)
    fileName := GenerateRandomString(12)
    kind, _ := filetype.Match(fileBytes)

    newPath := filepath.Join(path, fileName+"."+kind.Extension)
    newFile, err := os.Create(newPath)
    if err != nil {
        return "Couldn't create file" + err.Error(), false
    }
    defer newFile.Close() // idempotent, okay to call twice
    // writing image content.
    if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
        return "Couldn't write file" + err.Error(), false
    }
    return newPath, true
}
