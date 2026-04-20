# G1 - OWASP Top 10 Haritalama

## Amac

Orta olcekli bir web projesini OWASP Top 10 kategorilerine gore incelemek ve riskleri uygulama kodu/konfigurasyon mantigi ile eslestirmek.

## Secilen Proje

- Proje: Directus
- Neden: Node.js tabanli, API agirlikli, kimlik dogrulama/yetkilendirme, dosya yukleme ve veritabani erisimi gibi bircok web guvenligi konusunu iceriyor.
- Kanit: `ekran-goruntuleri/01-directus-repo.png` eklenmeli.

## Uyguladigim Adimlar

1. Projenin GitHub deposunu actim.
2. README, dependency listesi, API katmani, auth/yetki sistemi ve security sekmesini inceledim.
3. OWASP Top 10 listesinden 5 kategori sectim.
4. Her kategori icin projenin hangi parcasinda risk olabilecegini ve hangi savunmanin beklendigini yazdim.
5. Bulgulari tablo halinde ozetledim.

## OWASP Haritalama Tablosu

| OWASP Kategorisi | Projedeki Ilgili Alan | Olasilik | Savunma / Beklenen Kontrol |
|---|---|---:|---|
| A01 Broken Access Control | API permission checks, role based access | Yuksek etki | Her endpoint icin server-side authorization kontrolu, default deny |
| A02 Cryptographic Failures | Token/session secrets, password hashing | Orta | Secret'lar env ile gelmeli, guclu hashing, HTTPS zorunlu olmali |
| A03 Injection | Database query builder, filters, search params | Yuksek | Parametreli sorgular, ORM/query builder, input validation |
| A05 Security Misconfiguration | CORS, headers, admin panel settings | Orta | Guvenli default ayarlar, security headers, least privilege config |
| A06 Vulnerable and Outdated Components | npm dependencies | Orta | Dependabot, npm audit, Semgrep/SCA, duzenli update |

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] GitHub Security tabinda hangi uyarilari gordun?
- [DOLDUR] Kaynak kodda anlamakta zorlandigin bolum neydi?
- [DOLDUR] Bir kategori icin riskin gercek mi teorik mi oldugunu nasil ayirdin?

## Sonuc

Directus gibi API agirlikli projelerde en kritik riskler access control, injection ve dependency guvenligidir. Bir acigin var oldugunu iddia etmek icin yalnizca teknoloji kullanimi yeterli degildir; endpoint davranisi, kod akisi ve konfigurasyon birlikte incelenmelidir.

## Ogrendigim 3 Sey

- OWASP Top 10 bir checklist degil, riskleri siniflandirma modelidir.
- Access control hatalari genellikle tek bir dosyada degil, route + service + policy akisi icinde ortaya cikar.
- Dependency uyarisi her zaman exploit anlamina gelmez; versiyon, kullanim sekli ve erisilebilirlik birlikte degerlendirilmelidir.
