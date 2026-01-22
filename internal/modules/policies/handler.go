package policies

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// Handler handles policies HTTP requests
type Handler struct{}

// GeneratePrivacyPolicyPDF generates privacy policy PDF
func (h *Handler) GeneratePrivacyPolicyPDF(w http.ResponseWriter, r *http.Request) {
	pdfContent := h.generatePrivacyPolicy()
	
	filename := fmt.Sprintf("kebijakan_privasi_%s.pdf", time.Now().Format("20060102"))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfContent)))
	w.Write(pdfContent)
}

// GenerateTermsOfServicePDF generates terms of service PDF
func (h *Handler) GenerateTermsOfServicePDF(w http.ResponseWriter, r *http.Request) {
	pdfContent := h.generateTermsOfService()
	
	filename := fmt.Sprintf("terms_of_service_%s.pdf", time.Now().Format("20060102"))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfContent)))
	w.Write(pdfContent)
}

func (h *Handler) generatePrivacyPolicy() []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	
	// Header
	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(0, 0, 210, 40, "F")
	
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetY(12)
	pdf.CellFormat(0, 10, "KEBIJAKAN PRIVASI", "", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetY(24)
	pdf.CellFormat(0, 6, "Karir Nusantara - Platform Pencarian Kerja Terpercaya Indonesia", "", 1, "C", false, 0, "")
	
	pdf.Ln(15)
	
	// Body
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	
	sections := []struct {
		title   string
		content string
	}{
		{
			"1. PENDAHULUAN",
			"Karir Nusantara (\"kami\", \"platform\") berkomitmen untuk melindungi privasi dan keamanan data pengguna. Kebijakan Privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, menyimpan, dan melindungi informasi Anda.",
		},
		{
			"2. INFORMASI YANG KAMI KUMPULKAN",
			"Kami mengumpulkan informasi berikut dari pengguna platform:\n• Data Perusahaan: Nama, alamat, kontak, informasi industri\n• Data Akun: Email, password (terenkripsi), profil pengguna\n• Data Lowongan: Deskripsi pekerjaan, requirements, salary range\n• Data Calon Karyawan: CV, portfolio, riwayat pendidikan dan pekerjaan\n• Data Transaksi: Invoice, bukti pembayaran, riwayat pembayaran\n• Data Teknis: IP address, browser, device type, aktivitas pengguna",
		},
		{
			"3. KEAMANAN DATA",
			"Semua data disimpan dengan enkripsi menggunakan standar keamanan internasional:\n• Enkripsi SSL/TLS untuk semua transmisi data\n• Database dienkripsi dengan algoritma AES-256\n• Backup data dilakukan secara berkala\n• Audit keamanan dilakukan setiap 6 bulan\n• Akses data dibatasi hanya untuk staff yang berwenang\n• Sistem monitoring 24/7 untuk mendeteksi aktivitas mencurigakan",
		},
		{
			"4. PENGGUNAAN DATA",
			"Data Anda digunakan untuk:\n• Menyediakan layanan perekrutan yang optimal\n• Verifikasi identitas dan pencegahan fraud\n• Mengirimkan notifikasi dan komunikasi terkait layanan\n• Meningkatkan kualitas dan fitur platform\n• Keperluan administrasi dan kepatuhan hukum\n• Analisis statistik dan research (data anonimisasi)",
		},
		{
			"5. PEMBAGIAN DATA DENGAN PIHAK KETIGA",
			"Kami TIDAK membagikan data pribadi Anda kepada pihak ketiga tanpa izin eksplisit, kecuali:\n• Diperlukan untuk keperluan hukum atau order dari pengadilan\n• Payment gateway (PG Aman, Bank) untuk proses pembayaran\n• Email service provider untuk mengirim notifikasi\n• Jika ada merger atau akuisisi perusahaan\n• Partner teknis yang menandatangani data processing agreement",
		},
		{
			"6. HAK DAN KONTROL PENGGUNA",
			"Anda memiliki hak untuk:\n• Mengakses data pribadi Anda kapan saja\n• Meminta koreksi atau update data yang tidak akurat\n• Menghapus akun dan semua data terkait (subject to legal holds)\n• Meminta laporan data yang kami simpan tentang Anda\n• Menolak penggunaan data untuk tujuan tertentu\n• Mengunduh data Anda dalam format yang dapat dibaca mesin",
		},
		{
			"7. RETENSI DATA",
			"Data disimpan selama periode berikut:\n• Data Akun: Selamanya atau sampai akun dihapus\n• Data Lowongan: Sesuai kebutuhan atau sampai lowongan ditutup\n• Data Pembayaran: 7 tahun (sesuai regulasi pajak Indonesia)\n• Data Log/Teknis: 90 hari\n• Data Personal: Sesuai permintaan atau sampai akun dihapus",
		},
		{
			"8. COOKIES DAN TRACKING",
			"Platform menggunakan cookies untuk:\n• Menjaga sesi login Anda tetap aktif\n• Menyimpan preferensi pengguna\n• Analytics dan improvement platform\n\nAnda dapat menonaktifkan cookies melalui pengaturan browser, namun ini dapat mempengaruhi fungsionalitas platform.",
		},
		{
			"9. KEPATUHAN REGULASI",
			"Kami mematuhi peraturan privasi data berikut:\n• Undang-Undang Nomor 8 Tahun 1997 tentang Dokumen Perusahaan\n• Undang-Undang Nomor 36 Tahun 1999 tentang Telekomunikasi\n• Standar ISO 27001 untuk Information Security Management\n• Best practices internasional untuk data protection",
		},
		{
			"10. KONTAK & PENGADUAN",
			"Jika Anda memiliki pertanyaan atau pengaduan tentang privasi data:\nEmail: privacy@karirnusantara.com\nChat Support: Tersedia di dashboard platform\nAlamat: PT Karir Nusantara, Indonesia\n\nWaktu respons: 5-7 hari kerja",
		},
	}
	
	for _, section := range sections {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		
		// Section title
		pdf.SetFont("Arial", "B", 11)
		pdf.SetFillColor(41, 128, 185)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, 8, section.title, "", 1, "L", true, 0, "")
		
		// Section content
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, section.content, "", "L", false)
		pdf.Ln(3)
	}
	
	// Footer
	pdf.SetY(-25)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 5, fmt.Sprintf("Dokumen ini digenerate otomatis pada %s", time.Now().Format("02 January 2006 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Karir Nusantara - Kebijakan Privasi", "", 1, "C", false, 0, "")
	
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating privacy policy PDF: %v", err)
		return []byte{}
	}
	
	return buf.Bytes()
}

func (h *Handler) generateTermsOfService() []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	
	// Header
	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(0, 0, 210, 40, "F")
	
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetY(12)
	pdf.CellFormat(0, 10, "TERMS OF SERVICE", "", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetY(24)
	pdf.CellFormat(0, 6, "Karir Nusantara - Platform Pencarian Kerja Terpercaya Indonesia", "", 1, "C", false, 0, "")
	
	pdf.Ln(15)
	
	// Body
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	
	sections := []struct {
		title   string
		content string
	}{
		{
			"1. PENERIMAAN SYARAT & KETENTUAN",
			"Dengan mengakses dan menggunakan platform Karir Nusantara, Anda setuju untuk mematuhi semua syarat dan ketentuan yang tercantum dalam dokumen ini. Jika Anda tidak setuju, silakan tidak menggunakan platform kami.",
		},
		{
			"2. LAYANAN YANG DISEDIAKAN",
			"Karir Nusantara menyediakan platform digital untuk:\n• Perusahaan: Membuat dan mengelola lowongan pekerjaan\n• Pencari Kerja: Mencari dan melamar lowongan pekerjaan\n• Admin: Verifikasi pembayaran dan moderasi konten\n\nLayanan ini disediakan 'AS IS' tanpa jaminan apa pun.",
		},
		{
			"3. BIAYA & PEMBAYARAN",
			"• Harga Lowongan: Rp 10.000 per lowongan pekerjaan\n• Gratis 1 lowongan untuk perusahaan baru\n• Metode Pembayaran: Transfer Bank Indonesia\n• Verifikasi Pembayaran: 1-2 jam jam kerja\n• Pembatalan: Dapat dilakukan kapan saja tanpa biaya\n• Refund: Hanya jika ada kesalahan pembayaran dari sistem kami\n\nSemua transaksi bersifat FINAL dan tidak dapat dibatalkan setelah dikonfirmasi.",
		},
		{
			"4. KEBIJAKAN LOWONGAN PEKERJAAN",
			"Perusahaan setuju untuk:\n• Membuat lowongan yang jujur dan transparan\n• Menyediakan deskripsi pekerjaan yang jelas dan detail\n• Menampilkan salary range yang kompetitif dan akurat\n• Tidak melakukan diskriminasi dalam proses rekrutmen\n• Tidak memposting konten yang melanggar hukum\n• Mematuhi peraturan ketenagakerjaan Indonesia\n\nKarir Nusantara berhak menolak atau menghapus lowongan yang melanggar kebijakan ini.",
		},
		{
			"5. TANGGUNGJAWAB PENGGUNA",
			"Pengguna platform bertanggung jawab untuk:\n• Menjaga kerahasiaan login dan password\n• Semua aktivitas yang terjadi di akun Anda\n• Konten yang Anda posting atau upload\n• Kepatuhan terhadap hukum dan regulasi yang berlaku\n• Tidak melakukan aktivitas yang merugikan platform atau pengguna lain\n\nAnda setuju untuk tidak:\n• Melakukan hacking atau unauthorized access\n• Membuat bot atau automation yang tidak diizinkan\n• Membagikan data pengguna lain\n• Menggunakan platform untuk spam atau phishing\n• Melakukan aktivitas yang melanggar hak pihak ketiga",
		},
		{
			"6. PEMBATASAN TANGGUNG JAWAB",
			"Karir Nusantara TIDAK bertanggung jawab untuk:\n• Kerugian finansial atau bisnis yang diakibatkan penggunaan platform\n• Kualitas atau keakuratan lowongan pekerjaan\n• Kelakuan atau kredibilitas pengguna platform\n• Downtime atau gangguan teknis platform\n• Kehilangan atau kerusakan data (meski kami berusaha melindungi)\n• Konten pihak ketiga yang link-nya terdapat di platform\n\nTanggung jawab maksimal kami terbatas pada biaya yang Anda bayarkan.",
		},
		{
			"7. PEMUTUSAN LAYANAN",
			"Karir Nusantara dapat menghentikan layanan Anda jika:\n• Ada pelanggaran berulang terhadap syarat & ketentuan ini\n• Ada aktivitas fraud atau mencurigakan\n• Tidak ada aktivitas selama 12 bulan berturut-turut\n• Atas permintaan Anda sendiri\n\nPemutusan layanan dapat berakibat:\n• Penghapusan semua lowongan yang aktif\n• Akun tidak dapat diakses\n• Data pembayaran tetap disimpan sesuai hukum",
		},
		{
			"8. PERUBAHAN SYARAT & KETENTUAN",
			"Karir Nusantara berhak mengubah syarat dan ketentuan ini kapan saja. Perubahan signifikan akan diberitahu melalui email. Penggunaan platform setelah pemberitahuan perubahan berarti Anda menyetujui syarat yang baru.",
		},
		{
			"9. KESUKESAN PLATFORM & DOWNTIME",
			"Kami berkomitmen untuk menjaga uptime 99% namun tidak menjamin zero downtime. Downtime dapat terjadi untuk:\n• Maintenance dan update sistem\n• Improvement fitur platform\n• Backup data\n• Force majeure atau insiden teknis yang tidak terduga\n\nAkami akan memberitahu downtime yang terencana minimal 24 jam sebelumnya.",
		},
		{
			"10. DISPUTE & PENYELESAIAN SENGKETA",
			"Setiap sengketa harus diselesaikan melalui:\n1. Komunikasi langsung dengan support team (5-7 hari kerja)\n2. Mediasi untuk menemukan solusi yang adil\n3. Arbitrase sesuai hukum Indonesia\n\nHukum yang berlaku adalah Hukum Republik Indonesia.",
		},
		{
			"11. KONTAK SUPPORT",
			"Untuk pertanyaan atau komplain:\nEmail: support@karirnusantara.com\nChat Support: Tersedia di dashboard platform\nJam Kerja: Senin - Jumat, 09:00 - 17:00 WIB\nResponse Time: 24 jam kerja",
		},
	}
	
	for _, section := range sections {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		
		// Section title
		pdf.SetFont("Arial", "B", 11)
		pdf.SetFillColor(41, 128, 185)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, 8, section.title, "", 1, "L", true, 0, "")
		
		// Section content
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, section.content, "", "L", false)
		pdf.Ln(3)
	}
	
	// Footer
	pdf.SetY(-25)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 5, fmt.Sprintf("Dokumen ini digenerate otomatis pada %s", time.Now().Format("02 January 2006 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Karir Nusantara - Terms of Service", "", 1, "C", false, 0, "")
	
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating terms of service PDF: %v", err)
		return []byte{}
	}
	
	return buf.Bytes()
}
