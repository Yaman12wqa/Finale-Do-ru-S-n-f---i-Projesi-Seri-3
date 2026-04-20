# G10 - Security Headers A+

## Amac

Public web uygulamasinda securityheaders.com uzerinden A+ skoruna ulasmak.

## Uyguladigim Adimlar

1. Frontend icin `next.config.mjs` dosyasina guvenlik headerlari eklendi.
2. Ayrica basit Express `helmet-demo` hazirlandi; bu demo securityheaders.com A+ testi icin daha kolay deploy edilir.
3. Uygulama production modda build edildi.
4. Public deploy URL olusturulduktan sonra securityheaders.com uzerinden tarama yapilacak.
5. A+ sonucunun permalink'i ve ekran goruntusu bu rapora eklenecek.

## Eklenen Headerlar

- `Content-Security-Policy`
- `Strict-Transport-Security`
- `X-Frame-Options`
- `X-Content-Type-Options`
- `Referrer-Policy`
- `Permissions-Policy`
- `Cross-Origin-Opener-Policy`
- `Cross-Origin-Resource-Policy`

## Test Komutlari

Local header kontrolu:

```bash
curl -I http://localhost:3000
```

Helmet demo:

```bash
cd gorev-10/helmet-demo
npm install
npm start
curl -I http://localhost:4010
```

SecurityHeaders testi:

```text
https://securityheaders.com/?q=YOUR_PUBLIC_URL&followRedirects=on
```

## Kanit

- SecurityHeaders permalink: `[DOLDUR]`
- Ekran goruntusu: `ekran-goruntuleri/01-securityheaders-a-plus.png`

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Ilk taramada hangi header eksikti?
- [DOLDUR] CSP hangi resource'u engelledi?
- [DOLDUR] A+ almak icin hangi header degeri degisti?

## Sonuc

Security headers, tarayici tarafinda clickjacking, MIME sniffing, gereksiz izinler ve bazi XSS etkilerini azaltir. A+ skor guzel bir hedef olsa da gercek guvenlik icin backend kontrolleri, input validation ve auth kontrolleri de gerekir.

## Ogrendigim 3 Sey

- HSTS sadece HTTPS public domain uzerinde anlamlidir.
- CSP guclu bir savunmadir ama hatali yazilirsa uygulama kaynaklarini bozabilir.
- Security header skoru uygulamanin tamamini guvenli yapmaz, sadece bir katmani olcer.
