package main

import (
	"database/sql" //Veri tabanı için
	"fmt"          // yazı yazdırmak için
	"log"          // hata mesajlarını bastırmak için

	_ "github.com/lib/pq"
)

// Veritabanı bağlantı bilgilerim
// Const Sabit!!
const (
	sunucu     = "localhost"
	port       = 5432
	kullanici  = "postgres"
	sifre      = "12345" // Şifrem
	veritabani = "mydb"  // Veri Tabanımın Adı
)

//%s	String
//%d	int
//%f	Float
//%v	Değişkenin türüne göre

func main() {
	// Bağlantı cümlesi oluşturdum
	baglantiBilgisi := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		sunucu, port, kullanici, sifre, veritabani)

	// Veritabanına bağlan
	db, err := sql.Open("postgres", baglantiBilgisi) //Sql bağlantısı açtık
	if err != nil {
		log.Fatal("Bağlantı açma hatası:", err)
	}
	defer db.Close()

	// db.ping() ile Bağlantı testi
	err = db.Ping()
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}
	fmt.Println("PostgreSQL veritabanına başarıyla bağlanıldı!")

	//----------------------------------//
	//CREATE TABLE: Yeni tablo oluştur

	//IF NOT EXISTS: Eğer daha önce yoksa

	//id SERIAL PRIMARY KEY: Otomatik artan benzersiz ID

	//isim TEXT NOT NULL: Boş olmayan metin alanı
	//----------------------------------//

	// Tablo oluşturdum
	tabloSorgusu := `
    CREATE TABLE IF NOT EXISTS urunler (
        id SERIAL PRIMARY KEY,
        isim TEXT NOT NULL
    );`
	_, err = db.Exec(tabloSorgusu)
	if err != nil {
		log.Fatal("Tablo oluşturma hatası:", err)
	}
	fmt.Println("Tablo başarıyla oluşturuldu veya zaten var.")

	// Veri ekleme kısmı
	ekleSorgusu := `INSERT INTO urunler (isim) VALUES ($1) RETURNING id`
	var yeniID int

	err = db.QueryRow(ekleSorgusu, "Muz").Scan(&yeniID)

	if err != nil {
		log.Fatal("Veri ekleme hatası:", err)
	}
	fmt.Printf("Yeni ürün eklendi. ID: %d\n", yeniID)

	// Verileri listeleme kısmı
	satirlar, err := db.Query("SELECT id, isim FROM urunler")
	if err != nil {
		log.Fatal("Veri çekme hatası:", err)
	}
	defer satirlar.Close()

	fmt.Println("Ürünler tablosundaki veriler:")
	for satirlar.Next() {
		//Her satırdan Çekilen Veriyi Tutmak İçin iki değişken tanımlanır
		var id int
		var isim string

		err = satirlar.Scan(&id, &isim)

		if err != nil {
			log.Fatal("Satır okuma hatası:", err)
		}
		fmt.Printf("%d: %s\n", id, isim)
	}

	// Ürün güncelleme
	guncelleSorgusu := `UPDATE urunler SET isim = $1 WHERE id = $2`
	sonuc, err := db.Exec(guncelleSorgusu, "Güncellenmiş Ürün")
	if err != nil {
		log.Fatal("Veri güncelleme hatası:", err)
	}
	degisenSatirSayisi, err := sonuc.RowsAffected()
	if err != nil {
		log.Fatal("Güncellenen satır sayısı alınamadı:", err)
	}
	fmt.Printf("%d satır güncellendi.\n", degisenSatirSayisi)

	// Ürün silme
	silSorgusu := `DELETE FROM urunler WHERE id = $1`
	sonuc, err = db.Exec(silSorgusu, 1)
	if err != nil {
		log.Fatal("Veri silme hatası:", err)
	}
	silinenSatirSayisi, err := sonuc.RowsAffected()
	if err != nil {
		log.Fatal("Silinen satır sayısı alınamadı:", err)
	}
	fmt.Printf("%d satır silindi.\n", silinenSatirSayisi)

}
