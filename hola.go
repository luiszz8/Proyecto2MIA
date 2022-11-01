package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
Esta estructura nos ayuda a llevar el control
sobre cual es el primer espacio libre donde se puede agregar información,
la información se comienza a agregar despues de esta estructura
*/

type control struct {
	FirstSpace int64
}

/*
un nombre que lo identifica y 5 apuntadores que van hacia otras estructuras.
*/

type SuperBloque struct {
	S_filesystem_type   [10]byte
	S_inodes_count      [10]byte
	S_blocks_count      [10]byte
	S_free_blocks_count [10]byte
	S_free_inodes_count [10]byte
	S_mtime             [10]byte
	S_mnt_count         [10]byte
	S_magic             [10]byte
	S_inode_size        [10]byte
	S_block_size        [10]byte
	S_firts_ino         [10]byte
	S_first_blo         [10]byte
	S_bm_inode_start    [10]byte
	S_bm_block_start    [10]byte
	S_inode_start       [10]byte
	S_block_start       [10]byte
}

type Inodo struct {
	I_uid   [10]byte
	I_gid   [10]byte
	I_size  [10]byte
	I_atime [10]byte
	I_ctime [10]byte
	I_mtime [10]byte
	I_block [10]byte
	I_type  byte
	I_perm  [10]byte
}

type CarpetaBloque struct {
	B_content [4]Content
}

type ArchivoBloque struct {
	B_content [64]byte
}

type Content struct {
	B_name  [10]byte
	B_inodo [6]byte
}

type MBR struct {
	Mbr_tamano         [10]byte
	Mbr_fecha_creacion [10]byte
	Mbr_dsk_signature  [10]byte
	Dsk_fit            byte
	Mbr_partition_1    Partition
	Mbr_partition_2    Partition
	Mbr_partition_3    Partition
	Mbr_partition_4    Partition
}

type PartitionMontada struct {
	particion Partition
	id        string
	ruta      string
	montada   bool
	activa    bool
}

type usuarioLogueado struct {
	id     string
	user   string
	pass   string
	grupo  string
	active bool
	admin  bool
	actual PartitionMontada
}

type PartitionesMontadas struct {
	particiones []PartitionMontada
}

type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  [10]byte
	Part_size   [10]byte
	Part_name   [15]byte
}

type Traslado struct {
	idNum     int
	comienzo  int
	fin       int
	anterior  int
	siguiente int
}

type Traslados struct {
	Lista []Traslado
}

type Particiones struct {
	Lista []Partition
}

/*
	type EBR struct {
		part_status byte
		part_fit    byte
		part_start  [4]byte
		part_size   [4]byte
		part_next   [4]byte
		part_name   [10]byte
	}
*/
type MyStructure struct {
	Nombre    [10]byte
	Apuntador [5]int64
}

var listaMontadas PartitionesMontadas
var userActive usuarioLogueado

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func convertirImagen(ruta string) string {
	imagen := ""
	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(ruta)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	// Print the full base64 representation of the image
	//fmt.Println(base64Encoding)
	imagen = base64Encoding
	return imagen
}

func mainArchivos() {
	userActive.active = false
	userActive.admin = false

	//prueba := `mkdisk -Size=15 -unit=m -path=/home/luis/proyecto2/Disco4.dk -fit=bf`
	//prueba = strings.ToLower(prueba)
	//prueba := `rmdisk -path=/home/luis/proyecto2/Disco3.dk`
	//prueba2 := `fdisk -Size=300 -path=/home/luis/proyecto2/Disco4.dk -name=Particion2`
	//prueba2 = strings.ToLower(prueba2)
	//prueba := `mount -path=/home/luis/proyecto2/Disco3.dk -name=Particion1`
	//prueba = strings.ToLower(prueba)
	//prueba3 := `mount -path=/home/luis/proyecto2/Disco4.dk -name=Particion2`
	//prueba3 = strings.ToLower(prueba3)

	//prueba4 := `mkfs -id=391disco4`
	//prueba4 = strings.ToLower(prueba4)

	listaMontadas := PartitionesMontadas{}

	fmt.Println(len(listaMontadas.particiones))
	//cadena := "--Hola--Mundo--Ejemplo--Parzibyte.me"
	//ejecutar(prueba)
	//ejecutar(prueba2)
	//ejecutar(prueba3)
	//ejecutar(prueba4)
	//ejecutar(prueba5)
	//ejecutar(prueba6)
	//ejecutar(prueba7)
	//ejecutar(prueba8)

	//tree("/home/luis/proyecto2/arbol.dot", "391disco4")

	//repDisco("/home/luis/proyecto2/mbr.dot", "391disco4")
	//repFile("/home/luis/proyecto2/mbr.dot", "391disco4", "user.txt")
	//repSuper("/home/luis/proyecto2/mbr.dot", "391disco4")

}
func PrueComentarios(prueba string) {
	divido := strings.Split(prueba, "#")
	fmt.Println("Cometario")
	fmt.Println(divido[0])
	fmt.Println(divido[1])
	fmt.Println("FCometario")
}

func ejecutar(prueba string) string {
	imagen := ""
	prueba = strings.Split(prueba, "#")[0]
	prueba = strings.Split(prueba, "\r")[0]
	if prueba == "" {
		return imagen
	}
	prueba = strings.ToLower(prueba)
	if prueba == "pause" {
		pause()
		return imagen
	}
	fmt.Println(prueba)
	delimitador_mkdisk := "mkdisk "
	delimitador_rmdisk := "rmdisk "
	delimitador_fdisk := "fdisk "
	delimitador_mount := "mount "
	delimitador_mkfs := "mkfs "
	delimitador_mkdir := "mkdir "
	delimitador_mkfile := "mkfile "
	delimitador_mkgrp := "mkgrp "
	delimitador_mkusr := "mkusr "
	delimitador_login := "login "
	delimitador_logout := "logout"
	delimitador_rep := "rep "

	arreglo := strings.Split(prueba, delimitador_mkdisk)
	if arreglo[0] == "" {
		mkdisk(arreglo[1])
	}

	arreglo2 := strings.Split(prueba, delimitador_rmdisk)
	if arreglo2[0] == "" {
		rmdisk(arreglo2[1])
	}

	arreglo3 := strings.Split(prueba, delimitador_fdisk)
	if arreglo3[0] == "" {
		fdisk(arreglo3[1])
	}

	arreglo4 := strings.Split(prueba, delimitador_mount)
	if arreglo4[0] == "" {
		mount(arreglo4[1])
	}

	arreglo5 := strings.Split(prueba, delimitador_mkfs)
	if arreglo5[0] == "" {
		mkfs(arreglo5[1])
	}

	arreglo6 := strings.Split(prueba, delimitador_mkdir)
	if arreglo6[0] == "" {
		mkdir(arreglo6[1], false)
	}

	arreglo7 := strings.Split(prueba, delimitador_mkfile)
	if arreglo7[0] == "" {
		mkdir(arreglo7[1], true)
		mkfile(arreglo7[1])
	}

	arreglo8 := strings.Split(prueba, delimitador_mkgrp)
	if arreglo8[0] == "" {
		mkgrp(arreglo8[1])
	}

	arreglo9 := strings.Split(prueba, delimitador_mkusr)
	if arreglo9[0] == "" {
		mkusr(arreglo9[1])
	}

	arreglo10 := strings.Split(prueba, delimitador_login)
	if arreglo10[0] == "" {
		login(arreglo10[1])
	}

	arreglo11 := strings.Split(prueba, delimitador_logout)
	if arreglo11[0] == "" {
		logout()
	}

	arreglo12 := strings.Split(prueba, delimitador_rep)
	if arreglo12[0] == "" {
		imagen = reportes(arreglo12[1])
	}
	return imagen
}

func rmdisk(linea string) {
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	ruta := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				ruta = arregloCI[1]
				fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				ruta = arregloExec[1]
				fmt.Println(ruta)
			}

		}
	}
	eliminarArchivo(ruta)
}

func mkdisk(linea string) {
	disco := MBR{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	megas := false
	ruta := ""
	disco.Dsk_fit = 'f'
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-size" {
			copy(disco.Mbr_tamano[:], arregloExec[1])
			//disco.mbr_tamano = []byte(arregloExec[1])
			//fmt.Println(disco.mbr_tamano)
			posicionSeguir := posicionVacio(disco.Mbr_tamano[:])
			valor := string(disco.Mbr_tamano[:posicionSeguir])

			//valor = strings.Trim(valor, "5")
			//fmt.Println(len(valor))
			byteToInt, _ := strconv.Atoi(valor)
			fmt.Println(byteToInt)
		} else if arregloExec[0] == "-fit" {
			if arregloExec[1] == "bf" {
				disco.Dsk_fit = 'b'
			} else if arregloExec[1] == "ff" {
				disco.Dsk_fit = 'f'
			} else if arregloExec[1] == "wf" {
				disco.Dsk_fit = 'w'
			}
			fmt.Println(string(disco.Dsk_fit))
		} else if arregloExec[0] == "-unit" {
			if arregloExec[1] == "k" {

			} else if arregloExec[1] == "m" {
				megas = true
			}
			fmt.Println(string(arregloExec[1]))
		} else if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				ruta = arregloCI[1]
				fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				ruta = arregloExec[1]
				fmt.Println(ruta)
			}

		}
	}
	posicionSeguir := posicionVacio(disco.Mbr_tamano[:])
	byteToInt, _ := strconv.Atoi(string(disco.Mbr_tamano[:posicionSeguir]))
	numTexto := ""
	if megas {
		byteToInt = byteToInt * 1024 * 1024
		numTexto = strconv.Itoa(byteToInt)
		copy(disco.Mbr_tamano[:], []byte(numTexto))
		//disco.mbr_tamano = []byte(numTexto)
	} else {
		byteToInt = byteToInt * 1024
		//disco.mbr_tamano = []byte(string(byteToInt))
		numTexto = strconv.Itoa(byteToInt)
		copy(disco.Mbr_tamano[:], []byte(numTexto))
	}
	fmt.Println("Tamaño del disco")
	fmt.Println(byteToInt)
	arreglodeRutas := strings.Split(ruta, "/")
	rutaCarpeta := ""
	for i := 1; i < len(arreglodeRutas)-1; i++ {
		rutaCarpeta = rutaCarpeta + "/" + arreglodeRutas[i]
	}
	fmt.Println(rutaCarpeta)
	crearDirectorioSiNoExiste(rutaCarpeta)
	crearArchivo(ruta, disco)
}

func fdisk(linea string) {
	disco := Partition{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	megas := false
	kilo := true
	ruta := ""
	disco.Part_type = 'p'
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-size" {
			copy(disco.Part_size[:], arregloExec[1])
			//disco.part_size = []byte(arregloExec[1])
			//fmt.Println(disco.mbr_tamano)
			posicionSeguir := posicionVacio(disco.Part_size[:])
			byteToInt, _ := strconv.Atoi(string(disco.Part_size[:posicionSeguir]))
			fmt.Println(byteToInt)
		} else if arregloExec[0] == "-fit" {
			if arregloExec[1] == "bf" {
				disco.Part_fit = 'b'
			} else if arregloExec[1] == "ff" {
				disco.Part_fit = 'f'
			} else if arregloExec[1] == "wf" {
				disco.Part_fit = 'w'
			}
			fmt.Println(string(disco.Part_fit))
		} else if arregloExec[0] == "-unit" {
			if arregloExec[1] == "k" {
				kilo = true
			} else if arregloExec[1] == "m" {
				megas = true
			}
			//fmt.Println(string(disco.dsk_fit))
		} else if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				ruta = arregloCI[1]
				fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				ruta = arregloExec[1]
				fmt.Println(ruta)
			}

		} else if arregloExec[0] == "-type" {
			if arregloExec[1] == "p" {
				disco.Part_type = 'p'
			} else if arregloExec[1] == "e" {
				disco.Part_type = 'e'
			} else if arregloExec[1] == "l" {
				disco.Part_type = 'l'
			}
			fmt.Println(string(disco.Part_type))
		} else if arregloExec[0] == "-name" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(string(disco.Part_name[:posicionSeguir]))
		}
	}
	posicionSeguir := posicionVacio(disco.Part_size[:])
	byteToInt, _ := strconv.Atoi(string(disco.Part_size[:posicionSeguir]))
	numTexto := ""
	if megas {
		byteToInt = byteToInt * 1024 * 1024
		numTexto = strconv.Itoa(byteToInt)
		copy(disco.Part_size[:], []byte(numTexto))
		//disco.part_size = []byte(numTexto)
	} else if kilo {
		byteToInt = byteToInt * 1024
		numTexto = strconv.Itoa(byteToInt)
		copy(disco.Part_size[:], []byte(numTexto))
		//disco.part_size = []byte(string(byteToInt))
	} else {
		//byteToInt = byteToInt
		numTexto = strconv.Itoa(byteToInt)
		copy(disco.Part_size[:], []byte(numTexto))
		//disco.part_size = []byte(string(byteToInt))
	}
	fmt.Println("Tamaño del disco")
	fmt.Println(byteToInt)

	//Abrimos el archivo.
	file, err := os.Open(ruta)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	temporal := LeerDisco(file, size, 1)

	// tamanos
	posPT1 := posicionVacio(temporal.Mbr_partition_1.Part_size[:])
	posP1, _ := strconv.Atoi(string(temporal.Mbr_partition_1.Part_size[:posPT1]))
	posPT2 := posicionVacio(temporal.Mbr_partition_2.Part_size[:])
	posP2, _ := strconv.Atoi(string(temporal.Mbr_partition_2.Part_size[:posPT2]))
	posPT3 := posicionVacio(temporal.Mbr_partition_3.Part_size[:])
	posP3, _ := strconv.Atoi(string(temporal.Mbr_partition_3.Part_size[:posPT3]))
	posPT4 := posicionVacio(temporal.Mbr_partition_4.Part_size[:])
	posP4, _ := strconv.Atoi(string(temporal.Mbr_partition_4.Part_size[:posPT4]))
	tamparticiones := posP1 + posP2 + posP3 + posP4

	tamDisocT := posicionVacio(temporal.Mbr_tamano[:])
	tamDisco, _ := strconv.Atoi(string(temporal.Mbr_tamano[:tamDisocT]))

	if byteToInt+tamparticiones > tamDisco {
		fmt.Println("Tamaño no valida se desborda")
		return
	}

	if disco.Part_type == 'e' {
		if temporal.Mbr_partition_1.Part_type == 'e' || temporal.Mbr_partition_2.Part_type == 'e' || temporal.Mbr_partition_3.Part_type == 'e' || temporal.Mbr_partition_4.Part_type == 'e' {
			fmt.Println("Ya existe una extendida")
			return
		}
	}

	if disco.Part_type == 'l' {
		if temporal.Mbr_partition_1.Part_type == 'e' || temporal.Mbr_partition_2.Part_type == 'e' || temporal.Mbr_partition_3.Part_type == 'e' || temporal.Mbr_partition_4.Part_type == 'e' {

		} else {
			fmt.Println("No existe Extendida")
			return
		}
	}

	trasladosL := Traslados{}
	cantidad := 0
	comienzoPrimera := 0

	if posP1 != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		comienzoT := posicionVacio(temporal.Mbr_partition_1.Part_start[:])
		aux.comienzo, _ = strconv.Atoi(string(temporal.Mbr_partition_1.Part_start[:comienzoT]))
		comienzoPrimera = aux.comienzo
		aux.fin = posP1 + aux.comienzo
		aux.anterior = aux.comienzo
		if posP2 != 0 {
			//auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
			comienzoTSig := posicionVacio(temporal.Mbr_partition_2.Part_start[:])
			aux.siguiente, _ = strconv.Atoi(string(temporal.Mbr_partition_2.Part_start[:comienzoTSig]))
		} else {
			aux.siguiente = tamDisco + aux.comienzo
		}
		trasladosL.Lista = append(trasladosL.Lista, aux)
		//auxParticiones.push_back(aux)
		cantidad++
		if temporal.Mbr_partition_1.Part_type == 'e' {
			//extendida = discoAux.mbr_partition_1
		}
	}

	if posP2 != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		comienzoT := posicionVacio(temporal.Mbr_partition_2.Part_start[:])
		aux.comienzo, _ = strconv.Atoi(string(temporal.Mbr_partition_2.Part_start[:comienzoT]))
		aux.fin = posP2 + aux.comienzo
		aux.anterior = trasladosL.Lista[0].fin
		if posP3 != 0 {
			//auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
			comienzoTSig := posicionVacio(temporal.Mbr_partition_3.Part_start[:])
			aux.siguiente, _ = strconv.Atoi(string(temporal.Mbr_partition_3.Part_start[:comienzoTSig]))
		} else {
			aux.siguiente = tamDisco + comienzoPrimera
		}
		trasladosL.Lista = append(trasladosL.Lista, aux)
		//auxParticiones.push_back(aux)
		cantidad++
		if temporal.Mbr_partition_2.Part_type == 'e' {
			//extendida = discoAux.mbr_partition_1
		}
	}

	if posP3 != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		comienzoT := posicionVacio(temporal.Mbr_partition_3.Part_start[:])
		aux.comienzo, _ = strconv.Atoi(string(temporal.Mbr_partition_3.Part_start[:comienzoT]))
		aux.fin = posP3 + aux.comienzo
		aux.anterior = trasladosL.Lista[1].fin
		if posP4 != 0 {
			//auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
			comienzoTSig := posicionVacio(temporal.Mbr_partition_4.Part_start[:])
			aux.siguiente, _ = strconv.Atoi(string(temporal.Mbr_partition_4.Part_start[:comienzoTSig]))
		} else {
			aux.siguiente = tamDisco + comienzoPrimera
		}
		trasladosL.Lista = append(trasladosL.Lista, aux)
		//auxParticiones.push_back(aux)
		cantidad++
		if temporal.Mbr_partition_3.Part_type == 'e' {
			//extendida = discoAux.mbr_partition_1
		}
	}

	if posP4 != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		comienzoT := posicionVacio(temporal.Mbr_partition_4.Part_start[:])
		aux.comienzo, _ = strconv.Atoi(string(temporal.Mbr_partition_4.Part_start[:comienzoT]))
		aux.fin = posP4 + aux.comienzo
		aux.anterior = trasladosL.Lista[2].fin
		aux.siguiente = tamDisco + comienzoPrimera
		trasladosL.Lista = append(trasladosL.Lista, aux)
		//auxParticiones.push_back(aux)
		cantidad++
		if temporal.Mbr_partition_4.Part_type == 'e' {
			//extendida = discoAux.mbr_partition_1
		}
	}

	particonesL := Particiones{}
	particonesL.Lista = append(particonesL.Lista, temporal.Mbr_partition_1)
	particonesL.Lista = append(particonesL.Lista, temporal.Mbr_partition_2)
	particonesL.Lista = append(particonesL.Lista, temporal.Mbr_partition_3)
	particonesL.Lista = append(particonesL.Lista, temporal.Mbr_partition_4)

	nuevoMBR := AgregarParticionNueva(temporal, disco, trasladosL.Lista, particonesL.Lista, cantidad)

	EditarArchivo(ruta, nuevoMBR)
}

func menu() {
	var option = 0
	for {
		fmt.Println("----------------------------------------")
		fmt.Println("------------Escoja una opción-----------")
		fmt.Println("----------------------------------------")
		fmt.Println("------- 1. Crear archivo binario -------")
		fmt.Println("------ 2.Eliminar archivo binario ------")
		fmt.Println("---------  3. Crear estructura ---------")
		fmt.Println("---------- 4. Leer estructura ----------")
		fmt.Println("------------- 5. Renombrar -------------")
		fmt.Println("--------------- 6. Salir ---------------")
		fmt.Scanf("%d\n", &option)
		switch option {
		case 1:
			//crearArchivo()
			break
		case 2:
			//eliminarArchivo()
			break
		case 3:
			crearEstructura()
			break
		case 4:
			leerEstructuras()
			break
		case 5:
			Renombrar()
			break
		case 6:
			salir()
			break
		}
	}
}

func Renombrar() {
	lectura := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el nombre del Nodo a Renombrar")
	exNodo, _ := lectura.ReadString('\n')
	exNodo = strings.TrimSpace(exNodo)
	fmt.Println("Ingrese el nuevo nombre del nodo")
	newNodo, _ := lectura.ReadString('\n')
	newNodo = strings.TrimSpace(newNodo)
	fmt.Println("EL nommbre e abuscar es: ", exNodo)
	fmt.Println("El nuevo nombre es ", newNodo)

	//Abrimos el archivo.
	file, err := os.OpenFile("misEstructuras.bin", os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	control := control{}

	nodoOrigen := MyStructure{}

	nodoNuevo := MyStructure{}

	var sizeControl int64 = int64(unsafe.Sizeof(control))
	var sizeEstructura int64 = int64(unsafe.Sizeof(nodoNuevo))

	file.Seek(0, 0)
	dataControl := leerBytes(file, int(sizeControl))
	bufferControl := bytes.NewBuffer(dataControl)
	err = binary.Read(bufferControl, binary.BigEndian, &control)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var posicionEscribir = sizeControl + control.FirstSpace*sizeEstructura

	var posicionOrigen = buscarEstructura(file, sizeControl, exNodo)

	file.Seek(posicionOrigen, 0)
	dataOrigen := leerBytes(file, int(sizeEstructura))
	bufferOrigen := bytes.NewBuffer(dataOrigen)
	err = binary.Read(bufferOrigen, binary.BigEndian, &nodoOrigen)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//actualizamos el primer puntero vacio que encontremos del origen
	for i := 0; i < 5; i++ {
		if nodoOrigen.Apuntador[i] == -1 {
			copy(nodoOrigen.Nombre[:], newNodo)
			break
		}
	}

	//lo reescribimos
	file.Seek(posicionOrigen, 0)
	var bufferOrigenW bytes.Buffer
	binary.Write(&bufferOrigenW, binary.BigEndian, &nodoOrigen)
	escribirBytes(file, bufferOrigenW.Bytes())

	//escribimos el nuevo nodo
	file.Seek(posicionEscribir, 0)
	var bufferNuevo bytes.Buffer
	binary.Write(&bufferNuevo, binary.BigEndian, &nodoNuevo)
	escribirBytes(file, bufferNuevo.Bytes())

	//actualizamos el valor del control
	control.FirstSpace = control.FirstSpace + 1
	file.Seek(0, 0)
	var bufferControlW bytes.Buffer
	binary.Write(&bufferControlW, binary.BigEndian, &control)
	escribirBytes(file, bufferControlW.Bytes())

	fmt.Println("renombrado exitosamente!")

}

func Graficar() {
	fmt.Println("Graficar ")
}

func crearArchivo(ruta string, disco MBR) {
	fmt.Println("-----Crear archivo binario-----")

	//variable para llevar control del tamaño del disco
	//var size = 2

	//se procede a crear el archivo
	file, err := os.Create(ruta)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//se crea una variable temporal con un cero que nos ayudará a llenar nuestro archivo de ceros lógicos
	var temporal int8 = 0
	s := &temporal
	var binario bytes.Buffer
	posicionSeguir := posicionVacio(disco.Mbr_tamano[:])
	tam, _ := strconv.Atoi(string(disco.Mbr_tamano[:posicionSeguir]))
	for i := 0; i < tam; i++ {
		binary.Write(&binario, binary.BigEndian, s)
	}
	escribirBytes(file, binario.Bytes())
	/*
			se realiza un for para llenar el archivo completamente de ceros
			NOTA: Para esta parte se recomienda tener un buffer con 1024 ceros (ya que 1024 es la medida
			mínima a escribir) para que este ciclo sea más eficiente


		cont := 1
		for i := 0; i < tam; i++ {

			fmt.Println(cont)
			cont++
		}
	*/
	/*
		se escribira un estudiante por default para llevar el control.
		En el proyecto, el que nos ayuda a llevar el control de las
		particiones es el mbr
	*/

	//Creamos el struct de control y el struct de miEstructura
	miControl := control{FirstSpace: 1}

	nodoRaiz := MyStructure{}
	copy(nodoRaiz.Nombre[:], "201404106")
	//Inicializamos todos los apuntadores de la estructura en -1
	for i := 0; i < 5; i++ {
		nodoRaiz.Apuntador[i] = -1
	}

	//nos posicionamos al inicio del archivo usando la funcion Seek
	file.Seek(0, 0)

	//Escribimos struct de control
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &miControl)
	escribirBytes(file, bufferControl.Bytes())

	//movemos el puntero a donde ira nuestra primera estructura
	file.Seek(int64(unsafe.Sizeof(miControl)), 0)

	//Escribimos struct raiz
	var bufferNodo bytes.Buffer
	binary.Write(&bufferNodo, binary.BigEndian, &disco)
	escribirBytes(file, bufferNodo.Bytes())
	file.Close()
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0755)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			panic(err)
		}
	}
}

func EditarArchivo(ruta string, disco MBR) {
	fmt.Println("-----Crear archivo binario-----")

	//variable para llevar control del tamaño del disco
	//var size = 2

	//se procede a crear el archivo
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//se crea una variable temporal con un cero que nos ayudará a llenar nuestro archivo de ceros lógicos
	var temporal int8 = 0
	s := &temporal
	var binario bytes.Buffer
	posicionSeguir := posicionVacio(disco.Mbr_tamano[:])
	tam, _ := strconv.Atoi(string(disco.Mbr_tamano[:posicionSeguir]))
	for i := 0; i < tam; i++ {
		binary.Write(&binario, binary.BigEndian, s)
	}
	escribirBytes(file, binario.Bytes())
	/*
			se realiza un for para llenar el archivo completamente de ceros
			NOTA: Para esta parte se recomienda tener un buffer con 1024 ceros (ya que 1024 es la medida
			mínima a escribir) para que este ciclo sea más eficiente


		cont := 1
		for i := 0; i < tam; i++ {

			fmt.Println(cont)
			cont++
		}
	*/
	/*
		se escribira un estudiante por default para llevar el control.
		En el proyecto, el que nos ayuda a llevar el control de las
		particiones es el mbr
	*/

	//Creamos el struct de control y el struct de miEstructura
	miControl := control{FirstSpace: 1}

	nodoRaiz := MyStructure{}
	copy(nodoRaiz.Nombre[:], "201404106")
	//Inicializamos todos los apuntadores de la estructura en -1
	for i := 0; i < 5; i++ {
		nodoRaiz.Apuntador[i] = -1
	}

	//nos posicionamos al inicio del archivo usando la funcion Seek
	file.Seek(0, 0)

	//Escribimos struct de control
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &miControl)
	escribirBytes(file, bufferControl.Bytes())

	//movemos el puntero a donde ira nuestra primera estructura
	file.Seek(int64(unsafe.Sizeof(miControl)), 0)

	//Escribimos struct raiz
	var bufferNodo bytes.Buffer
	binary.Write(&bufferNodo, binary.BigEndian, &disco)
	escribirBytes(file, bufferNodo.Bytes())
	file.Close()
}

func eliminarArchivo(ruta string) {
	fmt.Println("-----Eliminar archivo binario-----")

	err := os.Remove(ruta)

	if err != nil {
		fmt.Println("Error al eliminar el archivo.")
	} else {
		fmt.Println("Archivo eliminado exitosamente!")
	}
}

func crearEstructura() {
	fmt.Println("-----Crear estudiante-----")

	//Toma los datos de entrada del usuario
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el origen:")
	origen, _ := reader.ReadString('\n')
	origen = strings.TrimSpace(origen)

	fmt.Print("Ingrese el nombre:")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	//Abrimos el archivo.
	file, err := os.OpenFile("misEstructuras.bin", os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//declaramos la variable control que nos ayudara a escribir el nuevo
	control := control{}

	//declaramos una estructura para leer el nodo de origen
	nodoOrigen := MyStructure{}

	//declaramos el nodo nuevo a crear con sus datos
	nodoNuevo := MyStructure{}
	copy(nodoNuevo.Nombre[:], nombre)
	for i := 0; i < 5; i++ {
		nodoNuevo.Apuntador[i] = -1
	}

	//obtenemor el size del control para empezar a leer desde ahi
	var sizeControl int64 = int64(unsafe.Sizeof(control))
	//obtenemos el size de la estructrua que nos servira en otras operaciones
	var sizeEstructura int64 = int64(unsafe.Sizeof(nodoNuevo))

	//leemos el control
	file.Seek(0, 0)
	dataControl := leerBytes(file, int(sizeControl))
	bufferControl := bytes.NewBuffer(dataControl)
	err = binary.Read(bufferControl, binary.BigEndian, &control)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//obtenemos la posicion donde toca escribir
	var posicionEscribir = sizeControl + control.FirstSpace*sizeEstructura

	//leemos el nodo origen para actualizar el apuntador
	var posicionOrigen = buscarEstructura(file, sizeControl, origen)

	file.Seek(posicionOrigen, 0)
	dataOrigen := leerBytes(file, int(sizeEstructura))
	bufferOrigen := bytes.NewBuffer(dataOrigen)
	err = binary.Read(bufferOrigen, binary.BigEndian, &nodoOrigen)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//actualizamos el primer puntero vacio que encontremos del origen
	for i := 0; i < 5; i++ {
		if nodoOrigen.Apuntador[i] == -1 {
			nodoOrigen.Apuntador[i] = posicionEscribir
			break
		}
	}

	//lo reescribimos
	file.Seek(posicionOrigen, 0)
	var bufferOrigenW bytes.Buffer
	binary.Write(&bufferOrigenW, binary.BigEndian, &nodoOrigen)
	escribirBytes(file, bufferOrigenW.Bytes())

	//escribimos el nuevo nodo
	file.Seek(posicionEscribir, 0)
	var bufferNuevo bytes.Buffer
	binary.Write(&bufferNuevo, binary.BigEndian, &nodoNuevo)
	escribirBytes(file, bufferNuevo.Bytes())

	//actualizamos el valor del control
	control.FirstSpace = control.FirstSpace + 1
	file.Seek(0, 0)
	var bufferControlW bytes.Buffer
	binary.Write(&bufferControlW, binary.BigEndian, &control)
	escribirBytes(file, bufferControlW.Bytes())

	fmt.Println("Creado exitosamente!")
}

/*
Metodo que buscara un nombre especifico dentro de la estructura
(origen representa el nombre que se buscara dentro de la esructura)
posicionEstructura nos ayuda hacer el arbol recursivo para poder
recorrer
*/
func buscarEstructura(file *os.File, posicionEstructura int64, origen string) int64 {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := MyStructure{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		return -1

	} else {
		temporal := string(estructuraTemporal.Nombre[:len(origen)])
		if origen == temporal {
			return posicionEstructura

		} else {

			var posicionInicio int64 = -1

			for i := 0; i < 5 && posicionInicio == -1; i++ {
				if estructuraTemporal.Apuntador[i] != -1 {
					posicionInicio = buscarEstructura(file, estructuraTemporal.Apuntador[i], origen)
				}
			}

			return posicionInicio
		}
	}
}

func leerEstructuras() {
	fmt.Println("-----Leer Estructuras-----")

	//Abrimos el archivo.
	file, err := os.Open("misEstructuras.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	fmt.Println(leerEstructura(file, size, 1))
}

/*
Este metodo será un método recursivo que recibira la posicion en la cual se encuentra
la estructura que se desea leer, y devolverá el string de la concatenación de nombres
*/
func leerEstructura(file *os.File, posicionEstructura int64, nivel int) string {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := MyStructure{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		return "Error!"

	} else {
		//Si se leyo exitosamente los datos
		//declaramos un string temporal para almacenar los nombres de las demas estructuras
		nombres := string(estructuraTemporal.Nombre[:]) + "\n"

		//hacemos un ciclo para recorrer todos los apuntadores y concatenamos el nombre
		for i := 0; i < 5; i++ {
			if estructuraTemporal.Apuntador[i] != -1 {
				for j := 0; j < nivel; j++ {
					nombres += "     "
				}
				nombres += leerEstructura(file, estructuraTemporal.Apuntador[i], nivel+1)
			}
		}

		return (nombres)
	}

}

func salir() {
	os.Exit(0)
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func LeerDisco(file *os.File, posicionEstructura int64, nivel int) MBR {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := MBR{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		fmt.Print("Error!")
		return estructuraTemporal

	} else {
		return estructuraTemporal
	}
}

func posicionVacio(a []byte) int {
	pos := 0
	for i := 0; i < len(a); i++ {
		if a[i] == byte(0) {
			return i
		}
	}
	return pos
}

func byteToString(a byte) string {

	if a == byte(0) {
		return ""
	}

	return string(a)
}

func firstFitPartition(disco MBR, t []Traslado, p []Partition, nueva Partition) MBR {
	//n := len(p)
	m := len(t)
	//MBR actulizado=MBR();
	contador := 0
	agregar := false

	for j := 0; j < m; j++ {
		pos := posicionVacio(nueva.Part_size[:])
		nuevaTam, _ := strconv.Atoi(string(nueva.Part_size[:pos]))
		if t[j].siguiente-t[j].fin >= nuevaTam {
			numTexto := strconv.Itoa(t[j].fin)
			//nueva.part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			contador = j
			agregar = true
			break
		}
	}
	if agregar {
		if contador == 0 {
			aux2 := Partition{}
			aux2 = disco.Mbr_partition_3
			disco.Mbr_partition_4 = aux2
			aux := Partition{}
			aux = disco.Mbr_partition_2
			disco.Mbr_partition_3 = aux
			disco.Mbr_partition_2 = nueva
			//extendidaInicio=disco.mbr_partition_2.part_start;
		} else if contador == 1 {
			aux := Partition{}
			aux = disco.Mbr_partition_3
			disco.Mbr_partition_4 = aux
			disco.Mbr_partition_3 = nueva
			//extendidaInicio=disco.mbr_partition_3.part_start;
		} else if contador == 2 {
			disco.Mbr_partition_4 = nueva
			//extendidaInicio=disco.mbr_partition_4.part_start;
		} else {
			fmt.Println("Particiones llenas")
		}

	} else {
		fmt.Println("Sin espacio para agregar particion")
	}
	return disco
}

func bestFitPartitton(disco MBR, t []Traslado, p []Partition, nueva Partition) MBR {
	//int n = p.size();
	m := len(t)
	agregar := false
	bestIdx := -1
	for j := 0; j < m; j++ {
		pos := posicionVacio(nueva.Part_size[:])
		nuevaTam, _ := strconv.Atoi(string(nueva.Part_size[:pos]))
		if t[j].siguiente-t[j].fin >= nuevaTam {
			if bestIdx == -1 {
				bestIdx = j
			} else if t[bestIdx].siguiente-t[bestIdx].fin > t[j].siguiente-t[j].fin {
				bestIdx = j
			}
		}
	}

	if bestIdx != -1 {
		if bestIdx == 0 {
			numTexto := strconv.Itoa(t[bestIdx].fin)
			//nueva.Part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			agregar = true
			aux := Partition{}
			aux = disco.Mbr_partition_3
			disco.Mbr_partition_4 = aux
			aux2 := Partition{}
			aux2 = disco.Mbr_partition_2
			disco.Mbr_partition_3 = aux2
			disco.Mbr_partition_2 = nueva
			//extendidaInicio=nueva.part_start;
		} else if bestIdx == 1 {
			numTexto := strconv.Itoa(t[bestIdx].fin)
			//nueva.Part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			agregar = true
			aux := Partition{}
			aux = disco.Mbr_partition_3
			disco.Mbr_partition_4 = aux
			disco.Mbr_partition_3 = nueva
			//extendidaInicio=nueva.part_start;
		} else if bestIdx == 2 {
			numTexto := strconv.Itoa(t[bestIdx].fin)
			//nueva.Part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			agregar = true
			disco.Mbr_partition_4 = nueva
			//extendidaInicio=nueva.part_start;
		}
	}

	if !agregar {
		fmt.Println("Sin espacio para agregar particion")
	}
	return disco
}

func worstFitPartitton(disco MBR, t []Traslado, p []Partition, nueva Partition) MBR {
	//int n = p.size();
	m := len(t)
	agregar := false
	bestIdx := -1
	for j := 0; j < m; j++ {
		pos := posicionVacio(nueva.Part_size[:])
		nuevaTam, _ := strconv.Atoi(string(nueva.Part_size[:pos]))
		if t[j].siguiente-t[j].fin >= nuevaTam {
			if bestIdx == -1 {
				bestIdx = j
			} else if t[bestIdx].siguiente-t[bestIdx].fin < t[j].siguiente-t[j].fin {
				bestIdx = j
			}
		}
	}

	if bestIdx != -1 {
		if bestIdx == 0 {
			numTexto := strconv.Itoa(t[bestIdx].fin)
			//nueva.part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			agregar = true
			aux := Partition{}
			aux = disco.Mbr_partition_3
			disco.Mbr_partition_4 = aux
			aux2 := Partition{}
			aux2 = disco.Mbr_partition_2
			disco.Mbr_partition_3 = aux2
			disco.Mbr_partition_2 = nueva
			//extendidaInicio=nueva.part_start
		} else if bestIdx == 1 {
			numTexto := strconv.Itoa(t[bestIdx].fin)
			//nueva.Part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			agregar = true
			aux := Partition{}
			aux = disco.Mbr_partition_3
			disco.Mbr_partition_4 = aux
			disco.Mbr_partition_3 = nueva
			//extendidaInicio=nueva.part_start;
		} else if bestIdx == 2 {
			numTexto := strconv.Itoa(t[bestIdx].fin)
			//nueva.Part_start = []byte(numTexto)
			copy(nueva.Part_start[:], numTexto)
			agregar = true
			disco.Mbr_partition_4 = nueva
			//extendidaInicio=nueva.part_start;
		}
	}
	if !agregar {
		fmt.Println("Sin espacio para agregar particion")
	}
	return disco
}

func AgregarParticionNueva(mbr MBR, p Partition, t []Traslado, ps []Partition, u int) MBR {
	if u == 0 {
		tamTexto := strconv.Itoa(int(unsafe.Sizeof(MBR{})))
		copy(p.Part_start[:], []byte(tamTexto))
		//p.Part_start =

		mbr.Mbr_partition_1 = p
		//extendidaInicio = p.part_start
		return mbr
	} else {
		if mbr.Dsk_fit == 'f' {
			return bestFitPartitton(mbr, t, ps, p)
		}
		if mbr.Dsk_fit == 'b' {
			return bestFitPartitton(mbr, t, ps, p)
		}
		if mbr.Dsk_fit == 'w' {
			return bestFitPartitton(mbr, t, ps, p)
		}
	}
	return mbr
}

func mount(linea string) {
	disco := Partition{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	ruta := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				ruta = arregloCI[1]
				fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				ruta = arregloExec[1]
				fmt.Println(ruta)
			}

		} else if arregloExec[0] == "-name" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(string(disco.Part_name[:posicionSeguir]))
		}
	}

	//Abrimos el archivo.
	file, err := os.Open(ruta)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	temporal := LeerDisco(file, size, 1)
	arregloSpit := strings.Split(ruta, "/")
	tamA := len(arregloSpit)
	nombreDisco := strings.Split(arregloSpit[tamA-1], ".")[0]
	fmt.Println(nombreDisco)
	numPart := "0"
	A_Montar := Partition{}
	if temporal.Mbr_partition_1.Part_name == disco.Part_name {
		temporal.Mbr_partition_1.Part_status = 's'
		numPart = "1"
		A_Montar = temporal.Mbr_partition_1
		//actual.particion=auxDisco.mbr_partition_1;
	} else if temporal.Mbr_partition_2.Part_name == disco.Part_name {
		temporal.Mbr_partition_2.Part_status = 's'
		numPart = "2"
		A_Montar = temporal.Mbr_partition_2
		//actual.particion=auxDisco.mbr_partition_2
	} else if temporal.Mbr_partition_3.Part_name == disco.Part_name {
		temporal.Mbr_partition_3.Part_status = 's'
		numPart = "3"
		A_Montar = temporal.Mbr_partition_3
		//actual.particion=auxDisco.mbr_partition_3
	} else if temporal.Mbr_partition_4.Part_name == disco.Part_name {
		temporal.Mbr_partition_4.Part_status = 's'
		numPart = "4"
		A_Montar = temporal.Mbr_partition_4
		//actual.particion=auxDisco.mbr_partition_4
	}

	EditarArchivo(ruta, temporal)
	fmt.Println("39" + numPart + "a")
	montada := PartitionMontada{}
	montada.id = "39" + numPart + "a"
	montada.montada = true
	montada.particion = A_Montar
	montada.ruta = ruta
	montada.activa = true
	listaMontadas.particiones = append(listaMontadas.particiones, montada)
	fmt.Println("Montada")
}

func mkfs(linea string) {
	disco := Partition{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	nombreId := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-type" {
			if arregloExec[1] == "full" {
				fmt.Println("Full")
			} else {
				disco.Part_type = 'e'
			}

		} else if arregloExec[0] == "-id" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
				nombreId = arregloCI[1]
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
				nombreId = arregloExec[1]
			}
			posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(string(disco.Part_name[:posicionSeguir]))
		}
	}
	tamMontadas := len(listaMontadas.particiones)
	auxMontada := PartitionMontada{}
	for i := 0; i < tamMontadas; i++ {
		if nombreId == listaMontadas.particiones[i].id {
			auxMontada = listaMontadas.particiones[i]
			break
		}
	}

	numEstructuras := 0
	superBlock := SuperBloque{}
	tamParticionT := posicionVacio(auxMontada.particion.Part_size[:])
	tamParticion, _ := strconv.Atoi(string(auxMontada.particion.Part_size[:tamParticionT]))
	numEstructuras = (tamParticion - int(unsafe.Sizeof(SuperBloque{})))
	numEstructuras = numEstructuras / (4 + int(unsafe.Sizeof(Inodo{})) + (3 * int(unsafe.Sizeof(CarpetaBloque{}))))
	numEstructuras = int(math.Floor(float64(numEstructuras)))
	copy(superBlock.S_filesystem_type[:], "2")
	copy(superBlock.S_inodes_count[:], strconv.Itoa(numEstructuras))
	copy(superBlock.S_blocks_count[:], strconv.Itoa(3*numEstructuras))
	copy(superBlock.S_free_blocks_count[:], strconv.Itoa(3*numEstructuras))
	copy(superBlock.S_free_inodes_count[:], strconv.Itoa(numEstructuras))
	copy(superBlock.S_mtime[:], string(time.Now().GoString()))
	copy(superBlock.S_mnt_count[:], "1")
	copy(superBlock.S_magic[:], "61267")
	copy(superBlock.S_firts_ino[:], "2")
	copy(superBlock.S_first_blo[:], "2")

	ceros := '0'
	grupo := "1,G,root\n1,U,root,root,123\n"

	//EXT2

	inodoStart, _ := strconv.Atoi(string(auxMontada.particion.Part_start[:posicionVacio(auxMontada.particion.Part_start[:])]))
	inicioParticion := inodoStart
	inodoStart = inodoStart + int(unsafe.Sizeof(SuperBloque{}))
	copy(superBlock.S_bm_inode_start[:], strconv.Itoa(inodoStart))
	blockStart := inodoStart + numEstructuras
	copy(superBlock.S_bm_block_start[:], strconv.Itoa(blockStart))
	copy(superBlock.S_inode_start[:], strconv.Itoa(blockStart+(3*numEstructuras)))
	copy(superBlock.S_block_start[:], strconv.Itoa(inodoStart+(int(unsafe.Sizeof(SuperBloque{}))*numEstructuras)))
	//inodo_Start := blockStart + (3 * numEstructuras)
	block_Start := inodoStart + (int(unsafe.Sizeof(SuperBloque{})) * numEstructuras)

	//Abrimos el archivo.
	file, err := os.OpenFile(auxMontada.ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	control := control{}

	//var sizeControl int64 = int64(unsafe.Sizeof(control))
	//var sizeEstructura int64 = int64(unsafe.Sizeof(superBlock))
	fmt.Print(control)
	file.Seek(int64(inicioParticion), 0)
	var bufferOrigenW bytes.Buffer
	binary.Write(&bufferOrigenW, binary.BigEndian, &superBlock)
	escribirBytes(file, bufferOrigenW.Bytes())

	file.Seek(int64(inodoStart), 0)
	var bufferCeros bytes.Buffer
	for i := 0; i < numEstructuras; i++ {
		binary.Write(&bufferCeros, binary.BigEndian, &ceros)
	}
	escribirBytes(file, bufferCeros.Bytes())

	var bufferCeros2 bytes.Buffer
	file.Seek(int64(blockStart), 0)
	for i := 0; i < (3 * numEstructuras); i++ {
		binary.Write(&bufferCeros2, binary.BigEndian, &ceros)
	}
	escribirBytes(file, bufferCeros2.Bytes())

	inodo := Inodo{}
	copy(inodo.I_gid[:], "1")
	copy(inodo.I_uid[:], "1")
	copy(inodo.I_atime[:], string(time.Now().GoString()))
	copy(inodo.I_ctime[:], string(time.Now().GoString()))
	copy(inodo.I_mtime[:], string(time.Now().GoString()))
	copy(inodo.I_block[:], "0")
	inodo.I_type = '0'
	copy(inodo.I_perm[:], "664")

	raiz := CarpetaBloque{}
	copy(raiz.B_content[0].B_name[:], ".")
	copy(raiz.B_content[0].B_inodo[:], "0")
	copy(raiz.B_content[1].B_name[:], "..")
	copy(raiz.B_content[1].B_inodo[:], "0")
	copy(raiz.B_content[2].B_name[:], "user.txt")
	copy(raiz.B_content[2].B_inodo[:], "1")

	unionInodo := Inodo{}
	copy(unionInodo.I_uid[:], "1")
	copy(unionInodo.I_gid[:], "1")
	//unionInodo.i_s = sizeof(grupo.c_str()) + sizeof(CarpetaBloque);
	copy(unionInodo.I_atime[:], string(time.Now().GoString()))
	copy(unionInodo.I_ctime[:], string(time.Now().GoString()))
	copy(unionInodo.I_mtime[:], string(time.Now().GoString()))
	copy(unionInodo.I_block[:], "1")
	unionInodo.I_type = '1'
	copy(unionInodo.I_perm[:], "664")

	//inodo.i_s = unionInodo.i_s + sizeof(CarpetaBloque) + sizeof(Inodo);

	archivoBloque := ArchivoBloque{}
	copy(archivoBloque.B_content[:], grupo)

	var bufferM_I bytes.Buffer
	file.Seek(int64(inodoStart), 0)
	unos := '1'
	for i := 0; i < 2; i++ {
		binary.Write(&bufferM_I, binary.BigEndian, &unos)
	}
	escribirBytes(file, bufferM_I.Bytes())

	var bufferM_B bytes.Buffer
	file.Seek(int64(blockStart), 0)

	for i := 0; i < 2; i++ {
		binary.Write(&bufferM_B, binary.BigEndian, &unos)
	}
	escribirBytes(file, bufferM_B.Bytes())

	var bufferI_S bytes.Buffer
	file.Seek(int64(inodoStart), 0)
	binary.Write(&bufferI_S, binary.BigEndian, &inodo)
	binary.Write(&bufferI_S, binary.BigEndian, &unionInodo)
	escribirBytes(file, bufferI_S.Bytes())

	var bufferB_S bytes.Buffer
	file.Seek(int64(block_Start), 0)
	binary.Write(&bufferB_S, binary.BigEndian, &raiz)
	binary.Write(&bufferB_S, binary.BigEndian, &archivoBloque)
	escribirBytes(file, bufferB_S.Bytes())

}

func escribir() {
	//Abrimos el archivo.
	file, err := os.OpenFile("misEstructuras.bin", os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	control := control{}

	var sizeControl int64 = int64(unsafe.Sizeof(control))
	var sizeEstructura int64 = int64(unsafe.Sizeof(control))

	file.Seek(0, 0)
	dataControl := leerBytes(file, int(sizeControl))
	bufferControl := bytes.NewBuffer(dataControl)
	err = binary.Read(bufferControl, binary.BigEndian, &control)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var posicionEscribir = sizeControl + control.FirstSpace*sizeEstructura

	var posicionOrigen = buscarEstructura(file, sizeControl, "con")

	file.Seek(posicionOrigen, 0)
	dataOrigen := leerBytes(file, int(sizeEstructura))
	bufferOrigen := bytes.NewBuffer(dataOrigen)
	err = binary.Read(bufferOrigen, binary.BigEndian, &control)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//actualizamos el primer puntero vacio que encontremos del origen

	//lo reescribimos
	file.Seek(posicionOrigen, 0)
	var bufferOrigenW bytes.Buffer
	binary.Write(&bufferOrigenW, binary.BigEndian, &control)
	escribirBytes(file, bufferOrigenW.Bytes())

	//escribimos el nuevo nodo
	file.Seek(posicionEscribir, 0)
	var bufferNuevo bytes.Buffer
	binary.Write(&bufferNuevo, binary.BigEndian, &control)
	escribirBytes(file, bufferNuevo.Bytes())

	//actualizamos el valor del control
	control.FirstSpace = control.FirstSpace + 1
	file.Seek(0, 0)
	var bufferControlW bytes.Buffer
	binary.Write(&bufferControlW, binary.BigEndian, &control)
	escribirBytes(file, bufferControlW.Bytes())

	fmt.Println("renombrado exitosamente!")
}

func mkdir(linea string, archivo bool) {

	delimitador := " "
	delimitador2 := "="
	delimitador4 := "/"
	arreglo := strings.Split(linea, delimitador)
	ruta := ""
	padre := false
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-p" {
			padre = true
			fmt.Println(padre)
		} else if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				ruta = arregloCI[1]
				fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				ruta = arregloExec[1]
				fmt.Println(ruta)
			}

		}
	}

	tamMontadas := len(listaMontadas.particiones)
	auxMontada := PartitionMontada{}
	for i := 0; i < tamMontadas; i++ {
		if listaMontadas.particiones[i].activa {
			auxMontada = listaMontadas.particiones[i]
			break
		}
	}
	//Abrimos el archivo.
	file, err := os.Open(auxMontada.ruta)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	disco := LeerDisco(file, size, 1)
	inicioParticion := 0
	if auxMontada.particion.Part_name == disco.Mbr_partition_1.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_1.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_1.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == disco.Mbr_partition_2.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_2.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_2.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == disco.Mbr_partition_3.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_3.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_3.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == disco.Mbr_partition_4.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_4.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_4.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	}

	super := LeerSuperBloque(file, int64(inicioParticion), 1)

	posicionInodoT := posicionVacio(super.S_inode_start[:])
	posicionInodo, _ := strconv.Atoi(string(super.S_inode_start[:posicionInodoT]))
	posicionInodo = inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	posicionBloqueT := posicionVacio(super.S_block_start[:])
	posicionBloque, _ := strconv.Atoi(string(super.S_block_start[:posicionBloqueT]))

	inodoRaiz := LeerInodo(file, int64(posicionInodo), 1)

	arregloCarpetas := strings.Split(ruta, delimitador4)
	tamRuta := len(arregloCarpetas)
	if archivo {
		tamRuta = tamRuta - 1
	}
	fmt.Println(posicionBloque + tamRuta)
	fmt.Println(arregloCarpetas)
	fmt.Println(inodoRaiz)
	file.Close()
	if tamRuta > 0 {
		LeerDirectorioCarpeta(inodoRaiz, posicionInodo, posicionBloque, auxMontada.ruta, arregloCarpetas, tamRuta, super, inicioParticion, 1, posicionBloque, 0, "")
	}

}

func LeerSuperBloque(file *os.File, posicionEstructura int64, nivel int) SuperBloque {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := SuperBloque{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		fmt.Print("Error!")
		return estructuraTemporal

	} else {
		return estructuraTemporal
	}
}

func LeerInodo(file *os.File, posicionEstructura int64, nivel int) Inodo {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := Inodo{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		fmt.Print("Error!")
		return estructuraTemporal

	} else {
		return estructuraTemporal
	}
}

func LeerBloqueCarpeta(file *os.File, posicionEstructura int64, nivel int) CarpetaBloque {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := CarpetaBloque{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		fmt.Print("Error!")
		return estructuraTemporal

	} else {
		return estructuraTemporal
	}
}

func LeerBloqueArchivo(file *os.File, posicionEstructura int64, nivel int) ArchivoBloque {

	//Se mueve el puntero hacia la nueva posicion
	file.Seek(posicionEstructura, 0)

	//Se declara una estructura temporal para leer la estructura del archivo
	estructuraTemporal := ArchivoBloque{}
	var size int = int(unsafe.Sizeof(estructuraTemporal))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &estructuraTemporal)

	if err != nil {
		log.Fatal("binary.Read failed", err)
		fmt.Print("Error!")
		return estructuraTemporal

	} else {
		return estructuraTemporal
	}
}

func LeerDirectorioCarpeta(inodo Inodo, posI int, posB int, ruta string, listaNombre []string, nivel int, super SuperBloque, inicioParticion int, numInodo int, posBI int, contador int, CarpetaAnterior string) {
	//tamRuta = len(listaNombre)
	nombreN := listaNombre[contador]
	//Abrimos el archivo.
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//inodo = LeerInodo(file, int64(posI), 1)
	tipoInodo := string(inodo.I_type)
	existeCarpeta := false
	iBlockLibre := 0
	iBlockLibreCarpeta := -1
	carpetaTemporal := CarpetaBloque{}
	//bloqueCarpetaLibreT := bitsAString(super.S_first_blo[:])
	//bloqueCarpetaLibre := bitsAInt(super.S_first_blo[:])
	primerVacio := false
	SiguienteInodo := 0
	if tipoInodo == "0" {
		for i := 0; i < 10; i++ {
			if primerVacio || existeCarpeta {
				break
			}
			if inodo.I_block[i] != 0 {
				iBlockLibre++
				posicionCarpetaT := byteToString(inodo.I_block[i])
				posicionCarpeta, _ := strconv.Atoi(posicionCarpetaT)
				carpeta := LeerBloqueCarpeta(file, int64(posicionCarpeta*int(unsafe.Sizeof(CarpetaBloque{})))+int64(posB), 1)
				//posB = posB + int(unsafe.Sizeof(carpeta))
				for j := 0; j < 4; j++ {
					posicionT := posicionVacio(carpeta.B_content[j].B_name[:])
					nombreCarpeta := string(carpeta.B_content[j].B_name[:posicionT])
					if nombreCarpeta == "" {
						iBlockLibreCarpeta = j
						primerVacio = true
						carpetaTemporal = carpeta
						break
					}
					if nombreCarpeta == "." || nombreCarpeta == ".." {

					} else {
						if nombreCarpeta == nombreN {
							existeCarpeta = true
							SiguienteInodo = bitsAInt(carpeta.B_content[j].B_inodo[:])
							break
						}
					}
				}

			}
		}
	}

	if !existeCarpeta {
		if iBlockLibreCarpeta != -1 {
			copy(carpetaTemporal.B_content[iBlockLibreCarpeta].B_name[:], nombreN)
			posicionILT := posicionVacio(super.S_firts_ino[:])
			posicionIL := string(super.S_firts_ino[:posicionILT])
			primerInodoLibre, _ := strconv.Atoi(string(posicionIL))
			SiguienteInodo = primerInodoLibre
			posicionBLT := posicionVacio(super.S_first_blo[:])
			posicionBL := string(super.S_first_blo[:posicionBLT])
			primerBloqueLibre, _ := strconv.Atoi(string(posicionBL))
			copy(carpetaTemporal.B_content[iBlockLibreCarpeta].B_inodo[:], posicionIL)

			copy(super.S_firts_ino[:], strconv.Itoa(primerInodoLibre+1))
			copy(super.S_first_blo[:], strconv.Itoa(primerBloqueLibre+1))

			inodoN := Inodo{}
			copy(inodoN.I_gid[:], "1")
			copy(inodoN.I_uid[:], "1")
			copy(inodoN.I_atime[:], string(time.Now().GoString()))
			copy(inodoN.I_ctime[:], string(time.Now().GoString()))
			copy(inodoN.I_mtime[:], string(time.Now().GoString()))
			copy(inodoN.I_block[:], posicionBL)
			inodoN.I_type = '0'
			copy(inodoN.I_perm[:], "664")
			copy(inodoN.I_size[:], "81")

			nuevaCarpeta := CarpetaBloque{}
			copy(nuevaCarpeta.B_content[0].B_name[:], ("."))
			copy(nuevaCarpeta.B_content[0].B_inodo[:], posicionIL)
			copy(nuevaCarpeta.B_content[1].B_name[:], (".."))
			copy(nuevaCarpeta.B_content[1].B_inodo[:], posicionIL)

			file.Seek(0, 0)
			file.Seek(int64(inicioParticion), 0)
			var bufferOrigenS bytes.Buffer
			binary.Write(&bufferOrigenS, binary.BigEndian, &super)
			escribirBytes(file, bufferOrigenS.Bytes())

			numBloqueT := string(inodo.I_block[iBlockLibre-1])
			numBloque, _ := strconv.Atoi(string(numBloqueT))
			axu1 := int64(64*numBloque) + int64(posBI)
			fmt.Println(axu1)
			tamInodo := int(unsafe.Sizeof(Inodo{}))
			axu2 := int64(tamInodo*primerInodoLibre) + int64(posI)
			fmt.Println(axu2)

			file.Seek(int64(64*numBloque)+int64(posBI), 0)
			var bufferBloqueN bytes.Buffer
			binary.Write(&bufferBloqueN, binary.BigEndian, &carpetaTemporal)
			escribirBytes(file, bufferBloqueN.Bytes())

			file.Seek(int64(64*primerBloqueLibre)+int64(posBI), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaCarpeta)
			escribirBytes(file, bufferBloque.Bytes())

			file.Seek(int64(tamInodo*primerInodoLibre)+int64(posI), 0)
			var bufferInodo bytes.Buffer
			binary.Write(&bufferInodo, binary.BigEndian, &inodoN)
			escribirBytes(file, bufferInodo.Bytes())

		} else {
			nuevaCarpeta := CarpetaBloque{}
			copy(nuevaCarpeta.B_content[0].B_name[:], nombreN)
			posicionILT := posicionVacio(super.S_firts_ino[:])
			posicionIL := string(super.S_firts_ino[:posicionILT])
			primerInodoLibre, _ := strconv.Atoi(string(posicionIL))
			SiguienteInodo = primerInodoLibre
			posicionBLT := posicionVacio(super.S_first_blo[:])
			posicionBL := string(super.S_first_blo[:posicionBLT])
			primerBloqueLibre, _ := strconv.Atoi(string(posicionBL))
			copy(nuevaCarpeta.B_content[0].B_inodo[:], strconv.Itoa(primerInodoLibre))
			inodo.I_block[iBlockLibre] = byte(posicionBL[0])

			copy(super.S_firts_ino[:], strconv.Itoa(primerInodoLibre+1))

			copy(super.S_first_blo[:], strconv.Itoa(primerBloqueLibre+2))

			inodoN := Inodo{}
			copy(inodoN.I_gid[:], "1")
			copy(inodoN.I_uid[:], "1")
			copy(inodoN.I_atime[:], string(time.Now().GoString()))
			copy(inodoN.I_ctime[:], string(time.Now().GoString()))
			copy(inodoN.I_mtime[:], string(time.Now().GoString()))
			copy(inodoN.I_block[:], strconv.Itoa(primerBloqueLibre+1))
			inodoN.I_type = '0'
			copy(inodoN.I_size[:], "81")
			copy(inodoN.I_perm[:], "664")

			nuevaCarpetaN := CarpetaBloque{}
			copy(nuevaCarpetaN.B_content[0].B_name[:], ("."))
			copy(nuevaCarpetaN.B_content[0].B_inodo[:], posicionIL)
			copy(nuevaCarpetaN.B_content[1].B_name[:], (".."))
			copy(nuevaCarpetaN.B_content[1].B_inodo[:], posicionIL)

			tamInodo := int(unsafe.Sizeof(Inodo{}))
			file.Seek(0, 0)
			file.Seek(int64(inicioParticion), 0)
			var bufferOrigenW bytes.Buffer
			binary.Write(&bufferOrigenW, binary.BigEndian, &super)
			escribirBytes(file, bufferOrigenW.Bytes())

			//numBloqueT := string(inodo.I_block[iBlockLibre-1])
			//numBloque, _ := strconv.Atoi(string(numBloqueT))

			file.Seek(int64(tamInodo*(numInodo-1))+int64(posI), 0)
			var bufferInodoActual bytes.Buffer
			binary.Write(&bufferInodoActual, binary.BigEndian, &inodo)
			escribirBytes(file, bufferInodoActual.Bytes())

			file.Seek(int64(64*primerBloqueLibre)+int64(posBI), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaCarpeta)
			escribirBytes(file, bufferBloque.Bytes())

			file.Seek(int64(64*(primerBloqueLibre+1))+int64(posBI), 0)
			var bufferBloqueN bytes.Buffer
			binary.Write(&bufferBloqueN, binary.BigEndian, &nuevaCarpetaN)
			escribirBytes(file, bufferBloqueN.Bytes())

			file.Seek(int64(tamInodo*(primerInodoLibre))+int64(posI), 0)
			var bufferInodo bytes.Buffer
			binary.Write(&bufferInodo, binary.BigEndian, &inodoN)
			escribirBytes(file, bufferInodo.Bytes())

		}
	}
	if nivel > 1 {
		inodoTem := LeerInodo(file, int64(posI)+(int64(unsafe.Sizeof(Inodo{}))*int64(SiguienteInodo)), 1)
		tipoTemInodo := string(inodoTem.I_type)
		if tipoTemInodo == "0" {
			LeerDirectorioCarpeta(inodoTem, int(posI), posB, ruta, listaNombre, (nivel - 1), super, inicioParticion, SiguienteInodo+1, posBI, contador+1, listaNombre[contador])
		} else {
			//bloqueArc := ArchivoBloque{}
			//posB = posB + (int(unsafe.Sizeof(bloqueArc) * 16))
		}
	}
}

func tree(p string, id string) {
	path := ""
	spr := SuperBloque{}
	inode := Inodo{}
	//inodeArchivo := Inodo{}

	//partition := Partition{}
	auxMontada := PartitionMontada{}
	for _, actual := range listaMontadas.particiones {
		if actual.id == id {
			auxMontada = actual
			//partition = auxMontada.particion
			path = auxMontada.ruta
		}
	}

	//Abrimos el archivo.
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	inicioParticion, _ := strconv.Atoi(string(auxMontada.particion.Part_start[:posicionVacio(auxMontada.particion.Part_start[:])]))
	inicioInodos := inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	spr = LeerSuperBloque(file, int64(inicioParticion), 1)

	inode = LeerInodo(file, int64(inicioInodos), 1)

	posT := posicionVacio(spr.S_firts_ino[:])
	inodoLibre, _ := strconv.Atoi(string(spr.S_firts_ino[:posT]))

	posBT := posicionVacio(spr.S_block_start[:])
	posBloques, _ := strconv.Atoi(string(spr.S_block_start[:posBT]))

	content := "digraph G{\n" +
		"rankdir=LR;\n" +
		"graph [ dpi = \"600\" ]; \n" +
		"forcelabels= true;\n" +
		"node [shape = plaintext];\n"

	for i := 0; i < inodoLibre; i++ {
		content += "inode" + strconv.Itoa(i) + "  [label = <<table BGCOLOR=\"#87CEEB\">\n" +
			"<tr><td COLSPAN = '2'><font color=\"black\">INODO " +
			strconv.Itoa(i) + "</font></td></tr>\n" +
			"<tr>\n" +
			"<td>i_uid</td>\n" +
			"<td>" +
			bitsAString(inode.I_uid[:]) + "</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td>i_gid</td>\n" +
			"<td>" +
			bitsAString(inode.I_gid[:]) + "</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td>i_s</td>\n" +
			"<td>" +
			bitsAString(inode.I_size[:]) + "</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td>i_atime</td>\n" +
			"<td>" +
			time.Now().Local().String() + "</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td>i_ctime</td>\n" +
			"<td>" +
			time.Now().Local().String() + "</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td>i_mtime</td>\n" +
			"<td>" +
			time.Now().Local().String() + "</td>\n" +
			"</tr>\n"

		for j := 0; j < 10; j++ {
			content += "<tr>\n" +
				"<td>i_block_" + strconv.Itoa(j+1) + "</td>\n" +
				"<td port=\"" + strconv.Itoa(j) + "\">" +
				byteToString(inode.I_block[j]) + "</td>\n" +
				"</tr>\n"
		}

		content += "<tr>\n" +
			"<td>i_type</td>\n" +
			"<td>" + string(inode.I_type) + "</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td>i_perm</td>\n" +
			"<td>" + bitsAString(inode.I_perm[:]) + "</td>\n" +
			"</tr>\n</table>>];\n"

		if inode.I_type == '0' {
			for j := 0; j < 10; j++ {
				if inode.I_block[j] != 0 {
					content += "inode" + strconv.Itoa(i) + ":" + strconv.Itoa(j) + "-> BLOCK" + string(inode.I_block[j]) +
						"\n"
					tem, _ := strconv.Atoi(string(inode.I_block[j]))
					//fmt.Println(tem)
					foldertmp := LeerBloqueCarpeta(file, int64(posBloques+(int(unsafe.Sizeof(CarpetaBloque{}))*tem)), 1)

					content += "BLOCK" + byteToString(inode.I_block[j]) + "  [label = <<table BGCOLOR=\"#f4a020\">\n" +
						"<tr><td COLSPAN = '2'><font color=\"black\">BLOCK " +
						byteToString(inode.I_block[j]) + "</font></td></tr>\n"

					for k := 0; k < 4; k++ {
						ctmp := ""
						ctmp += bitsAString(foldertmp.B_content[k].B_name[:])
						content += "<tr>\n" +
							"<td>" + bitsAString(foldertmp.B_content[k].B_name[:]) + "</td>\n" +
							"<td port=\"" + strconv.Itoa(k) + "\">" +
							bitsAString(foldertmp.B_content[k].B_inodo[:]) + "</td>\n" +
							"</tr>\n"
					}

					content += "</table>>];\n"

					for b := 0; b < 4; b++ {
						if string(foldertmp.B_content[b].B_inodo[0]) != "0" && foldertmp.B_content[b].B_inodo[0] != 0 {
							nm := bitsAString(foldertmp.B_content[b].B_name[:])
							if !((nm == ".") || (nm == "..")) {
								content +=
									"BLOCK" + byteToString(inode.I_block[j]) + ":" + strconv.Itoa(b) + " -> inode" +
										bitsAString(foldertmp.B_content[b].B_inodo[:]) + ";\n"
							}
						}
					}
				}
			}
		} else {
			for j := 0; j < 10; j++ {
				if inode.I_block[j] != 0 {
					tem, _ := strconv.Atoi(string(inode.I_block[j]))
					//fmt.Println(tem)
					filetmp := LeerBloqueArchivo(file, int64(posBloques+(int(unsafe.Sizeof(CarpetaBloque{}))*tem)), 1)

					content += "inode" + strconv.Itoa(i) + ":" + strconv.Itoa(j) + "-> BLOCK" +
						byteToString(inode.I_block[j]) +
						"\n"

					content += "BLOCK" + byteToString(inode.I_block[j]) + " [label = <<table BGCOLOR=\"#008000\">\n" +
						"<tr><td COLSPAN = '2'>BLOCK " +
						byteToString(inode.I_block[j]) +
						"</td></tr>\n <tr><td COLSPAN = '2'>" + bitsAString(filetmp.B_content[:]) +
						"</td></tr>\n</table>>];\n"
				}
			}
		}
		inode = LeerInodo(file, int64(inicioInodos)+(int64(unsafe.Sizeof(Inodo{}))*int64(i+1)), 1)

	}
	content += "\n\n}\n"
	fmt.Println(content)
	crearDot(content, p, 1)
}

func bitsAInt(arreglo []byte) int {
	posT := posicionVacio(arreglo[:])
	num, _ := strconv.Atoi(string(arreglo[:posT]))
	return num
}

func bitsAString(arreglo []byte) string {
	posT := posicionVacio(arreglo[:])
	tex := string(arreglo[:posT])
	return tex
}

func crearDot(contenido string, nombre string, formato int) {
	f, err := os.Create("/home/luis/proyecto2/mbr.dot")
	if err != nil {
		panic(err)
	}
	f.Close()

	b := []byte(contenido)
	er := ioutil.WriteFile("/home/luis/proyecto2/mbr.dot", b, 0644)
	if er != nil {
		log.Fatal(er)
	}
	arreglodeRutas := strings.Split(nombre, "/")
	rutaCarpeta := ""
	for i := 1; i < len(arreglodeRutas)-1; i++ {
		rutaCarpeta = rutaCarpeta + "/" + arreglodeRutas[i]
	}
	fmt.Println(rutaCarpeta)
	crearDirectorioSiNoExiste(rutaCarpeta)
	compilarDot(nombre, formato)
}

func compilarDot(ruta string, formato int) {
	if formato == 1 {
		out, err := exec.Command("dot", "-Tjpg", "/home/luis/proyecto2/mbr.dot", "-o", ruta).Output()
		if err != nil {
			log.Fatal(err)

		}
		fmt.Println(string(out))
	} else {
		out, err := exec.Command("dot", "-Tpdf", "/home/luis/proyecto2/mbr.dot", "-o", ruta).Output()
		if err != nil {
			log.Fatal(err)

		}
		fmt.Println(string(out))
	}

}

func mkfile(linea string) {

	delimitador := " "
	delimitador2 := "="
	delimitador4 := "/"
	arreglo := strings.Split(linea, delimitador)
	ruta := ""
	rutaContenido := ""
	padre := false
	tamArchivo := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-r" {
			padre = true
			fmt.Println(padre)
		} else if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				ruta = arregloCI[1]
				fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				ruta = arregloExec[1]
				fmt.Println(ruta)
			}

		} else if arregloExec[0] == "-size" {
			tamArchivo = arregloExec[1]
			fmt.Println(tamArchivo)
		} else if arregloExec[0] == "-cont" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				rutaContenido = arregloCI[1]
				fmt.Println(rutaContenido)
			} else {
				fmt.Println("Sin Comillas")
				rutaContenido = arregloExec[1]
				fmt.Println(rutaContenido)
			}

		}
	}
	contenidoArchivo := ""
	if rutaContenido != "" {
		// leer el arreglo de bytes del archivo
		datosComoBytes, err := ioutil.ReadFile(rutaContenido)
		if err != nil {
			log.Fatal(err)
		}
		// convertir el arreglo a string
		contenidoArchivo = string(datosComoBytes)
	}
	// imprimir el string
	//fmt.Println(datosComoString)

	if len(contenidoArchivo) == 0 {
		contenidoArchivo = "Vacio"
	}

	tamMontadas := len(listaMontadas.particiones)
	auxMontada := PartitionMontada{}
	for i := 0; i < tamMontadas; i++ {
		if listaMontadas.particiones[i].activa {
			auxMontada = listaMontadas.particiones[i]
			break
		}
	}
	//Abrimos el archivo.
	file, err := os.Open(auxMontada.ruta)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	disco := LeerDisco(file, size, 1)
	inicioParticion := 0
	if auxMontada.particion.Part_name == disco.Mbr_partition_1.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_1.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_1.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == disco.Mbr_partition_2.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_2.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_2.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == disco.Mbr_partition_3.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_3.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_3.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == disco.Mbr_partition_4.Part_name {
		posicionSeguir := posicionVacio(disco.Mbr_partition_4.Part_start[:])
		inicioP, _ := strconv.Atoi(string(disco.Mbr_partition_4.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	}

	super := LeerSuperBloque(file, int64(inicioParticion), 1)

	posicionInodoT := posicionVacio(super.S_inode_start[:])
	posicionInodo, _ := strconv.Atoi(string(super.S_inode_start[:posicionInodoT]))
	posicionInodo = inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	posicionBloqueT := posicionVacio(super.S_block_start[:])
	posicionBloque, _ := strconv.Atoi(string(super.S_block_start[:posicionBloqueT]))

	inodoRaiz := LeerInodo(file, int64(posicionInodo), 1)

	arregloCarpetas := strings.Split(ruta, delimitador4)
	tamRuta := len(arregloCarpetas)

	fmt.Println(posicionBloque + tamRuta)
	fmt.Println(arregloCarpetas)
	fmt.Println(inodoRaiz)
	file.Close()
	LeerDirectorioArchivo(inodoRaiz, posicionInodo, posicionBloque, auxMontada.ruta, arregloCarpetas, tamRuta, super, inicioParticion, 1, posicionBloque, 0, "", contenidoArchivo)
}

func LeerDirectorioArchivo(inodo Inodo, posI int, posB int, ruta string, listaNombre []string, nivel int, super SuperBloque, inicioParticion int, numInodo int, posBI int, contador int, CarpetaAnterior string, contenido string) {
	//tamRuta = len(listaNombre)
	nombreN := listaNombre[contador]
	//Abrimos el archivo.
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//inodo = LeerInodo(file, int64(posI), 1)
	tipoInodo := string(inodo.I_type)
	existeCarpeta := false
	iBlockLibre := 0
	iBlockLibreCarpeta := -1
	carpetaTemporal := CarpetaBloque{}
	//bloqueCarpetaLibreT := bitsAString(super.S_first_blo[:])
	//bloqueCarpetaLibre := bitsAInt(super.S_first_blo[:])
	primerVacio := false
	SiguienteInodo := 0
	if tipoInodo == "0" {
		for i := 0; i < 10; i++ {
			if primerVacio || existeCarpeta {
				break
			}
			if inodo.I_block[i] != 0 {
				iBlockLibre++
				posicionCarpetaT := byteToString(inodo.I_block[i])
				posicionCarpeta, _ := strconv.Atoi(posicionCarpetaT)
				carpeta := LeerBloqueCarpeta(file, int64(posicionCarpeta*int(unsafe.Sizeof(CarpetaBloque{})))+int64(posB), 1)
				//posB = posB + int(unsafe.Sizeof(carpeta))
				for j := 0; j < 4; j++ {
					posicionT := posicionVacio(carpeta.B_content[j].B_name[:])
					nombreCarpeta := string(carpeta.B_content[j].B_name[:posicionT])
					if nombreCarpeta == "" {
						iBlockLibreCarpeta = j
						primerVacio = true
						carpetaTemporal = carpeta
						break
					}
					if nombreCarpeta == "." || nombreCarpeta == ".." {

					} else {
						if nombreCarpeta == nombreN {
							existeCarpeta = true
							SiguienteInodo = bitsAInt(carpeta.B_content[j].B_inodo[:])
							break
						}
					}
				}

			}
		}
	}

	if !existeCarpeta {
		tamContenido := len(contenido)
		contadorBloques := 1
		for tamContenido > 60 {
			contadorBloques++
			tamContenido = tamContenido - 60
		}

		if iBlockLibreCarpeta != -1 {
			copy(carpetaTemporal.B_content[iBlockLibreCarpeta].B_name[:], nombreN)
			posicionILT := posicionVacio(super.S_firts_ino[:])
			posicionIL := string(super.S_firts_ino[:posicionILT])
			primerInodoLibre, _ := strconv.Atoi(string(posicionIL))
			SiguienteInodo = primerInodoLibre
			posicionBLT := posicionVacio(super.S_first_blo[:])
			posicionBL := string(super.S_first_blo[:posicionBLT])
			primerBloqueLibre, _ := strconv.Atoi(string(posicionBL))
			copy(carpetaTemporal.B_content[iBlockLibreCarpeta].B_inodo[:], posicionIL)

			copy(super.S_firts_ino[:], strconv.Itoa(primerInodoLibre+1))
			copy(super.S_first_blo[:], strconv.Itoa(primerBloqueLibre+contadorBloques))

			inodoN := Inodo{}
			copy(inodoN.I_gid[:], "1")
			copy(inodoN.I_uid[:], "1")
			copy(inodoN.I_atime[:], string(time.Now().GoString()))
			copy(inodoN.I_ctime[:], string(time.Now().GoString()))
			copy(inodoN.I_mtime[:], string(time.Now().GoString()))
			copy(inodoN.I_block[:], posicionBL)
			copy(inodoN.I_size[:], strconv.Itoa(tamContenido))
			inodoN.I_type = '1'
			copy(inodoN.I_perm[:], "664")
			punteroInicio := 0
			punteroFinal := 60
			for i := 0; i < contadorBloques; i++ {
				nuevaArchivo := ArchivoBloque{}
				if punteroFinal > len(contenido) {
					punteroFinal = len(contenido) - 1
				}
				copy(nuevaArchivo.B_content[:], (contenido[punteroInicio:punteroFinal]))
				inodoN.I_block[i] = strconv.Itoa(primerBloqueLibre + i)[0]
				file.Seek(int64(64*(primerBloqueLibre+i))+int64(posBI), 0)
				var bufferBloque bytes.Buffer
				binary.Write(&bufferBloque, binary.BigEndian, &nuevaArchivo)
				escribirBytes(file, bufferBloque.Bytes())
				punteroFinal = punteroFinal + 61
				punteroInicio = punteroInicio + 61
			}

			file.Seek(0, 0)
			file.Seek(int64(inicioParticion), 0)
			var bufferOrigenS bytes.Buffer
			binary.Write(&bufferOrigenS, binary.BigEndian, &super)
			escribirBytes(file, bufferOrigenS.Bytes())

			numBloqueT := string(inodo.I_block[iBlockLibre-1])
			numBloque, _ := strconv.Atoi(string(numBloqueT))
			axu1 := int64(64*numBloque) + int64(posBI)
			fmt.Println(axu1)
			tamInodo := int(unsafe.Sizeof(Inodo{}))
			axu2 := int64(tamInodo*primerInodoLibre) + int64(posI)
			fmt.Println(axu2)

			file.Seek(int64(64*numBloque)+int64(posBI), 0)
			var bufferBloqueN bytes.Buffer
			binary.Write(&bufferBloqueN, binary.BigEndian, &carpetaTemporal)
			escribirBytes(file, bufferBloqueN.Bytes())

			file.Seek(int64(tamInodo*primerInodoLibre)+int64(posI), 0)
			var bufferInodo bytes.Buffer
			binary.Write(&bufferInodo, binary.BigEndian, &inodoN)
			escribirBytes(file, bufferInodo.Bytes())

		} else {
			nuevaCarpeta := CarpetaBloque{}
			copy(nuevaCarpeta.B_content[0].B_name[:], nombreN)
			posicionILT := posicionVacio(super.S_firts_ino[:])
			posicionIL := string(super.S_firts_ino[:posicionILT])
			primerInodoLibre, _ := strconv.Atoi(string(posicionIL))
			SiguienteInodo = primerInodoLibre
			posicionBLT := posicionVacio(super.S_first_blo[:])
			posicionBL := string(super.S_first_blo[:posicionBLT])
			primerBloqueLibre, _ := strconv.Atoi(string(posicionBL))
			copy(nuevaCarpeta.B_content[0].B_inodo[:], strconv.Itoa(primerInodoLibre))
			inodo.I_block[iBlockLibre] = byte(posicionBL[0])

			copy(super.S_firts_ino[:], strconv.Itoa(primerInodoLibre+1))

			copy(super.S_first_blo[:], strconv.Itoa(primerBloqueLibre+1+contadorBloques))

			inodoN := Inodo{}
			copy(inodoN.I_gid[:], "1")
			copy(inodoN.I_uid[:], "1")
			copy(inodoN.I_atime[:], string(time.Now().GoString()))
			copy(inodoN.I_ctime[:], string(time.Now().GoString()))
			copy(inodoN.I_mtime[:], string(time.Now().GoString()))
			copy(inodoN.I_block[:], strconv.Itoa(primerBloqueLibre+1))
			inodoN.I_type = '1'
			copy(inodoN.I_size[:], strconv.Itoa(tamContenido))
			copy(inodoN.I_perm[:], "664")

			punteroInicio := 0
			punteroFinal := 60
			for i := 0; i < contadorBloques; i++ {
				nuevaArchivo := ArchivoBloque{}
				if punteroFinal > len(contenido) {
					punteroFinal = len(contenido) - 1
				}
				copy(nuevaArchivo.B_content[:], (contenido[punteroInicio:punteroFinal]))
				inodoN.I_block[i] = strconv.Itoa(primerBloqueLibre + i + 1)[0]
				file.Seek(int64(64*(primerBloqueLibre+i))+int64(posBI), 0)
				var bufferBloque bytes.Buffer
				binary.Write(&bufferBloque, binary.BigEndian, &nuevaArchivo)
				escribirBytes(file, bufferBloque.Bytes())
				punteroFinal = punteroFinal + 61
				punteroInicio = punteroInicio + 61
			}

			tamInodo := int(unsafe.Sizeof(Inodo{}))
			file.Seek(0, 0)
			file.Seek(int64(inicioParticion), 0)
			var bufferOrigenW bytes.Buffer
			binary.Write(&bufferOrigenW, binary.BigEndian, &super)
			escribirBytes(file, bufferOrigenW.Bytes())

			//numBloqueT := string(inodo.I_block[iBlockLibre-1])
			//numBloque, _ := strconv.Atoi(string(numBloqueT))

			file.Seek(int64(tamInodo*(numInodo-1))+int64(posI), 0)
			var bufferInodoActual bytes.Buffer
			binary.Write(&bufferInodoActual, binary.BigEndian, &inodo)
			escribirBytes(file, bufferInodoActual.Bytes())

			file.Seek(int64(64*primerBloqueLibre)+int64(posBI), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaCarpeta)
			escribirBytes(file, bufferBloque.Bytes())

			file.Seek(int64(tamInodo*(primerInodoLibre))+int64(posI), 0)
			var bufferInodo bytes.Buffer
			binary.Write(&bufferInodo, binary.BigEndian, &inodoN)
			escribirBytes(file, bufferInodo.Bytes())

		}
	}
	if nivel > 1 {
		inodoTem := LeerInodo(file, int64(posI)+(int64(unsafe.Sizeof(Inodo{}))*int64(SiguienteInodo)), 1)
		tipoTemInodo := string(inodoTem.I_type)
		if tipoTemInodo == "0" {
			LeerDirectorioArchivo(inodoTem, int(posI), posB, ruta, listaNombre, (nivel - 1), super, inicioParticion, SiguienteInodo+1, posBI, contador+1, listaNombre[contador], contenido)
		} else {
			//bloqueArc := ArchivoBloque{}
			//posB = posB + (int(unsafe.Sizeof(bloqueArc) * 16))
		}
	}
}

func mkgrp(linea string) {
	disco := Partition{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	name := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-name" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				name = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				name = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(string(disco.Part_name[:posicionSeguir]))
		}
	}

	tamMontadas := len(listaMontadas.particiones)
	auxMontada := PartitionMontada{}
	for i := 0; i < tamMontadas; i++ {
		if listaMontadas.particiones[i].activa {
			auxMontada = listaMontadas.particiones[i]
			break
		}
	}
	//Abrimos el archivo.
	file, err := os.OpenFile(auxMontada.ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	discoA := LeerDisco(file, size, 1)
	inicioParticion := 0
	if auxMontada.particion.Part_name == discoA.Mbr_partition_1.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_1.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_1.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_2.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_2.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_2.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_3.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_3.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_3.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_4.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_4.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_4.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	}

	super := LeerSuperBloque(file, int64(inicioParticion), 1)

	posicionInodoT := posicionVacio(super.S_inode_start[:])
	posicionInodo, _ := strconv.Atoi(string(super.S_inode_start[:posicionInodoT]))
	posicionInodo = inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	posicionBloqueT := posicionVacio(super.S_block_start[:])
	posicionBloque, _ := strconv.Atoi(string(super.S_block_start[:posicionBloqueT]))

	inodoUsuarios := LeerInodo(file, int64(posicionInodo)+int64(unsafe.Sizeof(Inodo{})), 1)
	contenido := ""
	for i := 0; i < 10; i++ {
		if inodoUsuarios.I_block[i] != 0 {
			numBloqueT := byteToString(inodoUsuarios.I_block[i])
			numBloque, _ := strconv.Atoi(numBloqueT)
			ArchivoTem := LeerBloqueArchivo(file, int64(posicionBloque)+(int64(unsafe.Sizeof(ArchivoBloque{}))*int64(numBloque)), 1)
			contenido = contenido + bitsAString(ArchivoTem.B_content[:])
		}
	}

	arregloDatos := strings.Split(contenido, "\n")
	//arreglo[len(arreglo)-1]
	contadorGrupos := 0
	for ident, linea := range arregloDatos {
		listaG := strings.Split(linea, ",")
		if ident != len(arregloDatos)-1 {
			if listaG[1] == "G" || listaG[1] == "g" {
				contadorGrupos++
				if listaG[2] == name {
					return
				}
			}
		}
	}

	contenido = contenido + strconv.Itoa(contadorGrupos+1) + ",G," + name + "\n"
	tamGrupos := len(contenido)

	cantidadBloques := 1

	for tamGrupos > 60 {
		cantidadBloques++
		tamGrupos = tamGrupos - 60
	}

	punteroInicio := 0
	punteroFinal := 60
	for i := 0; i < cantidadBloques; i++ {
		if inodoUsuarios.I_block[i] != 0 {
			nuevaArchivo := ArchivoBloque{}
			if punteroFinal > len(contenido) {
				punteroFinal = len(contenido)
			}
			copy(nuevaArchivo.B_content[:], (contenido[punteroInicio:punteroFinal]))
			//inodoN.I_block[i] = strconv.Itoa(primerBloqueLibre + i + 1)[0]
			apuntador, _ := strconv.Atoi(byteToString(inodoUsuarios.I_block[i]))
			file.Seek(int64(64*(apuntador))+int64(posicionBloque), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaArchivo)
			escribirBytes(file, bufferBloque.Bytes())
			punteroFinal = punteroFinal + 61
			punteroInicio = punteroInicio + 61

		} else {
			nuevaArchivo := ArchivoBloque{}
			if punteroFinal > len(contenido) {
				punteroFinal = len(contenido)
			}
			copy(nuevaArchivo.B_content[:], (contenido[punteroInicio:punteroFinal]))
			//inodoN.I_block[i] = strconv.Itoa(primerBloqueLibre + i + 1)[0]
			posicionBLT := posicionVacio(super.S_first_blo[:])
			posicionBL := string(super.S_first_blo[:posicionBLT])
			primerBloqueLibre, _ := strconv.Atoi(string(posicionBL))
			inodoUsuarios.I_block[i] = strconv.Itoa(primerBloqueLibre)[0]
			apuntador, _ := strconv.Atoi(byteToString(inodoUsuarios.I_block[i]))

			file.Seek(int64(posicionInodo)+int64(unsafe.Sizeof(Inodo{})), 0)
			var bufferInodoActual bytes.Buffer
			binary.Write(&bufferInodoActual, binary.BigEndian, &inodoUsuarios)
			escribirBytes(file, bufferInodoActual.Bytes())

			file.Seek(int64(64*(apuntador))+int64(posicionBloque), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaArchivo)
			escribirBytes(file, bufferBloque.Bytes())
			punteroFinal = punteroFinal + 61
			punteroInicio = punteroInicio + 61

			copy(super.S_first_blo[:], strconv.Itoa(primerBloqueLibre+1))
			file.Seek(0, 0)
			file.Seek(int64(inicioParticion), 0)
			var bufferOrigenW bytes.Buffer
			binary.Write(&bufferOrigenW, binary.BigEndian, &super)
			escribirBytes(file, bufferOrigenW.Bytes())
		}
	}

	file.Close()

}

func login(linea string) bool {
	disco := Partition{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	usuario := ""
	pass := ""
	id := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-usuario" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				usuario = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				usuario = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(usuario)
		} else if arregloExec[0] == "-password" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				pass = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				pass = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(pass)
		} else if arregloExec[0] == "-id" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				id = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				id = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(id)
		}
	}

	tamMontadas := len(listaMontadas.particiones)
	auxMontada := PartitionMontada{}
	for i := 0; i < tamMontadas; i++ {
		if listaMontadas.particiones[i].activa {
			auxMontada = listaMontadas.particiones[i]
			break
		}
	}
	//Abrimos el archivo.
	file, err := os.OpenFile(auxMontada.ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	discoA := LeerDisco(file, size, 1)
	inicioParticion := 0
	if auxMontada.particion.Part_name == discoA.Mbr_partition_1.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_1.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_1.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_2.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_2.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_2.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_3.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_3.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_3.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_4.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_4.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_4.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	}

	super := LeerSuperBloque(file, int64(inicioParticion), 1)

	posicionInodoT := posicionVacio(super.S_inode_start[:])
	posicionInodo, _ := strconv.Atoi(string(super.S_inode_start[:posicionInodoT]))
	posicionInodo = inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	posicionBloqueT := posicionVacio(super.S_block_start[:])
	posicionBloque, _ := strconv.Atoi(string(super.S_block_start[:posicionBloqueT]))

	inodoUsuarios := LeerInodo(file, int64(posicionInodo)+int64(unsafe.Sizeof(Inodo{})), 1)
	contenido := ""
	for i := 0; i < 10; i++ {
		if inodoUsuarios.I_block[i] != 0 {
			numBloqueT := byteToString(inodoUsuarios.I_block[i])
			numBloque, _ := strconv.Atoi(numBloqueT)
			ArchivoTem := LeerBloqueArchivo(file, int64(posicionBloque)+(int64(unsafe.Sizeof(ArchivoBloque{}))*int64(numBloque)), 1)
			contenido = contenido + bitsAString(ArchivoTem.B_content[:])
		}
	}

	arregloDatos := strings.Split(contenido, "\n")
	//arreglo[len(arreglo)-1]

	//existeU := false
	for ident, linea := range arregloDatos {
		listaG := strings.Split(linea, ",")
		if ident != len(arregloDatos)-1 {
			if listaG[1] == "U" || listaG[1] == "u" {
				if listaG[3] == usuario && listaG[4] == pass {
					userActive.id = id
					userActive.user = usuario
					userActive.pass = pass
					userActive.grupo = listaG[0]
					userActive.active = true
					userActive.actual = auxMontada
					if usuario == "root" {
						userActive.admin = true
					}
					file.Close()
					fmt.Println("Bienvenido al sistema")
					return true
				}
			}
		}
	}
	fmt.Println("Error de credenciales")
	file.Close()
	return false

}

func mkusr(linea string) {
	disco := Partition{}
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	usuario := ""
	pass := ""
	group := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-usuario" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				usuario = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				usuario = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(usuario)
		} else if arregloExec[0] == "-pwd" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				pass = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				pass = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(pass)
		} else if arregloExec[0] == "-grp" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]
				copy(disco.Part_name[:], arregloCI[1])
				group = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]
				copy(disco.Part_name[:], arregloExec[1])
				group = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(group)
		}
	}

	tamMontadas := len(listaMontadas.particiones)
	auxMontada := PartitionMontada{}
	for i := 0; i < tamMontadas; i++ {
		if listaMontadas.particiones[i].activa {
			auxMontada = listaMontadas.particiones[i]
			break
		}
	}
	//Abrimos el archivo.
	file, err := os.OpenFile(auxMontada.ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	discoA := LeerDisco(file, size, 1)
	inicioParticion := 0
	if auxMontada.particion.Part_name == discoA.Mbr_partition_1.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_1.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_1.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_2.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_2.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_2.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_3.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_3.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_3.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	} else if auxMontada.particion.Part_name == discoA.Mbr_partition_4.Part_name {
		posicionSeguir := posicionVacio(discoA.Mbr_partition_4.Part_start[:])
		inicioP, _ := strconv.Atoi(string(discoA.Mbr_partition_4.Part_start[:posicionSeguir]))
		inicioParticion = inicioP
	}

	super := LeerSuperBloque(file, int64(inicioParticion), 1)

	posicionInodoT := posicionVacio(super.S_inode_start[:])
	posicionInodo, _ := strconv.Atoi(string(super.S_inode_start[:posicionInodoT]))
	posicionInodo = inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	posicionBloqueT := posicionVacio(super.S_block_start[:])
	posicionBloque, _ := strconv.Atoi(string(super.S_block_start[:posicionBloqueT]))

	inodoUsuarios := LeerInodo(file, int64(posicionInodo)+int64(unsafe.Sizeof(Inodo{})), 1)
	contenido := ""
	for i := 0; i < 10; i++ {
		if inodoUsuarios.I_block[i] != 0 {
			numBloqueT := byteToString(inodoUsuarios.I_block[i])
			numBloque, _ := strconv.Atoi(numBloqueT)
			ArchivoTem := LeerBloqueArchivo(file, int64(posicionBloque)+(int64(unsafe.Sizeof(ArchivoBloque{}))*int64(numBloque)), 1)
			contenido = contenido + bitsAString(ArchivoTem.B_content[:])
		}
	}

	arregloDatos := strings.Split(contenido, "\n")
	//arreglo[len(arreglo)-1]
	contadorGrupos := 0
	existeG := false
	//existeU := false
	for ident, linea := range arregloDatos {
		listaG := strings.Split(linea, ",")
		if ident != len(arregloDatos)-1 {
			if listaG[1] == "G" || listaG[1] == "g" {

				if listaG[2] == group {
					existeG = true
				}
			} else if listaG[1] == "U" || listaG[1] == "u" {
				contadorGrupos++
				if listaG[3] == usuario {
					return
				}
			}
		}
	}
	if !existeG {
		return
	}
	contenido = contenido + strconv.Itoa(contadorGrupos+1) + ",U," + group + "," + usuario + "," + pass + "\n"

	tamGrupos := len(contenido)

	cantidadBloques := 1

	for tamGrupos > 60 {
		cantidadBloques++
		tamGrupos = tamGrupos - 60
	}

	punteroInicio := 0
	punteroFinal := 60
	for i := 0; i < cantidadBloques; i++ {
		if inodoUsuarios.I_block[i] != 0 {
			nuevaArchivo := ArchivoBloque{}
			if punteroFinal > len(contenido) {
				punteroFinal = len(contenido)
			}
			copy(nuevaArchivo.B_content[:], (contenido[punteroInicio:punteroFinal]))
			//inodoN.I_block[i] = strconv.Itoa(primerBloqueLibre + i + 1)[0]
			apuntador, _ := strconv.Atoi(byteToString(inodoUsuarios.I_block[i]))
			file.Seek(int64(64*(apuntador))+int64(posicionBloque), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaArchivo)
			escribirBytes(file, bufferBloque.Bytes())
			punteroFinal = punteroFinal + 60
			punteroInicio = punteroInicio + 60

		} else {
			nuevaArchivo := ArchivoBloque{}
			if punteroFinal > len(contenido) {
				punteroFinal = len(contenido)
			}
			copy(nuevaArchivo.B_content[:], (contenido[punteroInicio:punteroFinal]))
			//inodoN.I_block[i] = strconv.Itoa(primerBloqueLibre + i + 1)[0]
			posicionBLT := posicionVacio(super.S_first_blo[:])
			posicionBL := string(super.S_first_blo[:posicionBLT])
			primerBloqueLibre, _ := strconv.Atoi(string(posicionBL))
			inodoUsuarios.I_block[i] = strconv.Itoa(primerBloqueLibre)[0]
			apuntador, _ := strconv.Atoi(byteToString(inodoUsuarios.I_block[i]))

			file.Seek(int64(posicionInodo)+int64(unsafe.Sizeof(Inodo{})), 0)
			var bufferInodoActual bytes.Buffer
			binary.Write(&bufferInodoActual, binary.BigEndian, &inodoUsuarios)
			escribirBytes(file, bufferInodoActual.Bytes())

			file.Seek(int64(64*(apuntador))+int64(posicionBloque), 0)
			var bufferBloque bytes.Buffer
			binary.Write(&bufferBloque, binary.BigEndian, &nuevaArchivo)
			escribirBytes(file, bufferBloque.Bytes())
			punteroFinal = punteroFinal + 60
			punteroInicio = punteroInicio + 60

			copy(super.S_first_blo[:], strconv.Itoa(primerBloqueLibre+1))
			file.Seek(0, 0)
			file.Seek(int64(inicioParticion), 0)
			var bufferOrigenW bytes.Buffer
			binary.Write(&bufferOrigenW, binary.BigEndian, &super)
			escribirBytes(file, bufferOrigenW.Bytes())
		}
	}

	file.Close()

}

func logout() {
	if !userActive.active {
		fmt.Println("Sin usuario logueado")
		return
	}
	userActive.id = ""
	userActive.user = ""
	userActive.pass = ""
	userActive.grupo = ""
	userActive.active = false
	userActive.admin = false
	fmt.Println("Adios")
}

func repDisco(p string, id string) {
	path := ""
	//spr := SuperBloque{}
	//inode := Inodo{}
	//inodeArchivo := Inodo{}

	//partition := Partition{}
	auxMontada := PartitionMontada{}
	for _, actual := range listaMontadas.particiones {
		if actual.id == id {
			auxMontada = actual
			//partition = auxMontada.particion
			path = auxMontada.ruta
		}
	}

	//Abrimos el archivo.
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//inicioParticion, _ := strconv.Atoi(string(auxMontada.particion.Part_start[:posicionVacio(auxMontada.particion.Part_start[:])]))

	//obtenemor el size del control para empezar a leer desde ahi
	var size int64 = int64(unsafe.Sizeof(control{}))
	disk := LeerDisco(file, size, 1)

	var partitions [4]Partition
	var extended Partition
	fmt.Println(extended)
	ext := false
	fmt.Println(ext)
	partitions[0] = disk.Mbr_partition_1
	partitions[1] = disk.Mbr_partition_2
	partitions[2] = disk.Mbr_partition_3
	partitions[3] = disk.Mbr_partition_4

	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == '1' {
			if partitions[i].Part_type == 'E' {
				ext = true
				extended = partitions[i]
			}
		}
	}

	content := ""

	content += "digraph G{\n" +
		"rankdir=TB;\n" +
		"forcelabels= true;\n" +
		"graph [ dpi = \"600\" ]; \n" +
		"node [shape = plaintext];\n"
	content += "nodo1 [label = <<table>\n"
	content += "<tr>\n"

	content += "<td ROWSPAN='2'> MBR </td>\n"

	var auxParticiones []Traslado
	cantidad := 0
	var extendida Partition

	if disk.Mbr_partition_1.Part_size[0] != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		aux.comienzo = bitsAInt(disk.Mbr_partition_1.Part_start[:])
		aux.fin = bitsAInt(disk.Mbr_partition_1.Part_start[:]) + bitsAInt(disk.Mbr_partition_1.Part_size[:])
		aux.anterior = bitsAInt(disk.Mbr_partition_1.Part_start[:])
		if disk.Mbr_partition_2.Part_size[0] != 0 {
			//auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
			aux.siguiente = bitsAInt(disk.Mbr_partition_2.Part_start[:])
		} else {
			aux.siguiente = bitsAInt(disk.Mbr_tamano[:]) + bitsAInt(disk.Mbr_partition_1.Part_start[:])
		}
		auxParticiones = append(auxParticiones, aux)
		cantidad++
		if disk.Mbr_partition_1.Part_type == 'e' {
			extendida = disk.Mbr_partition_1
		}
	}
	if disk.Mbr_partition_2.Part_size[0] != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		aux.comienzo = bitsAInt(disk.Mbr_partition_2.Part_start[:])
		aux.fin = bitsAInt(disk.Mbr_partition_2.Part_start[:]) + bitsAInt(disk.Mbr_partition_2.Part_size[:])
		aux.anterior = auxParticiones[0].fin //discoAux.mbr_partition_2.part_start-libre;
		if disk.Mbr_partition_3.Part_size[0] != 0 {
			//auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
			aux.siguiente = bitsAInt(disk.Mbr_partition_3.Part_start[:])
		} else {
			aux.siguiente = bitsAInt(disk.Mbr_tamano[:]) + bitsAInt(disk.Mbr_partition_1.Part_start[:])
		}
		auxParticiones = append(auxParticiones, aux)
		cantidad++
		if disk.Mbr_partition_2.Part_type == 'e' {
			extendida = disk.Mbr_partition_2
		}
	}
	if disk.Mbr_partition_3.Part_size[0] != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		aux.comienzo = bitsAInt(disk.Mbr_partition_3.Part_start[:])
		aux.fin = bitsAInt(disk.Mbr_partition_3.Part_start[:]) + bitsAInt(disk.Mbr_partition_3.Part_size[:])
		aux.anterior = auxParticiones[1].fin //discoAux.mbr_partition_3.part_start-libre;
		if disk.Mbr_partition_4.Part_size[0] != 0 {
			//auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
			aux.siguiente = bitsAInt(disk.Mbr_partition_4.Part_start[:])
		} else {
			aux.siguiente = bitsAInt(disk.Mbr_tamano[:]) + bitsAInt(disk.Mbr_partition_1.Part_start[:])
		}
		auxParticiones = append(auxParticiones, aux)
		cantidad++
		if disk.Mbr_partition_3.Part_type == 'e' {
			extendida = disk.Mbr_partition_3
		}
	}
	if disk.Mbr_partition_4.Part_size[0] != 0 {
		aux := Traslado{}
		aux.idNum = cantidad + 1
		aux.comienzo = bitsAInt(disk.Mbr_partition_4.Part_start[:])
		aux.fin = bitsAInt(disk.Mbr_partition_4.Part_start[:]) + bitsAInt(disk.Mbr_partition_4.Part_size[:])
		aux.anterior = auxParticiones[2].fin //discoAux.mbr_partition_4.part_start-libre;
		//if(cantidad !=0)
		//{
		//  auxParticiones.at(cantidad-1).siguiente = aux.comienzo - (auxParticiones.at(cantidad-1).fin);
		//}
		aux.siguiente = bitsAInt(disk.Mbr_tamano[:]) + bitsAInt(disk.Mbr_partition_1.Part_start[:])
		auxParticiones = append(auxParticiones, aux)
		cantidad++
		if disk.Mbr_partition_4.Part_type == 'e' {
			extendida = disk.Mbr_partition_4
		}
	}

	for i := 0; i < cantidad; i++ {
		if partitions[i].Part_size[0] != 0 {
			num1 := bitsAInt(partitions[i].Part_size[:])
			num2 := bitsAInt(disk.Mbr_tamano[:])
			var num float64 = float64(num1) / float64(num2)
			num = num * 100
			if partitions[i].Part_type == 'p' {
				content += "<td ROWSPAN='2'>Primaria " + fmt.Sprintf("%f", num) + "%</td>\n"
			} else if partitions[i].Part_type == 'e' {
				//content +="<tr>\n";
				content += "<td ROWSPAN='2'>Extendida " + fmt.Sprintf("%f", num) + "%</td>\n"
				//content += "</tr>\n";
				//vector<EBR> logicas=ObtenerLogicas(partitions[i],path.c_str());
				libreE := 0
				//for _, logica := range logicas{
				//	float num1=logica.part_s;
				//	float num2=partitions[i].part_s;
				//	float num=num1/num2;
				//	content += "<td ROWSPAN='2'>Logica "+ to_string(num)+"%</td>\n";
				//	libreE=libreE+num;
				//}
				if libreE != 0 {
					result := 100 - libreE
					content += "<td ROWSPAN='2'>Libre " + strconv.Itoa(result) + "%</td>\n"
				}
				//if(logicas.size()==0){
				//content +="<tr>\n";
				//	content += "<td ROWSPAN='2'>EBR</td>\n";
				//content += "</tr>\n";
				//}

			}
			if auxParticiones[i].siguiente == auxParticiones[i].fin {

			} else {
				num1 := auxParticiones[i].siguiente - auxParticiones[i].fin
				num2 := bitsAInt(disk.Mbr_tamano[:])
				var num float64 = float64(num1) / float64(num2)
				num = num * 100
				content += "<td ROWSPAN='2'>Libre " + fmt.Sprintf("%f", num) + "%</td>\n"
			}

		}
	}
	fmt.Println(extendida)
	content += "</tr>\n\n"

	content += "</table>>];\n}\n"
	crearDot(content, p, 1)
}

func repFile(p string, id string, nombre string) {
	path := ""
	spr := SuperBloque{}
	inode := Inodo{}
	//inodeArchivo := Inodo{}

	//partition := Partition{}
	auxMontada := PartitionMontada{}
	for _, actual := range listaMontadas.particiones {
		if actual.id == id {
			auxMontada = actual
			//partition = auxMontada.particion
			path = auxMontada.ruta
		}
	}

	//Abrimos el archivo.
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	posicionNombre := len(strings.Split(nombre, "/"))
	nombreActual := strings.Split(nombre, "/")[posicionNombre-1]
	inicioParticion, _ := strconv.Atoi(string(auxMontada.particion.Part_start[:posicionVacio(auxMontada.particion.Part_start[:])]))
	inicioInodos := inicioParticion + int(unsafe.Sizeof(SuperBloque{}))
	spr = LeerSuperBloque(file, int64(inicioParticion), 1)

	inode = LeerInodo(file, int64(inicioInodos), 1)

	posT := posicionVacio(spr.S_firts_ino[:])
	inodoLibre, _ := strconv.Atoi(string(spr.S_firts_ino[:posT]))

	posBT := posicionVacio(spr.S_block_start[:])
	posBloques, _ := strconv.Atoi(string(spr.S_block_start[:posBT]))

	content := "digraph G{\n" +
		"rankdir=LR;\n" +
		"graph [ dpi = \"600\" ]; \n" +
		"forcelabels= true;\n" +
		"node [shape = plaintext];\n"
	content += "nodo1 [label = <<table>\n"
	content += "<tr>\n"

	content += "<td>" + nombre + "</td>\n"
	content += "<td>"
	inodoArchvio := 0
	for i := 0; i < inodoLibre; i++ {

		if inode.I_type == '0' {
			for j := 0; j < 10; j++ {
				if inode.I_block[j] != 0 {

					tem, _ := strconv.Atoi(string(inode.I_block[j]))
					//fmt.Println(tem)
					foldertmp := LeerBloqueCarpeta(file, int64(posBloques+(int(unsafe.Sizeof(CarpetaBloque{}))*tem)), 1)

					for k := 0; k < 4; k++ {
						ctmp := ""
						ctmp += bitsAString(foldertmp.B_content[k].B_name[:])
						if ctmp == nombreActual {
							inodoArchvio = bitsAInt(foldertmp.B_content[k].B_inodo[:])
							inode = LeerInodo(file, int64(inicioInodos)+(int64(unsafe.Sizeof(Inodo{}))*int64(inodoArchvio)), 1)

							contenido := ""
							for k := 0; k < 10; k++ {
								if inode.I_block[k] != 0 {
									numBloqueT := byteToString(inode.I_block[i])
									numBloque, _ := strconv.Atoi(numBloqueT)
									ArchivoTem := LeerBloqueArchivo(file, int64(posBloques)+(int64(unsafe.Sizeof(ArchivoBloque{}))*int64(numBloque)), 1)
									contenido = contenido + bitsAString(ArchivoTem.B_content[:])
								}
							}
							content = content + contenido
						}
					}
				}
			}
		}
		inode = LeerInodo(file, int64(inicioInodos)+(int64(unsafe.Sizeof(Inodo{}))*int64(i+1)), 1)

	}
	content += "</td>\n"
	content += "</tr>\n\n"
	content += "</table>>];\n}\n"
	//content += "\n\n}\n"
	fmt.Println(content)
	crearDot(content, p, 1)
}

func repSuper(p string, id string) {
	path := ""
	spr := SuperBloque{}

	//inodeArchivo := Inodo{}

	//partition := Partition{}
	auxMontada := PartitionMontada{}
	for _, actual := range listaMontadas.particiones {
		if actual.id == id {
			auxMontada = actual
			//partition = auxMontada.particion
			path = auxMontada.ruta
		}
	}

	//Abrimos el archivo.
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	inicioParticion, _ := strconv.Atoi(string(auxMontada.particion.Part_start[:posicionVacio(auxMontada.particion.Part_start[:])]))

	spr = LeerSuperBloque(file, int64(inicioParticion), 1)

	content := "digraph G{\n" +
		"rankdir=LR;\n" +
		"graph [ dpi = \"600\" ]; \n" +
		"forcelabels= true;\n" +
		"node [shape = plaintext];\n"
	content += "nodo1 [label = <<table>\n"
	content += "<tr><td>SuperBloque"
	content += "</td></tr>\n"
	content += "<tr><td>s_filesystem_type/ " + bitsAString(spr.S_filesystem_type[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_inodes_count/ " + bitsAString(spr.S_inodes_count[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_blocks_count/ " + bitsAString(spr.S_blocks_count[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_free_blocks_count/ " + bitsAString(spr.S_free_blocks_count[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_free_inodes_count/ " + bitsAString(spr.S_free_inodes_count[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_mtime/ " + time.Now().Local().String()
	content += "</td></tr>\n"
	content += "<tr><td>s_mnt_count/ " + bitsAString(spr.S_mnt_count[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_magic/ " + bitsAString(spr.S_magic[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_inode_size/ " + strconv.Itoa(int(unsafe.Sizeof(Inodo{})))
	content += "</td></tr>\n"
	content += "<tr><td>s_block_size/ " + strconv.Itoa(int(unsafe.Sizeof(CarpetaBloque{})))
	content += "</td></tr>\n"
	content += "<tr><td>s_firts_ino/ " + bitsAString(spr.S_firts_ino[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_first_blo/ " + bitsAString(spr.S_first_blo[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_bm_inode_start/ " + bitsAString(spr.S_bm_inode_start[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_bm_block_start/ " + bitsAString(spr.S_bm_block_start[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_inode_start/ " + bitsAString(spr.S_inode_start[:])
	content += "</td></tr>\n"
	content += "<tr><td>s_block_start/ " + bitsAString(spr.S_block_start[:])
	content += "</td></tr>\n"
	content += "</table>>];\n}\n"
	//content += "\n\n}\n"
	fmt.Println(content)
	crearDot(content, p, 1)
}

func pause() {
	var eleccion string //Declarar variable y tipo antes de escanear, esto es obligatorio
	fmt.Println("Pausa")
	fmt.Scanln(&eleccion)
	fmt.Println("Salio de pausa")

}

func reportes(linea string) string {

	imagen := ""
	delimitador := " "
	delimitador2 := "="
	arreglo := strings.Split(linea, delimitador)
	name := ""
	path := ""
	id := ""
	ruta := ""
	for _, separado := range arreglo {
		arregloExec := strings.Split(separado, delimitador2)
		fmt.Println(arregloExec[0])
		//fmt.Println(arregloExec[1])
		if arregloExec[0] == "-name" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]

				name = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]

				name = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(name)
		} else if arregloExec[0] == "-path" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]

				path = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]

				path = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(path)
		} else if arregloExec[0] == "-id" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]

				id = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]

				id = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(id)
		} else if arregloExec[0] == "-ruta" {
			if arregloExec[1][0] == '"' {
				fmt.Println("Comillas")
				delimitador3 := '"'
				arregloCI := strings.Split(linea, string(delimitador3))
				//ruta = arregloCI[1]

				ruta = arregloCI[1]
				//disco.part_name = []byte(arregloCI[1])
				//fmt.Println(ruta)
			} else {
				fmt.Println("Sin Comillas")
				//ruta = arregloExec[1]

				ruta = arregloExec[1]
				//disco.part_name = []byte(arregloExec[1])
				//fmt.Println(ruta)
			}
			//posicionSeguir := posicionVacio(disco.Part_name[:])
			fmt.Println(ruta)
		}
	}

	if name == "disk" {
		repDisco(path, id)
	} else if name == "tree" {
		tree(path, id)
	} else if name == "file" {
		repFile(path, id, ruta)
	} else if name == "sb" {
		repSuper(path, id)
	}
	imagen = convertirImagen(path)
	return imagen
}
