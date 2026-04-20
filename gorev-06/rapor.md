# G6 - OAuth 2.0 + PKCE Demo

## Amac

Google OAuth Authorization Code + PKCE akisini uygulamak ve `code_verifier` / `code_challenge` mantigini anlamak.

## Uyguladigim Adimlar

1. Google Cloud Console uzerinde OAuth Client olusturulacak.
2. Redirect URI olarak `http://localhost:4006/callback` eklenecek.
3. Client ID ve secret `.env` veya terminal environment variable olarak verilecek.
4. `pkce-demo` uygulamasi calistirilacak.
5. `/login` ile Google'a yonlendirme, `/callback` ile token exchange test edilecek.

## Calistirma

PowerShell:

```powershell
$env:GOOGLE_CLIENT_ID="your-client-id"
$env:GOOGLE_CLIENT_SECRET="your-client-secret"
$env:OAUTH_REDIRECT_URI="http://localhost:4006/callback"
node gorev-06/pkce-demo/server.js
```

Tarayici:

```text
http://localhost:4006/login
```

## Guvenlik Notlari

- Client secret repoya commit edilmez.
- `code_verifier` random uretilir.
- `code_challenge` SHA-256 ve base64url ile uretilir.
- `state` parametresi CSRF benzeri OAuth saldirilarini azaltir.

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Google Console redirect URI hatasi aldin mi?
- [DOLDUR] `code_verifier` ve `code_challenge` nerede goruldu?
- [DOLDUR] Callback sonrasi hangi profil bilgisi ekranda basildi?

## Sonuc

PKCE, authorization code ele gecse bile token exchange icin ek bir kanit ister. Bu ozellikle public client ve SPA/native app senaryolarinda onemli bir savunmadir.

## Ogrendigim 3 Sey

- OAuth login sadece redirect degil, state ve PKCE dogrulama akisi da gerektirir.
- Client secret `.env` veya secret manager disinda tutulmamalidir.
- Redirect URI birebir eslesmediginde Google token akisini reddeder.
