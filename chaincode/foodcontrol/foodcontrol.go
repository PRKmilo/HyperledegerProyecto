package main

import (
	//modulo para convertir json a bytes
	"encoding/json"
	//modulo para mostrar por pantalla
	"fmt"
	//modulo de hyperledger
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for control the food
type SmartContract struct {
	contractapi.Contract
}

//IMPORTANTE!!!! EL ACTIVO QUE QUEREMOS PERSISITIR

type Escritura struct {
	Nmatricula string `json:"numeroMatricula"`
	NombreTitular string `json:"nombreTitular"`
	DescripcionActo string `json:"descripcionActo"`
	ObjetoEscritura string `json:"objetoEscritura"`
	Fecha string `json:"fecha"`
	FormaPago string `json:"formaPago"`
	Contenido string `json:"contenido"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, nmatricula string, nombreTitular string, descripcionActo string, objetoEscritura string ,fecha string, formaPago string, contenido string) error {

	//Validaciones de sintaxis

	//validaciones de negocio

	//validamos primero si el alumno ya existe en la blockchain
	resEscritura, err :=s.Query(ctx, Nmatricula)
	if resEscritura != nil {
		fmt.Printf("Alumno already exist error: %s", err.Error())
		return err
	}

	escritura := Escritura{
		Nmatricula:  nmatricula,
		NombreTitular: nombreTitular,
		descripcionActo: descripcionActo,
		Objetoescritura: objetoEscritura,
		Fecha: time.Now().Format(time.RFC3339),
		FormaPago: formaPago,
		Contenido: contenido
	}

	//transformo alumno a bytes
	//lo resivo como json, opero en go como structura y lo persisito en bytes.
	escrituraAsBytes, err := json.Marshal(escritura)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	//PutState es el que nos permite guardar en el libro distribuido
	return ctx.GetStub().PutState(nmatricula, escrituraAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, nmatricula string) (*Escritura, error) {

	//busco el alumno por su id en la blockchain  
	escrituraAsBytes, err := ctx.GetStub().GetState(nmatricula)

	//validamos si tenemos un error
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	//validamos si no encontro el alumno
	if escrituraAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", nmatricula)
	}

	//declaramos una estructura Alumno
	escritura := new(Escritura)

	//convertimos de bytes a la estructura alumno
	err = json.Unmarshal(escrituraAsBytes, escritura)
	//validamos si existio algun error
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
	}

	//retornamos el alumno o null
	return escritura, nil
}


func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, nmatricula string, nombreTitular string, descripcionActo string, objetoEscritura string, fecha string, formaPago string, contenido string) error {
    // Obtén el registro actual
    escritura, err := s.Query(ctx, nmatricula)
    if err != nil {
        return err
    }

    // Parse la fecha de creación
    fechaCreacion, err := time.Parse(time.RFC3339, escritura.FechaCreacion)
    if err != nil {
        return fmt.Errorf("Error parsing fechaCreacion: %s", err.Error())
    }

    // Compara la fecha actual con la fecha de creación más 2 años
    dosAnos := time.Now().Sub(fechaCreacion)
    if dosAnos < 2*365*24*time.Hour { // 2 años en horas
        return fmt.Errorf("El objeto no puede ser modificado antes de 2 años desde su creación")
    }

    // Actualiza el registro
    escritura.NombreTitular = nombreTitular
    escritura.DescripcionActo = descripcionActo
    escritura.ObjetoEscritura = objetoEscritura
    escritura.Fecha = fecha
    escritura.FormaPago = formaPago
    escritura.Contenido = contenido

    // Convierte el registro actualizado a bytes
    escrituraAsBytes, err := json.Marshal(escritura)
    if err != nil {
        return fmt.Errorf("Marshal error: %s", err.Error())
    }

    // Guarda el registro actualizado
    return ctx.GetStub().PutState(nmatricula, escrituraAsBytes)
}

//Food describes basic details of what makes up a food
/*
type Food struct {
	Farmer  string `json:"farmer"`
	Variety string `json:"variety"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, foodId string, farmer string, variety string) error {

	//Validaciones de sintaxis

	//validaciones de negocio

	food := Food{
		Farmer:  farmer,
		Variety: variety,
	}

	foodAsBytes, err := json.Marshal(food)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(foodId, foodAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, foodId string) (*Food, error) {

	foodAsBytes, err := ctx.GetStub().GetState(foodId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if foodAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", foodId)
	}

	food := new(Food)

	err = json.Unmarshal(foodAsBytes, food)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
	}

	return food, nil
}
*/

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create foodcontrol chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting foodcontrol chaincode: %s", err.Error())
	}
}
