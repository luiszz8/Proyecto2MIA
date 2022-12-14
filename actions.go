package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getHola(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mainArchivos()
	json.NewEncoder(w).Encode("comando Exitoso")
}

type inicioS struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	IdG      string `json:"Id"`
}

type Instrucion struct {
	Comando string `json:"Comando"`
}

var lgn = logn{}

type logn []inicioS

func Login(w http.ResponseWriter, r *http.Request) {
	var lgn = logn{}
	var inicio inicioS
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error datos no validos")
	}
	json.Unmarshal(reqBody, &inicio) //La función Unmarshal() en la codificación del paquete / json se utiliza para descomprimir o decodificar los datos de JSON a la estructura
	lgn = append(lgn, inicio)
	fmt.Println(inicio)
	fmt.Println("USUARIO: ", inicio.Username)
	inicio.Username = strings.Split(inicio.Username, "\r")[0]
	inicio.Password = strings.Split(inicio.Password, "\r")[0]
	inicio.IdG = strings.Split(inicio.IdG, "\r")[0]
	fmt.Println("PASSWORD: ", inicio.Password)
	Bandera := login("Login -password=" + inicio.Password + " -usuario=" + inicio.Username + " -id=" + inicio.IdG)
	w.Header().Set("Content-type", "application/json")
	Mensaje := "Crendeciales Incorrectas"
	if Bandera {
		Mensaje = "Bienvenido"
	}
	json.NewEncoder(w).Encode(Mensaje)

}

func Comandos(w http.ResponseWriter, r *http.Request) {
	var comand Instrucion
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error datos no validos")
	}
	json.Unmarshal(reqBody, &comand) //La función Unmarshal() en la codificación del paquete / json se utiliza para descomprimir o decodificar los datos de JSON a la estructura
	regreso := ""
	regreso = ejecutar(comand.Comando)

	w.Header().Set("Content-type", "application/json")
	if regreso == "" {
		json.NewEncoder(w).Encode("comando Exitoso")
	} else {
		json.NewEncoder(w).Encode(regreso)
	}

}

type Archivo struct {
	Contenido string `json:"Username"`
}

func Carga(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error datos no validos")
	}
	content := string(reqBody)
	fmt.Println(content)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode("nice")
}
