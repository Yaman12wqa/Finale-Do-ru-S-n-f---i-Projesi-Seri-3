# G3 - XSS + CSP Korumasi

## Amac

Reflected XSS davranisini gormek, sonra output encoding ve Content-Security-Policy ile riski azaltmak.

## Uyguladigim Adimlar

1. `express-csp-demo` klasorunde Express uygulamasini hazirladim.
2. `/vulnerable?name=...` endpointinde kullanici girdisini HTML icine ham bastim.
3. XSS payload ile alert calisip calismadigini test ettim.
4. `/safe?name=...` endpointinde HTML escaping uyguladim.
5. CSP header ekledim: `default-src 'self'; script-src 'self'; object-src 'none'; base-uri 'self'`.
6. Tarayici DevTools Console ve Network sekmelerinden sonucu kontrol ettim.

## Calistirma

```bash
cd gorev-03/express-csp-demo
npm install
npm start
```

Test:

```text
http://localhost:4003/vulnerable?name=<script>alert(1)</script>
http://localhost:4003/safe?name=<script>alert(1)</script>
```

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Ilk endpointte payload nasil davrandi?
- [DOLDUR] CSP console'da hangi mesaji verdi?
- [DOLDUR] Safe endpointte HTML nasil encode edildi?

## Sonuc

Reflected XSS, kullanici girdisi HTML/JS contextine guvenli encode edilmeden yazildiginda olusur. CSP tek basina yeterli degildir; asil savunma dogru output encoding ve framework guvenli varsayimlaridir.

## Ogrendigim 3 Sey

- XSS korumasi icin context-aware output encoding gereklidir.
- CSP ikinci savunma katmanidir, hatali rendering'i tamamen telafi etmez.
- DevTools Console CSP ihlallerini anlamak icin pratik bir debug aracidir.
