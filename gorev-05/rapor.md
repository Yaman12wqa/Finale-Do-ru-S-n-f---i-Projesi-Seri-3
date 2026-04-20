# G5 - JWT Guvenligi Audit

## Amac

Kendi projedeki JWT kullanimi veya JWT eklenirse uygulanacak guvenlik kriterlerini denetlemek.

## Proje Durumu

SecScan su anda kullanici hesabi veya JWT tabanli login icermiyor. Bu nedenle hard-coded JWT secret, uzun omurlu token veya localStorage token riski mevcut degil. Bu gorevde mevcut durum audit edildi ve JWT eklenecekse uygulanacak kurallar yazildi.

## Kontrol Listesi

| Kriter | SecScan Durumu | Dogru Yaklasim |
|---|---|---|
| Secret en az 256-bit random mi? | JWT yok | `openssl rand -base64 32` ile uretilip env'de tutulmali |
| Secret `.env` icinde mi, hard-coded mi? | Hard-coded secret yok | `.env`, `.gitignore`, production secret manager |
| Token expiration var mi? | JWT yok | Access token en fazla 30 dakika |
| Token nerede saklaniyor? | JWT yok | HttpOnly, Secure, SameSite cookie tercih edilmeli |
| Logout tokeni gecersiz kiliyor mu? | JWT yok | Cookie silinmeli, gerekirse denylist/session version kullanilmali |

## Uyguladigim Adimlar

1. Projede `jwt`, `secret`, `token`, `localStorage` anahtar kelimeleri arandi.
2. `.env.example` dosyalari kontrol edildi.
3. `.gitignore` dosyasinda `.env` dosyalarinin ignore edildigi dogrulandi.
4. JWT eklenirse uygulanacak guvenli varsayimlar tabloya yazildi.

## Komutlar

```bash
rg -n "jwt|secret|token|localStorage|session" .
```

Guvenli secret uretme:

```bash
openssl rand -base64 32
```

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] `rg` komutu hangi dosyalarda token kelimesi buldu?
- [DOLDUR] Bunlardan hangileri gercek secret degildi?
- [DOLDUR] JWT eklenirse hangi endpointlerde auth gerekir?

## Sonuc

SecScan'in mevcut halinde JWT riski yoktur. Auth eklenecekse token localStorage yerine HttpOnly cookie ile saklanmali, secret env'den okunmali ve access token suresi kisa tutulmalidir.

## Ogrendigim 3 Sey

- JWT kullanmamak da bazen dogru tasarim karari olabilir.
- En buyuk JWT hatalari hard-coded secret, uzun expiration ve localStorage saklamadir.
- Logout JWT'de sadece client tarafindan token silmek degil, sunucu tarafinda strateji gerektirir.
