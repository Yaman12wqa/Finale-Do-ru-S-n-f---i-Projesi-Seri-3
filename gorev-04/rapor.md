# G4 - CSRF Token Sistemi

## Amac

Express uzerinde CSRF token mantigini anlamak ve token olmayan POST istegini engellemek.

## Uyguladigim Adimlar

1. `csrf-demo` klasorunde Express uygulamasini hazirladim.
2. Session cookie ve CSRF token ureten middleware ekledim.
3. `/transfer` formuna hidden CSRF input koydum.
4. Token yokken veya yanlisken POST isteginin `403 Forbidden` dondugunu test ettim.
5. Token dogruyken transfer simulasyonunun basarili oldugunu gordum.

## Calistirma

```bash
cd gorev-04/csrf-demo
npm install
npm start
```

Tarayici:

```text
http://localhost:4004/transfer
```

Token olmadan test:

```bash
curl -i -X POST http://localhost:4004/transfer -d "to=alice&amount=100"
```

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Token olmadan gelen yanit neydi?
- [DOLDUR] Cookie ve hidden input nasil eslesti?
- [DOLDUR] SameSite cookie ayarinin etkisi ne oldu?

## Sonuc

CSRF savunmasi icin state-changing POST isteklerinde tahmin edilemeyen token kontrolu ve SameSite cookie ayari birlikte kullanilmalidir.

## Ogrendigim 3 Sey

- CSRF, kullanicinin oturum cookie'sinin otomatik gonderilmesini suistimal eder.
- Token sunucu tarafinda oturumla iliskili dogrulanmalidir.
- SameSite cookie faydalidir ama token kontrolunun yerine gecmez.
