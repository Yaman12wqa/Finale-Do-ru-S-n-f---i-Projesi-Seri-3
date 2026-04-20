# G7 - Nmap + ZAP Taramasi

## Amac

Kendi local uygulamamizi Nmap ve OWASP ZAP ile tarayip bulgulari analiz etmek.

## Hedef

- SecScan frontend: `http://localhost:3000`
- SecScan backend: `http://localhost:8080`

## Komutlar

```bash
nmap -sV -p 3000,8080 localhost
```

ZAP baseline:

```bash
docker run --rm -t -v "$(pwd)/gorev-07:/zap/wrk" ghcr.io/zaproxy/zaproxy:stable zap-baseline.py -t http://host.docker.internal:3000 -r zap-report.html
```

## Uyguladigim Adimlar

1. SecScan uygulamasini `docker compose up -d` ile baslattim.
2. Nmap ile 3000 ve 8080 portlarini taradim.
3. ZAP baseline ile frontend uzerinden pasif tarama yaptim.
4. HTML rapordaki medium/high bulgulari listeledim.
5. Bulgularin OWASP kategorisini not ettim.

## Bulgular

| Arac | Bulgu | Seviye | OWASP Kategorisi | Not |
|---|---|---|---|---|
| Nmap | [DOLDUR] | [DOLDUR] | [DOLDUR] | [DOLDUR] |
| ZAP | [DOLDUR] | [DOLDUR] | [DOLDUR] | [DOLDUR] |

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Docker icinden localhost yerine neden `host.docker.internal` kullandim?
- [DOLDUR] ZAP raporu hangi dosyada olustu?
- [DOLDUR] Hangi bulgu false positive olabilir?

## Sonuc

Nmap servis kesfi icin, ZAP ise HTTP guvenlik basliklari ve pasif web bulgulari icin faydalidir. ZAP bulgulari otomatik olarak kesin acik kabul edilmemeli, uygulama davranisi ile dogrulanmalidir.

## Ogrendigim 3 Sey

- `localhost` Docker container icinde host makineyi gostermeyebilir.
- ZAP baseline pasif ve daha guvenli bir tarama modudur.
- Tarama sonucu kadar false positive analizi de raporun parcasidir.
