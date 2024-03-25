package main

import (
	"fmt"
	"os"
)

type Biodata struct {
	Nama       string
	Alamat    string
	Pekerjaan string
	Alasan     string
}

var biodataFriends = []Biodata{
 	{
 		Nama: "Rahmi",
 		Alamat: "Malang",
 		Pekerjaan: "Frontend Engineer",
 		Alasan: "Untuk meningkatkan pemahaman mengenai cara kerja backend dengan lebih baik",
 	},
 	{
 		Nama: "Asma",
 		Alamat: "Banjarmasin",
 		Pekerjaan: "Mahasiswa",
 		Alasan: "Ingin mempelajari bahasa pemrograman yang relevan dengan industri teknologi saat ini",
 	},
 	{
 		Nama: "Tom",
 		Alamat: "Jakarta",
 		Pekerjaan: "Web Developer",
 		Alasan: "Ingin mempelajari bahasa pemrograman yang efisien dan cocok untuk pengembangan web",
 	},
 	{
 		Nama: "Anya",
 		Alamat: "Surabaya",
 		Pekerjaan: "Data Scientist",
 		Alasan: "Ingin mempelajari bahasa pemrograman yang efisien untuk analisis data dan machine learning",
 	},
}

func main() {
	args := os.Args
	number := args[1]
	ShowData(number)
}

func ShowData(number string) {
	var num int
	_, err := fmt.Sscanf(number, "%d", &num)

	if err != nil {
		fmt.Println("Harus integer")
		return
	}

	if num < 1 || num > len(biodataFriends) {
		fmt.Println("Tidak ada data")
		return
	}

	biodata := biodataFriends[num-1]
	fmt.Println("Nama:", biodata.Nama)
	fmt.Println("Alamat:", biodata.Alamat)
	fmt.Println("Pekerjaan:", biodata.Pekerjaan)
	fmt.Println("Alasan memilih kelas Golang:", biodata.Alasan)
}