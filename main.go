package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"map/animalsandpictures/redisclient"
	"net/http"

	"github.com/gin-gonic/gin"
)

const namespace string = "animals"

var client = redisclient.Default()

type Animal struct {
	Id       string
	Race     string
	Birthday string
	Name     string
}

func main() {
	// Gin routing setup
	router := gin.Default()

	// Define the API routes
	router.GET("/animals/:id", getAnimal)
	router.POST("/animals", createAnimal)
	router.PUT("/animals/:id", updateAnimal)
	router.DELETE("/animals/:id", deleteAnimal)

	// Pictures
	router.GET("/animalpicture/:id", getPicture)
	router.POST("/animalpicture", postPicture)

	// Start the HTTP server on port 8090
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start the server: ", err)
	}
	fmt.Println("Go Animals & Pictures CRUD API initialized successfully..")
}

func set(animal *Animal) {
	fmt.Printf("Persisting id: %s, name: %s, race: %s, birthday: %s \n", animal.Id, animal.Name, animal.Race, animal.Birthday)

	json, err := json.Marshal(animal)
	if err != nil {
		panic(err)
	}

	err = client.SetAnimal(&animal.Id, string(json))
	if err != nil {
		panic(err)
	}
}

// Handler to retrieve a specific animal
func getAnimal(c *gin.Context) {
	id := c.Param("id")
	animal, animalPresent := client.GetAnimal(&id)
	if !animalPresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{id: json.RawMessage(animal)})
}

// Handler to create a new animal
func createAnimal(c *gin.Context) {
	var animal Animal

	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	set(&animal)
	c.JSON(http.StatusCreated, animal)
}

// Handler to update an existing animal
func updateAnimal(c *gin.Context) {
	id := c.Param("id")
	if !client.IsAnimalPresent(&id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}

	var animal Animal
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	set(&animal)
	c.JSON(http.StatusOK, animal)
}

// Handler to delete an animal
func deleteAnimal(c *gin.Context) {
	id := c.Param("id")
	if animalDeleted := client.DeleteAnimal(&id); !animalDeleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}
	c.JSON(http.StatusOK, id)
}

// Handle uploading a picture as form file with key "img"
func postPicture(c *gin.Context) {
	file, header, err := c.Request.FormFile("img")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}

	err = client.SetPicture(&header.Filename, (base64.StdEncoding.EncodeToString(buf.Bytes())))
	if err != nil {
		panic(err)
	}
	c.Status(http.StatusCreated)
}

// Handle retrieving a picture
func getPicture(c *gin.Context) {
	id := c.Param("id")
	img, imagePresent := client.GetPicture(&id)

	if !imagePresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "Picture not found"})
		return
	}

	unbased, _ := base64.StdEncoding.DecodeString(img)
	res := bytes.NewReader(unbased)
	c.Header("Content-Type", "image/jpeg")
	_, err := io.Copy(c.Writer, res)
	if err != nil {
		panic(err)
	}
}
