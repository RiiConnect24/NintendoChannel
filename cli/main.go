package main

import (
    "NintendoChannel/csdata"
    "NintendoChannel/dllist"
    "NintendoChannel/thumbnail"
    "fmt"
    "os"
    "strconv"
)

func main() {
    fmt.Println("WiiLink Nintendo Channel File Generator")
    fmt.Println()

    if len(os.Args) < 2 {
        fmt.Println("Usage: ", os.Args[0], " <operation>")
        fmt.Println("Available operations:")
        fmt.Println("1 - DLList and game info")
		fmt.Println("2 - DLList and game info (force)")
        fmt.Println("3 - Thumbnails")
        fmt.Println("4 - CSData")
        return
    }

    selection, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Invalid input. Please provide a valid numeric operation.")
        return
    }

    switch selection {
    case 1:
        dllist.MakeDownloadList(false)
    case 2:
        dllist.MakeDownloadList(true)
    case 3:
        thumbnail.WriteThumbnail()
    case 4:
        csdata.CreateCSData()
    default:
        fmt.Println("\nInvalid Selection")
    }
}