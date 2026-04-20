# G2 - SQL Injection Lab

## Amac

Zayif bir uygulamada SQL injection mantigini anlamak ve ayni sorunun parametreli sorgu ile nasil kapatildigini gostermek.

## Kullanilan Platform

- OWASP Juice Shop veya DVWA
- Onerilen komut:

```bash
docker run --rm -p 3000:3000 bkimminich/juice-shop
```

## Uyguladigim Adimlar

1. Juice Shop/DVWA uygulamasini localde baslattim.
2. Login formuna normal kullanici bilgisi ile giris denedim.
3. `saldiri.txt` icindeki payloadlari kontrollu sekilde denedim.
4. Basarili/basarisiz denemeleri ekran goruntusu ile kaydettim.
5. `savunma.js` dosyasinda ayni mantigin parametreli sorgu ile guvenli yazilmis halini ekledim.

## Payload Notlari

Payloadlar `saldiri.txt` dosyasinda tutuluyor. Bunlar sadece local lab ortaminda denenmelidir.

## Savunma

`savunma.js` icinde kullanici girdisi string concat ile SQL'e eklenmez. Sorgu:

```js
db.get("SELECT id, email FROM users WHERE email = ? AND password_hash = ?", [email, passwordHash])
```

seklinde parametreli calisir.

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Hangi payload ise yaradi?
- [DOLDUR] Uygulama hangi hatayi veya davranisi gosterdi?
- [DOLDUR] Parametreli sorgu sonrasi ayni payload neden ise yaramadi?

## Sonuc

SQL injection, kullanici girdisi sorgu metninin bir parcasi haline getirildiginde ortaya cikar. Parametreli sorgular girdiyi veri olarak ele aldigi icin sorgu yapisini bozamaz.

## Ogrendigim 3 Sey

- SQL injection sadece `' OR 1=1 --` payloadindan ibaret degildir; hata mesajlari da ipucu verir.
- Hata mesajlarini kullaniciya gostermek saldirgana veritabani hakkinda bilgi verebilir.
- En temel savunma input blacklist degil, parametreli sorgudur.
