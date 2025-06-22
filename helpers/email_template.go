package helpers

func GetEmailVerificationTemplate() (subject string, body string) {
	subject = "Registrasi PFL berhasil - Aktivasi akun anda"
	body = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>PFL Email Verification</title>
		</head>
		<body style="font-family: Arial, Helvetica, sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;">

			<div style="width: 100%; max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px;">
				
				<!-- PFL Text -->
				<div style="font-size: 32px; color: #00009C; font-weight: bold; text-align: center; margin-bottom: 20px;">
					PFL
				</div>

				<!-- Greeting -->
				<h1 style="font-size: 24px; font-weight: 600; text-align: center; margin-bottom: 20px;">
					Halo, {{user_name}}
				</h1>

				<!-- Thank you message -->
				<p style="font-size: 18px; text-align: center; margin-bottom: 30px;">
					Terima kasih telah melakukan registrasi di Pro Futsal League
				</p>

				<!-- Instructions -->
				<div style="background-color: #f9f9f9; padding: 20px; border-radius: 8px; text-align: center;">
					<h2 style="font-size: 20px; font-weight: bold; margin-bottom: 15px;">Petunjuk Selanjutnya:</h2>
					<p style="font-size: 16px; margin-bottom: 20px;">
						Klik tombol di bawah ini untuk mengonfirmasi alamat email Anda.
					</p>
					<a href="{{link_verification}}" style="display: inline-block; background-color: #0000aa; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: 500; font-size: 16px;">
						Verifikasi Email Sekarang
					</a>
				</div>

				<!-- Disclaimer -->
				<p style="font-size: 14px; text-align: center; margin-top: 30px;">
					Jika Anda tidak merasa melakukan registrasi, silakan abaikan email ini.
				</p>

				<!-- Footer -->
				<div style="background-image: linear-gradient(to right, #00009C, #000022); color: white; text-align: center; padding: 20px;">
					<p>Mempunyai kendala terkait registrasi?</p>
					<p>Silahkan kontak email CS kami di <a href="mailto:cs@profutsaleague" style="color: white;">cs@profutsaleague</a></p>
				</div>

			</div>

		</body>
		</html>
	`
	return
}

func GetEmailTicketPurchaseQRTemplate() (subject string, body string) {
	subject = "Pembelian Tiket Pro Futsal League"
	body = `
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>QR Ticket - PFL</title>
		</head>
		<body style="font-family: Arial, Helvetica, sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;">
			<div style="max-width: 680px; margin: 0 auto; background-color: #ffffff;">
				<div style="margin: 0 auto; padding: 20px; max-width: 624px;">
					<div style="text-align: center; margin-bottom: 20px;">
						<img src="logo-blue.png" alt="PFL Logo" style="height: 96px;">
						<p style="font-size: 20px; font-weight: bold; margin: 10px 0;">Tiket Anda Siap Digunakan</p>
						<p style="font-size: 14px; margin: 0;">Terimakasih telah melakukan pembelian tiket Pro Futsal League</p>
					</div>
					<p style="font-size: 14px; margin-top: 20px;">Tiket Anda terlampir dalam email ini. Jika Anda mengalami masalah saat membuka lampiran, Anda tetap dapat mengakses tiket Anda kapan saja melalui halaman tiket atau dengan mengklik tautan berikut:</p>
					<p style="font-size: 14px; text-align: center; padding: 20px 0;"><a
							href="{{ticket_purchase_url}}"
							style="color: #2b51c0;">{{ticket_purchase_url}}</a></p>

					<div style="background-color: #FAFAFA; padding: 15px; border-radius: 8px;">
						<h4 style="margin-top: 0;">Informasi Penting</h4>
						<ul style="padding-left: 20px; font-size: 14px;">
							<li>Setiap tiket hanya berlaku untuk satu hari pertandingan sesuai tanggal yang tertera</li>
							<li>Pastikan untuk menyimpan tiket digital Anda dengan baik</li>
							<li>Direkomendasikan hadir 30 menit sebelum pertandingan dimulai</li>
							<li>Tiket tidak dapat dipindahkan atau dijual kembali</li>
							<li>Anda dapat menunjukan QR tiket pada panitia untuk ditukar gelang</li>
						</ul>
					</div>
				</div>
				<div
					style="margin-top: 30px; text-align: center; font-size: 13px; background: linear-gradient(to right, #00009B, #000035); color: #fff; padding: 15px;">
					Memunyai kendala terkait pembelian tiket?<br>
					Hubungi kami via email: <a style="color:#fff;" href="mailto:cs@profutsalleague">cs@profutsalleague</a>
				</div>
			</div>
		</body>

		</html>
		<!-- End of One Day Ticket HTML -->
	`

	return subject, body
}
