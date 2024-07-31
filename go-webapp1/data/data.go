package data

import (
    "encoding/json"
    "os"
    "sync"
)

var mu sync.Mutex

type Item struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

func SaveToFile(items []Item, filename string) error {
    mu.Lock()
    defer mu.Unlock()

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(items)
}
