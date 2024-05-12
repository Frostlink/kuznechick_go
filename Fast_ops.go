package main

import "fmt"

func fast_operations() {
	//record()
	converter()
	fmt.Println("Progress 10%")
	fmt.Println("Converted done")
	upload_to_storage()
	fmt.Println("Progress 20%")
	fmt.Println("Upload to storage done")
	recognize()
	fmt.Println("Progress 70%")
	fmt.Println("Recognize audio done")
	Text_encrypt()
	fmt.Println("Progress 100%")
	fmt.Println("Text Encrypted")
	Text_decrypt()
	fmt.Println("Text Decrypted")

}
