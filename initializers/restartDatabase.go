package initializers

import (
	"log"

	"gorm.io/gorm"
)

// Reinicializa o banco de dados: DropTables e AutoMigrate
func ResetDatabase(db *gorm.DB, models []interface{}) error {
	// 1. Excluir todas as tabelas (na ordem inversa para lidar com chaves estrangeiras)
	// O GORM fornece o Migrator para operações de esquema.
	log.Println("Excluindo todas as tabelas...")

	// O GORM tentará lidar com a ordem correta, mas pode ser mais seguro
	// fazer a exclusão na ordem inversa das dependências.
	// db.Migrator().DropTable(models...) funciona na maioria dos casos.
	if err := db.Migrator().DropTable(models...); err != nil {
		log.Printf("Erro ao excluir tabelas: %v\n", err)
		return err
	}

	log.Println("Tabelas excluídas com sucesso.")

	// 2. Recriar as tabelas (AutoMigrate)
	log.Println("Recriando tabelas...")
	if err := db.AutoMigrate(models...); err != nil {
		log.Printf("Erro ao rodar AutoMigrate: %v\n", err)
		return err
	}
	log.Println("Banco de dados reinicializado com sucesso!")
	return nil
}
