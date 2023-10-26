package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "github.com/tidwall/gjson"
)

const (
    API_KEY   = "2b10I9c4DHkPJWMzrvYIDCOTu"
    PROJECT   = "all"
    apiURL    = "https://my-api.plantnet.org/v2/identify/%s?lang=pt-br&api-key=%s"
)

type PlantResponse struct {
    CommonNames    string
    ScientificName string
    Score          float64
}

func RecognizePlant(image multipart.File) (string, error) {


    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    data := map[string][]string{
        "organs": {"leaf"},
    }

    for key, values := range data {
        for _, value := range values {
            _ = writer.WriteField(key, value)
        }
    }

    part, err := writer.CreateFormFile("images", "image.jpg")  // Set a default file name
    if err != nil {
        return "", fmt.Errorf("error creating form file: %v", err)
    }

    _, err = io.Copy(part, image)
    if err != nil {
        return "", fmt.Errorf("error copying file: %v", err)
    }

    writer.Close()

    req, err := http.NewRequest("POST", fmt.Sprintf(apiURL, PROJECT, API_KEY), body)
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&result)
    if err != nil {
        return "", fmt.Errorf("error decoding JSON: %v", err)
    }

    return FormatResponse(result)
}


func FormatResponse(result map[string]interface{}) (string, error) {
    prettyJSON, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        return "", fmt.Errorf("error formatting JSON: %v", err)
    }

    content := prettyJSON

    bestMatch := gjson.Get(string(content), "bestMatch")
    commonNames := gjson.Get(string(content), "results.0.species.commonNames")


    finalResult := fmt.Sprintf("Best Match: %s\nCommon Names:\n", bestMatch)
    for _, r := range commonNames.Array() {
        stringor := r.String()
        finalResult += fmt.Sprintf("- %s\n", stringor)
    }

    return finalResult, nil
}
