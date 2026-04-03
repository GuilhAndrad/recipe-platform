package database

import (
	"log"
 
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
 
// Connect abre a conexão com o banco e executa as migrations.
func Connect(databaseURL string, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v", err)
	}
 
	log.Println("Banco de dados conectado com sucesso!")
 
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}
 
	log.Println("Migrations executadas com sucesso!")
 
	return db
}