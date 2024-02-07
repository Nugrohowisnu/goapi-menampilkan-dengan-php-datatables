package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Kritiks struct {
	NoID        int    `json:"no-id"`
	NIK         string `json:"nik"`
	Nama        string `json:"nama"`
	Departement string `json:"departement"`
	Kritik      string `json:"kritik"`
	Saran       string `json:"saran"`
}

func main() {
	// Koneksi ke database MySQL
	dsn := "root:@tcp(127.0.0.1:3306)/kritik-saran?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Inisialisasi aplikasi Fiber
	app := fiber.New()

	// Endpoint untuk menampilkan semua kritik
	app.Get("/kritik", getKritikList)
	app.Post("/kritik", createKritik)
	app.Put("/kritik/:no_id", updateKritik)
	app.Delete("/kritik/:no_id", deleteKritik)

	// Jalankan server pada port 8000
	err = app.Listen(":8000")
	if err != nil {
		log.Fatal("Error starting Fiber app:", err)
	}
}

// Handler untuk endpoint GET /kritik
func getKritikList(c *fiber.Ctx) error {
	var kritikList []Kritiks
	if err := db.Find(&kritikList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(kritikList)
}

// Handler untuk endpoint POST /kritik
func createKritik(c *fiber.Ctx) error {
	// Membuat instance dari model Kritiks untuk menyimpan data yang diterima dari request
	var kritik Kritiks

	// Parse JSON request body dan bind ke variabel kritik
	if err := c.BodyParser(&kritik); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Membuat entry baru di database
	if err := db.Create(&kritik).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(kritik)
}

// Handler untuk endpoint PUT /kritik/:id (Edit kritik berdasarkan ID)
func updateKritik(c *fiber.Ctx) error {
	// Mendapatkan ID dari parameter URL
	id := c.Params("no_id")

	// Mencari kritik berdasarkan ID
	var kritik Kritiks
	if err := db.First(&kritik, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kritik not found"})
	}

	// Parse JSON request body dan bind ke variabel kritik
	if err := c.BodyParser(&kritik); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Simpan perubahan ke dalam database
	if err := db.Save(&kritik).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(&kritik)
}

func deleteKritik(c *fiber.Ctx) error {
	// Mendapatkan ID dari parameter URL
	id := c.Params("no_id")

	// Menghapus kritik dari database berdasarkan ID
	if err := db.Delete(&Kritiks{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}


