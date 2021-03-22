package helpers

import (
	"strconv"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

func SendMailVerification(Code int, Name string, Email string) error {
	h := hermes.Hermes{
		// Optional Theme
		Theme: new(hermes.Flat),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Pegibelanja",
			Link: "https://pegibelanja.com/",
			// Optional product logo
			Logo:        "https://pegibelanja.com/static/img/logo_light.png",
			Copyright:   "Copyright © 2020 Pegibelanja. All rights reserved.",
			TroubleText: "Kalau tombol {ACTION} tidak berfungsi, silahkan copy dan paste URL dibawah ini ke web browser kamu.",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: Name,
			Intros: []string{
				"Yay! Kamu telah berhasil mendaftar di Pegibelanja. Terima kasih telah memilih Pegibelanja sebagai tempat belanja kebutuhan pokok Kamu.",
			},
			Greeting: "Dear",
			Actions: []hermes.Action{
				{
					Instructions: "Untuk mengaktifkan akun Pegibelanja Anda, silahkan masukkan kode verifikasi berikut ini di halaman verifikasi akun:",
					InviteCode:   strconv.Itoa(Code),
				},
				// {
				// 	Instructions: "Atau klik tombol di bawah ini.",
				// 	Button: hermes.Button{
				// 		Color: "#e80909", // Optional action button color
				// 		Text:  "Verifikasi Sekarang",
				// 		Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
				// 	},
				// },
			},
			Outros: []string{
				"Dapatkan promo menarik untuk produk-produk berkualitas hanya di Pegibelanja!",
				"",
				"Mengapa Anda menerima email ini?",
				`Pegibelanja akan melakukan verifikasi alamat email untuk setiap email yang diregistrasikan di Pegibelanja.com. Email Kamu tidak dapat digunakan untuk masuk ke akun Pegibelanja Kamu 
				sebelum verifikasi dilakukan. Apabila Kamu tidak merasa melakukan registrasi email di Pegibelanja.com, silahkan abaikan email ini.`,
				"",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress("noreply@pegibelanja.com", "Pegibelanja"))
	m.SetHeader("To", Email)
	m.SetHeader("Subject", "Selamat Bergabung di Pegibelanja.com")
	m.SetBody("text/html", emailBody)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("redfin.id.rapidplex.com", 465, "info@pegibelanja.store", "gmaildgm100100")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendMailTest() error {
	h := hermes.Hermes{
		// Optional Theme
		Theme: new(hermes.Flat),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Pegibelanja",
			Link: "https://pegibelanja.com/",
			// Optional product logo
			Logo:        "https://pegibelanja.com/static/img/logo_light.png",
			Copyright:   "Copyright © 2020 Pegibelanja. All rights reserved.",
			TroubleText: "Kalau tombol {ACTION} tidak berfungsi, silahkan copy dan paste URL dibawah ini ke web browser kamu.",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: "John Doe",
			Intros: []string{
				"Yay! Kamu telah berhasil mendaftar di Pegibelanja. Terima kasih telah memilih Pegibelanja sebagai tempat belanja kebutuhan pokok Kamu.",
			},
			Greeting: "Dear",
			Actions: []hermes.Action{
				{
					Instructions: "Untuk mengaktifkan akun Pegibelanja Anda, silahkan masukkan kode verifikasi berikut ini di halaman verifikasi akun:",
					InviteCode:   "5455",
				},
				{
					Instructions: "Atau klik tombol di bawah ini.",
					Button: hermes.Button{
						Color: "#e80909", // Optional action button color
						Text:  "Verifikasi Sekarang",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Dapatkan promo menarik untuk produk-produk berkualitas hanya di Pegibelanja!",
				"",
				"Mengapa Anda menerima email ini?",
				`Pegibelanja akan melakukan verifikasi alamat email untuk setiap email yang diregistrasikan di Pegibelanja.com. Email Kamu tidak dapat digunakan untuk masuk ke akun Pegibelanja Kamu 
				sebelum verifikasi dilakukan. Apabila Kamu tidak merasa melakukan registrasi email di Pegibelanja.com, silahkan abaikan email ini.`,
				"",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress("noreply@pegibelanja.com", "Pegibelanja"))
	m.SetHeader("To", "m.ardhasy@gmail.com", "naufalikhsan448@gmail.com")
	m.SetHeader("Subject", "Selamat Bergabung di Pegibelanja.com")
	m.SetBody("text/html", emailBody)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("redfin.id.rapidplex.com", 465, "info@pegibelanja.store", "gmaildgm100100")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
